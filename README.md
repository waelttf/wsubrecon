# wsubrecon 🔍


`wsubrecon` is a powerful subdomain enumeration and live probing tool written in Go. It combines multiple reconnaissance tools and data sources to comprehensively discover subdomains and identify live hosts.

## Features ✨

- **Comprehensive Subdomain Enumeration** using:
  - [Subfinder](https://github.com/projectdiscovery/subfinder)
  - [Assetfinder](https://github.com/tomnomnom/assetfinder)
  - crt.sh certificate transparency logs
  - [Shosubgo](https://github.com/incogbyte/shosubgo) (Shodan integration)
  
- **Smart Processing**:
  - Merges and deduplicates results from all sources
  - Filters live subdomains using [httpx](https://github.com/projectdiscovery/httpx)
  
- **Organized Output**:
  - Creates dedicated output directory
  - Saves raw and processed results
  - Generates clean final list of live subdomains

## Installation 📥

### Prerequisites

Ensure these tools are installed and available in your `PATH`:

- [Go](https://golang.org/dl/) (v1.20+ recommended)
- [Subfinder](https://github.com/projectdiscovery/subfinder#installation)
- [Assetfinder](https://github.com/tomnomnom/assetfinder#installation)
- [Shosubgo](https://github.com/incogbyte/shosubgo#installation)
- [httpx](https://github.com/projectdiscovery/httpx#installation)

**Note:** You'll need a valid [Shodan API key](https://account.shodan.io/) for Shosubgo functionality.

### Building from Source

```bash
git clone https://github.com/waelttf/wsubrecon.git
cd wsubrecon
go build -o wsubrecon main.go
```

## Usage 🚀

### Basic Command

```bash
go mod init wsubrecon
go build -o wsubrecon wsubrecon.go
./wsubrecon -d example.com
```


