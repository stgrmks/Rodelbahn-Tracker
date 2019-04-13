package main

import (
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
	Result    []RbData
}

func (cc *CrawlerControl) Start(config *Config) {

	defer func() {
		cc.EndTime = time.Now()
	}()

	cc.StartTime = time.Now()
	log.Debugf("Crawl initiated at: %s", cc.StartTime)

	// Get RB Links if there are none
	if len(cc.Links) < 1 {
		cc.ExtractRbLinks(config)
	}

	// add wg for each link
	wg.Add(len(cc.Links[:2]))
	for _, Rb := range cc.Links[:2] {
		go cc.ExtractRbData(Rb)
	}
	wg.Wait()
	log.Debug("Crawler Finished")
}

func (cc *CrawlerControl) ExtractRbLinks(config *Config) {

	res, err := http.Get(config.BaseURL + config.ExtURL)
	if err != nil {
		log.Debug(err)
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Debug(err)
		return
	}

	doc.Find("a[href]").Each(func(_ int, link *goquery.Selection) {
		href, _ := link.Attr("href")
		if relDoc := strings.Contains(href, "/rodelbahnen-alpen/rodeltour/"); relDoc {
			newLink := config.BaseURL + href
			cc.Links = append(cc.Links, newLink)
		}
	})
}

func (cc *CrawlerControl) ExtractRbData(rbUrl string) {

	defer wg.Done()
	rbRes, err := http.Get(rbUrl)
	if err != nil {
		log.Debug(err)
		wg.Done()
	}

	defer rbRes.Body.Close()

	doc, err := goquery.NewDocumentFromReader(rbRes.Body)
	if err != nil {
		log.Debug(err)
		wg.Done()
	}

	location := strings.TrimSpace(doc.Find("h1").First().Text())

	// find table
	doc.Find(".table-striped").Each(func(_ int, table *goquery.Selection) {

		// iterate over entries
		table.Find("tr").Each(func(entryIdx int, tr *goquery.Selection) {
			rbEntry := new(RbData)

			// iterate over single cells
			tr.Find("td").EachWithBreak(func(dataIdx int, td *goquery.Selection) bool {

				tdValue := strings.TrimSpace(td.Text())

				switch {
				case dataIdx == 0:
					rbEntry.Time, err = time.Parse("2006-01-02", tdValue)
					if err != nil {
						log.Debug("Could not Parse Date!", tdValue)
						wg.Done()
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
				cc.Result = append(cc.Result, *rbEntry)
				log.Debugf("Found rating: %s", rbEntry)
			} else {
				log.Debug("No Rating!")
			}

		})

	})
}
