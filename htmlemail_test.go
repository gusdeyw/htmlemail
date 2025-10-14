package htmlemail

import (
	"strings"
	"testing"
)

func TestLoadTemplateFromString(t *testing.T) {
	content := "<h1>Hello $name$</h1>"
	template := LoadTemplateFromString(content)

	if template == nil {
		t.Fatal("LoadTemplateFromString returned nil")
	}

	if template.content != content {
		t.Errorf("Expected content '%s', got '%s'", content, template.content)
	}

	if template.variables == nil {
		t.Error("Variables map should be initialized")
	}
}

func TestTemplateSetVariable(t *testing.T) {
	template := LoadTemplateFromString("<h1>Hello $name$</h1>")

	// Test method chaining
	result := template.SetVariable("name", "John")
	if result != template {
		t.Error("SetVariable should return the same template instance for chaining")
	}

	// Test value
	if template.variables["name"] != "John" {
		t.Errorf("Expected variable 'name' to be 'John', got '%v'", template.variables["name"])
	}
}

func TestTemplateSetVariables(t *testing.T) {
	template := LoadTemplateFromString("<h1>Hello $name$, age $age$</h1>")

	vars := map[string]interface{}{
		"name": "Alice",
		"age":  25,
	}

	template.SetVariables(vars)

	if template.variables["name"] != "Alice" {
		t.Errorf("Expected variable 'name' to be 'Alice', got '%v'", template.variables["name"])
	}

	if template.variables["age"] != 25 {
		t.Errorf("Expected variable 'age' to be 25, got '%v'", template.variables["age"])
	}
}

func TestTemplateSetStruct(t *testing.T) {
	type TestData struct {
		FirstName string
		LastName  string
		Age       int
	}

	data := TestData{
		FirstName: "Bob",
		LastName:  "Smith",
		Age:       30,
	}

	template := LoadTemplateFromString("<h1>$first_name$ $last_name$, age $age$</h1>")
	template.SetStruct(data)

	if template.variables["first_name"] != "Bob" {
		t.Errorf("Expected 'first_name' to be 'Bob', got '%v'", template.variables["first_name"])
	}

	if template.variables["last_name"] != "Smith" {
		t.Errorf("Expected 'last_name' to be 'Smith', got '%v'", template.variables["last_name"])
	}

	if template.variables["age"] != 30 {
		t.Errorf("Expected 'age' to be 30, got '%v'", template.variables["age"])
	}
}

func TestTemplateRender(t *testing.T) {
	template := LoadTemplateFromString("<h1>Hello $name$, you are $age$ years old</h1>")
	template.SetVariable("name", "Charlie").SetVariable("age", 35)

	result, err := template.Render()
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	expected := "<h1>Hello Charlie, you are 35 years old</h1>"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestTemplateRenderWithUnresolvedPlaceholders(t *testing.T) {
	template := LoadTemplateFromString("<h1>Hello $name$, you are $age$ years old</h1>")
	template.SetVariable("name", "Charlie")
	// Note: not setting 'age'

	_, err := template.Render()
	if err == nil {
		t.Error("Expected error for unresolved placeholders")
	}

	if !strings.Contains(err.Error(), "unresolved placeholders") {
		t.Errorf("Expected error about unresolved placeholders, got: %v", err)
	}
}

func TestTemplateRenderWithDifferentStyles(t *testing.T) {
	// Test different placeholder styles
	tests := []struct {
		template string
		style    PlaceholderStyle
		expected string
	}{
		{"Hello $name$", DollarStyle, "Hello John"},
		{"Hello {{name}}", BraceStyle, "Hello John"},
		{"Hello %name%", PercentStyle, "Hello John"},
	}

	for _, test := range tests {
		tmpl := LoadTemplateFromString(test.template)
		tmpl.SetVariable("name", "John")

		result, err := tmpl.RenderWithStyle(test.style)
		if err != nil {
			t.Errorf("RenderWithStyle failed for %s: %v", test.template, err)
			continue
		}

		if result != test.expected {
			t.Errorf("For template '%s' with style %d, expected '%s', got '%s'",
				test.template, test.style, test.expected, result)
		}
	}
}

func TestNewEmailBuilder(t *testing.T) {
	builder := NewEmailBuilder()

	if builder == nil {
		t.Fatal("NewEmailBuilder returned nil")
	}

	if builder.data == nil {
		t.Error("Data map should be initialized")
	}
}

func TestEmailBuilderSetHTML(t *testing.T) {
	builder := NewEmailBuilder()
	html := "<h1>Test HTML</h1>"

	result := builder.SetHTML(html)
	if result != builder {
		t.Error("SetHTML should return the same builder instance for chaining")
	}

	if builder.htmlContent != html {
		t.Errorf("Expected HTML content '%s', got '%s'", html, builder.htmlContent)
	}
}

func TestEmailBuilderSetData(t *testing.T) {
	builder := NewEmailBuilder()

	builder.SetData("name", "Test User").SetData("age", 25)

	if builder.data["name"] != "Test User" {
		t.Errorf("Expected data 'name' to be 'Test User', got '%v'", builder.data["name"])
	}

	if builder.data["age"] != 25 {
		t.Errorf("Expected data 'age' to be 25, got '%v'", builder.data["age"])
	}
}

func TestEmailBuilderBuild(t *testing.T) {
	builder := NewEmailBuilder()

	html := "<h1>Hello $name$</h1><p>Age: $age$</p>"
	builder.SetHTML(html).SetData("name", "Alice").SetData("age", 30)

	result, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	expected := "<h1>Hello Alice</h1><p>Age: 30</p>"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestEmailBuilderWithCSS(t *testing.T) {
	builder := NewEmailBuilder()

	html := "<div class=\"test\">$content$</div>"
	css := ".test { color: red; }"

	builder.SetHTML(html).AddCSS(css).SetData("content", "Hello World")

	result, err := builder.Build()
	if err != nil {
		t.Fatalf("Build with CSS failed: %v", err)
	}

	if !strings.Contains(result, "Hello World") {
		t.Error("Result should contain the replaced content")
	}

	if !strings.Contains(result, ".test { color: red; }") {
		t.Error("Result should contain the CSS")
	}
}

func TestInsertHTMLInformation(t *testing.T) {
	html := "<h1>Hello $name$</h1>"

	result, err := InsertHTMLInformation(html, "name", "World")
	if err != nil {
		t.Fatalf("InsertHTMLInformation failed: %v", err)
	}

	expected := "<h1>Hello World</h1>"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestInsertHTMLInformationNotFound(t *testing.T) {
	html := "<h1>Hello World</h1>"

	_, err := InsertHTMLInformation(html, "name", "Test")
	if err == nil {
		t.Error("Expected error when placeholder not found")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected error about placeholder not found, got: %v", err)
	}
}

func TestBatchInsertHTMLInformation(t *testing.T) {
	html := "<h1>$title$</h1><p>Hello $name$, age $age$</p>"

	replacements := map[string]string{
		"title": "Welcome",
		"name":  "Bob",
		"age":   "25",
	}

	result, err := BatchInsertHTMLInformation(html, replacements)
	if err != nil {
		t.Fatalf("BatchInsertHTMLInformation failed: %v", err)
	}

	expected := "<h1>Welcome</h1><p>Hello Bob, age 25</p>"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestNewEmailTemplate(t *testing.T) {
	content := "<h1>Hello $name$</h1>"
	template := NewEmailTemplateFromString(content)

	if template == nil {
		t.Fatal("NewEmailTemplateFromString returned nil")
	}

	if template.htmlContent != content {
		t.Errorf("Expected content '%s', got '%s'", content, template.htmlContent)
	}
}

func TestEmailTemplateSetPlaceholder(t *testing.T) {
	template := NewEmailTemplateFromString("<h1>$greeting$ $name$</h1>")

	result := template.SetPlaceholder("greeting", "Hello").SetPlaceholder("name", "World")
	if result != template {
		t.Error("SetPlaceholder should return the same template for chaining")
	}

	if template.placeholders["greeting"] != "Hello" {
		t.Errorf("Expected placeholder 'greeting' to be 'Hello', got '%s'", template.placeholders["greeting"])
	}
}

func TestEmailTemplateRender(t *testing.T) {
	template := NewEmailTemplateFromString("<h1>$greeting$ $name$</h1>")
	template.SetPlaceholder("greeting", "Hi").SetPlaceholder("name", "Alice")

	result, err := template.Render()
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	expected := "<h1>Hi Alice</h1>"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestEmailTemplateRenderSafe(t *testing.T) {
	template := NewEmailTemplateFromString("<h1>$greeting$ $name$</h1>")
	template.SetPlaceholder("greeting", "Hello")
	// Note: not setting 'name'

	result := template.RenderSafe()
	expected := "<h1>Hello $name$</h1>" // 'name' should remain unresolved

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestConvertToString(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{"hello", "hello"},
		{42, "42"},
		{3.14, "3.14"},
		{true, "true"},
		{false, "false"},
		{nil, ""},
	}

	for _, test := range tests {
		result := convertToString(test.input)
		if result != test.expected {
			t.Errorf("convertToString(%v) = '%s', expected '%s'", test.input, result, test.expected)
		}
	}
}

func TestCamelToSnake(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"FirstName", "first_name"},
		{"lastName", "last_name"},
		{"Age", "age"},
		{"EmailAddress", "email_address"},
		{"ID", "i_d"},
	}

	for _, test := range tests {
		result := camelToSnake(test.input)
		if result != test.expected {
			t.Errorf("camelToSnake('%s') = '%s', expected '%s'", test.input, result, test.expected)
		}
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected float64
		hasError bool
	}{
		{float64(123.45), 123.45, false},
		{float32(67.89), 67.89, false}, // Note: float32 precision will be approximate
		{int(100), 100.0, false},
		{int64(200), 200.0, false},
		{int32(300), 300.0, false},
		{"456.78", 456.78, false},
		{"invalid", 0, true},
		{true, 0, true},
	}

	for _, test := range tests {
		result, err := ToFloat64(test.input)

		if test.hasError {
			if err == nil {
				t.Errorf("ToFloat64(%v) expected error but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("ToFloat64(%v) unexpected error: %v", test.input, err)
			}
			// Use approximate comparison for float32 inputs due to precision differences
			if _, isFloat32 := test.input.(float32); isFloat32 {
				if diff := result - test.expected; diff < -0.001 || diff > 0.001 {
					t.Errorf("ToFloat64(%v) = %f, expected approximately %f", test.input, result, test.expected)
				}
			} else {
				if result != test.expected {
					t.Errorf("ToFloat64(%v) = %f, expected %f", test.input, result, test.expected)
				}
			}
		}
	}
}

func TestFormatMoneyFromFloat(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{123.45, "$123.45"},
		{0.0, "$0.00"},
		{1000.99, "$1000.99"},
		{42, "$42.00"},
		{99.999, "$100.00"}, // Should round
	}

	for _, test := range tests {
		result := FormatMoneyFromFloat(test.input)
		if result != test.expected {
			t.Errorf("FormatMoneyFromFloat(%f) = '%s', expected '%s'", test.input, result, test.expected)
		}
	}
}

func TestBuildFixedStyledHTMLTable(t *testing.T) {
	data := []map[string]interface{}{
		{
			"date":      "2024-12-01",
			"rate_type": "Standard",
			"amount":    150.00,
		},
		{
			"date":      "2024-12-02",
			"rate_type": "Weekend",
			"amount":    200.50,
		},
	}

	result := BuildFixedStyledHTMLTable(data)

	// Check if result contains expected elements
	if !strings.Contains(result, "<table") {
		t.Error("Result should contain table tag")
	}
	if !strings.Contains(result, "Date") || !strings.Contains(result, "Rate Type") || !strings.Contains(result, "Amount") {
		t.Error("Result should contain header cells")
	}
	if !strings.Contains(result, "2024-12-01") || !strings.Contains(result, "Standard") {
		t.Error("Result should contain data from first row")
	}
	if !strings.Contains(result, "Total") {
		t.Error("Result should contain total row")
	}
	if !strings.Contains(result, "$350.50") { // 150.00 + 200.50
		t.Error("Result should contain correct total amount")
	}
	if !strings.Contains(result, "</table>") {
		t.Error("Result should contain closing table tag")
	}
}

func TestBuildHTMLTable(t *testing.T) {
	data := []map[string]interface{}{
		{"name": "Product A", "price": 25.99, "qty": 2},
		{"name": "Product B", "price": 15.50, "qty": 1},
	}

	options := TableOptions{
		Columns: []TableColumn{
			{Key: "name", Header: "Product"},
			{Key: "price", Header: "Price"},
			{Key: "qty", Header: "Quantity"},
		},
		ShowTotal:   true,
		TotalColumn: "price",
	}

	result := BuildHTMLTable(data, options)

	// Check basic structure
	if !strings.Contains(result, "<table") {
		t.Error("Result should contain table tag")
	}
	if !strings.Contains(result, "Product") || !strings.Contains(result, "Price") || !strings.Contains(result, "Quantity") {
		t.Error("Result should contain headers")
	}
	if !strings.Contains(result, "Product A") || !strings.Contains(result, "Product B") {
		t.Error("Result should contain data")
	}
	if !strings.Contains(result, "Total") {
		t.Error("Result should contain total row when ShowTotal is true")
	}
}

func TestBuildHTMLTableNoTotal(t *testing.T) {
	data := []map[string]interface{}{
		{"item": "Test", "value": 100},
	}

	options := TableOptions{
		Columns: []TableColumn{
			{Key: "item", Header: "Item"},
			{Key: "value", Header: "Value"},
		},
		ShowTotal: false,
	}

	result := BuildHTMLTable(data, options)

	if strings.Contains(result, "Total") {
		t.Error("Result should not contain total row when ShowTotal is false")
	}
}

func TestBuildHTMLTableEmptyData(t *testing.T) {
	data := []map[string]interface{}{}

	options := TableOptions{
		Columns: []TableColumn{
			{Key: "test", Header: "Test"},
		},
	}

	result := BuildHTMLTable(data, options)

	// Should still contain table structure and headers
	if !strings.Contains(result, "<table") {
		t.Error("Result should contain table tag even with empty data")
	}
	if !strings.Contains(result, "Test") {
		t.Error("Result should contain header even with empty data")
	}
}

func TestBuildHTMLTableWithCustomStyles(t *testing.T) {
	data := []map[string]interface{}{
		{"test": "value"},
	}

	customTableStyle := `style="border:2px solid red"`
	customHeaderStyle := `style="background:blue"`
	customCellStyle := `style="color:green"`

	options := TableOptions{
		TableStyle:  customTableStyle,
		HeaderStyle: customHeaderStyle,
		CellStyle:   customCellStyle,
		Columns: []TableColumn{
			{Key: "test", Header: "Test Column"},
		},
	}

	result := BuildHTMLTable(data, options)

	if !strings.Contains(result, customTableStyle) {
		t.Error("Result should contain custom table style")
	}
	if !strings.Contains(result, customHeaderStyle) {
		t.Error("Result should contain custom header style")
	}
	if !strings.Contains(result, customCellStyle) {
		t.Error("Result should contain custom cell style")
	}
}

func TestMinifyHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple HTML",
			input:    "<html>  <body>  <h1>Hello</h1>  </body>  </html>",
			expected: "<html><body><h1>Hello</h1></body></html>",
		},
		{
			name:     "remove comments",
			input:    "<!-- comment --><div>Hello</div>",
			expected: "<div>Hello</div>",
		},
		{
			name:     "preserve conditional comments",
			input:    "<!--[if IE]><div>IE only</div><![endif]--><div>All browsers</div>",
			expected: "<!--[if IE]><div>IE only</div><![endif]--><div>All browsers</div>",
		},
		{
			name:     "preserve pre content",
			input:    "<pre>  formatted\n  text  </pre>",
			expected: "<pre>  formatted\n  text  </pre>",
		},
		{
			name:     "preserve textarea content",
			input:    "<textarea>  user\n  input  </textarea>",
			expected: "<textarea>  user\n  input  </textarea>",
		},
		{
			name:     "preserve script content",
			input:    "<script>  var x = 1;\n  console.log(x);  </script>",
			expected: "<script>  var x = 1;\n  console.log(x);  </script>",
		},
		{
			name:     "collapse multiple spaces",
			input:    "<p>Hello   world</p>",
			expected: "<p>Hello world</p>",
		},
		{
			name:     "remove spaces around equals",
			input:    `<input type="text" name="test" value="123">`,
			expected: `<input type="text" name="test" value="123">`,
		},
		{
			name:     "complex HTML",
			input:    `<!DOCTYPE html>
<html>
<head>
    <title>Test</title>
</head>
<body>
    <!-- This is a comment -->
    <h1>  Hello   World  </h1>
    <p>This is a <strong>test</strong> paragraph.</p>
    <pre>
        Preserved formatting
        with multiple lines
    </pre>
</body>
</html>`,
			expected: `<!DOCTYPE html><html><head><title>Test</title></head><body><h1> Hello World </h1><p>This is a <strong>test</strong> paragraph.</p>
<pre>
        Preserved formatting
        with multiple lines
    </pre>
</body></html>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MinifyHTML(tt.input)
			if result != tt.expected {
				t.Errorf("MinifyHTML() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestTemplateRenderMinified(t *testing.T) {
	template := LoadTemplateFromString("<html>  <body>  <h1>Hello $name$</h1>  <p>Welcome!</p>  </body>  </html>")
	template.SetVariable("name", "World")

	result, err := template.RenderMinified()
	if err != nil {
		t.Fatalf("RenderMinified() error = %v", err)
	}

	expected := "<html><body><h1>Hello World</h1><p>Welcome!</p></body></html>"
	if result != expected {
		t.Errorf("RenderMinified() = %q, expected %q", result, expected)
	}
}

func TestTemplateRenderWithStyleMinified(t *testing.T) {
	template := LoadTemplateFromString("<html><body><h1>Hello {{name}}</h1><p>Welcome!</p></body></html>")
	template.SetVariable("name", "World")

	result, err := template.RenderWithStyleMinified(BraceStyle)
	if err != nil {
		t.Fatalf("RenderWithStyleMinified() error = %v", err)
	}

	expected := "<html><body><h1>Hello World</h1><p>Welcome!</p></body></html>"
	if result != expected {
		t.Errorf("RenderWithStyleMinified() = %q, expected %q", result, expected)
	}
}

func TestTemplateRenderSafeMinified(t *testing.T) {
	template := LoadTemplateFromString("<html>  <body>  <h1>Hello $name$</h1>  <p>Welcome!</p>  </body>  </html>")
	template.SetVariable("name", "World")

	result := template.RenderSafeMinified()

	expected := "<html><body><h1>Hello World</h1><p>Welcome!</p></body></html>"
	if result != expected {
		t.Errorf("RenderSafeMinified() = %q, expected %q", result, expected)
	}
}

func TestTemplateRenderWithGoTemplateMinified(t *testing.T) {
	template := LoadTemplateFromString("<html><body><h1>Hello {{.Name}}</h1><p>Welcome!</p></body></html>")
	template.SetVariable("Name", "World")

	result, err := template.RenderWithGoTemplateMinified()
	if err != nil {
		t.Fatalf("RenderWithGoTemplateMinified() error = %v", err)
	}

	expected := "<html><body><h1>Hello World</h1><p>Welcome!</p></body></html>"
	if result != expected {
		t.Errorf("RenderWithGoTemplateMinified() = %q, expected %q", result, expected)
	}
}

func TestEmailBuilderEnableMinification(t *testing.T) {
	builder := NewEmailBuilder().
		SetHTML("<html>  <body>  <h1>Hello $name$</h1>  </body>  </html>").
		SetData("name", "World").
		EnableMinification()

	result, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	expected := "<html><body><h1>Hello World</h1></body></html>"
	if result != expected {
		t.Errorf("Build() with minification = %q, expected %q", result, expected)
	}
}

func TestEmailBuilderDisableMinification(t *testing.T) {
	builder := NewEmailBuilder().
		SetHTML("<html>  <body>  <h1>Hello $name$</h1>  </body>  </html>").
		SetData("name", "World").
		EnableMinification().
		DisableMinification()

	result, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Should not be minified
	if !strings.Contains(result, "  <body>") {
		t.Error("Expected unminified HTML with extra spaces")
	}
}

func TestTemplateRenderSafe(t *testing.T) {
	template := LoadTemplateFromString("<h1>Hello $name$</h1><p>$missing$ placeholder</p>")
	template.SetVariable("name", "World")

	result := template.RenderSafe()

	expected := "<h1>Hello World</h1><p>$missing$ placeholder</p>"
	if result != expected {
		t.Errorf("RenderSafe() = %q, expected %q", result, expected)
	}
}
