package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"designs.capital/dogepool/config"
	"github.com/gorilla/mux"
)

type EnhancedAPIServer struct {
	router *mux.Router
	config *config.Config
	port   string
}

type BlockchainSyncStatus struct {
	Chain              string  `json:"chain"`
	Blocks             int64   `json:"blocks"`
	Headers            int64   `json:"headers"`
	SyncProgress       float64 `json:"sync_progress"`
	IsSyncing          bool    `json:"is_syncing"`
	VerificationProgress float64 `json:"verification_progress"`
}

type PoolStats struct {
	PoolName       string                  `json:"pool_name"`
	Hashrate       float64                 `json:"hashrate"`
	Miners         int                     `json:"miners"`
	Workers        int                     `json:"workers"`
	BlocksFound    int                     `json:"blocks_found"`
	LastBlockTime  string                  `json:"last_block_time"`
	NetworkDiff    map[string]float64      `json:"network_difficulty"`
	PoolDiff       float64                 `json:"pool_difficulty"`
	SyncStatus     []BlockchainSyncStatus  `json:"sync_status"`
	Chains         []string                `json:"chains"`
	MergedMining   bool                    `json:"merged_mining"`
}

type MinerStats struct {
	Address        string             `json:"address"`
	Hashrate       float64            `json:"hashrate"`
	SharesValid    int64              `json:"shares_valid"`
	SharesInvalid  int64              `json:"shares_invalid"`
	LastShare      string             `json:"last_share"`
	Balance        map[string]float64 `json:"balance"`
	Paid           map[string]float64 `json:"paid"`
	Workers        []WorkerStats      `json:"workers"`
}

type WorkerStats struct {
	Name          string  `json:"name"`
	Hashrate      float64 `json:"hashrate"`
	SharesValid   int64   `json:"shares_valid"`
	SharesInvalid int64   `json:"shares_invalid"`
	LastShare     string  `json:"last_share"`
}

type Payment struct {
	TxID      string  `json:"txid"`
	Amount    float64 `json:"amount"`
	Chain     string  `json:"chain"`
	Address   string  `json:"address"`
	Timestamp string  `json:"timestamp"`
}

type BlockInfo struct {
	Height       uint64  `json:"height"`
	Hash         string  `json:"hash"`
	Chain        string  `json:"chain"`
	Difficulty   float64 `json:"difficulty"`
	Reward       float64 `json:"reward"`
	Confirmations int    `json:"confirmations"`
	Status       string  `json:"status"`
	Time         string  `json:"time"`
	Miner        string  `json:"miner"`
}

func NewEnhancedAPIServer(cfg *config.Config) *EnhancedAPIServer {
	router := mux.NewRouter()
	server := &EnhancedAPIServer{
		router: router,
		config: cfg,
		port:   cfg.API.Port,
	}

	server.setupRoutes()
	return server
}

func (s *EnhancedAPIServer) setupRoutes() {
	// Enable CORS
	s.router.Use(corsMiddleware)

	// Pool endpoints
	s.router.HandleFunc("/api/pool/stats", s.GetPoolStats).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/pool/blocks", s.GetPoolBlocks).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/pool/payments", s.GetPoolPayments).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/pool/miners", s.GetPoolMiners).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/pool/sync", s.GetSyncStatus).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/pool/hashrate", s.GetPoolHashrate).Methods("GET", "OPTIONS")

	// Miner endpoints
	s.router.HandleFunc("/api/miner/{address}/stats", s.GetMinerStats).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/miner/{address}/payments", s.GetMinerPayments).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/miner/{address}/blocks", s.GetMinerBlocks).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/miner/{address}/workers", s.GetMinerWorkers).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/miner/{address}/balance", s.GetMinerBalance).Methods("GET", "OPTIONS")

	// Chain-specific endpoints
	s.router.HandleFunc("/api/chain/{chain}/stats", s.GetChainStats).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/chain/{chain}/blocks", s.GetChainBlocks).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/chain/{chain}/difficulty", s.GetChainDifficulty).Methods("GET", "OPTIONS")

	// Network endpoints
	s.router.HandleFunc("/api/network/hashrate", s.GetNetworkHashrate).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/network/difficulty", s.GetNetworkDifficulty).Methods("GET", "OPTIONS")

	// Live data endpoints (for dashboard)
	s.router.HandleFunc("/api/live/hashrate", s.GetLiveHashrate).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/live/workers", s.GetLiveWorkers).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/live/shares", s.GetLiveShares).Methods("GET", "OPTIONS")

	// Configuration endpoints
	s.router.HandleFunc("/api/config/chains", s.GetConfiguredChains).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/config/fees", s.GetPoolFees).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/api/config/ports", s.GetPoolPorts).Methods("GET", "OPTIONS")

	log.Printf("Enhanced API routes configured")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *EnhancedAPIServer) GetPoolStats(w http.ResponseWriter, r *http.Request) {
	stats := PoolStats{
		PoolName:     s.config.PoolName,
		Hashrate:     getCurrentHashrate(),
		Miners:       getActiveMinerCount(),
		Workers:      getActiveWorkerCount(),
		BlocksFound:  getTotalBlocksFound(),
		LastBlockTime: getLastBlockTime(),
		NetworkDiff:  getNetworkDifficulties(s.config),
		PoolDiff:     s.config.PoolDifficulty,
		SyncStatus:   getSyncStatus(s.config),
		Chains:       s.config.BlockChainOrder,
		MergedMining: len(s.config.BlockChainOrder) > 1,
	}

	respondJSON(w, stats)
}

func (s *EnhancedAPIServer) GetSyncStatus(w http.ResponseWriter, r *http.Request) {
	syncStatus := getSyncStatus(s.config)
	respondJSON(w, syncStatus)
}

func (s *EnhancedAPIServer) GetPoolBlocks(w http.ResponseWriter, r *http.Request) {
	limit := getQueryInt(r, "limit", 50)
	offset := getQueryInt(r, "offset", 0)

	blocks := getRecentBlocks(limit, offset)
	respondJSON(w, blocks)
}

func (s *EnhancedAPIServer) GetPoolPayments(w http.ResponseWriter, r *http.Request) {
	limit := getQueryInt(r, "limit", 50)
	offset := getQueryInt(r, "offset", 0)

	payments := getRecentPayments(limit, offset)
	respondJSON(w, payments)
}

func (s *EnhancedAPIServer) GetPoolMiners(w http.ResponseWriter, r *http.Request) {
	miners := getActiveMiners()
	respondJSON(w, miners)
}

func (s *EnhancedAPIServer) GetPoolHashrate(w http.ResponseWriter, r *http.Request) {
	duration := r.URL.Query().Get("duration")
	if duration == "" {
		duration = "1h"
	}

	hashrate := getHashrateHistory(duration)
	respondJSON(w, hashrate)
}

func (s *EnhancedAPIServer) GetMinerStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	stats := getMinerStats(address)
	respondJSON(w, stats)
}

func (s *EnhancedAPIServer) GetMinerPayments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	limit := getQueryInt(r, "limit", 50)

	payments := getMinerPayments(address, limit)
	respondJSON(w, payments)
}

func (s *EnhancedAPIServer) GetMinerBlocks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	limit := getQueryInt(r, "limit", 50)

	blocks := getMinerBlocks(address, limit)
	respondJSON(w, blocks)
}

func (s *EnhancedAPIServer) GetMinerWorkers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	workers := getMinerWorkers(address)
	respondJSON(w, workers)
}

func (s *EnhancedAPIServer) GetMinerBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	balance := getMinerBalance(address)
	respondJSON(w, balance)
}

func (s *EnhancedAPIServer) GetChainStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chain := vars["chain"]

	stats := getChainStats(chain)
	respondJSON(w, stats)
}

func (s *EnhancedAPIServer) GetChainBlocks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chain := vars["chain"]
	limit := getQueryInt(r, "limit", 50)

	blocks := getChainBlocks(chain, limit)
	respondJSON(w, blocks)
}

func (s *EnhancedAPIServer) GetChainDifficulty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chain := vars["chain"]

	difficulty := getChainDifficulty(chain)
	respondJSON(w, map[string]interface{}{
		"chain":      chain,
		"difficulty": difficulty,
	})
}

func (s *EnhancedAPIServer) GetNetworkHashrate(w http.ResponseWriter, r *http.Request) {
	hashrates := getNetworkHashrates(s.config)
	respondJSON(w, hashrates)
}

func (s *EnhancedAPIServer) GetNetworkDifficulty(w http.ResponseWriter, r *http.Request) {
	difficulties := getNetworkDifficulties(s.config)
	respondJSON(w, difficulties)
}

func (s *EnhancedAPIServer) GetLiveHashrate(w http.ResponseWriter, r *http.Request) {
	hashrate := getCurrentHashrate()
	respondJSON(w, map[string]float64{"hashrate": hashrate})
}

func (s *EnhancedAPIServer) GetLiveWorkers(w http.ResponseWriter, r *http.Request) {
	workers := getActiveWorkerCount()
	respondJSON(w, map[string]int{"workers": workers})
}

func (s *EnhancedAPIServer) GetLiveShares(w http.ResponseWriter, r *http.Request) {
	shares := getLiveShares()
	respondJSON(w, shares)
}

func (s *EnhancedAPIServer) GetConfiguredChains(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, s.config.BlockChainOrder)
}

func (s *EnhancedAPIServer) GetPoolFees(w http.ResponseWriter, r *http.Request) {
	fees := make(map[string]interface{})
	for chain, payoutConfig := range s.config.Payouts.Chains {
		fees[chain] = payoutConfig.PoolRewardRecipients
	}
	respondJSON(w, fees)
}

func (s *EnhancedAPIServer) GetPoolPorts(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{
		"stratum": s.config.Port,
		"api":     s.config.API.Port,
	})
}

func (s *EnhancedAPIServer) Start() error {
	addr := ":" + s.port
	log.Printf("Enhanced API server starting on %s", addr)
	return http.ListenAndServe(addr, s.router)
}

// Helper functions

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func getQueryInt(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getSyncStatus(cfg *config.Config) []BlockchainSyncStatus {
	// This would query each daemon for sync status
	// Placeholder implementation
	status := make([]BlockchainSyncStatus, 0)
	for _, chain := range cfg.BlockChainOrder {
		status = append(status, BlockchainSyncStatus{
			Chain:                chain,
			Blocks:               100000,
			Headers:              100000,
			SyncProgress:         100.0,
			IsSyncing:            false,
			VerificationProgress: 100.0,
		})
	}
	return status
}

func getCurrentHashrate() float64 {
	// Query from persistence
	return 1000000.0 // 1 MH/s placeholder
}

func getActiveMinerCount() int {
	// Query from persistence
	return 10 // Placeholder
}

func getActiveWorkerCount() int {
	// Query from persistence
	return 25 // Placeholder
}

func getTotalBlocksFound() int {
	// Query from persistence
	return 100 // Placeholder
}

func getLastBlockTime() string {
	// Query from persistence
	return time.Now().Add(-time.Hour).Format(time.RFC3339)
}

func getNetworkDifficulties(cfg *config.Config) map[string]float64 {
	diff := make(map[string]float64)
	for _, chain := range cfg.BlockChainOrder {
		diff[chain] = 1000000.0 // Placeholder
	}
	return diff
}

func getNetworkHashrates(cfg *config.Config) map[string]float64 {
	hashrates := make(map[string]float64)
	for _, chain := range cfg.BlockChainOrder {
		hashrates[chain] = 10000000000.0 // 10 GH/s placeholder
	}
	return hashrates
}

func getRecentBlocks(limit, offset int) []BlockInfo {
	// Query from persistence
	return []BlockInfo{} // Placeholder
}

func getRecentPayments(limit, offset int) []Payment {
	// Query from persistence
	return []Payment{} // Placeholder
}

func getActiveMiners() []string {
	// Query from persistence
	return []string{} // Placeholder
}

func getHashrateHistory(duration string) map[string]float64 {
	// Query from persistence
	return make(map[string]float64) // Placeholder
}

func getMinerStats(address string) MinerStats {
	// Query from persistence
	return MinerStats{
		Address:       address,
		Hashrate:      100000.0,
		SharesValid:   1000,
		SharesInvalid: 10,
		LastShare:     time.Now().Format(time.RFC3339),
		Balance:       make(map[string]float64),
		Paid:          make(map[string]float64),
		Workers:       []WorkerStats{},
	}
}

func getMinerPayments(address string, limit int) []Payment {
	return []Payment{} // Placeholder
}

func getMinerBlocks(address string, limit int) []BlockInfo {
	return []BlockInfo{} // Placeholder
}

func getMinerWorkers(address string) []WorkerStats {
	return []WorkerStats{} // Placeholder
}

func getMinerBalance(address string) map[string]float64 {
	return make(map[string]float64) // Placeholder
}

func getChainStats(chain string) map[string]interface{} {
	return map[string]interface{}{
		"chain":      chain,
		"hashrate":   1000000.0,
		"difficulty": 1000000.0,
		"blocks":     100,
	}
}

func getChainBlocks(chain string, limit int) []BlockInfo {
	return []BlockInfo{} // Placeholder
}

func getChainDifficulty(chain string) float64 {
	return 1000000.0 // Placeholder
}

func getLiveShares() map[string]int64 {
	return map[string]int64{
		"valid":   1000,
		"invalid": 10,
	}
}
