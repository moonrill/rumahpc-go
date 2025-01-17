package templates

import (
	"fmt"
	"html/template"
)

// FormatIDR converts a number to IDR format
func FormatIDR(amount int) string {
	formatted := fmt.Sprintf("%.0f", float64(amount))
	length := len(formatted)
	var result string

	for i := length - 1; i >= 0; i-- {
		if (length-i-1)%3 == 0 && i != length-1 {
			result = "." + result
		}
		result = string(formatted[i]) + result
	}

	return "Rp " + result
}

var TemplateFuncs = template.FuncMap{
	"formatIDR": FormatIDR,
}

const Delivered = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Order Delivered</title>
    <style>
        body {
            font-family: 'Segoe UI', Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f8f9fa;
            -webkit-font-smoothing: antialiased;
        }
        
        .container {
            max-width: 600px;
            margin: 20px auto;
            background: #ffffff;
            border-radius: 12px;
            overflow: hidden;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
        }
        
        .header {
            background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%);
            color: white;
            text-align: center;
            padding: 30px 20px;
        }
        
        .header h1 {
            margin: 0;
            font-size: 28px;
            font-weight: 600;
            letter-spacing: -0.5px;
        }
        
        .content {
            padding: 30px;
            color: #2c3e50;
            line-height: 1.6;
        }
        
        .order-details {
            margin: 25px 0;
            background: #f8f9fa;
            border-radius: 8px;
            padding: 20px;
            border: 1px solid #e9ecef;
        }
        
        .order-details h3 {
            margin: 0 0 20px 0;
            color: #2c3e50;
            font-size: 20px;
            font-weight: 600;
        }
        
        .order-item {
            margin-bottom: 20px;
            padding-bottom: 20px;
            border-bottom: 1px solid #e9ecef;
        }
        
        .order-item:last-child {
            border-bottom: none;
            margin-bottom: 10px;
            padding-bottom: 0;
        }
        
        .price-tag {
            color: #2c3e50;
            font-weight: 600;
            font-size: 16px;
        }
        
        .total-price {
            margin-top: 20px;
            padding-top: 20px;
            border-top: 2px solid #e9ecef;
            font-size: 18px;
            font-weight: 600;
        }
        
        .button {
            display: inline-block;
            padding: 12px 24px;
            background: #4CAF50;
            color: white;
            text-decoration: none;
            border-radius: 6px;
            font-weight: 500;
            margin: 20px 0;
            text-align: center;
            transition: background-color 0.2s;
        }
        
        .button:hover {
            background: #45a049;
        }
        
        .shipping-info {
            background: #fff;
            padding: 15px;
            border-radius: 6px;
            margin-top: 20px;
            border: 1px solid #e9ecef;
        }
        
        .footer {
            text-align: center;
            font-size: 14px;
            color: #6c757d;
            padding: 20px;
            background: #f8f9fa;
            border-top: 1px solid #e9ecef;
        }
        
        .product-grid {
            display: table;
            width: 100%;
        }
        
        .product-details {
            margin: 10px 0;
        }
        
        .product-name {
            font-weight: 600;
            color: #2c3e50;
            font-size: 16px;
        }
        
        .highlight {
            color: #4CAF50;
            font-weight: 600;
        }
        
        @media only screen and (max-width: 600px) {
            .container {
                margin: 10px;
                border-radius: 8px;
            }
            
            .content {
                padding: 20px;
            }
            
            .header h1 {
                font-size: 24px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Your Order Has Been Delivered! 🎉</h1>
        </div>
        <div class="content">
            <p>Dear <span class="highlight">{{.User.Name}}</span>,</p>
            <p>Great news! Your order <span class="highlight">#{{.ID}}</span> has been successfully delivered to your doorstep.</p>
            
            <div class="order-details">
                <h3>📦 Order Summary</h3>
                {{range .OrderItems}}
                <div class="order-item">
                    <div class="product-grid">
                        <div class="product-details">
                            <div class="product-name">{{.Product.Name}}</div>
                            <p style="margin: 5px 0;">Quantity: {{.Quantity}}</p>
                            <p class="price-tag" style="margin: 5px 0;">Subtotal: {{formatIDR .SubTotal}}</p>
                        </div>
                    </div>
                </div>
                {{end}}
                
                <div class="total-price">
                    Total Amount: <span class="highlight">{{formatIDR .TotalPrice}}</span>
                </div>
                
                <div class="shipping-info">
                    <strong>📍 Shipping Address:</strong><br>
                    {{.Address.Address}},<br>
                    {{.Address.City}}, {{.Address.Province}},<br>
                    {{.Address.ZipCode}}
                </div>
            </div>
            
            <p>We hope you're delighted with your purchase! If you have any questions or need assistance, our support team is here to help.</p>
            
            <a href="https://example.com/support" class="button">Contact Support</a>
            
            <p style="font-size: 14px; color: #6c757d; margin-top: 20px;">
                Note: If you did not receive your order or have any concerns, please contact us immediately.
            </p>
        </div>
        
        <div class="footer">
            <p style="margin: 5px 0;">Thank you for shopping with RumahPC! 💚</p>
            <p style="margin: 5px 0;">&copy; 2025 RumahPC. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`
