# DNSToolkit

## Current Functionality:
- Using command-line parameters, user is able to make dns requests to the specified domain
- Basic error checking is implemented, such as checking if the right cli parameters are used, the right domain name is valid, etc.

## To-do
- Add a dns record type as the --type flag (A, AAAA, MX, CNAME, TXT, etc.)

## Tests (Depending on what functionality you are testing)
- go run .\main.go --query A google.com   
- go run .\main.go --query AAAA example.com 
- go run .\main.go --query MX google.com
- go run .\main.go --query CNAME example.com 
- go run .\main.go --query TXT google.com
- go run .\main.go --query a google.com
- go run .\main.go --query mx example.com
- go run .\main.go --query invalid-command
- go run .\main.go --query MX
- go run .\main.go

> You could compile and run the binary file, but for testing purpses I am using **go run**       
