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

	// export MONGODB_URI="mongodb://username:password@localhost:27017"
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	username := os.Getenv("BASICAUTH_USERNAME")
	password := os.Getenv("BASICAUTH_PASSWORD")

	rssTitle := os.Getenv("RSS_TITLE")
	rssDescription := os.Getenv("RSS_DESCRIPTION")

	rssInit := app.RSSInit{
		DatabaseUri:    uri,
		Username:       username,
		Password:       password,
		RssTitle:       rssTitle,
		RssDescription: rssDescription,
	}

	fmt.Println("The RSS feed is running!")

	rootPath := fmt.Sprintf("/%s", os.Getenv("ROOT_PATH"))

	http.HandleFunc(rootPath, rssInit.CreateIndex)
	http.HandleFunc(fmt.Sprintf("/api%s", rootPath), rssInit.CreateItemAPI)
	http.HandleFunc("/health", rssInit.HealthCheck)
	http.ListenAndServe(":8082", nil)

}
