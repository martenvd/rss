// GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o rssfilter main.go

package app

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	DatabaseUri string
}

func (db *DB) CreateIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

	feed := db.CreateRSSFeed()

	w.Write(feed)
}

func (db *DB) CreateItemAPI(w http.ResponseWriter, r *http.Request) {

	var jsonItem ItemJSON

	err := json.NewDecoder(r.Body).Decode(&jsonItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bsonItem := ItemBSON(jsonItem)

	db.WriteToDatabase(bsonItem, "rss", "feeditems")
}

func (db *DB) WriteToDatabase(item ItemBSON, database string, collection string) error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(db.DatabaseUri).SetServerAPIOptions(serverAPI)
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
	filter := bson.D{{Key: "title", Value: item.Title}, {Key: "link", Value: item.Link}, {Key: "description", Value: item.Description}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "title", Value: item.Title}}}}
	options := options.Update().SetUpsert(true)

	result, err := coll.UpdateOne(context.TODO(), filter, update, options)
	if err != nil {
		return err
	}

	fmt.Printf("Document inserted with ID: %s\n", result.UpsertedID)

	return nil
}

func (db *DB) CreateRSSFeed() []byte {
	rssFeed := Rss{
		XMLName:  xml.Name{Space: "rss"},
		Version:  "2.0",
		Channels: []Channel{},
	}

	items := db.GetAllFromDatabaseAndConvert()

	rssFeed.Channels = []Channel{
		{
			XMLName:     xml.Name{Space: "channel"},
			Title:       "Martens test RSS Feed",
			Link:        "https://var.tf/",
			Description: "Martens test RSS Feed Channel",
			Atom:        "http://www.w3.org/2005/Atom",
			AtomLink: AtomLink{
				Href: "https://var.tf/rss",
				Rel:  "self",
				Type: "application/rss+xml",
			},
			Item: items,
		},
	}

	feed, _ := xml.MarshalIndent(rssFeed, "", " ")

	return feed

}

func (db *DB) GetAllFromDatabaseAndConvert() []Item {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(db.DatabaseUri).SetServerAPIOptions(serverAPI)
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
		})
	}

	return items

}
