package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nerock/goblockchain/internal/blockchain"
	"github.com/nerock/goblockchain/internal/wallet"
)

type HTTPError struct {
	Error string `json:"error"`
}

type BlockchainCache interface {
	Get(string) *blockchain.Blockchain
	Put(string, *blockchain.Blockchain)
}

type Server struct {
	port  int
	cache BlockchainCache
}

func New(port int, cache BlockchainCache) *Server {
	return &Server{port: port, cache: cache}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/blockchain", s.Blockchain)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), mux)
}

func (s *Server) getBlockchain() (any, int) {
	bc, err := func(c BlockchainCache) (*blockchain.Blockchain, error) {
		bc := c.Get("blockchain")
		if bc == nil {
			minersWallet, err := wallet.New()
			if err != nil {
				return nil, fmt.Errorf("could not create new wallet: %w", err)
			}

			bc, err = blockchain.New(minersWallet.Address())
			if err != nil {
				return nil, fmt.Errorf("could not create new blockchain: %w", err)
			}

			c.Put("blockchain", bc)
		}

		return bc, nil
	}(s.cache)

	if err != nil {
		return HTTPError{err.Error()}, http.StatusInternalServerError
	}

	return bc, http.StatusOK
}

func (s *Server) Blockchain(w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var res any
	switch r.Method {
	case http.MethodGet:
		res, statusCode = s.getBlockchain()
	default:
		statusCode = http.StatusMethodNotAllowed
	}

	if err := write(w, res, statusCode); err != nil {
		log.Println(err)
	}
}

func write(w http.ResponseWriter, res any, statusCode int) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if res != nil {
		return json.NewEncoder(w).Encode(res)
	}

	return nil
}
