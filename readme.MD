# CSRF, XSRF, "Sea Surf"

### /etc/hosts
`127.0.0.1       localhost domain1.com domain2.com`

### Start Vulnerable Server without CORS
`go run server/main.go`

### Start Vulnerable Server with CORS
`go run server/main.go -cors`

### Start Vulnerable Server with CORS and credential
`go run server/main.go -cors -cred`

### STart Attack Server
`go run attacjer/main.go`

### Refernces
https://www.owasp.org/index.php/Cross-Site_Request_Forgery_%28CSRF%29
https://www.html5rocks.com/en/tutorials/cors/