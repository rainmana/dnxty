// main.go
package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

// DomainTXT holds the DNS TXT record result for a domain.
type DomainTXT struct {
	Domain string `json:"domain" yaml:"domain"`
	TXT    string `json:"txt" yaml:"txt"`
	Key    string `json:"key" yaml:"key"`
	Value  string `json:"value" yaml:"value"`
}

func main() {
	// Define flags with descriptive names and detailed help messages.
	filePath := flag.String("file", "", "Path to a text file containing domain names (one domain per line).")
	outputFormat := flag.String("format", "pretty", "Output format. Options: pretty (default), json, yaml, csv.")

	// Override the default Usage function to provide more context.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n dnxty - A DNS TXT Record Extraction Utility\n\n")
		fmt.Fprintf(os.Stderr, " Usage:\n")
		fmt.Fprintf(os.Stderr, "   %s [options] domain1 [domain2 ...]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, " Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n Examples:\n")
		fmt.Fprintf(os.Stderr, "   %s google.com facebook.com\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "   %s --file domains.txt --format json\n\n", os.Args[0])
	}

	flag.Parse()

	// Gather domains from file (if provided) and from command-line arguments.
	var domains []string
	if *filePath != "" {
		file, err := os.Open(*filePath)
		if err != nil {
			color.Red("Error opening file %s: %v", *filePath, err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				domains = append(domains, line)
			}
		}
		if err := scanner.Err(); err != nil {
			color.Red("Error reading file %s: %v", *filePath, err)
			os.Exit(1)
		}
	}

	// Append domains provided as positional arguments.
	domains = append(domains, flag.Args()...)
	if len(domains) == 0 {
		color.Yellow("No domains provided. Please supply domains as arguments or via the --file flag.\n")
		flag.Usage()
		os.Exit(1)
	}

	// Prepare a slice to store results.
	var results []DomainTXT

	// Regex to capture key=value pairs (commonly used for domain verification).
	re := regexp.MustCompile(`([\w\.\-]+)=([A-Za-z0-9\+\/=]+)`)

	// For each domain, perform a DNS TXT lookup.
	for _, domain := range domains {
		txtRecords, err := net.LookupTXT(domain)
		if err != nil {
			color.Red("Error looking up TXT records for %s: %v", domain, err)
			continue
		}

		// Process each TXT record.
		for _, txt := range txtRecords {
			key := ""
			value := ""
			match := re.FindStringSubmatch(txt)
			if len(match) == 3 {
				key = match[1]
				value = match[2]
			}
			results = append(results, DomainTXT{
				Domain: domain,
				TXT:    txt,
				Key:    key,
				Value:  value,
			})
		}
	}

	// Output the results in the chosen format.
	switch strings.ToLower(*outputFormat) {
	case "pretty":
		printPretty(results)
	case "json":
		printJSON(results)
	case "yaml":
		printYAML(results)
	case "csv":
		printCSV(results)
	default:
		color.Yellow("Unknown output format '%s'. Defaulting to pretty.", *outputFormat)
		printPretty(results)
	}
}

// printPretty outputs results in a color-coded, well-aligned table.
func printPretty(results []DomainTXT) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Domain", "TXT Record", "Key", "Value"})

	headerColors := []tablewriter.Colors{
		{tablewriter.FgHiBlueColor, tablewriter.Bold},
		{tablewriter.FgHiBlueColor, tablewriter.Bold},
		{tablewriter.FgHiBlueColor, tablewriter.Bold},
		{tablewriter.FgHiBlueColor, tablewriter.Bold},
	}
	table.SetHeaderColor(headerColors...)

	for _, r := range results {
		table.Append([]string{r.Domain, r.TXT, r.Key, r.Value})
	}
	table.Render()
}

// printJSON outputs results in JSON format.
func printJSON(results []DomainTXT) {
	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		color.Red("Error marshalling JSON: %v", err)
		return
	}
	fmt.Println(string(b))
}

// printYAML outputs results in YAML format.
func printYAML(results []DomainTXT) {
	b, err := yaml.Marshal(results)
	if err != nil {
		color.Red("Error marshalling YAML: %v", err)
		return
	}
	fmt.Println(string(b))
}

// printCSV outputs results in CSV format.
func printCSV(results []DomainTXT) {
	writer := csv.NewWriter(os.Stdout)
	if err := writer.Write([]string{"Domain", "TXT Record", "Key", "Value"}); err != nil {
		color.Red("Error writing CSV header: %v", err)
		return
	}
	for _, r := range results {
		if err := writer.Write([]string{r.Domain, r.TXT, r.Key, r.Value}); err != nil {
			color.Red("Error writing CSV row: %v", err)
			return
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		color.Red("Error flushing CSV: %v", err)
	}
}
