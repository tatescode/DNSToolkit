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

	// Checks if the number of arguments is correct.
	// If you add a CLI parameter as a variable, make sure to adjust this accordingly
	if len(arguments) != 4 {
		fmt.Printf("%s\nError: Invalid number of arguments\n\n", breaker)
		fmt.Printf("Example Usage: %s --query <type> <domain-name>\n", programName)
		fmt.Printf("\nSupported record types: A, AAAA, MX\n")
		fmt.Printf("%s\n", breaker)
		return
	}

	command := arguments[1]
	recordTypeStr := arguments[2]
	domainName := arguments[3]

	// Checks if the command is --query
	if command != "--query" {
		fmt.Printf("%s\n", breaker)
		fmt.Printf("Unknown command: %s\n\nPlease use: --query\n", command)
		fmt.Printf("Usage: %s --query <record-type> <domain-name>\n", programName)
		fmt.Printf("Supported record types: A, AAAA, MX\n")
		fmt.Printf("%s\n", breaker)
		return
	}

	// Input validation
	recordTypeStrUpper := strings.ToUpper(recordTypeStr)
	recordType := recordTypeStrUpper

	// map of supported types of DNS records
	supportedRecordTypes := map[string]bool{
		"A":	true,
		"AAAA":	true,
		"MX":	true,
	}

	if !supportedRecordTypes[recordType] {
		fmt.Printf("%s\nError: Unsupported record type: '%s'\n\n", breaker, recordTypeStr)
		fmt.Printf("Supported record types are: A, AAAA, MX\n")
		fmt.Printf("Usage: %s --query <record-type> <domain-name>\n", programName)
		fmt.Printf("%s\n", breaker)
		return
	}

	fmt.Printf("%s Results for %s query of %s %s\n", breaker, recordType, domainName, breaker)

	switch recordType {
	case "A", "AAAA": // Added case for "A" and "AAAA" (grouped together)
		ips, err := net.LookupIP(domainName)
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
		if len(ips) == 0 {
			fmt.Printf("%s\nNo %s records found for %s", breaker, recordType, domainName)
			fmt.Printf("%s\n", breaker)
			return
		}
		fmt.Printf("%s %s records for %s %s\n", breaker, recordType, domainName, breaker)
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
		if !foundIPv4 && !foundIPv6 {
			fmt.Println("\nNo IPv4 or IPv4 addresses found.")
		}
		fmt.Printf("%s\n", breaker)

	case "MX":
		mxRecords, err := net.LookupMX(domainName)
		if err != nil {
			dnsError, ok := err.(*net.DNSError)
			if ok && dnsError.IsNotFound {
				fmt.Printf("%s\nError: Invalid domain name: '%s' for MX lookup\n", breaker, domainName) 
				fmt.Println("Please check the domain name and try again.")
				fmt.Printf("%s\n", breaker)
				return
			}
			fmt.Printf("%s\nError looking up MX records for '%s': %v\n", breaker, domainName, err) 
			fmt.Printf("Please check your network connection or DNS settings for MX records.\n")
			fmt.Printf("%s\n", breaker)
			return
		}
		fmt.Printf("%s MX Records for %s %s\n", breaker, domainName, breaker) 
		fmt.Println("\nMail Exchange Records:")
		for _, mx := range mxRecords {
			fmt.Printf("- Host: %-30s Priority: %-5d\n", mx.Host, mx.Pref)
		}
		fmt.Printf("%s\n", breaker)

	default:
		fmt.Printf("%s\nError: Internal error - Unsupported record type '%s' reached printing stage\n", breaker, recordType)
		fmt.Printf("%s\n", breaker)
		return
	}
}


