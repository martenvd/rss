// GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o rssfilter main.go

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/martenvd/rss/internal/app"
)

func main() {

	uri := os.Getenv("MONGODB_URI")
	username := os.Getenv("BASICAUTH_USERNAME")
	password := os.Getenv("BASICAUTH_PASSWORD")
	// export MONGODB_URI="mongodb://username:password@localhost:27017"
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	rssInit := app.RSSInit{
		DatabaseUri: uri,
		Username:    username,
		Password:    password,
	}

	fmt.Println("The RSS feed is running!")

	http.HandleFunc("/", rssInit.CreateIndex)
	http.HandleFunc("/api", rssInit.CreateItemAPI)
	http.ListenAndServe(":8082", nil)

}
