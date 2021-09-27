package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/mike/inv/router"
)

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == " " {
		err := os.Setenv("PORT", "9001")
		if err != nil {
			return "environment variable not set", err
		}
		err = os.Setenv("DATABASE_URL","postgres://evgllcxlyxdzru:bba19fa32f9668a5785cb38bd8ed822b00790b1b82e973b641f9848e28ca42a9@ec2-54-157-4-216.compute-1.amazonaws.com:5432/deefo83mbp8mg8")
		if err != nil {
			return "environment variable $DATABASE_URL not set",err
		}

	}
	fmt.Println(port)
	fmt.Println(os.Getenv("DATABASE_URL"))
	return ":" + port, nil
}

func handleRequests() {
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(`Error determining listen address: check environment variables \n
		(If running on heroku check env vars)`)
	}
	r := router.Router(addr, err)


	_ = http.ListenAndServe(addr, r)
}

/* func handleDevRequests() {

	r := router.Router(":9001", nil)
	log.Fatal()


	http.ListenAndServe(":9001", r)
} */
func main() {
	//handleDevRequests()

	handleRequests()
}
