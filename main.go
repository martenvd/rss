// GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o rssfilter main.go

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/martenvd/rss/internal/app"
)

func main() {

	uri := os.Getenv("MONGODB_URI")
	// export MONGODB_URI="mongodb://username:password@localhost:27017"
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	db := app.DB{
		DatabaseUri: uri,
	}

	http.HandleFunc("/", db.CreateIndex)
	http.HandleFunc("/api", db.CreateItemAPI)
	http.ListenAndServe(":8082", nil)

}
