package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func printARecordsOutput(domainName, recordType string, ips []net.IP, breaker string) {
	fmt.Printf("%s %s Records for %s %s\n", breaker, recordType, domainName, breaker)

	foundIPv4 := false
	foundIPv6 := false
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			if !foundIPv4 {
				fmt.Println("\nIPV4 Addresses:")
				foundIPv4 = true
			}
			fmt.Printf("  - IP Address: %-20s\n", ipv4.String())
		} else if ipv6 := ip.To16(); ipv6 != nil {
			if !foundIPv6 {
				fmt.Println("\nIPv6 Addresses:")
				foundIPv6 = true
			}
			fmt.Printf("  - IPv6 Address: %-40s\n", ipv6.String())
		}
	}
	if !foundIPv4 && !foundIPv6 {
		fmt.Println("\n  No IPv4 or IPv6 addresses found.")
	}
	fmt.Printf("%s\n", breaker)
}

func printMXRecordsOutput(domainName string, mxRecords []*net.MX, breaker string) {
	fmt.Printf("%s MX Records for %s %s\n", breaker, domainName, breaker)

	if len(mxRecords) > 0 {
		fmt.Println("\nMail Exchange Records:")
	}

	for _, mx := range mxRecords {
		fmt.Printf("  - Host: %-30s  Priority: %-5d\n", mx.Host, mx.Pref)
	}
	if len(mxRecords) == 0 {
		fmt.Println("\n  No MX records found.")
	}

	fmt.Printf("%s\n", breaker)
}

func printNoRecordsFoundOutput(recordType, domainName, breaker string) {
	fmt.Printf("%s\nNo %s records found for %s\n", breaker, recordType, domainName)
	fmt.Printf("%s\n", breaker)
}

func printErrorOutput(breaker, errorMessage string) {
	fmt.Printf("%s\nError: %s\n%s\n", breaker, errorMessage, breaker)
}

func printInvalidDomainErrorOutput(breaker, domainName string) {
	fmt.Printf("%s\nError: Invalid domain name: '%s'\n", breaker, domainName)
	fmt.Println("Please check the domain name and try again.")
	fmt.Printf("%s\n", breaker)
}

func printUsageErrorOutput(breaker, programName string, supportedRecordTypes []string, unknownCommand string) {
	fmt.Printf("%s\nError: Invalid number of arguments\n\n", breaker)
	fmt.Printf("Example Usage: %s --query <record-type> <domain-name>\n", programName)
	fmt.Printf("\nSupported record types: %s\n", strings.Join(supportedRecordTypes, ", "))
	fmt.Printf("%s\n", breaker)

	if unknownCommand != "" {
		fmt.Printf("%s\n", breaker)
		fmt.Printf("Unknown command: %s\n\nPlease use: --query\n", unknownCommand)
		fmt.Printf("Usage: %s --query <record-type> <domain-name>\n", programName)
		fmt.Printf("Supported record types: %s\n", strings.Join(supportedRecordTypes, ", "))
		fmt.Printf("%s\n", breaker)
	}
}

func printUnsupportedRecordTypeErrorOutput(breaker string, recordTypeStr string, supportedRecordTypes []string, programName string) {
	fmt.Printf("%s\nError: Unsupported record type: '%s'\n\n", breaker, recordTypeStr)
	fmt.Printf("Supported record types are %s\n", strings.Join(supportedRecordTypes, ""))
	fmt.Printf("Usage: %s --query <record-type> <domain-name>\n", programName)
	fmt.Printf("%s\n", breaker)
}

func printInternalErrorOutput(breaker, recordType string) {
	fmt.Printf("%s\nError: Internal error - Unsupported record type '%s' reached printing stage\n", breaker, recordType)
	fmt.Printf("%s\n", breaker)
}

func parseAndValidateArgs(arguments []string, programName string, breaker string) (command string, recordTypeStr string, domainName string, isValid bool) {
	
	if len(arguments) != 4 {
		printUsageErrorOutput(breaker, programName, []string{"A", "AAAA", "MX", "CNAME"}, "")
		return "", "", "", false
	}

	command = arguments[1]
	recordTypeStr = arguments[2]
	domainName = arguments[3]

	if command != "--query" {
		printUsageErrorOutput(breaker, programName, []string{"A", "AAAA", "MX", "CNAME"}, command)
		return "", "", "", false
	}

	recordTypeStrUpper := strings.ToUpper(recordTypeStr)
	recordType := recordTypeStrUpper

	supportedRecordTypes := map[string]bool{
		"A":		true,
		"AAAA":		true,
		"MX":		true,
		"CNAME":	true,
	}

	if !supportedRecordTypes[recordType] {
		printUnsupportedRecordTypeErrorOutput(breaker, recordTypeStr, []string{"A", "AAAA", "MX", "CNAME"}, programName)
		return "", "", "", false
	}

	return command, recordType, domainName, true
}

func main () {
	arguments := os.Args
	programName := filepath.Base(arguments[0])
	breaker := strings.Repeat("-", 30)
	_, recordTypeStr, domainName, isValidArgs := parseAndValidateArgs(arguments, programName, breaker)

	if !isValidArgs {
		return
	}

	fmt.Printf("%s Results for %s query of %s %s\n", breaker, recordTypeStr, domainName, breaker)

	switch recordTypeStr {
	case "A", "AAAA":
		ips, err := net.LookupIP(domainName)
		if err != nil {
			dnsError, ok := err.(*net.DNSError)
			if ok && dnsError.IsNotFound {
				printInvalidDomainErrorOutput(breaker, domainName)
				return
			}
			printErrorOutput(breaker, fmt.Sprintf("Error looking up %s records for '%s': %v. Please check your network connection or DNS settings.", recordTypeStr, domainName, err))
			return
		}
		if len(ips) == 0 {
			printNoRecordsFoundOutput(recordTypeStr, domainName, breaker)
			return
		}
		printARecordsOutput(domainName, recordTypeStr, ips, breaker)
	
	case "MX":
		mxRecords, err := net.LookupMX(domainName)
		if err != nil {
			dnsError, ok := err.(*net.DNSError)
			if ok && dnsError.IsNotFound {
				printInvalidDomainErrorOutput(breaker, domainName)
				return
			}
			printErrorOutput(breaker, fmt.Sprintf("Error looking up MX records for '%s': %v. Please check your network connection or DNS settings for MX records.", domainName, err))
			return
		}
		if len(mxRecords) == 0 {
			printNoRecordsFoundOutput(recordTypeStr, domainName, breaker)
			return
		}
		printMXRecordsOutput(domainName, mxRecords, breaker)
	
	case "CNAME":
		cname, err := net.LookupCNAME(domainName)
		if err != nil {
			dnsError, ok := err.(*net.DNSError)
			if ok && dnsError.IsNotFound {
				printNoRecordsFoundOutput(recordTypeStr, domainName, breaker)
				return
			}
			printErrorOutput(breaker, fmt.Sprintf("Error looking up CNAME record for '%s': %v. Please check your network connection or DNS settings for CNAME records.", domainName, err))
			return
		}
		fmt.Printf("%s CNAME Record for %s %s\n", breaker, domainName, breaker)
		fmt.Println("\nCanonical Name (CNAME):")
		fmt.Printf("- Target: %s\n", cname)
		fmt.Printf("%s\n", breaker)
	
	default:
		printInternalErrorOutput(breaker, recordTypeStr)
		return
	}
	
}
