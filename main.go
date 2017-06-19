package main

import (
	"net/http"
	"os"
	"time"
    "log"
    "github.com/garyburd/redigo/redis"
 	"crypto/tls"
	
)


var redisPool = &redis.Pool{
    MaxIdle: 5,
    MaxActive: 10,
    Wait: true,
    IdleTimeout: 10 * time.Second,
    Dial: func() (conn redis.Conn, err error) {
      
        conn, err = redis.DialURL(os.Getenv("REDIS_URL"), redis.DialTLSConfig(&tls.Config{InsecureSkipVerify: true}))
    
        if err != nil {
            panic(err)
        }

        return 
    },
}

var (
	waitTimeout = time.Minute * 10
    rr          redisReceiver
	rw          redisWriter
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("PORT : %v  - Error -> Post must be set.", port)
        
	}
    
    
	rr = newRedisReceiver(redisPool)
	rw = newRedisWriter(redisPool)

	go func() {
		for {
    		rr.broadcast(availableMessage)
            err := rr.run()
			if err == nil {
				break
			}
			log.Println(err)		
        }
	}()

	go func() {
		for {
            err := rw.run()
			if err == nil {
				break
			}
			log.Println(err)
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", handleWebsocket)
	log.Println(http.ListenAndServe(":"+port, nil))
}
