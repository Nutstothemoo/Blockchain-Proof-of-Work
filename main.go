package main

import (
	"log"
	"net/http"
	"spew"

	"github.com/joho/godotenv"
)
type Block struct {
	index int
	Timestamp string
	Data int 
	Hash string
	PrevHash string
	Diffficulty int
	Nonce string
}

const difficulty = 1

var Blockchain []Block
var mutex sync.Mutex

func main(){
	err:= godotenv.Load()
	if err!= nil {
		log.Fatal(err)
	}
	go func() {
		t:= time.now()
		genesisBlock:= Block{0, t.String(), 0, calculateHash(genesisBlock), "", difficulty, ""}
		spew.Dump(genesisBlock)
		mutex.Lock()
	}	()
	log.Fatal(run())	
}

func run() error {

}

func makeMuxRouter() http.Handler {
	HandleFunc("/", handleGetBlockchain).Methods("GET")
	HandleFunc("/", handleWriteBlock).Methods("POST")
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {

}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {

}
func blockIsValid(newBlock, oldBlock Block) bool {

}
func generateBlock(oldBlock Block, BPM int) (Block, error) {

}
func calculateHash(block Block) string {

}

func isHashValid(hash string, difficulty int) bool {

}
func responseWithJSON(w http.ResponseWriter, json []byte, code int) {

}