# dnxty 🔍

[![Go Version](https://img.shields.io/badge/Go-1.23.3%2B-blue?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge)](LICENSE)

**dnxty** is a lightweight, command‑line DNS TXT Record Extraction Utility designed for OSINT analysts, security researchers, and beginners alike. It performs DNS lookups on domains and extracts key/value pairs from TXT records (commonly used for domain verification), with options for simplified output, colorized output, multiple formats, and advanced filtering.

---

## ✨ Features

- **DNS TXT Record Lookup**: Query domains for TXT records using Go’s native DNS libraries.
- **Key/Value Extraction**: Automatically extract common verification strings (e.g. `google-site-verification`) into user‑friendly keys (e.g. `google`).
- **Simplified Mode**: Use the `--simple` flag to output only the domain and a deduplicated, simplified key.
- **Multiple Output Formats**: Print results as a pretty table, JSON, YAML, or CSV.
- **Color & Syntax Highlighting**: Enjoy vibrant, color‑coded output by default (with the option to disable via `--no-color`).
- **Advanced Filtering**: Skip SPF records by default (unless overridden with `--include-spf`) and choose to output all TXT records if desired.
- **OSINT & Automation Friendly**: Easily combine with other Linux command‑line utilities for advanced filtering and analysis.

---

## 🚀 Installation

### Prerequisites

- [Go](https://golang.org/doc/install) 
- Tested on Go version `go1.23.3 linux/amd64`

### Building from Source

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/yourusername/dnxty.git
   cd dnxty
   ```

2. **Download Dependencies:**

   The dependencies are managed via Go modules. If needed, run:

   ```bash
   go mod tidy
   ```

3. **Compile the Binary:**

   To compile and produce an executable named `dnxty`, run:

   ```bash
   go build -o dnxty main.go
   ```
    
> [!NOTE]
> On Windows, this will produce `dnxty.exe`


---

## ⚙️ Usage

Run **dnxty** with one or more domains as positional arguments, or provide a file of domains. Below are some examples:

### Basic Lookup

```bash
./dnxty google.com
```

### Multiple Domains

```bash
./dnxty google.com facebook.com
```

### Using a Domain List File

```bash
./dnxty --file domains.txt
```

### Output in JSON Format

```bash
./dnxty --file domains.txt --format json
```

### Include All TXT Records (Even Without a Valid Key/Value)

```bash
./dnxty --all google.com
```

### Include SPF Records (Disabled by Default)

```bash
./dnxty --include-spf google.com
```

### Simplified Output (Domain + Simplified Key)

```bash
./dnxty --simple google.com
```

### Advanced Usage with Linux CLI Tools

Pipe the JSON output into [`jq`](https://stedolan.github.io/jq/) for further filtering:

```bash
./dnxty --format json google.com | jq '.[].domain'
```

Combine with `grep` and `awk` for custom filtering:

```bash
./dnxty google.com | grep google | awk '{print $1, $3}'
```

---

## 💡 Use Cases for OSINT Analysts

- **Domain Reconnaissance**: Quickly retrieve and parse TXT records from multiple domains for potential indicators.
- **Verification Check**: Identify and extract domain verification tokens for services like Google, Facebook, and HubSpot.
- **Automation**: Integrate **dnxty** into larger OSINT pipelines (e.g., combining with `curl`, `grep`, `jq`, etc.) for automated monitoring and reconnaissance.

---

## 🛠️ Development

- **Code Style**: The project is written in Go and *attempts* follows standard Go conventions (this is my first, "real", Go project).
- **Contributing**: Feel free to open issues or pull requests if you have ideas for new features, bug fixes, or improvements.
- **Testing**: We welcome tests—if you add features, please include tests to help maintain the code quality.

---

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## 🤝 Acknowledgements

- Inspired by OSINT and security research tools.
- Thanks to the Go community for robust libraries like [`net`](https://golang.org/pkg/net/), [`flag`](https://golang.org/pkg/flag/), and third‑party packages for colorful CLI output.

---

Happy hunting! 🕵️‍♂️  
*— Rainmana --*
