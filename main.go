// main.go
package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/quick"
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
	// Define command-line flags.
	filePath := flag.String("file", "", "Path to a text file containing domain names (one domain per line).")
	outputFormat := flag.String("format", "pretty", "Output format. Options: pretty (default), json, yaml, csv.")
	noColor := flag.Bool("no-color", false, "Disable colored output and syntax highlighting.")
	allRecords := flag.Bool("all", false, "Include all TXT records, even those without a valid key/value pair.")
	// Default behavior: ignore SPF records unless --include-spf is provided.
	includeSPF := flag.Bool("include-spf", false, "Include SPF TXT records (records starting with 'v=spf1'). By default, SPF records are ignored.")

	// Override the default Usage function with a Typer-inspired help interface.
	flag.Usage = func() {
		// If no-color is enabled, force no colors for the help.
		if *noColor {
			color.NoColor = true
		}
		// Define styled printers.
		header := color.New(color.FgCyan, color.Bold)
		example := color.New(color.FgYellow)
		fmt.Fprintf(os.Stderr, "\n")
		header.Fprintf(os.Stderr, "dnxty - A DNS TXT Record Extraction Utility\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  %s [options] domain1 [domain2 ...]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		example.Fprintf(os.Stderr, "  %s google.com facebook.com\n", os.Args[0])
		example.Fprintf(os.Stderr, "  %s --file domains.txt --format json\n", os.Args[0])
		example.Fprintf(os.Stderr, "  %s --all google.com\n", os.Args[0])
		example.Fprintf(os.Stderr, "  %s --include-spf google.com\n\n", os.Args[0])
	}

	flag.Parse()

	// Set color.NoColor based on the flag.
	color.NoColor = *noColor

	// Gather domains from file (if provided) and from positional arguments.
	var domains []string
	if *filePath != "" {
		f, err := os.Open(*filePath)
		if err != nil {
			color.Red("Error opening file %s: %v", *filePath, err)
			os.Exit(1)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
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
	// Append any domains provided as positional arguments.
	domains = append(domains, flag.Args()...)
	if len(domains) == 0 {
		color.Yellow("No domains provided. Please supply domains as arguments or via the --file flag.\n")
		flag.Usage()
		os.Exit(1)
	}

	// Prepare to store results.
	var results []DomainTXT

	// Compile a regex to capture key=value pairs (commonly used for domain verification).
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
			// By default, ignore SPF records (those starting with "v=spf1") unless --include-spf is set.
			if !*includeSPF && strings.HasPrefix(strings.ToLower(txt), "v=spf1") {
				continue
			}
			key := ""
			value := ""
			match := re.FindStringSubmatch(txt)
			if len(match) == 3 {
				key = match[1]
				value = match[2]
			}
			// If the --all flag is not set, only include records with a valid key/value pair.
			if !*allRecords && key == "" && value == "" {
				continue
			}
			results = append(results, DomainTXT{
				Domain: domain,
				TXT:    txt,
				Key:    key,
				Value:  value,
			})
		}
	}

	// Output results in the chosen format.
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

// printPretty outputs the results as a formatted table with colored headers.
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

// printJSON outputs results in JSON format with syntax highlighting.
func printJSON(results []DomainTXT) {
	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		color.Red("Error marshalling JSON: %v", err)
		return
	}
	jsonStr := string(b)
	if !color.NoColor {
		// Use chroma to syntax-highlight the JSON.
		if err := quick.Highlight(os.Stdout, jsonStr, "json", "terminal", "monokai"); err != nil {
			// Fallback to plain printing if highlighting fails.
			fmt.Println(jsonStr)
		}
	} else {
		fmt.Println(jsonStr)
	}
}

// printYAML outputs results in YAML format with syntax highlighting.
func printYAML(results []DomainTXT) {
	b, err := yaml.Marshal(results)
	if err != nil {
		color.Red("Error marshalling YAML: %v", err)
		return
	}
	yamlStr := string(b)
	if !color.NoColor {
		if err := quick.Highlight(os.Stdout, yamlStr, "yaml", "terminal", "monokai"); err != nil {
			fmt.Println(yamlStr)
		}
	} else {
		fmt.Println(yamlStr)
	}
}

// printCSV outputs results in CSV format. We generate CSV output into a buffer,
// then highlight it if colors are enabled.
func printCSV(results []DomainTXT) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	// Write CSV header.
	if err := writer.Write([]string{"Domain", "TXT Record", "Key", "Value"}); err != nil {
		color.Red("Error writing CSV header: %v", err)
		return
	}
	// Write CSV rows.
	for _, r := range results {
		if err := writer.Write([]string{r.Domain, r.TXT, r.Key, r.Value}); err != nil {
			color.Red("Error writing CSV row: %v", err)
			return
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		color.Red("Error flushing CSV: %v", err)
		return
	}
	csvStr := buf.String()
	if !color.NoColor {
		// Attempt syntax highlighting using the "csv" lexer.
		if err := quick.Highlight(os.Stdout, csvStr, "csv", "terminal", "monokai"); err != nil {
			fmt.Println(csvStr)
		}
	} else {
		fmt.Println(csvStr)
	}
}

