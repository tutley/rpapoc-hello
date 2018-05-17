# helloUser

This is a part of the reverse proxy authentication proof of concept.

This is a simple application to listen on localhost:port and respond to http requests with a hello json message, including the user information, which will be fetched from a database based on the username, which is retrieved from the http header as inserted by a proxy server.

so http header "Token-Claim-Username" will contain the username, which will be used to fetch data from the database.

Run this app with go installed:
LISTEN_PORT=9000 SERVICE_NAME=TestService DB_HOST=localhost go run main.go

The idea is to contanerize this:
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w' .

### API

GET / (must have Token-Claim-Username header)
