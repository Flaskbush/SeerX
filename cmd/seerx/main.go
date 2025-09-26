package main

import (
	"flag"
	"fmt"

	"github.com/flaskbush/seerx/exploits"
	"github.com/flaskbush/seerx/exploits/vsftpd"
	"github.com/flaskbush/seerx/recon"
)

func main() {
	// Command line flags
	target := flag.String("target", "192.168.88.112", "Target IP")
	exploitName := flag.String("exploit", "", "Exploit name")
	scan := flag.Bool("scan", false, "Perform port scanning")
	portRange := flag.String("ports", "1-1000", "Port range to scan (e.g., 1-1000)")
	flag.Parse()

	// Handle port scanning
	if *scan {
		scanner := recon.NewPortScanner()
		fmt.Printf("Scanning ports %s on %s...\n", *portRange, *target)
		openPorts, err := scanner.Scan(*target, *portRange)
		if err != nil {
			fmt.Printf("Scanning error: %v\n", err)
			return
		}
		fmt.Println("Open ports:")
		for _, port := range openPorts {
			fmt.Println(port)
		}
		return
	}

	// Handle exploit execution
	if *exploitName != "" {
		manager := exploits.NewExploitManager()
		manager.Register(vsftpd.NewVsftpdExploit("21", "6200"))

		err := manager.RunExploit(*exploitName, *target)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}

	// If no mode specified
	fmt.Println("Please specify either -scan or -exploit")
	flag.Usage()
}
