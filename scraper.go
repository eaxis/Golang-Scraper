package main

import (
	"os"
	"fmt"
	//"time"
	"bufio"
	"strings"
	"strconv"
	"io/ioutil"
	"golang.org/x/text/transform"
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"golang.org/x/text/encoding/charmap"
)

type result struct {
    article string
    title string
    vendor string
    inStock bool
    price float32
}

func main() {
	fmt.Printf("Input an article:")

	reader := bufio.NewReader(os.Stdin)

	article, _ := reader.ReadString('\n');
	
	html := get(article)
	
	results := scrape(html)

	//fmt.Println(results)
	
	if (len(results) == 0) {
		fmt.Printf("Nothing found\n\n\n\n")
	} else {
		for c, result := range results {
			fmt.Printf("\n------ Position #%d ------\n", c + 1)
			fmt.Printf("Article: %s\n", result.article)
			fmt.Printf("Title: %s\n", result.title)
			fmt.Printf("Vendor: %s\n", result.vendor)
			fmt.Printf("In stock: %s\n", boolToString(result.inStock))
			fmt.Printf("Price: %.2f rub.\n", result.price)
			fmt.Printf("\n")
			
			if (len(results) == c + 1) {
				fmt.Printf("\n\n\n")
			}
		}
	}
	
	main()
}

func get(article string) string {
	url := fmt.Sprintf("http://www.tr-auto.ru/catalog/search/?search=%s", article)
	
	request := gorequest.New()
	
	_, body, _ := request.Get(url).End()
	
	body = decode(body)
	
	return body
}

func scrape(html string) []result {
	results := []result{}

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	
	doc.Find(".catalog_table tr[id]").Each(func(i int, s *goquery.Selection) {
		result := result{
			article: trim(s.Find(".article").Text()),
			title: trim(s.Find(".title").Text()),
			vendor: trim(s.Find(".vender").Text()),
			inStock: trim(s.Find(".quantity").Text()) == "В наличии",
			price: priceToFloat(trim(s.Find(".price div.bx_price").Text())),
		}
		
		results = append(results, result)
	})
	
	return results
}

func trim(s string) string {
	return strings.TrimSpace(s)
}

func noSpaces(s string) string {
    return strings.Join(strings.Fields(s), "")
}

func boolToString(b bool) string {
	string := "Yes"
	
	if (b == false) {
		string = "No"
	}
	
	return string
}

func decode(s string) string {
	sr := strings.NewReader(s)
	tr := transform.NewReader(sr, charmap.Windows1251.NewDecoder())
	buf, _ := ioutil.ReadAll(tr)

	string := string(buf)
	
	return string
}

func priceToFloat(s string) float32 {
	s = strings.Replace(s, "руб.", "", -1)
	s = strings.Replace(s, ",", ".", -1)
	
	s = noSpaces(s)

	//fmt.Printf("%s \n", s)

	converted, _ := strconv.ParseFloat(s, 32)
	
	//fmt.Printf("%f\n\r", converted)

	return float32(converted)
}