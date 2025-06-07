package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const shodanKey = "" // Hardcoded API key

func runCommand(name string, args []string, outputFile string) error {
	cmd := exec.Command(name, args...)
	outfile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outfile.Close()
	cmd.Stdout = outfile

	// Suppress stderr (banner, logs, etc.)
	cmd.Stderr = io.Discard

	return cmd.Run()
}

func normalizeDomain(url string) string {
	url = strings.TrimSpace(url)
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	return strings.Split(url, "/")[0]
}

func fetchCRTSh(domain, outputFile string) error {
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data []map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}

	domains := make(map[string]bool)
	for _, entry := range data {
		if name, ok := entry["name_value"].(string); ok {
			for _, d := range strings.Split(name, "\n") {
				d = strings.ReplaceAll(d, "*.", "")
				domains[d] = true
			}
		}
	}

	var sorted []string
	for d := range domains {
		sorted = append(sorted, d)
	}
	sort.Strings(sorted)

	return os.WriteFile(outputFile, []byte(strings.Join(sorted, "\n")), 0644)
}

func mergeFiles(files []string, outputFile string) error {
	all := make(map[string]bool)
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			domain := normalizeDomain(scanner.Text())
			if domain != "" {
				all[domain] = true
			}
		}
		f.Close()
	}

	var final []string
	for d := range all {
		final = append(final, d)
	}
	sort.Strings(final)
	return os.WriteFile(outputFile, []byte(strings.Join(final, "\n")), 0644)
}

func main() {
	domain := flag.String("d", "", "Target domain (e.g., example.com)")
	flag.Parse()

	if *domain == "" {
		fmt.Println("[-] Please provide a domain using -d")
		return
	}

	outputDir := *domain + "_output"
	os.MkdirAll(outputDir, 0755)

	fmt.Println("[*] Running Subfinder...")
	runCommand("subfinder", []string{"-d", *domain}, filepath.Join(outputDir, "subfinder.txt"))

	fmt.Println("[*] Running Assetfinder...")
	runCommand("assetfinder", []string{"--subs-only", *domain}, filepath.Join(outputDir, "assetfinder.txt"))

	fmt.Println("[*] Fetching from crt.sh...")
	fetchCRTSh(*domain, filepath.Join(outputDir, "crtsh.txt"))

	fmt.Println("[*] Running Shosubgo...")
	runCommand("shosubgo", []string{"-d", *domain, "-s", shodanKey}, filepath.Join(outputDir, "shosubgo.txt"))

	fmt.Println("[*] Merging results and removing duplicates...")
	allSubFile := filepath.Join(outputDir, "all_subdomains.txt")
	mergeFiles([]string{
		filepath.Join(outputDir, "subfinder.txt"),
		filepath.Join(outputDir, "assetfinder.txt"),
		filepath.Join(outputDir, "crtsh.txt"),
		filepath.Join(outputDir, "shosubgo.txt"),
	}, allSubFile)

	fmt.Println("[*] Probing live subdomains using httpx...")
	httpxCmd := exec.Command("bash", "-c", fmt.Sprintf(`cat %s | httpx -silent -mc 200,301,302,401,403,500`, allSubFile))
	output, err := httpxCmd.Output()
	if err != nil {
		fmt.Println("httpx failed:", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	seen := make(map[string]bool)
	var clean []string
	for _, line := range lines {
		domain := normalizeDomain(line)
		if domain != "" && !seen[domain] {
			seen[domain] = true
			clean = append(clean, domain)
		}
	}

	finalOut := filepath.Join(outputDir, "FINAL_subdomains.txt")
	os.WriteFile(finalOut, []byte(strings.Join(clean, "\n")), 0644)

	fmt.Println("[✓] Done! Live subdomains saved to:", finalOut)
}
