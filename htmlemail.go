// Package htmlemail provides utilities for generating dynamic HTML email content from templates.
package htmlemail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// EmailBuilder provides a fluent interface for building HTML emails.
type EmailBuilder struct {
	htmlContent string
	cssStyles   []string
	inlineCSS   bool
	data        map[string]interface{}
	errors      []error
}

// Template represents an HTML email template with dynamic content support.
type Template struct {
	content      string
	filePath     string
	variables    map[string]interface{}
	conditionals []ConditionalBlock
	loops        []LoopBlock
}

// ConditionalBlock represents conditional content in templates.
type ConditionalBlock struct {
	Condition   string
	Content     string
	ElseContent string
}

// LoopBlock represents repeating content in templates.
type LoopBlock struct {
	Variable   string
	Collection string
	Template   string
}

// PlaceholderStyle defines how placeholders are formatted in templates.
type PlaceholderStyle int

const (
	// DollarStyle uses $variable$ format (like your original code)
	DollarStyle PlaceholderStyle = iota
	// BraceStyle uses {{variable}} format (Go template style)
	BraceStyle
	// PercentStyle uses %variable% format
	PercentStyle
)

// NewEmailBuilder creates a new EmailBuilder instance.
func NewEmailBuilder() *EmailBuilder {
	return &EmailBuilder{
		data: make(map[string]interface{}),
	}
}

// LoadTemplate loads an HTML template from a file.
func LoadTemplate(filePath string) (*Template, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file %s: %w", filePath, err)
	}

	return &Template{
		content:   string(content),
		filePath:  filePath,
		variables: make(map[string]interface{}),
	}, nil
}

// LoadTemplateFromString creates a template from a string.
func LoadTemplateFromString(content string) *Template {
	return &Template{
		content:   content,
		variables: make(map[string]interface{}),
	}
}

// SetVariable sets a variable for template replacement.
func (t *Template) SetVariable(key string, value interface{}) *Template {
	if t.variables == nil {
		t.variables = make(map[string]interface{})
	}
	t.variables[key] = value
	return t
}

// SetVariables sets multiple variables at once.
func (t *Template) SetVariables(vars map[string]interface{}) *Template {
	if t.variables == nil {
		t.variables = make(map[string]interface{})
	}
	for k, v := range vars {
		t.variables[k] = v
	}
	return t
}

// SetStruct sets variables from a struct (similar to your BookingEmailStruct approach).
func (t *Template) SetStruct(data interface{}) *Template {
	if t.variables == nil {
		t.variables = make(map[string]interface{})
	}

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return t
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typ.Field(i)
		value := v.Field(i)

		// Convert field name to snake_case for template placeholders
		key := camelToSnake(field.Name)
		t.variables[key] = value.Interface()
	}

	return t
}

// Render renders the template with the provided variables using dollar-style placeholders.
func (t *Template) Render() (string, error) {
	return t.RenderWithStyle(DollarStyle)
}

// RenderWithStyle renders the template with a specific placeholder style.
func (t *Template) RenderWithStyle(style PlaceholderStyle) (string, error) {
	content := t.content

	for key, value := range t.variables {
		placeholder := formatPlaceholder(key, style)
		stringValue := convertToString(value)
		content = strings.ReplaceAll(content, placeholder, stringValue)
	}

	// Check for unresolved placeholders
	unresolved := findUnresolvedPlaceholders(content, style)
	if len(unresolved) > 0 {
		return content, fmt.Errorf("unresolved placeholders found: %v", unresolved)
	}

	return content, nil
}

// RenderPartial renders only specific variables, leaving others untouched.
func (t *Template) RenderPartial(keys []string) (string, error) {
	content := t.content

	for _, key := range keys {
		if value, exists := t.variables[key]; exists {
			placeholder := "$" + key + "$"
			stringValue := convertToString(value)
			content = strings.ReplaceAll(content, placeholder, stringValue)
		}
	}

	return content, nil
}

// RenderWithGoTemplate uses Go's template engine for more advanced features.
func (t *Template) RenderWithGoTemplate() (string, error) {
	tmpl, err := template.New("email").Parse(t.content)
	if err != nil {
		return "", fmt.Errorf("failed to parse Go template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, t.variables)
	if err != nil {
		return "", fmt.Errorf("failed to execute Go template: %w", err)
	}

	return buf.String(), nil
}

// AddConditional adds conditional content to the template.
func (t *Template) AddConditional(condition, content, elseContent string) *Template {
	t.conditionals = append(t.conditionals, ConditionalBlock{
		Condition:   condition,
		Content:     content,
		ElseContent: elseContent,
	})
	return t
}

// AddLoop adds repeating content to the template.
func (t *Template) AddLoop(variable, collection, template string) *Template {
	t.loops = append(t.loops, LoopBlock{
		Variable:   variable,
		Collection: collection,
		Template:   template,
	})
	return t
}

// ProcessConditionals processes conditional blocks in the template.
func (t *Template) ProcessConditionals() (string, error) {
	content := t.content

	for _, conditional := range t.conditionals {
		conditionResult, err := evaluateCondition(conditional.Condition, t.variables)
		if err != nil {
			return "", fmt.Errorf("failed to evaluate condition '%s': %w", conditional.Condition, err)
		}

		replacement := conditional.ElseContent
		if conditionResult {
			replacement = conditional.Content
		}

		// Replace the conditional block with the appropriate content
		conditionPlaceholder := fmt.Sprintf("{{#if %s}}", conditional.Condition)
		content = strings.ReplaceAll(content, conditionPlaceholder, replacement)
	}

	return content, nil
}

// EmailBuilder methods

// SetHTML sets the base HTML content for the email.
func (eb *EmailBuilder) SetHTML(html string) *EmailBuilder {
	eb.htmlContent = html
	return eb
}

// LoadHTMLFromFile loads HTML content from a file.
func (eb *EmailBuilder) LoadHTMLFromFile(filePath string) *EmailBuilder {
	content, err := os.ReadFile(filePath)
	if err != nil {
		eb.errors = append(eb.errors, fmt.Errorf("failed to load HTML file %s: %w", filePath, err))
		return eb
	}
	eb.htmlContent = string(content)
	return eb
}

// AddCSS adds CSS styles to the email.
func (eb *EmailBuilder) AddCSS(css string) *EmailBuilder {
	eb.cssStyles = append(eb.cssStyles, css)
	return eb
}

// AddCSSFromFile loads CSS from a file.
func (eb *EmailBuilder) AddCSSFromFile(filePath string) *EmailBuilder {
	content, err := os.ReadFile(filePath)
	if err != nil {
		eb.errors = append(eb.errors, fmt.Errorf("failed to load CSS file %s: %w", filePath, err))
		return eb
	}
	eb.cssStyles = append(eb.cssStyles, string(content))
	return eb
}

// SetInlineCSS enables or disables CSS inlining.
func (eb *EmailBuilder) SetInlineCSS(inline bool) *EmailBuilder {
	eb.inlineCSS = inline
	return eb
}

// SetData sets template data.
func (eb *EmailBuilder) SetData(key string, value interface{}) *EmailBuilder {
	eb.data[key] = value
	return eb
}

// SetDataFromStruct sets data from a struct.
func (eb *EmailBuilder) SetDataFromStruct(data interface{}) *EmailBuilder {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return eb
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typ.Field(i)
		value := v.Field(i)
		key := camelToSnake(field.Name)
		eb.data[key] = value.Interface()
	}

	return eb
}

// Build builds the final HTML email.
func (eb *EmailBuilder) Build() (string, error) {
	if len(eb.errors) > 0 {
		return "", fmt.Errorf("builder has errors: %v", eb.errors)
	}

	html := eb.htmlContent

	// Replace placeholders
	for key, value := range eb.data {
		placeholder := "$" + key + "$"
		stringValue := convertToString(value)
		html = strings.ReplaceAll(html, placeholder, stringValue)
	}

	// Add CSS
	if len(eb.cssStyles) > 0 {
		cssContent := strings.Join(eb.cssStyles, "\n")
		if eb.inlineCSS {
			// Basic CSS inlining (you could use a more sophisticated library)
			html = inlineCSS(html, cssContent)
		} else {
			// Add CSS to head
			html = addCSSToHead(html, cssContent)
		}
	}

	return html, nil
}

// Utility functions

// formatPlaceholder formats a placeholder according to the style.
func formatPlaceholder(key string, style PlaceholderStyle) string {
	switch style {
	case DollarStyle:
		return "$" + key + "$"
	case BraceStyle:
		return "{{" + key + "}}"
	case PercentStyle:
		return "%" + key + "%"
	default:
		return "$" + key + "$"
	}
}

// convertToString converts any value to its string representation.
func convertToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%.2f", v)
	case bool:
		return strconv.FormatBool(v)
	case nil:
		return ""
	default:
		// Try JSON encoding for complex types
		if jsonBytes, err := json.Marshal(v); err == nil {
			return string(jsonBytes)
		}
		return fmt.Sprintf("%v", v)
	}
}

// camelToSnake converts CamelCase to snake_case.
func camelToSnake(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

// findUnresolvedPlaceholders finds placeholders that haven't been resolved.
func findUnresolvedPlaceholders(content string, style PlaceholderStyle) []string {
	var pattern string
	switch style {
	case DollarStyle:
		pattern = `\$([a-zA-Z_][a-zA-Z0-9_]*)\$`
	case BraceStyle:
		pattern = `\{\{([a-zA-Z_][a-zA-Z0-9_]*)\}\}`
	case PercentStyle:
		pattern = `%([a-zA-Z_][a-zA-Z0-9_]*)%`
	default:
		pattern = `\$([a-zA-Z_][a-zA-Z0-9_]*)\$`
	}

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(content, -1)

	var unresolved []string
	for _, match := range matches {
		if len(match) > 1 {
			unresolved = append(unresolved, match[1])
		}
	}

	return unresolved
}

// evaluateCondition evaluates a simple condition.
func evaluateCondition(condition string, variables map[string]interface{}) (bool, error) {
	// Simple condition evaluation (can be extended)
	parts := strings.Fields(condition)
	if len(parts) < 3 {
		return false, fmt.Errorf("invalid condition format: %s", condition)
	}

	variable := parts[0]
	operator := parts[1]
	expected := parts[2]

	value, exists := variables[variable]
	if !exists {
		return false, nil
	}

	switch operator {
	case "==":
		return convertToString(value) == expected, nil
	case "!=":
		return convertToString(value) != expected, nil
	case ">":
		if num, err := strconv.ParseFloat(convertToString(value), 64); err == nil {
			if expectedNum, err := strconv.ParseFloat(expected, 64); err == nil {
				return num > expectedNum, nil
			}
		}
		return false, fmt.Errorf("cannot compare non-numeric values with >")
	default:
		return false, fmt.Errorf("unsupported operator: %s", operator)
	}
}

// inlineCSS performs basic CSS inlining (simplified version).
func inlineCSS(html, css string) string {
	// This is a basic implementation. For production, consider using a proper CSS inlining library
	if !strings.Contains(html, "<style>") {
		// Add CSS to head
		html = addCSSToHead(html, css)
	}
	return html
}

// addCSSToHead adds CSS to the HTML head section.
func addCSSToHead(html, css string) string {
	style := fmt.Sprintf("<style type=\"text/css\">\n%s\n</style>", css)

	if strings.Contains(html, "</head>") {
		return strings.Replace(html, "</head>", style+"\n</head>", 1)
	} else if strings.Contains(html, "<body>") {
		return strings.Replace(html, "<body>", style+"\n<body>", 1)
	} else {
		return style + "\n" + html
	}
}

// Legacy functions for backward compatibility (similar to your original approach)

// ReadHTMLFile reads an HTML file and returns its content as a string.
// This is equivalent to your original ReadHTMLFile function.
func ReadHTMLFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read HTML file %s: %w", filePath, err)
	}

	// Convert to string and normalize quotes
	htmlString := string(content)
	htmlString = strings.ReplaceAll(htmlString, "'", "\"")

	return htmlString, nil
}

// InsertHTMLInformation replaces placeholders in HTML content.
// This is equivalent to your original InsertHTMLInformation function but enhanced.
func InsertHTMLInformation(html, target, replace string) (string, error) {
	placeholder := "$" + target + "$"
	if strings.Contains(html, placeholder) {
		return strings.ReplaceAll(html, placeholder, replace), nil
	}
	return html, fmt.Errorf("placeholder $%s$ not found in HTML content", target)
}

// BatchInsertHTMLInformation replaces multiple placeholders at once.
func BatchInsertHTMLInformation(html string, replacements map[string]string) (string, error) {
	var errors []string
	result := html

	for target, replace := range replacements {
		placeholder := "$" + target + "$"
		if strings.Contains(result, placeholder) {
			result = strings.ReplaceAll(result, placeholder, replace)
		} else {
			errors = append(errors, fmt.Sprintf("placeholder $%s$ not found", target))
		}
	}

	if len(errors) > 0 {
		return result, fmt.Errorf("some placeholders not found: %s", strings.Join(errors, ", "))
	}

	return result, nil
}

// CreateEmailTemplate creates a reusable email template similar to your booking function.
type EmailTemplate struct {
	templatePath string
	htmlContent  string
	placeholders map[string]string
}

// NewEmailTemplate creates a new email template from a file.
func NewEmailTemplate(templatePath string) (*EmailTemplate, error) {
	content, err := ReadHTMLFile(templatePath)
	if err != nil {
		return nil, err
	}

	return &EmailTemplate{
		templatePath: templatePath,
		htmlContent:  content,
		placeholders: make(map[string]string),
	}, nil
}

// NewEmailTemplateFromString creates a new email template from a string.
func NewEmailTemplateFromString(content string) *EmailTemplate {
	return &EmailTemplate{
		htmlContent:  content,
		placeholders: make(map[string]string),
	}
}

// SetPlaceholder sets a single placeholder value.
func (et *EmailTemplate) SetPlaceholder(key, value string) *EmailTemplate {
	et.placeholders[key] = value
	return et
}

// SetPlaceholders sets multiple placeholder values.
func (et *EmailTemplate) SetPlaceholders(placeholders map[string]string) *EmailTemplate {
	for k, v := range placeholders {
		et.placeholders[k] = v
	}
	return et
}

// SetStructData sets placeholders from a struct (like your BookingEmailStruct).
func (et *EmailTemplate) SetStructData(data interface{}) *EmailTemplate {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return et
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typ.Field(i)
		value := v.Field(i)

		// Use the field name as-is or convert to snake_case
		key := camelToSnake(field.Name)
		et.placeholders[key] = convertToString(value.Interface())
	}

	return et
}

// Render renders the template with all set placeholders.
func (et *EmailTemplate) Render() (string, error) {
	return BatchInsertHTMLInformation(et.htmlContent, et.placeholders)
}

// RenderSafe renders the template and ignores missing placeholders.
func (et *EmailTemplate) RenderSafe() string {
	result := et.htmlContent
	for key, value := range et.placeholders {
		placeholder := "$" + key + "$"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// GetUnresolvedPlaceholders returns a list of placeholders that haven't been set.
func (et *EmailTemplate) GetUnresolvedPlaceholders() []string {
	pattern := `\$([a-zA-Z_][a-zA-Z0-9_]*)\$`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(et.htmlContent, -1)

	var unresolved []string
	for _, match := range matches {
		if len(match) > 1 {
			placeholder := match[1]
			if _, exists := et.placeholders[placeholder]; !exists {
				unresolved = append(unresolved, placeholder)
			}
		}
	}

	return unresolved
}

// ToFloat64 converts various types to float64.
func ToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

// FormatMoneyFromFloat formats a float64 as a currency string.
func FormatMoneyFromFloat(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

// TableColumn represents a table column configuration.
type TableColumn struct {
	Key    string // The key to look for in row data
	Header string // The header text to display
	Style  string // CSS styles for this column
}

// TableOptions defines options for HTML table generation.
type TableOptions struct {
	TableStyle  string        // CSS styles for the table element
	HeaderStyle string        // CSS styles for header cells
	CellStyle   string        // CSS styles for data cells
	Columns     []TableColumn // Column definitions
	ShowTotal   bool          // Whether to show a total row
	TotalColumn string        // Which column to calculate total for
}

// BuildHTMLTable creates an HTML table from data with configurable options.
func BuildHTMLTable(data []map[string]interface{}, options TableOptions) string {
	var html strings.Builder
	var totalAmount float64

	// Default table style if none provided
	tableStyle := options.TableStyle
	if tableStyle == "" {
		tableStyle = `cellpadding="0" cellspacing="0" class="es-table" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;width:100%" role="presentation"`
	}

	html.WriteString(fmt.Sprintf(`<table %s>`+"\n", tableStyle))

	// Header row
	if len(options.Columns) > 0 {
		html.WriteString(`<tr>` + "\n")
		for _, col := range options.Columns {
			headerStyle := options.HeaderStyle
			if headerStyle == "" {
				headerStyle = `style="padding:0;Margin:0;border-width:1px;border-style:solid;border-color:#cccccc"`
			}
			html.WriteString(fmt.Sprintf(`<td %s>%s</td>`+"\n", headerStyle, col.Header))
		}
		html.WriteString(`</tr>` + "\n")
	}

	// Data rows
	for _, row := range data {
		html.WriteString(`<tr>` + "\n")
		for _, col := range options.Columns {
			cellStyle := col.Style
			if cellStyle == "" {
				cellStyle = options.CellStyle
				if cellStyle == "" {
					cellStyle = `style="padding:0;Margin:0;border-width:1px;border-style:solid;border-color:#cccccc"`
				}
			}

			value := ""
			if val, ok := row[col.Key]; ok && val != nil {
				// Check if this is the total column and it's a number
				if options.ShowTotal && col.Key == options.TotalColumn {
					if f, err := ToFloat64(val); err == nil {
						totalAmount += f
						value = FormatMoneyFromFloat(f)
					} else {
						value = fmt.Sprintf("%v", val)
					}
				} else {
					value = fmt.Sprintf("%v", val)
				}
			}

			html.WriteString(fmt.Sprintf(`<td %s>%s</td>`+"\n", cellStyle, value))
		}
		html.WriteString(`</tr>` + "\n")
	}

	// Total row
	if options.ShowTotal && options.TotalColumn != "" {
		html.WriteString(`<tr>` + "\n")
		for i, col := range options.Columns {
			cellStyle := col.Style
			if cellStyle == "" {
				cellStyle = options.CellStyle
				if cellStyle == "" {
					cellStyle = `style="padding:0;Margin:0;border-width:1px;border-style:solid;border-color:#cccccc"`
				}
			}

			if i == 0 {
				// First column shows "Total"
				html.WriteString(fmt.Sprintf(`<td %s; font-weight:bold;">Total</td>`+"\n", cellStyle))
			} else if col.Key == options.TotalColumn {
				// Total column shows the calculated total
				html.WriteString(fmt.Sprintf(`<td %s; font-weight:bold;">%s</td>`+"\n", cellStyle, FormatMoneyFromFloat(totalAmount)))
			} else {
				// Other columns are empty
				html.WriteString(fmt.Sprintf(`<td %s></td>`+"\n", cellStyle))
			}
		}
		html.WriteString(`</tr>` + "\n")
	}

	html.WriteString(`</table>` + "\n")
	return html.String()
}

// BuildFixedStyledHTMLTable creates an HTML table with your specific styling for date/rate_type/amount data.
func BuildFixedStyledHTMLTable(data []map[string]interface{}) string {
	var html strings.Builder
	var totalAmount float64

	html.WriteString(`<table cellpadding="0" cellspacing="0" class="es-table" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;width:100%" role="presentation">` + "\n")

	// Fixed header
	html.WriteString(`<tr>` + "\n")
	html.WriteString(`<td style="padding:0;Margin:0;border-width:1px;border-style:solid;border-color:#cccccc">Date</td>` + "\n")
	html.WriteString(`<td style="padding:0;Margin:0;border-width:1px;border-style:solid;border-color:#cccccc">Rate Type</td>` + "\n")
	html.WriteString(`<td style="padding:0;Margin:0;border-width:1px;border-style:solid;border-color:#cccccc">Amount</td>` + "\n")
	html.WriteString(`</tr>` + "\n")

	// Data rows
	for _, row := range data {
		date := fmt.Sprintf("%v", row["date"])
		rateType := fmt.Sprintf("%v", row["rate_type"])
		amount := ""

		if val, ok := row["amount"]; ok && val != nil {
			if f, err := ToFloat64(val); err == nil {
				totalAmount += f
				amount = FormatMoneyFromFloat(f)
			} else {
				amount = fmt.Sprintf("%v", val)
			}
		}

		html.WriteString(`<tr>` + "\n")
		html.WriteString(fmt.Sprintf(`<td style="padding:0;Margin:0;border-width:1px;border-style:solid;border-color:#cccccc">%s</td>`+"\n", date))
		html.WriteString(fmt.Sprintf(`<td style="padding:0;Margin:0;border-width:1px;border-style:solid;border-color:#cccccc">%s</td>`+"\n", rateType))
		html.WriteString(fmt.Sprintf(`<td style="padding:0;Margin:0;border-width:1px;border-style:solid;border-color:#cccccc">%s</td>`+"\n", amount))
		html.WriteString(`</tr>` + "\n")
	}

	// Total row
	html.WriteString(`<tr>` + "\n")
	html.WriteString(`<td style="padding:0;Margin:0;border-width:1px;border-style:solid;border-color:#cccccc; font-weight:bold;">Total</td>` + "\n")
	html.WriteString(`<td style="padding:0;Margin:0;border-width:1px;border-style:solid;border-color:#cccccc"></td>` + "\n")
	html.WriteString(fmt.Sprintf(`<td style="padding:0;Margin:0;border-width:1px;border-style:solid;border-color:#cccccc; font-weight:bold;">%s</td>`+"\n", FormatMoneyFromFloat(totalAmount)))
	html.WriteString(`</tr>` + "\n")

	html.WriteString(`</table>` + "\n")
	return html.String()
}
