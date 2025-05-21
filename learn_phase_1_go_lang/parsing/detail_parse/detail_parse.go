package detail_parse

import (
	"encoding/json"
	"fmt"
	"go_lang/parsing/main_parse"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type MealDetail struct {
	Receipt     string       `json:"receipt"`
	Flag        string       `json:"flag"`
	Ingredents  []Ingredient `json:"ingredents"`
	ReceiptName string       `json:"receipt_name"`
}
type Ingredient struct {
	ImageURL string `json:"image_url"`
	Caption  string `json:"caption"`
}

var fullDetailMeal []MealDetail

func DetailParse() {
	file, _ := os.ReadFile("output/items.json")
	var loadedMeals []main_parse.Meal
	_ = json.Unmarshal(file, &loadedMeals)
	fmt.Printf("Loaded %d meals from output/items.json\n", len(loadedMeals))
	for _, meal := range loadedMeals {
		parseMeal(meal)
	}
	saveFinalMeals(fullDetailMeal)
}

func parseMeal(meal main_parse.Meal) {
	c := colly.NewCollector()

	c.OnHTML("section#feature", func(e *colly.HTMLElement) {
		html, _ := e.DOM.Html()
		start := strings.Index(html, "Instructions")
		end := strings.Index(html, "Browse More")
		if start == -1 || end == -1 || end <= start {
			fmt.Println("Could not find Instructions or Browse More section")
			return
		}
		start += len("Instructions")
		instructions := html[start:end]
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(instructions))

		text := strings.TrimSpace(doc.Text())

		flagURL := ""
		e.DOM.Find("img[src*='/images/icons/flags/big/64/']").Each(func(_ int, s *goquery.Selection) {
			src, exists := s.Attr("src")
			if exists && strings.Contains(src, "/images/icons/flags/big/64/") {
				flagURL = src
			}
		})
		// Extract country code from flagURL, e.g. "/images/icons/flags/big/64/gb.png" -> "gb"
		flagCode := getFlag(flagURL)
		// Parse ingredients: extract src and figcaption from each <figure>

		var ingredients []Ingredient
		// Use the full section DOM, not just the instructions fragment
		e.DOM.Find("figure").Each(func(_ int, s *goquery.Selection) {
			img := s.Find("img")
			src, _ := img.Attr("src")
			caption := s.Find("figcaption").Text()
			ingredients = append(ingredients, Ingredient{
				ImageURL: strings.TrimSpace(src),
				Caption:  strings.TrimSpace(caption),
			})
		})
		receiptName := ""
		//	<meta property="og:title" content="Apple Frangipan Tart Recipe">
		e.DOM.Closest("html").Find("meta[property='og:title']").Each(func(_ int, meta *goquery.Selection) {
			content, exists := meta.Attr("content")
			if exists {
				receiptName = strings.TrimSpace(content)
			}
		})
		finalResult := MealDetail{
			Receipt:     text,
			Flag:        flagCode,
			Ingredents:  ingredients,
			ReceiptName: receiptName,
		}
		fullDetailMeal = append(fullDetailMeal, finalResult)
	})
	err := c.Visit(meal.Href)
	if err != nil {
		fmt.Println("Error visiting page:", err)
	}
}

func saveFinalMeals(finalResult []MealDetail) bool {
	jsonData, err := json.MarshalIndent(finalResult, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return true
	}
	// Create output directory if it doesn't exist
	err = os.MkdirAll("output", 0755)
	if err != nil {
		fmt.Println("Error creating output directory:", err)
		return true
	}

	// Write JSON to output/items.json
	err = os.WriteFile("output/final_items.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return true
	}
	return false
}
func getFlag(flagURL string) string {
	if flagURL != "" {
		parts := strings.Split(flagURL, "/")
		if len(parts) > 0 {
			filename := parts[len(parts)-1]
			if strings.HasSuffix(filename, ".png") {
				flagURL = strings.TrimSuffix(filename, ".png")
			}
		}
	}
	return flagURL
}
