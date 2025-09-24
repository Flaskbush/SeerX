# SeerX
SeersX is a modular Command &amp; Control (C2) framework built for stealth, scalability, and precision. Featuring the SeerAgent, it enables secure remote operations, dynamic payload delivery, and adaptive control across diverse environments.
Designed for red team simulations and advanced adversary emulation, SeersX is currently tested and demonstrated in a controlled lab environment using Metasploitable 2 virtual machines as primary targets.



SeerX/
├── cmd/                    # CLI entry points
│   └── seerx/main.go       # Main CLI (launches scans, exploits, C2)
├── c2/                     # C2 server and agent management
│   ├── server/             # C2 server code
│   │   ├── server.go       # Main server logic
│   │   └── api.go          # HTTP API for agents and dashboard
│   ├── agent/              # Future agents (placeholders for C/C++)
│   │   └── agent.cpp       # C++ Agent (coming soon)
│   └── models.go           # Data structures (e.g., Agent, Session)
├── exploits/               # Exploit modules
│   ├── vsftpd/             # vsftpd exploit
│   │   └── vsftpd.go
│   └── interface.go        # Exploit interface
├── recon/                  # Reconnaissance tools
│   └── portscanner.go
├── go.mod                  # Go dependencies
└── README.md               # Project documentation
