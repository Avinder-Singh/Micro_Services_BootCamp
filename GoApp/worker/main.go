package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/garyburd/redigo/redis"
)

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func main() {

	i := 0
	var pool= newPool()
	for i=1;i<=50 ;i++ {
		log.Println(i)
		//time.Sleep(2* time.Second)
		resp, err := http.Get("http://localhost:8080/stg/tokens/20")
		if err != nil {
			log.Fatalln(err)
		}

		var result map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&result)

		log.Println(result)
		//token := result["token"]
		//log.Println(token)

		bytesRepresentation, err := json.Marshal(result)
		if err != nil {
			log.Fatalln(err)
		}

		response, err := http.Post("http://localhost:8082/hasher", "application/json", bytes.NewBuffer(bytesRepresentation))
		if err != nil {
			log.Fatalln(err)
		}

		var finalresult map[string]string

		json.NewDecoder(response.Body).Decode(&finalresult)

		log.Println(finalresult)
		log.Println(finalresult["hash"])
		hash := finalresult["hash"]
		//fmt.Println(hash)

		if (hash[0] == '0') { // if lucky hash
			//fmt.Println("true")
			st:="lucky-hash-found"
			c := pool.Get()
			n, err := c.Do("LPUSH", "list-8", hash)
			if err != nil {
				log.Fatalln(err)
			}
			m, err := c.Do("PUBLISH", "mychannel-1", st)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("Lucky hash found")
			fmt.Println(n)
			fmt.Println(m)
		}
	}
}

