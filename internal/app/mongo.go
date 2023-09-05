package app

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"time"

	_ "time/tzdata"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (rss *RSSInit) WriteToMongoDatabase(item ItemBSON, database string, collection string) error {
	var ctx context.Context

	if rss.DatabaseType == "mongodb4" {
		ctx = context.Background()
		clientOptions := options.Client().ApplyURI(rss.DatabaseUri).SetDirect(true)

		c, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatalf("unable to initialize connection %v", err)
		}

		err = c.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("unable to connect %v", err)
		}

		err = rss.AddToMongoDatabase(item, c, ctx)
		if err != nil {
			log.Printf("Adding item to database resulted in the following error: %s", err)
		}

	} else if rss.DatabaseType == "mongodb6" {
		ctx = context.TODO()
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(rss.DatabaseUri).SetServerAPIOptions(serverAPI)
		c, err := mongo.Connect(ctx, opts)
		if err != nil {
			return err
		}
		defer func() {
			if err = c.Disconnect(ctx); err != nil {
				fmt.Println(err)
			}
		}()

		err = rss.AddToMongoDatabase(item, c, ctx)
		if err != nil {
			log.Printf("Adding item to database resulted in the following error: %s", err)
		}
	}

	return nil
}

func (rss *RSSInit) GetAllFromMongoDatabaseAndConvert() []Item {

	var items []Item

	if rss.DatabaseType == "mongodb4" {
		c := &mongo.Client{}

		ctx := context.Background()
		clientOptions := options.Client().ApplyURI(rss.DatabaseUri).SetDirect(true)
		defer c.Disconnect(ctx)

		c, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatalf("unable to initialize connection %v", err)
		}

		err = c.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("unable to connect %v", err)
		}

		items = rss.GetMongoFindResults(c, ctx)

	} else if rss.DatabaseType == "mongodb6" {
		c := &mongo.Client{}

		ctx := context.TODO()
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(rss.DatabaseUri).SetServerAPIOptions(serverAPI)
		c, err := mongo.Connect(ctx, opts)
		if err != nil {
			fmt.Println(err)
		}
		defer func() {
			if err = c.Disconnect(ctx); err != nil {
				fmt.Println(err)
			}
		}()

		items = rss.GetMongoFindResults(c, ctx)
	}

	return items
}

func (rss *RSSInit) GetMongoFindResults(c *mongo.Client, ctx context.Context) []Item {

	coll := c.Database("rss").Collection("feeditems")

	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println(err)
	}

	var results []ItemBSON
	if err = cursor.All(ctx, &results); err != nil {
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

func (rss *RSSInit) AddToMongoDatabase(item ItemBSON, c *mongo.Client, ctx context.Context) error {
	loc, _ := time.LoadLocation("Europe/Amsterdam")
	now := time.Now().In(loc).Format("Mon, 02 Jan 2006 15:04:05 -0700")

	coll := c.Database("rss").Collection("feeditems")
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
			{Key: "pubDate", Value: now},
		}}}
	options := options.Update().SetUpsert(true)

	result, err := coll.UpdateOne(ctx, filter, insert, options)
	if err != nil {
		return err
	}

	fmt.Printf("Document inserted with ID: %s\n", result.UpsertedID)

	return nil
}
