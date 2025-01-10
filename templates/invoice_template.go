package templates

const Invoice = `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8" />
		<style>
		p {
			margin: 0;
		}
		body {
			padding: 2rem;
			color: #2d3748;
			font-family: Arial, sans-serif;
		}

		.header {
			margin-bottom: 2rem;
			display: flex;
			justify-content: space-between;
			align-items: center;
		}

		.header img {
			width: 8rem;
			height: 8rem;
		}

		.header h1 {
			font-size: 2rem;
			font-weight: bold;
			letter-spacing: 2px;
			margin: 0;
		}

		.header p {
			font-size: 1.125rem;
			color: #38a169;
			margin: 0;
		}

		.grid {
			display: grid;
			grid-template-columns: 1fr 1fr;
			gap: 1rem;
			margin-bottom: 2rem;
		}

		h2 {
			font-size: 1.125rem;
			font-weight: bold;
			text-transform: uppercase;
			margin-bottom: 0.5rem;
		}

		.table {
			width: 100%;
			border-collapse: collapse;
			margin-bottom: 2rem;
		}

		.table thead tr {
			border-top: 2px solid #2d3748;
			border-bottom: 2px solid #2d3748;
		}

		.table th,
		.table td {
			padding: 1rem;
			text-align: right;
		}

		.table th {
			font-size: 0.875rem;
			font-weight: bold;
			text-transform: uppercase;
		}

		.table td {
			font-size: 0.875rem;
		}

		.table td:first-child,
		.table th:first-child {
			text-align: left;
		}

		.table td p {
			margin: 0;
		}

		.table .product-name {
			color: #38a169;
			font-weight: bold;
		}

		.totals {
			text-align: right;
			margin-bottom: 2rem;
		}

		.totals table {
			display: inline-block;
			font-size: 0.875rem;
		}

		.totals th {
			padding-right: 1rem;
			color: #718096;
			text-align: left;
		}

		.totals td {
			font-weight: bold;
			color: #2d3748;
			text-align: right;
		}

		.payment-info {
			font-size: 0.875rem;
		}

		.payment-info strong {
			font-weight: bold;
		}

		.text-right {
			text-align: right;
		}

		.top {
			vertical-align: top;
		}

		.font-bold {
			font-weight: bold;
		}
		</style>
	</head>
	<body>
		<!-- Header -->
		<div class="header">
		<img src="https://i.imgur.com/cnRXoSU.png" alt="Company Logo" />
		<div class="text-right">
			<h1>INVOICE</h1>
			<p class="font-bold">{{.InvoiceId}}</p>
		</div>
		</div>

		<div class="grid">
		<div>
			<h2>Diterbitkan atas nama</h2>
			<p>Merchant: <span class="font-bold">{{.Merchant}}</span></p>
		</div>
		<div>
			<h2>Untuk</h2>
			<table>
			<tr>
				<td class="top">Pembeli</td>
				<td class="top">:</td>
				<td class="top font-bold">{{.User}}</td>
			</tr>
			<tr>
				<td class="top">Tanggal Pembelian</td>
				<td class="top">:</td>
				<td class="top font-bold">{{.Date}}</td>
			</tr>
			<tr>
				<td class="top">Alamat</td>
				<td class="top">:</td>
				<td class="top">
				<p>
					<span class="font-bold">{{ .ContactName }}</span> ({{
					.ContactNumber }})
				</p>
				<p style="white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ .Address }}</p>
				</td>
			</tr>
			</table>
		</div>
		</div>

		<!-- Product Table -->
		<table class="table">
		<thead>
			<tr>
			<th>Product</th>
			<th>Jumlah</th>
			<th>Harga Satuan</th>
			<th>Total Harga</th>
			</tr>
		</thead>
		<tbody>
			{{range .InvoiceItems}}
			<tr>
			<td>
				<p class="product-name">{{ .Name }}</p>
				<p>Berat: {{ .Weight }} kg</p>
			</td>
			<td>{{.Quantity}}</td>
			<td>{{.Price}}</td>
			<td>{{.SubTotal}}</td>
			</tr>
			{{end}}
		</tbody>
		</table>

		<!-- Totals -->
		<div class="totals">
		<table>
			<tr>
			<th>Total Price:</th>
			<td>{{.TotalPrice}}</td>
			</tr>
		</table>
		</div>

		<!-- Payment Info -->
		<div class="payment-info">
		<p><strong>Payment Method:</strong> {{.PaymentMethod}}</p>
		<p><strong>Courier:</strong> {{.Courier}}</p>
		</div>
	</body>
	</html>



`
