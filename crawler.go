package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

type CrawlerControl struct {
	StartTime time.Time
	EndTime   time.Time
	Links     []string
	Data      []RbData
}

func (cc *CrawlerControl) Start(config *Config) {

	cc.StartTime = time.Now()
	log.Debugf("Crawl initiated at: %s", cc.StartTime)

	if len(cc.Links) < 1 {
		cc.ExtractRbLinks(config)
	}

	wg.Add(len(MyConfig.RbList))
	for _, Rb := range MyConfig.RbList {
		go ExtractRbData(Rb)
	}
	wg.Wait()

}

func (cc *CrawlerControl) ExtractRbLinks(config *Config) {

	res, err := http.Get(config.BaseURL + config.ExtURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a[href]").Each(func(_ int, link *goquery.Selection) {
		href, _ := link.Attr("href")
		if relDoc := strings.Contains(href, "/rodelbahnen-alpen/rodeltour/"); relDoc {
			newLink := config.BaseURL + href
			cc.Links = append(cc.Links, newLink)
		}
	})
}

func ExtractRbData(rbUrl string) {

	defer wg.Done()
	rbRes, err := http.Get(rbUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer rbRes.Body.Close()

	doc, err := goquery.NewDocumentFromReader(rbRes.Body)
	if err != nil {
		log.Fatal(err)
	}

	location := strings.TrimSpace(doc.Find("h1").First().Text())

	// find table
	doc.Find(".table-striped").Each(func(_ int, table *goquery.Selection) {

		// iterate over entries
		table.Find("tr").Each(func(entryIdx int, tr *goquery.Selection) {
			rbEntry := &RbData{}

			// iterate over single cells
			tr.Find("td").EachWithBreak(func(dataIdx int, td *goquery.Selection) bool {

				tdValue := strings.TrimSpace(td.Text())

				switch {
				case dataIdx == 0:
					rbEntry.Time, err = time.Parse("2006-01-02", tdValue)
					if err != nil {
						log.Fatal("Could not Parse Date!", tdValue)
					}
				case dataIdx == 1:
					rbEntry.User = tdValue
				case dataIdx == 2:
					if tdValue != "" {
						rbEntry.Rating = tdValue
					} else {
						return false
					}
				case dataIdx == 3:
					rbEntry.Comment = tdValue
				case (dataIdx > 3) && (dataIdx < 10):
					rbEntry.Comment += tdValue
				}
				return true
			})

			if rbEntry.Rating != "" {
				rbEntry.Link = rbUrl
				rbEntry.Location = location
				if err := rbEntry.Commit(ActiveDbSession.Collection); err == nil {
					if MyConfig.Notify == true {
						//SlackNotifier(rbEntry)
					}
				}
			} else {
				log.Info("No Rating!")
			}

		})

	})
}
