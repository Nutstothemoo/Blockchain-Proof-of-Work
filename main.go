package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Block struct {
	Index        int
	Timestamp    string
	Data         int
	Hash         string
	PrevHash     string
	Difficulty   int
	Nonce        string
	Signature    []byte
}

const difficulty = 1

var Blockchain []Block
var mutex sync.Mutex
var privateKey *rsa.PrivateKey


func main() {
	err := godotenv.Load()
	if err != nil {
			log.Fatal(err)
	}

	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
			log.Fatal(err)
	}

	go func() {
			t := time.Now()
			genesisBlock := Block{0, t.String(), 0, "", "", difficulty, "", nil}
			hash, err := calculateHash(genesisBlock)
			if err != nil {
					log.Fatal(err)
			}
			genesisBlock.Hash = hash

			genesisBlock.Signature, err = rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, []byte(genesisBlock.Hash), nil)
			if err != nil {
					log.Fatal(err)
			}

			spew.Dump(genesisBlock)
			mutex.Lock()
			Blockchain = append(Blockchain, genesisBlock)
			mutex.Unlock()
	}()

	log.Fatal(run())
}

func run() error {
	mux := makeMuxRouter()
	httpAddr := "8080"
	log.Println("Listening on ", httpAddr)
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var newBlock Block
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newBlock); err != nil {
			respondWithJSON(w, r, http.StatusBadRequest, r.Body)
			return
	}
	defer r.Body.Close()

	mutex.Lock()
	defer mutex.Unlock()
	publicKey := &privateKey.PublicKey

	valid, err := blockIsValid(newBlock, Blockchain[len(Blockchain)-1], publicKey)
		if err != nil {
			respondWithJSON(w, r, http.StatusInternalServerError, "Error validating the block: "+err.Error())
			return
	}

	if valid {
			Blockchain = append(Blockchain, newBlock)
			respondWithJSON(w, r, http.StatusCreated, newBlock)
	} else {
			respondWithJSON(w, r, http.StatusConflict, "Block is not valid")
	}
}

func proofOfWork(block Block) string {
	nonce := 0
	hash, _ := calculateHash(block)
	for !isHashValid(hash, difficulty) {
			nonce++
			block.Nonce = strconv.Itoa(nonce)
			hash, _ = calculateHash(block)
	}
	return strconv.Itoa(nonce)
}

func generateBlock(oldBlock Block, Data int, privateKey *rsa.PrivateKey) (Block, error) {
	var newBlock Block

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = time.Now().String()
	newBlock.Data = Data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty

	newBlock.Nonce = proofOfWork(newBlock)
	hash, err := calculateHash(newBlock)
	if err != nil {
			return Block{}, err
	}
	newBlock.Hash = hash

	newBlock.Signature, err = rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, []byte(newBlock.Hash), nil)
	if err != nil {
			return Block{}, err
	}

	return newBlock, nil
}

func blockIsValid(newBlock, oldBlock Block, publicKey *rsa.PublicKey) (bool, error) {
	if oldBlock.Index+1 != newBlock.Index {
			return false, nil
	}

	if oldBlock.Hash != newBlock.PrevHash {
			return false, nil
	}

	hash, err := calculateHash(newBlock)
	if err != nil {
			return false, err
	}

	if hash != newBlock.Hash {
			return false, nil
	}

	err = rsa.VerifyPSS(publicKey, crypto.SHA256, []byte(newBlock.Hash), newBlock.Signature, nil)
	if err != nil {
			return false, err
	}

	return true, nil
}

func calculateHash(block Block) (string, error) {
	record := string(rune(block.Index)) + block.Timestamp + string(rune(block.Data)) + block.PrevHash + block.Nonce
	h := sha256.New()
	_, err := h.Write([]byte(record))
	if err != nil {
			return "", err
	}
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed), nil
}

func isHashValid(hash string, difficulty int) bool {
	prefix := "0000" 
	return hash[:difficulty] == prefix
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
