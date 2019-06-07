package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	//	"io"

	"github.com/gorilla/mux"
)

type response struct {
	Token string `json:"token"`
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/stg/tokens/{size}", GenerateToken).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func GenerateToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	size := vars["size"]
	ln, err := strconv.Atoi(size)
	if err == nil {
		fmt.Println(ln)
	}

	bytes := make([]byte, ln)
	for i := 0; i < ln; i++ {
		if i%2 == 0 {
			bytes[i] = byte(97 + rand.Intn(25))
		} else {
			bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
		}
	}
	//fmt.Println(bytes)
	encodedStr := hex.EncodeToString(bytes)

	res := response{
		encodedStr,
	}
	//fmt.Println(encodedStr)
	fmt.Println(res)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}