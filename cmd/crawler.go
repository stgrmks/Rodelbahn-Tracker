package cmd

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

var crawl = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl Data",
	Long:  `Crawler is executed.`,
	Run: func(cmd *cobra.Command, args []string) {
		RunStartCrawler()
	},
}

var periodicCrawl = &cobra.Command{
	Use:   "periodic-crawl",
	Short: "Periodically Crawl Data",
	Long:  `Crawler is executed periodically. Intended for usage as a service.`,
	Run: func(cmd *cobra.Command, args []string) {
		RunPeriodicCrawler()
	},
}

func RunPeriodicCrawler() {
	c := cron.New()
	c.AddFunc(config.Cron, RunStartCrawler)
	fmt.Printf("Periodical Crawl initiated: %s\n", config.Cron)
	c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func RunStartCrawler() {

	RunInitDB()

	fmt.Printf("Crawl initiated at: %s\n", time.Now().String())

	if len(config.RbList) == 0 {
		ExtractRbLinks()
	}

	wg.Add(len(config.RbList))
	for _, Rb := range config.RbList {
		go ExtractRbData(Rb)
	}
	wg.Wait()

}

func ExtractRbLinks() {

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
			config.RbList = append(config.RbList, newLink)
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
				if err := rbEntry.Commit(config.ActiveCollection); err == nil {
					if config.Notify == true {
						SlackNotifier(rbEntry)
					}
				}
			} else {
				log.Println("No Rating!")
			}

		})

	})
}
