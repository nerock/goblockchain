package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"

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
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/blockchain", func(r chi.Router) {
		r.Get("/", s.getBlockchain)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), r)
}

func (s *Server) getBlockchain(w http.ResponseWriter, r *http.Request) {
	bc := s.cache.Get("blockchain")
	if bc == nil {
		minersWallet, err := wallet.New()
		if err != nil {
			http.Error(w, fmt.Sprintf("could not create new wallet: %s", err), http.StatusInternalServerError)
			return
		}

		bc, err = blockchain.New(minersWallet.Address())
		if err != nil {
			http.Error(w, fmt.Sprintf("could not create new blockchain: %s", err), http.StatusInternalServerError)
			return
		}

		s.cache.Put("blockchain", bc)
	}

	if err := json.NewEncoder(w).Encode(bc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
