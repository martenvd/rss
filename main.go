// GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o rssfilter main.go

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/martenvd/rss/internal/app"
)

func main() {

	uri := os.Getenv("MONGODB_URI")

	// export MONGODB_URI="mongodb://username:password@localhost:27017"
	username := os.Getenv("BASICAUTH_USERNAME")
	password := os.Getenv("BASICAUTH_PASSWORD")
	databaseType := os.Getenv("DATABASE_TYPE")
	connectionString := os.Getenv("MSSQL_CONNECTION_STRING")

	if len(databaseType) == 0 {
		databaseType = "mongodb6"
	}

	rssTitle := os.Getenv("RSS_TITLE")
	rssDescription := os.Getenv("RSS_DESCRIPTION")
	rootPath := os.Getenv("ROOT_PATH")

	rssInit := app.RSSInit{
		DatabaseType:     databaseType,
		ConnectionString: connectionString,
		DatabaseUri:      uri,
		Username:         username,
		Password:         password,
		RssTitle:         rssTitle,
		RssDescription:   rssDescription,
		RootPath:         rootPath,
	}

	fmt.Printf("Rss feed starting with %s driver.\n", databaseType)
	fmt.Println("The RSS feed is running!")

	http.HandleFunc(fmt.Sprintf("/%s", rootPath), rssInit.BasicAuth(rssInit.CreateIndex))
	http.HandleFunc(fmt.Sprintf("/api/%s", rootPath), rssInit.CreateItemAPI)
	http.HandleFunc("/health", rssInit.HealthCheck)
	http.ListenAndServe(":8082", nil)

}
