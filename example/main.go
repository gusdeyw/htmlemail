package main

import (
	"fmt"
	"log"

	"htmlemail"
)

// BookingEmailStruct represents booking data (similar to your original structure)
type BookingEmailStruct struct {
	HotelName     string
	HotelAddress  string
	HotelEmail    string
	HotelStreet   string
	Name          string
	Email         string
	Phone         string
	BookingCode   string
	BookingDate   string
	BookingStatus string
	ArrivalDate   string
	DepartureDate string
	Guest         string
	RoomType      string
	RoomRate      string
	RoomPrice     string
	NumberOfRooms string
}

func main() {
	fmt.Println("🎨 HTML Email Template Generator Examples")
	fmt.Println("=========================================")

	// Example 1: Your Original Approach (Enhanced)
	fmt.Println("\n1️⃣ Enhanced Version of Your Original Approach:")
	demonstrateOriginalApproach()

	// Example 2: Template-based Approach
	fmt.Println("\n2️⃣ Template-based Approach:")
	demonstrateTemplateApproach()

	// Example 3: Fluent EmailBuilder
	fmt.Println("\n3️⃣ Advanced EmailBuilder with CSS:")
	demonstrateEmailBuilder()

	// Example 4: Struct-based Template Population
	fmt.Println("\n4️⃣ Struct-based Template Population:")
	demonstrateStructBasedTemplate()

	// Example 5: Go Template Engine
	fmt.Println("\n5️⃣ Go Template Engine (Advanced):")
	demonstrateGoTemplateEngine()

	// Example 6: HTML Table Generation
	fmt.Println("\n6️⃣ HTML Table Generation:")
	demonstrateTableGeneration()

	fmt.Println("\n✅ All examples completed successfully!")
}

// demonstrateOriginalApproach shows an enhanced version of your original approach
func demonstrateOriginalApproach() {
	// Create a simple HTML template as string (simulating reading from file)
	htmlTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Booking Confirmation</title>
	</head>
	<body>
		<h1>Hotel: $hotel_name$</h1>
		<p>Address: $hotel_address$</p>
		<p>Email: $hotel_email$</p>
		<hr>
		<h2>Booking Details</h2>
		<p>Guest Name: $name$</p>
		<p>Email: $email$</p>
		<p>Phone: $phone$</p>
		<p>Booking Code: $booking_code$</p>
		<p>Check-in: $arrival_date$</p>
		<p>Check-out: $departure_date$</p>
		<p>Room Type: $room_type$</p>
		<p>Room Price: $room_price$</p>
		<p>Number of Rooms: $number_of_rooms$</p>
	</body>
	</html>
	`

	// Your original approach but enhanced
	html := htmlTemplate

	// Using the enhanced InsertHTMLInformation function
	html, _ = htmlemail.InsertHTMLInformation(html, "hotel_name", "Grand Hotel")
	html, _ = htmlemail.InsertHTMLInformation(html, "hotel_address", "123 Main St, City")
	html, _ = htmlemail.InsertHTMLInformation(html, "hotel_email", "info@grandhotel.com")
	html, _ = htmlemail.InsertHTMLInformation(html, "name", "John Doe")
	html, _ = htmlemail.InsertHTMLInformation(html, "email", "john@example.com")
	html, _ = htmlemail.InsertHTMLInformation(html, "phone", "+1234567890")
	html, _ = htmlemail.InsertHTMLInformation(html, "booking_code", "BK001")
	html, _ = htmlemail.InsertHTMLInformation(html, "arrival_date", "2024-12-25")
	html, _ = htmlemail.InsertHTMLInformation(html, "departure_date", "2024-12-28")
	html, _ = htmlemail.InsertHTMLInformation(html, "room_type", "Deluxe Suite")
	html, _ = htmlemail.InsertHTMLInformation(html, "room_price", "$299/night")
	html, _ = htmlemail.InsertHTMLInformation(html, "number_of_rooms", "2")

	fmt.Printf("✅ Generated HTML email (%d characters)\n", len(html))
	fmt.Println("   Using your original approach (enhanced)")
}

// demonstrateTemplateApproach shows the new template-based approach
func demonstrateTemplateApproach() {
	htmlTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>$subject$</title>
		<style>
			.header { background: #4CAF50; color: white; padding: 20px; }
			.content { padding: 20px; }
		</style>
	</head>
	<body>
		<div class="header">
			<h1>$title$</h1>
		</div>
		<div class="content">
			<p>Dear $customer_name$,</p>
			<p>$message$</p>
			<p>Amount: $amount$</p>
			<p>Status: $status$</p>
		</div>
	</body>
	</html>
	`

	// Using the new EmailTemplate approach
	template := htmlemail.NewEmailTemplateFromString(htmlTemplate)

	// Set data using the fluent interface
	html, err := template.
		SetPlaceholder("subject", "Invoice Notification").
		SetPlaceholder("title", "Invoice #INV-001").
		SetPlaceholder("customer_name", "Jane Smith").
		SetPlaceholder("message", "Your invoice is ready for review.").
		SetPlaceholder("amount", "$1,250.00").
		SetPlaceholder("status", "Paid").
		Render()

	if err != nil {
		log.Printf("Error rendering template: %v", err)
		return
	}

	fmt.Printf("✅ Generated HTML email (%d characters)\n", len(html))
	fmt.Println("   Using template-based approach")
}

// demonstrateEmailBuilder shows the advanced EmailBuilder
func demonstrateEmailBuilder() {
	htmlContent := `
	<div class="container">
		<h1>$title$</h1>
		<p>Hello $name$!</p>
		<p>Your order #$order_id$ is ready.</p>
		<div class="total">Total: $total$</div>
	</div>
	`

	css := `
	.container {
		max-width: 600px;
		margin: 0 auto;
		font-family: Arial, sans-serif;
	}
	.total {
		background: #f0f0f0;
		padding: 10px;
		font-weight: bold;
		border-radius: 5px;
	}
	`

	// Using EmailBuilder for advanced features
	html, err := htmlemail.NewEmailBuilder().
		SetHTML(htmlContent).
		AddCSS(css).
		SetInlineCSS(false).
		SetData("title", "Order Confirmation").
		SetData("name", "Alice Johnson").
		SetData("order_id", "ORD-12345").
		SetData("total", "$89.99").
		Build()

	if err != nil {
		log.Printf("Error building email: %v", err)
		return
	}

	fmt.Printf("✅ Generated HTML email (%d characters)\n", len(html))
	fmt.Println("   Using EmailBuilder with CSS")
}

// demonstrateStructBasedTemplate shows struct-based template population
func demonstrateStructBasedTemplate() {
	// Create booking data (like your original BookingEmailStruct)
	bookingData := BookingEmailStruct{
		HotelName:     "Luxury Resort",
		HotelAddress:  "456 Beach Ave, Miami",
		HotelEmail:    "reservations@luxuryresort.com",
		Name:          "Bob Wilson",
		Email:         "bob@example.com",
		Phone:         "+1987654321",
		BookingCode:   "BK002",
		BookingDate:   "2024-12-20",
		BookingStatus: "Confirmed",
		ArrivalDate:   "2024-12-30",
		DepartureDate: "2025-01-03",
		RoomType:      "Ocean View Suite",
		RoomPrice:     "$450/night",
		NumberOfRooms: "1",
	}

	htmlTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Booking Confirmation - $booking_code$</title>
	</head>
	<body>
		<h1>$hotel_name$</h1>
		<p>📧 $hotel_email$ | 📍 $hotel_address$</p>
		<hr>
		<h2>Booking Confirmation</h2>
		<p><strong>Guest:</strong> $name$ ($email$)</p>
		<p><strong>Phone:</strong> $phone$</p>
		<p><strong>Booking Code:</strong> $booking_code$</p>
		<p><strong>Status:</strong> $booking_status$</p>
		<p><strong>Dates:</strong> $arrival_date$ to $departure_date$</p>
		<p><strong>Room:</strong> $room_type$ - $room_price$</p>
		<p><strong>Rooms:</strong> $number_of_rooms$</p>
	</body>
	</html>
	`

	// Using struct-based population (automatic field mapping)
	template := htmlemail.NewEmailTemplateFromString(htmlTemplate)
	html := template.SetStructData(bookingData).RenderSafe() // Use RenderSafe to ignore missing placeholders

	fmt.Printf("✅ Generated HTML email (%d characters)\n", len(html))
	fmt.Println("   Using struct-based template population")

	// Show unresolved placeholders if any
	unresolved := template.GetUnresolvedPlaceholders()
	if len(unresolved) > 0 {
		fmt.Printf("   ⚠️  Unresolved placeholders: %v\n", unresolved)
	}
}

// demonstrateGoTemplateEngine shows Go template engine usage
func demonstrateGoTemplateEngine() {
	goTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>{{.Title}}</title>
	</head>
	<body>
		<h1>Hello {{.CustomerName}}!</h1>
		<p>Order Details:</p>
		<ul>
		{{range .Items}}
			<li>{{.Name}} - ${{.Price}} (Qty: {{.Quantity}})</li>
		{{end}}
		</ul>
		<p><strong>Total: ${{printf "%.2f" .Total}}</strong></p>
		{{if gt .Total 100.0}}
		<p>🎉 You qualify for free shipping!</p>
		{{end}}
	</body>
	</html>
	`

	// Data for Go template
	data := map[string]interface{}{
		"Title":        "Order Summary",
		"CustomerName": "Charlie Brown",
		"Items": []map[string]interface{}{
			{"Name": "Product A", "Price": 25.99, "Quantity": 2},
			{"Name": "Product B", "Price": 45.50, "Quantity": 1},
		},
		"Total": 97.48,
	}

	// Using Go template engine for advanced features
	template := htmlemail.LoadTemplateFromString(goTemplate)
	html, err := template.SetVariables(data).RenderWithGoTemplate()

	if err != nil {
		log.Printf("Error rendering Go template: %v", err)
		return
	}

	fmt.Printf("✅ Generated HTML email (%d characters)\n", len(html))
	fmt.Println("   Using Go template engine (with loops and conditionals)")
}

// demonstrateTableGeneration shows different ways to generate HTML tables
func demonstrateTableGeneration() {
	// Example data for tables
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
		{
			"date":      "2024-12-03",
			"rate_type": "Holiday Rate",
			"amount":    275.25,
		},
	}

	// Method 1: Your original fixed-styled table
	fmt.Println("   📊 Fixed-Styled Table (Your Original Approach):")
	fixedTable := htmlemail.BuildFixedStyledHTMLTable(tableData)
	fmt.Printf("   ✅ Generated table with %d rows (%d characters)\n", len(tableData)+2, len(fixedTable)) // +2 for header and total

	// Method 2: Configurable table with custom styling
	fmt.Println("   📊 Configurable Table with Custom Options:")
	customOptions := htmlemail.TableOptions{
		TableStyle:  `cellpadding="5" cellspacing="0" style="border-collapse:collapse;width:100%;font-family:Arial,sans-serif"`,
		HeaderStyle: `style="background-color:#4CAF50;color:white;padding:10px;border:1px solid #ddd;text-align:left"`,
		CellStyle:   `style="padding:8px;border:1px solid #ddd"`,
		Columns: []htmlemail.TableColumn{
			{Key: "date", Header: "Date", Style: ""},
			{Key: "rate_type", Header: "Rate Type", Style: ""},
			{Key: "amount", Header: "Amount", Style: `style="padding:8px;border:1px solid #ddd;text-align:right"`},
		},
		ShowTotal:   true,
		TotalColumn: "amount",
	}

	customTable := htmlemail.BuildHTMLTable(tableData, customOptions)
	fmt.Printf("   ✅ Generated custom styled table (%d characters)\n", len(customTable))

	// Method 3: Simple table without totals
	fmt.Println("   📊 Simple Table (No Totals):")
	simpleOptions := htmlemail.TableOptions{
		Columns: []htmlemail.TableColumn{
			{Key: "date", Header: "Date"},
			{Key: "rate_type", Header: "Description"},
			{Key: "amount", Header: "Price"},
		},
		ShowTotal: false,
	}

	simpleTable := htmlemail.BuildHTMLTable(tableData, simpleOptions)
	fmt.Printf("   ✅ Generated simple table (%d characters)\n", len(simpleTable))

	// Method 4: Using table in template
	fmt.Println("   📊 Table Embedded in Email Template:")
	tableHTML := htmlemail.BuildFixedStyledHTMLTable(tableData)

	emailTemplate := htmlemail.LoadTemplateFromString(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Invoice</title>
		</head>
		<body>
			<h1>Invoice for $customer_name$</h1>
			<p>Invoice Date: $invoice_date$</p>
			
			<h2>Charges:</h2>
			$charges_table$
			
			<p>Thank you for your business!</p>
		</body>
		</html>
	`)

	finalEmail, err := emailTemplate.
		SetVariable("customer_name", "John Doe").
		SetVariable("invoice_date", "2024-12-04").
		SetVariable("charges_table", tableHTML).
		Render()

	if err != nil {
		log.Printf("Error creating email with table: %v", err)
		return
	}

	fmt.Printf("   ✅ Generated complete email with embedded table (%d characters)\n", len(finalEmail))

	// Method 5: Multiple different tables
	fmt.Println("   📊 Multiple Table Types:")

	// Product inventory table
	inventoryData := []map[string]interface{}{
		{"product": "Widget A", "stock": 150, "price": 25.99},
		{"product": "Widget B", "stock": 75, "price": 45.50},
		{"product": "Widget C", "stock": 200, "price": 15.25},
	}

	inventoryOptions := htmlemail.TableOptions{
		Columns: []htmlemail.TableColumn{
			{Key: "product", Header: "Product Name"},
			{Key: "stock", Header: "In Stock"},
			{Key: "price", Header: "Unit Price"},
		},
		ShowTotal:   false,
		CellStyle:   `style="padding:5px;border:1px solid #ccc"`,
		HeaderStyle: `style="padding:5px;border:1px solid #ccc;background-color:#f0f0f0;font-weight:bold"`,
	}

	inventoryTable := htmlemail.BuildHTMLTable(inventoryData, inventoryOptions)
	fmt.Printf("   ✅ Generated inventory table (%d characters)\n", len(inventoryTable))
}
