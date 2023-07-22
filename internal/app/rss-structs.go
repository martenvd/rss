package app

import "encoding/xml"

type Rss struct {
	XMLName  xml.Name  `xml:"rss"`
	Version  string    `xml:"version,attr"`
	Channels []Channel `xml:"channel"`
}

type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Atom        string   `xml:"xmlns:atom,attr"`
	AtomLink    AtomLink `xml:"atom:link"`
	Item        []Item   `xml:"item"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type Image struct {
	XMLName xml.Name `xml:"image"`
	Title   string   `xml:"title"`
	Url     string   `xml:"url"`
	Link    string   `xml:"link"`
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link,omitempty"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
}

type ItemJSON struct {
	Title       string `json:"title"`
	Link        string `json:"link,omitempty"`
	Description string `json:"description"`
	PubDate     string `json:"pubDate"`
}

type ItemBSON struct {
	Title       string `bson:"title"`
	Link        string `bson:"link,omitempty"`
	Description string `bson:"description"`
	PubDate     string `bson:"pubDate"`
}

type GUID struct {
	IsPermaLink string `xml:"isPermaLink,attr"`
	InnerXML    string `xml:",innerxml"`
}

type URL struct {
	URL     string
	Filters []string
}
