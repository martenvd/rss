package app

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"log"

	_ "github.com/microsoft/go-mssqldb"
)

var db *sql.DB

func (rss *RSSInit) WriteToMSSQLDatabase(item ItemJSON, table string) error {
	var err error
	db, err = sql.Open("sqlserver", rss.ConnectionString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()

	if db == nil {
		err = errors.New("db is null")
		return err
	}

	// Check if database is alive.
	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	tsql := fmt.Sprintf(`
      INSERT INTO %s (title, description, link, pubDate) VALUES (@title, @description, @link, @pubDate);
      select isNull(SCOPE_IDENTITY(), -1);
    `, table)

	stmt, err := db.Prepare(tsql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("title", item.Title),
		sql.Named("description", item.Description),
		sql.Named("link", item.Link),
		sql.Named("pubDate", item.PubDate))
	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return err
	}

	return nil
}

func (rss *RSSInit) GetAllFromMSSQLDatabaseAndConvert(item ItemJSON, table string) []Item {
	var err error
	// Create connection pool
	db, err = sql.Open("sqlserver", rss.ConnectionString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	tableStatement := fmt.Sprintf("if not exists (select * from sysobjects where name='%s' and xtype='U') create table feeditems (title varchar(64) not null, description varchar(1024), link varchar(256), pubDate varchar(256))", table)

	_, err = db.Exec(tableStatement)
	if err != nil {
		panic(err)
	}

	tsqlRead := fmt.Sprintf("SELECT * FROM %s;", table)

	// Execute query
	rows, err := db.QueryContext(ctx, tsqlRead)
	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	var count int
	var items []Item

	// Iterate through the result set.
	for rows.Next() {
		var title, description, link, pubDate string

		// Get values from row.
		err := rows.Scan(&title, &description, &link, &pubDate)
		if err != nil {
			fmt.Println(err)
		}

		items = append(items, Item{
			XMLName:     xml.Name{Local: "item"},
			Title:       title,
			Description: description,
			Link:        link,
			PubDate:     pubDate,
		})

		count++
	}

	return items
}
