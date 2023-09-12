package server

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	"github.com/nerock/goblockchain/internal/signature"
	"github.com/nerock/goblockchain/internal/transaction"

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

	r.Get("/blockchain", s.getBlockchain)
	r.Post("/wallet", s.createWallet)

	r.Route("/transaction", func(r chi.Router) {
		r.Get("/", s.retrieveTransactions)
		r.Post("/", s.createTransaction)
	})

	r.Post("/mine", s.mine)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), r)
}

func (s *Server) getBlockchain(w http.ResponseWriter, r *http.Request) {
	bc, err := s.getOrCreateBlockchain()
	if err != nil {
		httpRes(w, fmt.Errorf("could not initialize blockchain: %w", err), http.StatusInternalServerError)
	}

	httpRes(w, bc, http.StatusOK)
}

func (s *Server) getOrCreateBlockchain() (*blockchain.Blockchain, error) {
	bc := s.cache.Get("blockchain")
	if bc == nil {
		minersWallet, err := wallet.New()
		if err != nil {
			return nil, err
		}

		bc, err = blockchain.New(minersWallet.Address())
		if err != nil {
			return nil, err
		}

		s.cache.Put("blockchain", bc)
	}

	return bc, nil
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

func (s *Server) retrieveTransactions(w http.ResponseWriter, r *http.Request) {
	bc, err := s.getOrCreateBlockchain()
	if err != nil {
		httpRes(w, fmt.Errorf("could not retrieve blockchain: %w", err), http.StatusInternalServerError)
		return
	}

	trs := bc.CopyTransactionPool()
	httpRes(w, struct {
		Transactions []*transaction.Transaction `json:"transactions"`
		Length       int                        `json:"length"`
	}{
		Transactions: bc.CopyTransactionPool(),
		Length:       len(trs),
	}, http.StatusOK)
}

func (s *Server) createTransaction(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body struct {
		SenderPrivateKey           string  `json:"sender_private_key"`
		SenderPublicKey            string  `json:"sender_public_key"`
		SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
		RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
		Value                      float32 `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpRes(w, fmt.Errorf("could not decode transaction body: %w", err), http.StatusBadRequest)
		return
	}

	if body.SenderPrivateKey == "" || body.SenderPublicKey == "" ||
		body.SenderBlockchainAddress == "" || body.RecipientBlockchainAddress == "" ||
		body.Value == 0 {
		httpRes(w, errors.New("invalid request"), http.StatusBadRequest)
		return
	}

	publicKey, err := strToECDSAPublicKey(body.SenderPublicKey)
	if err != nil {
		httpRes(w, fmt.Errorf("could not decode public key: %w", err), http.StatusBadRequest)
		return
	}

	privateKey, err := strToECDSAPrivateKey(body.SenderPrivateKey, *publicKey)
	if err != nil {
		httpRes(w, fmt.Errorf("could not decode private key: %w", err), http.StatusBadRequest)
		return
	}

	trs := transaction.New(body.SenderBlockchainAddress, body.RecipientBlockchainAddress, body.Value)
	sig, err := signature.Sign(trs, privateKey)
	if err != nil {
		httpRes(w, fmt.Errorf("could not sign transaction: %w", err), http.StatusBadRequest)
		return
	}

	bc, err := s.getOrCreateBlockchain()
	if err != nil {
		httpRes(w, fmt.Errorf("could not retrieve blockchain: %w", err), http.StatusInternalServerError)
		return
	}

	if err := bc.AddTransaction(body.SenderBlockchainAddress, body.RecipientBlockchainAddress, body.Value, publicKey, sig); err != nil {
		httpRes(w, fmt.Errorf("could not add transaction: %w", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) mine(w http.ResponseWriter, r *http.Request) {
	bc, err := s.getOrCreateBlockchain()
	if err != nil {
		httpRes(w, fmt.Errorf("could not retrieve blockchain: %w", err), http.StatusInternalServerError)
		return
	}

	if err := bc.Mining(); err != nil {
		httpRes(w, fmt.Errorf("could not mine: %w", err), http.StatusInternalServerError)
		return
	}

	httpRes(w, nil, http.StatusOK)
}

func strToECDSAPrivateKey(s string, publicKey ecdsa.PublicKey) (*ecdsa.PrivateKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("could not decode bytes: %w", err)
	}

	var d big.Int
	d.SetBytes(b)

	return &ecdsa.PrivateKey{PublicKey: publicKey, D: &d}, nil
}

func strToECDSAPublicKey(s string) (*ecdsa.PublicKey, error) {
	if len(s) != 128 {
		return nil, errors.New("invalid key length")
	}

	bx, err := hex.DecodeString(s[:64])
	if err != nil {
		return nil, fmt.Errorf("could not decode x bytes from public key: %w", err)
	}
	by, err := hex.DecodeString(s[64:])
	if err != nil {
		return nil, fmt.Errorf("could not decode y bytes from public key: %w", err)
	}

	var x, y big.Int
	x.SetBytes(bx)
	y.SetBytes(by)

	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}, nil
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
