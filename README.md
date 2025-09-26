# SeerX
SeersX is a modular Command &amp; Control (C2) framework built for stealth, scalability, and precision. Featuring the SeerAgent, it enables secure remote operations, dynamic payload delivery, and adaptive control across diverse environments.
Designed for red team simulations and advanced adversary emulation, SeersX is currently tested and demonstrated in a controlled lab environment using Metasploitable 2 virtual machines as primary targets.

## Installation
build:
    go mod tidy
    go build -o seerx ./cmd/seerx

run:
    ./seerx -scan
    ./seerx -exploit vsftpd
