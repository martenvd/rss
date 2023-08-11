// GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o rssfilter main.go

package app

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RSSInit struct {
	DatabaseUri    string
	Username       string
	Password       string
	RssTitle       string
	RssDescription string
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

	rss.WriteToDatabase(bsonItem, "rss", "feeditems")
}

func (rss *RSSInit) WriteToDatabase(item ItemBSON, database string, collection string) error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(rss.DatabaseUri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			fmt.Println(err)
		}
	}()

	coll := client.Database(database).Collection(collection)
	filter := bson.D{
		{Key: "title", Value: item.Title},
		{Key: "link", Value: item.Link},
		{Key: "description", Value: item.Description},
	}
	insert := bson.D{
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "title", Value: item.Title},
			{Key: "link", Value: item.Link},
			{Key: "description", Value: item.Description},
			{Key: "pubDate", Value: item.PubDate},
		}}}
	options := options.Update().SetUpsert(true)

	result, err := coll.UpdateOne(context.TODO(), filter, insert, options)
	if err != nil {
		return err
	}

	fmt.Printf("Document inserted with ID: %s\n", result.UpsertedID)

	return nil
}

func (rss *RSSInit) CreateRSSFeed() []byte {
	rssFeed := Rss{
		XMLName:  xml.Name{Space: "rss"},
		Version:  "2.0",
		Channels: []Channel{},
	}

	items := rss.GetAllFromDatabaseAndConvert()

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

func (rss *RSSInit) GetAllFromDatabaseAndConvert() []Item {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(rss.DatabaseUri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			fmt.Println(err)
		}
	}()

	coll := client.Database("rss").Collection("feeditems")

	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println(err)
	}

	var results []ItemBSON
	if err = cursor.All(context.TODO(), &results); err != nil {
		fmt.Println(err)
	}

	var items []Item

	for _, value := range results {
		items = append(items, Item{
			XMLName:     xml.Name{Local: "item"},
			Title:       value.Title,
			Link:        value.Link,
			Description: value.Description,
			PubDate:     value.PubDate,
		})
	}

	return items

}

func (rss *RSSInit) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Healthy"))
}
