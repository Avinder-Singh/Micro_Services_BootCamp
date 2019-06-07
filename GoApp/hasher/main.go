package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//	"io"

	"github.com/gorilla/mux"
)

type request_body struct {
	Token string
}

type response_body struct {
	Hash string `json:"hash"`
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/hasher", GenerateHash).Methods("POST")
	log.Fatal(http.ListenAndServe(":8082", router))
}

func GenerateHash(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("inside")

	decoder := json.NewDecoder(r.Body)
	var body request_body
	err := decoder.Decode(&body)
	if err != nil {
		panic(err)
	}

	decodedToken, err := hex.DecodeString(body.Token)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(decodedToken)
	hash := sha256.Sum256([]byte(decodedToken))
	//fmt.Println(hash)
	encodedHash := hex.EncodeToString(hash[:])
	fmt.Println(encodedHash)
	res := response_body{
		encodedHash,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}

	//vars := mux.Vars(r)
	//size := vars["size"]
	//ln, err := strconv.Atoi(size)
	//if err == nil {
	//	fmt.Println(ln)
	//}
}
