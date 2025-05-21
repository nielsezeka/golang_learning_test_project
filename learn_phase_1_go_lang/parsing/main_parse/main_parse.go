package main_parse

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type Meal struct {
	Name string `json:"name"`
	Href string `json:"href"`
}

var Meals []Meal

func ThumnailPageParse() {
	c := colly.NewCollector()
	c.OnHTML("section#feature .row .col-sm-3", func(e *colly.HTMLElement) {
		meal := e.Text
		href := e.ChildAttr("a", "href")
		Meals = append(Meals, Meal{Name: meal, Href: href})
	})
	for ch := 'a'; ch <= 'z'; ch++ {
		err := c.Visit(fmt.Sprintf("https://www.themealdb.com/browse/letter/%c", ch))
		if err != nil {
			fmt.Println("Error visiting page:", err)
		}
		fmt.Printf("Processed complete char:[%c]\n", ch)
	}

	jsonData, err := json.MarshalIndent(Meals, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}
	err = os.MkdirAll("output", 0755)
	if err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	// Write JSON to output/items.json
	err = os.WriteFile("output/items.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return
	}
	fmt.Println("Data written to output/items.json")
}
