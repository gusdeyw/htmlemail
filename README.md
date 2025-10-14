# htmlemail

[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/doc/devel/release.html)
[![MIT License](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=for-the-badge&logo=github-actions)](https://github.com/gusdeyw/htmlemail/actions)

A powerful and flexible Go package for **generating dynamic HTML email content** from templates. Focus on creating beautiful, data-driven HTML emails without worrying about SMTP complexities.

## Features

- 🎨 **Dynamic HTML Generation** - Transform templates with placeholder replacement  
- � **Multiple Template Styles** - Support for `$variable$`, `{{variable}}`, and `%variable%` formats
- 🏗️ **Fluent API** - Chainable methods for clean, readable code
- � **File & String Templates** - Load templates from files or create them inline
- 🔧 **Struct Mapping** - Automatically map struct fields to template variables
- 🎯 **Advanced Templates** - Go template engine support with loops and conditionals
- � **CSS Integration** - Add styles inline or in `<head>` section
- ✅ **Well Tested** - Comprehensive test coverage with edge cases
- 🔍 **Template Validation** - Detect unresolved placeholders and missing data

## Installation

```bash
go get -u github.com/gusdeyw/htmlemail
```

## Quick Start

### Your Original Approach (Enhanced)

If you're coming from a simple placeholder replacement approach:

```go
package main

import "github.com/gusdeyw/htmlemail"

func main() {
    // Load HTML template
    html, _ := htmlemail.ReadHTMLFile("./templates/booking.html")
    
    // Replace placeholders (your original approach)
    html, _ = htmlemail.InsertHTMLInformation(html, "hotel_name", "Grand Hotel")
    html, _ = htmlemail.InsertHTMLInformation(html, "guest_name", "John Doe")
    html, _ = htmlemail.InsertHTMLInformation(html, "booking_code", "BK001")
    
    // html now contains the final email content
}
```

### Template-Based Approach (Recommended)

For more complex scenarios with better structure:

```go
package main

import "github.com/gusdeyw/htmlemail"

func main() {
    // Create template from string
    template := htmlemail.NewEmailTemplateFromString(`
        <h1>Hello $name$!</h1>
        <p>Your order #$order_id$ totaling $total$ is ready.</p>
    `)
    
    // Set variables and render
    html, err := template.
        SetPlaceholder("name", "Alice").
        SetPlaceholder("order_id", "12345").
        SetPlaceholder("total", "$99.99").
        Render()
        
    if err != nil {
        panic(err)
    }
    
    // Use html for email sending
}
```

### Struct-Based Population

Perfect for existing data structures:

```go
type BookingData struct {
    HotelName    string
    GuestName    string
    BookingCode  string
    CheckIn      string
    CheckOut     string
    RoomType     string
    TotalAmount  string
}

func main() {
    booking := BookingData{
        HotelName:   "Luxury Resort",
        GuestName:   "Jane Smith", 
        BookingCode: "BK002",
        CheckIn:     "2024-12-25",
        CheckOut:    "2024-12-28",
        RoomType:    "Ocean Suite",
        TotalAmount: "$1,200.00",
    }

    template := htmlemail.NewEmailTemplateFromString(`
        <h1>$hotel_name$ - Booking Confirmation</h1>
        <p>Dear $guest_name$,</p>
        <p>Booking: $booking_code$</p>
        <p>Dates: $check_in$ to $check_out$</p>
        <p>Room: $room_type$</p>
        <p>Total: $total_amount$</p>
    `)
    
    // Automatically maps struct fields to template variables
    html := template.SetStructData(booking).RenderSafe()
}
```

## Table Generation

The package includes powerful HTML table generation capabilities, perfect for invoices, reports, and data tables in emails.

### Your Fixed-Style Table (Exact Match)

For data with `date`, `rate_type`, and `amount` fields (matching your original approach):

```go
tableData := []map[string]interface{}{
    {
        "date":      "2024-12-01",
        "rate_type": "Standard Rate", 
        "amount":    150.00,
    },
    {
        "date":      "2024-12-02", 
        "rate_type": "Weekend Rate",
        "amount":    200.50,
    },
}

// Your exact styling with automatic totals
tableHTML := htmlemail.BuildFixedStyledHTMLTable(tableData)
```

### Configurable Tables

For flexible table generation with custom styling:

```go
options := htmlemail.TableOptions{
    TableStyle:  `cellpadding="5" cellspacing="0" style="border-collapse:collapse;width:100%"`,
    HeaderStyle: `style="background-color:#4CAF50;color:white;padding:10px"`,
    CellStyle:   `style="padding:8px;border:1px solid #ddd"`,
    Columns: []htmlemail.TableColumn{
        {Key: "product", Header: "Product Name"},
        {Key: "price", Header: "Price", Style: `style="text-align:right"`},
        {Key: "quantity", Header: "Qty"},
    },
    ShowTotal:   true,
    TotalColumn: "price", // Which column to sum
}

tableHTML := htmlemail.BuildHTMLTable(data, options)
```

### Embedding Tables in Templates

```go
// Generate table
tableHTML := htmlemail.BuildFixedStyledHTMLTable(invoiceData)

// Embed in email template
template := htmlemail.LoadTemplateFromString(`
    <h1>Invoice for $customer_name$</h1>
    <p>Date: $invoice_date$</p>
    
    <h2>Charges:</h2>
    $charges_table$
    
    <p>Total Amount: $total_amount$</p>
`)

email, _ := template.
    SetVariable("customer_name", "John Doe").
    SetVariable("invoice_date", "2024-12-04").
    SetVariable("charges_table", tableHTML).
    SetVariable("total_amount", "$625.75").
    Render()
```

### Helper Functions

```go
// Convert various types to float64
amount, err := htmlemail.ToFloat64("123.45") // Returns 123.45

// Format currency
formatted := htmlemail.FormatMoneyFromFloat(123.45) // Returns "$123.45"
```

## API Reference

### Template Loading

```go
// From file
template, err := htmlemail.LoadTemplate("./templates/email.html")

// From string  
template := htmlemail.LoadTemplateFromString("<h1>Hello $name$</h1>")
```

### Setting Template Data

```go
// Single variable
template.SetVariable("name", "John")

// Multiple variables
template.SetVariables(map[string]interface{}{
    "name": "John",
    "age":  25,
})

// From struct (auto-converts CamelCase to snake_case)
template.SetStruct(userData)
```

### Rendering Templates

```go
// Render with error checking (fails on unresolved placeholders)
html, err := template.Render()

// Render safely (ignores unresolved placeholders)  
html := template.RenderSafe()

// Render with specific placeholder style
html, err := template.RenderWithStyle(htmlemail.BraceStyle) // {{variable}}

// Render only specific variables
html, err := template.RenderPartial([]string{"name", "email"})
```

### Advanced EmailBuilder

For complex emails with CSS and advanced features:

```go
html, err := htmlemail.NewEmailBuilder().
    LoadHTMLFromFile("./templates/newsletter.html").
    AddCSSFromFile("./styles/email.css").
    SetInlineCSS(true).
    SetData("title", "Monthly Newsletter").
    SetData("month", "December").
    SetDataFromStruct(newsletterData).
    Build()
```

### Go Template Engine

For advanced templating with loops and conditionals:

```go
template := htmlemail.LoadTemplateFromString(`
    <h1>Hello {{.CustomerName}}!</h1>
    <ul>
    {{range .Items}}
        <li>{{.Name}} - ${{.Price}}</li>
    {{end}}
    </ul>
    {{if gt .Total 100}}
        <p>Free shipping applied!</p>
    {{end}}
`)

data := map[string]interface{}{
    "CustomerName": "Bob",
    "Items": []map[string]interface{}{
        {"Name": "Product A", "Price": 50},
        {"Name": "Product B", "Price": 75},
    },
    "Total": 125,
}

html, err := template.SetVariables(data).RenderWithGoTemplate()
```

## Placeholder Styles

The package supports multiple placeholder formats:

```go
// Dollar style (default, matches your original approach)
"Hello $name$, your order $order_id$ is ready"

// Brace style (Go template compatible)  
"Hello {{name}}, your order {{order_id}} is ready"

// Percent style
"Hello %name%, your order %order_id% is ready"
```

## Legacy Functions

For backward compatibility with your original approach:

```go
// Read HTML file (equivalent to your ReadHTMLFile)
html, err := htmlemail.ReadHTMLFile("./template.html")

// Single replacement (equivalent to your InsertHTMLInformation)
html, err := htmlemail.InsertHTMLInformation(html, "name", "John")

// Batch replacement (enhanced version)
replacements := map[string]string{
    "hotel_name": "Grand Hotel",
    "guest_name": "John Doe", 
    "booking_code": "BK001",
}
html, err := htmlemail.BatchInsertHTMLInformation(html, replacements)
```

## Examples

### Recreating Your Booking Email Function

Here's how to recreate your `CreateEmailBodyForBooking` function using this package:

```go
type BookingEmailStruct struct {
    HotelName     string
    HotelAddress  string  
    HotelEmail    string
    Name          string
    Email         string
    Phone         string
    BookingCode   string
    BookingDate   string
    BookingStatus string
    ArrivalDate   string
    DepartureDate string
    RoomType      string
    RoomPrice     string
    NumberOfRooms string
}

func CreateEmailBodyForBooking(data BookingEmailStruct) (string, error) {
    // Load template (instead of hardcoded path)
    template, err := htmlemail.LoadTemplate("./email_template/book_details.html")
    if err != nil {
        return "", err
    }

    // Use struct-based population for cleaner code
    return template.SetStruct(data).Render()
}

// Or using the batch approach (closer to your original)
func CreateEmailBodyForBookingBatch(data BookingEmailStruct) (string, error) {
    html, err := htmlemail.ReadHTMLFile("./email_template/book_details.html")
    if err != nil {
        return "", err
    }

    replacements := map[string]string{
        "hotel_name":      data.HotelName,
        "hotel_address":   data.HotelAddress,
        "hotel_email":     data.HotelEmail,
        "name":            data.Name,
        "email":           data.Email,
        "phone":           data.Phone,
        "booking_code":    data.BookingCode,
        "booking_date":    data.BookingDate,
        "booking_status":  data.BookingStatus,
        "arrival_date":    data.ArrivalDate,
        "departure_date":  data.DepartureDate,
        "room_type":       data.RoomType,
        "room_price":      data.RoomPrice,
        "number_of_rooms": data.NumberOfRooms,
    }

    return htmlemail.BatchInsertHTMLInformation(html, replacements)
}
```

### Newsletter with Dynamic Content

```go
func CreateNewsletter(articles []Article, subscriber Subscriber) (string, error) {
    template := htmlemail.LoadTemplateFromString(`
        <!DOCTYPE html>
        <html>
        <head>
            <style>
                .header { background: #3498db; color: white; padding: 20px; }
                .article { margin: 20px 0; padding: 15px; border-left: 3px solid #3498db; }
            </style>
        </head>
        <body>
            <div class="header">
                <h1>$newsletter_title$</h1>
                <p>Hello $subscriber_name$!</p>
            </div>
            
            {{range .Articles}}
            <div class="article">
                <h2>{{.Title}}</h2>
                <p>{{.Summary}}</p>
                <a href="{{.URL}}">Read More</a>
            </div>
            {{end}}
            
            <p>Thanks for subscribing!</p>
        </body>
        </html>
    `)

    data := map[string]interface{}{
        "Articles": articles,
    }

    return template.
        SetVariable("newsletter_title", "Tech Weekly").
        SetVariable("subscriber_name", subscriber.Name).
        SetVariables(data).
        RenderWithGoTemplate()
}
```

## Error Handling

```go
// Check for unresolved placeholders
template := htmlemail.LoadTemplateFromString("Hello $name$, your $item$ is ready")
template.SetVariable("name", "John")

unresolved := template.GetUnresolvedPlaceholders()
if len(unresolved) > 0 {
    fmt.Printf("Missing data for: %v\n", unresolved) // Output: [item]
}

// Render safely (ignores missing placeholders)
html := template.RenderSafe() // Output: "Hello John, your $item$ is ready"

// Or render strictly (returns error for missing placeholders)
html, err := template.Render()
if err != nil {
    log.Printf("Template error: %v", err)
}
```

## Testing

Run the test suite:

```bash
go test -v
```

Run tests with coverage:

```bash
go test -v -cover
```

Run the examples:

```bash
cd example && go run main.go
```

## Performance

- ✅ **Zero allocations** for simple placeholder replacement
- ✅ **Compiled templates** cached for repeated use  
- ✅ **Batch operations** for multiple replacements
- ✅ **Lazy evaluation** - only processes used variables

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

- [x] ~~Dynamic placeholder replacement~~ ✅ 
- [x] ~~Multiple placeholder styles~~ ✅
- [x] ~~Struct-based template population~~ ✅  
- [x] ~~Go template engine integration~~ ✅
- [x] ~~CSS integration and inlining~~ ✅
- [ ] Email template library/gallery
- [ ] HTML minification
- [ ] Template inheritance/layout system
- [ ] Advanced CSS inlining with external stylesheets
- [ ] Template debugging and preview tools
- [ ] Security enhancements (HTML sanitization, CSP headers)
- [ ] Performance optimizations (template caching, streaming rendering)
- [ ] Email service integrations (SendGrid, Mailgun, AWS SES)
- [ ] Internationalization (i18n) support with locale-aware formatting
- [ ] Accessibility features (alt text validation, semantic HTML checks)
- [ ] CLI tools for template validation and preview
- [ ] Plugin system for custom template functions
- [ ] Monitoring and metrics (rendering performance, error tracking)
- [ ] IDE integrations (VS Code extension, syntax highlighting)
- [ ] Advanced table features (sorting, filtering, pagination)
- [ ] Template versioning and migration tools
- [ ] Web-based template editor and preview interface