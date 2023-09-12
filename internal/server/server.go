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

	r.Route("/wallet", func(r chi.Router) {
		r.Post("/", s.createWallet)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), r)
}

func (s *Server) getBlockchain(w http.ResponseWriter, r *http.Request) {
	bc := s.cache.Get("blockchain")
	if bc == nil {
		minersWallet, err := wallet.New()
		if err != nil {
			httpRes(w, fmt.Errorf("could not initialize wallet: %w", err), http.StatusInternalServerError)
			return
		}

		bc, err = blockchain.New(minersWallet.Address())
		if err != nil {
			httpRes(w, fmt.Errorf("could not create new blockchain: %w", err), http.StatusInternalServerError)
			return
		}

		s.cache.Put("blockchain", bc)
	}

	httpRes(w, bc, http.StatusOK)
}

func (s *Server) createWallet(w http.ResponseWriter, r *http.Request) {
	wlt, err := wallet.New()
	if err != nil {
		httpRes(w, fmt.Errorf("could not create new wallet: %w", err), http.StatusInternalServerError)
		return
	}

	httpRes(w, struct {
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
		Address    string `json:"address"`
	}{
		PrivateKey: fmt.Sprintf("%x", wlt.PrivateKey().D.Bytes()),
		PublicKey:  fmt.Sprintf("%x%x", wlt.PublicKey().X.Bytes(), wlt.PublicKey().Y.Bytes()),
		Address:    wlt.Address(),
	}, http.StatusCreated)
}

func httpRes(w http.ResponseWriter, res any, code int) {
	if err, ok := res.(error); ok {
		res = HTTPError{Error: err.Error()}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
