package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"path/filepath"
)

func main() {
	arguments := os.Args
	programName := filepath.Base(arguments[0])
	breaker := strings.Repeat("-", 30)
	// Checks if the number of arguments is correct
	if len(arguments) != 3 {
		fmt.Printf("%s\nError: Invalid number of arguments\n\n", breaker)
		fmt.Printf("Example Usage: %s dns-query <domain-name>\n", programName)
		fmt.Printf("%s\n", breaker)
		return
	}
	command := arguments[1]
	domainName := arguments[2]
	ips, err := net.LookupIP(domainName)

	// Checks if the command is dns-query
	if command != "--query" {
		fmt.Printf("%s\n", breaker)
		fmt.Printf("Unknown command: %s\n\nPlease use: --query\n", command)
		fmt.Printf("%s\n", breaker)
		return
	}

	// net.LookupIP Success (err == nil): Program proceeds to print IP addresses
	// net.LookupIP Failure (err != nil):
	// 				Is it a "Domain Not Found" error specifically?
	// 				YES Print the specific "Invalid domain name message"
	// 				NO (Some other DNS error): Print the general errord
	if err != nil {
		dnsError, ok := err.(*net.DNSError)
		if ok && dnsError.IsNotFound {
			fmt.Printf("%s\nError: Invalid domain name: '%s'\n\n", breaker, domainName)
			fmt.Println("Please check the domain name and try again.")
			fmt.Printf("%s", breaker)
			return
		}
		fmt.Printf("%s\nError looking up IP addresses for '%s: %v\n", breaker, domainName, err)
		fmt.Printf("Please check your network connection or DNS settings.\n")
		fmt.Printf("%s\n", breaker)
		return
	}

	// Prints the IP addresses
	fmt.Printf("\n%s Results for %s %s\n", breaker, domainName, breaker) // Header with domain
		foundIPv4 := false
		foundIPv6 := false

		for _, ip := range ips {
			if ipv4 := ip.To4(); ipv4 != nil {
				if !foundIPv4 {
					fmt.Println("\nIPv4 Addresses:")
					foundIPv4 = true
				}
				fmt.Println("-", ipv4.String())
			} else if ipv6 := ip.To16(); ipv6 != nil {
				if !foundIPv6 {
					fmt.Println("\nIPv6 Addresses:")
					foundIPv6 = true
				}
				fmt.Println("-", ipv6.String())
			}
		}
}


