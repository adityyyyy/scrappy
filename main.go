package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Product struct {
	url, image, name, price string
}

func main() {
	pageToScrap := []string{
		"https://www.scrapingcourse.com/ecommerce/page/1",
		"https://www.scrapingcourse.com/ecommerce/page/2",
		"https://www.scrapingcourse.com/ecommerce/page/3",
		"https://www.scrapingcourse.com/ecommerce/page/4",
	}

	var products []Product

	c := colly.NewCollector(
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		Parallelism: 5,
	})

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited", r.Request.URL)
	})

	// c.OnHTML("a.page-numbers", func(h *colly.HTMLElement) {
	// 	newPaginationLink := h.Attr("href")
	//
	// 	if !contains(pageToScrap, newPaginationLink) {
	// 		if !contains(pageDiscovered, newPaginationLink) {
	// 			pageToScrap = append(pageToScrap, newPaginationLink)
	// 		}
	// 		pageDiscovered = append(pageDiscovered, newPaginationLink)
	// 	}
	// })

	for _, urls := range pageToScrap {
		c.Visit(urls)
	}

	c.Wait()

	c.OnHTML("li.product", func(h *colly.HTMLElement) {
		product := Product{}

		product.url = h.ChildAttr("a", "href")
		product.image = h.ChildAttr("img", "href")
		product.name = h.ChildText("h2")
		product.price = h.ChildText(".price")

		products = append(products, product)
	})

	c.OnScraped(func(r *colly.Response) {
		file, err := os.Create("product.csv")
		if err != nil {
			log.Fatalln("Failed to create output  CSV file", err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)

		headers := []string{
			"url",
			"image",
			"name",
			"price",
		}

		writer.Write(headers)

		for _, product := range products {
			record := []string{
				product.url,
				product.image,
				product.name,
				product.price,
			}

			writer.Write(record)
		}

		defer writer.Flush()
	})
}

func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}
