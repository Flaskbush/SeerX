package exploits

import (
	"fmt"
)

// Exploit définit l'interface pour tous les modules d'exploit
type Exploit interface {
	Run(target string) error
	Name() string // Nom de l'exploit (ex. : "vsftpd-2.3.4")
}

// ExploitManager gère les exploits enregistrés
type ExploitManager struct {
	exploits map[string]Exploit
}

// NewExploitManager crée un nouveau manager
func NewExploitManager() *ExploitManager {
	return &ExploitManager{
		exploits: make(map[string]Exploit),
	}
}

// Register enregistre un exploit
func (m *ExploitManager) Register(exploit Exploit) {
	m.exploits[exploit.Name()] = exploit
}

// RunExploit exécute un exploit par son nom
func (m *ExploitManager) RunExploit(name, target string) error {
	exploit, exists := m.exploits[name]
	if !exists {
		return fmt.Errorf("exploit %s non trouvé", name)
	}
	return exploit.Run(target)
}
