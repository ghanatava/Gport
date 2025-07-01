package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ghanatava/Gport/internal/scanner"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "Gport",
	Short: "Gport is a CLI tool for managing ports.",
	Run:   runScanner,
}

var (
	host        string
	portRange   string
	timeout     time.Duration
	concurrency int
)

func init() {
	RootCmd.Flags().StringVarP(&host, "host", "H", "", "target host to connect (request)")
	RootCmd.Flags().StringVarP(&portRange, "ports", "p", "1-1024", "Port range to scan (e.g., 22,80,443 or 1-65535)")
	RootCmd.Flags().DurationVarP(&timeout, "timeout", "t", 1*time.Second, "Timeout per port scan")
	RootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 100, "Number of concurrent workers")
	RootCmd.MarkFlagRequired("host")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runScanner(cmd *cobra.Command, args []string) {
	ports, err := parsePortRange(portRange)
	if err != nil {
		fmt.Printf("Invalid port range: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting scan on %s for ports %s...\n", host, portRange)
	results := scanner.ScanPorts(host, ports, timeout, concurrency)
	fmt.Println("\n--- Scan Results ---")
	for _, res := range results {
		if res.Open {
			fmt.Printf("[+] Port %d OPEN\n", res.Port)
		}
	}
}

func parsePortRange(input string) ([]int, error) {
	var ports []int

	ranges := strings.Split(input, ",")
	for _, r := range ranges {
		if strings.Contains(r, "-") {
			parts := strings.SplitN(r, "-", 2)
			start, err1 := strconv.Atoi(parts[0])
			end, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil || start > end {
				return nil, fmt.Errorf("invalid range: %s", r)
			}
			for p := start; p <= end; p++ {
				ports = append(ports, p)
			}
		} else {
			p, err := strconv.Atoi(r)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", r)
			}
			ports = append(ports, p)
		}
	}

	return ports, nil
}
