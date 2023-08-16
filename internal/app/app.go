// GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o rssfilter main.go

package app

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

type RSSInit struct {
	DatabaseType     string
	ConnectionString string
	DatabaseUri      string
	Username         string
	Password         string
	RssTitle         string
	RssDescription   string
}

func (rss *RSSInit) CreateIndex(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(200)
	feed := rss.CreateRSSFeed()

	w.Write(feed)
}

func (rss *RSSInit) BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(rss.Username))
			expectedPasswordHash := sha256.Sum256([]byte(rss.Password))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func (rss *RSSInit) CreateItemAPI(w http.ResponseWriter, r *http.Request) {

	var jsonItem ItemJSON

	err := json.NewDecoder(r.Body).Decode(&jsonItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bsonItem := ItemBSON(jsonItem)

	if strings.Contains(rss.DatabaseType, "mongo") {
		rss.WriteToMongoDatabase(bsonItem, "rss", "feeditems")
	} else {
		rss.WriteToMSSQLDatabase(jsonItem, "feeditems")
	}
}

func (rss *RSSInit) CreateRSSFeed() []byte {
	rssFeed := Rss{
		XMLName:  xml.Name{Space: "rss"},
		Version:  "2.0",
		Channels: []Channel{},
	}

	var items []Item

	if strings.Contains(rss.DatabaseType, "mongo") {
		items = rss.GetAllFromMongoDatabaseAndConvert()
	} else {
		items = rss.GetAllFromMSSQLDatabaseAndConvert(ItemJSON{Title: "test", Link: "link", Description: "Description", PubDate: "vandaag"}, "feeditems")
		fmt.Println(items)
	}

	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}

	var pubDate string

	if len(items) == 0 {
		pubDate = ""
	} else {
		pubDate = items[0].PubDate
	}

	rssFeed.Channels = []Channel{
		{
			XMLName:     xml.Name{Space: "channel"},
			Title:       rss.RssTitle,
			Link:        "",
			Description: rss.RssDescription,
			Atom:        "http://www.w3.org/2005/Atom",
			AtomLink: AtomLink{
				Href: "",
				Rel:  "self",
				Type: "application/rss+xml",
			},
			LastBuildDate: pubDate,
			PubDate:       pubDate,
			Item:          items,
		},
	}

	feed, _ := xml.MarshalIndent(rssFeed, "", " ")

	return feed

}

func (rss *RSSInit) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Healthy"))
}
