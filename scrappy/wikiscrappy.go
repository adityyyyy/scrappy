package scrappy

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/gocolly/colly"
)

func wikiLink() {
	var links []string

	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(2),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
	})

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		link := h.Attr("href")
		links = append(links, link)
		h.Request.Visit(link)
	})

	c.OnScraped(func(r *colly.Response) {
		f, err := os.Create("wikiLinks.csv")
		if err != nil {
			log.Fatalln("err: ", err)
		}
		defer f.Close()

		writer := csv.NewWriter(f)
		defer writer.Flush()

		headers := []string{
			"S.No",
			"Link",
		}

		writer.Write(headers)

		for index, link := range links {
			writer.Write([]string{
				strconv.Itoa(index),
				link,
			})
		}
	})

	c.Visit("https://en.wikipedia.org/")

	c.Wait()
}
