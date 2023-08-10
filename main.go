package main

import (
	
)
type Block struct {

}

var Blockchain []Block

func main(){

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