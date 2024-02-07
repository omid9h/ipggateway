wrk.method = "POST"
wrk.body   = '{"InvoiceNumber": "INV123456","InvoiceDate": "2024-01-06","TerminalCode": "1234567893","MerchantCode": "789012","Amount": 12345,"RedirectAddress": "https://example.com/redirect","Timestamp": "2024-01-06T12:34:56Z","Action": 1}'
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Sign"] = "your-signature-header-value"
