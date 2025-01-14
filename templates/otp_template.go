package templates

const OTP = `
<!DOCTYPE html>
<html>
<head>
	<style>
		.container {
			max-width: 600px;
			margin: 0 auto;
			padding: 20px;
			font-family: Arial, sans-serif;
		}
		.header {
			background-color: #4CAF50;
			color: white;
			padding: 20px;
			text-align: center;
			border-radius: 5px 5px 0 0;
		}
		.content {
			background-color: #f9f9f9;
			padding: 20px;
			border: 1px solid #ddd;
			border-radius: 0 0 5px 5px;
		}
		.otp-code {
			font-size: 32px;
			font-weight: bold;
			color: #333;
			text-align: center;
			padding: 20px;
			background-color: #fff;
			border: 2px dashed #4CAF50;
			border-radius: 5px;
			margin: 20px 0;
		}
		.footer {
			text-align: center;
			margin-top: 20px;
			color: #666;
			font-size: 12px;
		}
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Your OTP Code</h1>
		</div>
		<div class="content">
			<p>Hello, </p>
			<p>Your one-time password (OTP) for authentication is:</p>
			<div class="otp-code">%s</div>
			<p>This code will expire in 5 minutes. Please do not share this code with anyone.</p>
			<p>If you didn't request this code, please ignore this email.</p>
		</div>
		<div class="footer">
			<p>This is an automated message, please do not reply.</p>
			<p>&copy; %d RumahPC. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
`
