package cli

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/flaskbush/seerx/c2/server"
)

type CLI struct {
	server *server.Server
}

func NewCLI(server *server.Server) *CLI {
	return &CLI{
		server: server,
	}
}

func (c *CLI) executor(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")

	switch blocks[0] {
	case "scan":
		if len(blocks) < 2 {
			fmt.Println("Usage: scan <target>")
			return
		}
		target := blocks[1]
		ports, err := c.server.GetPortScanner().Scan(target, "1-1024")
		if err != nil {
			fmt.Printf("Scan failed: %v\n", err)
			return
		}
		fmt.Printf("Open ports on %s: %v\n", target, ports)

	case "exploit":
		if len(blocks) < 3 {
			fmt.Println("Usage: exploit <name> <traget>")
			return
		}
		name := blocks[1]
		target := blocks[2]
		if err := c.server.GetExploitManager().RunExploit(name, target); err != nil {
			fmt.Printf("Exploit failed: %v\n", err)
			return
		}
		fmt.Printf("Exploit %s\n", name, target)

	case "agents":
		agents := c.server.GetAgents()
		if len(agents) == 0 {
			fmt.Println("No active agents")
			return
		}
		for _, agent := range agents {
			fmt.Printf("Agent ID: %s, IP: %s, Last seen: %s\n",
				agent.ID, agent.IP, agent.LastSeen)
		}
	case "help":
		fmt.Println("Available commands:")
		fmt.Println("  scan <target               - Scan target for open ports")
		fmt.Println("  exploit <name> <target>    - Run exploit against target")
		fmt.Println("  agents                     - List active agents")
		fmt.Println("  help                       - List commands")
		fmt.Println("  exit                       - Exit the CLI")

	case "exit":
		fmt.Println("Goodbye!")
		return

	default:
		fmt.Println("Unknown command. Type 'help' for available commands")
	}
}

func (c *CLI) completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "scan", Description: "Scan target for open ports"},
		{Text: "exploit", Description: "Run exploit against target"},
		{Text: "agents", Description: "List active agents"},
		{Text: "help", Description: "Show help"},
		{Text: "exit", Description: "Exit the CLI"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func (c *CLI) Start() {
	fmt.Println("Welcome to SeerX C2 Server")
	fmt.Println("Type 'help' for available commands")
	p := prompt.New(
		c.executor,
		c.completer,
		prompt.OptionPrefix("seerx"),
		prompt.OptionTitle("Seerx CLI"),
	)
	p.Run()
}
