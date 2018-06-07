package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// TrendData qiitaのトレンド
type TrendData struct {
	Trend struct {
		Edges []struct {
			Node struct {
				CreatedAt string `json:"createdAt"`
				Title     string `json:"title"`
				UUID      string `json:"uuid"`
				Author    struct {
					ProfileImageURL string `json:"profileImageUrl"`
					URLName         string `json:"urlName"`
				} `json:"author"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"trend"`
}

const (
	trendURL            = "https://qiita.com"
	internalServerError = "internal server error"
)

var port = flag.Int("port", 1234, "port mumber")

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", handleFeed)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}

func handleFeed(c echo.Context) error {
	res, err := http.Get(trendURL)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	feedItem := []*feeds.Item{}

	doc.Find("div[data-hyperapp-app=Trend]").Each(func(i int, s *goquery.Selection) {
		text, exists := s.Attr("data-hyperapp-props")
		if !exists {
			return
		}

		trends := &TrendData{}
		err := json.Unmarshal([]byte(text), trends)
		if err != nil {
			return
		}

		for _, item := range trends.Trend.Edges {
			url := url.URL{
				Scheme: "https",
				Host:   "qiita.com",
				Path:   path.Join(item.Node.Author.URLName, "items", item.Node.UUID),
			}
			created, _ := time.Parse(time.RFC3339, item.Node.CreatedAt)
			feedItem = append(feedItem, &feeds.Item{
				Title:       item.Node.Title,
				Link:        &feeds.Link{Href: url.String()},
				Description: item.Node.Title,
				Author:      &feeds.Author{Name: item.Node.Author.URLName, Email: ""},
				Created:     created,
			})
		}
	})

	feed := &feeds.Feed{
		Title:       "Qiita",
		Link:        &feeds.Link{Href: "https://qiita.com/trend"},
		Description: "Qiita trends",
		Updated:     time.Now(),
		Items:       feedItem,
	}
	text, err := feed.ToRss()
	if err != nil {
		c.Logger().Printf("%#v", err)
		return c.String(http.StatusInternalServerError, internalServerError)
	}

	c.Response().Header().Set("Content-Type", "text/xml; charset=utf-8")
	return c.String(http.StatusOK, text)
}
