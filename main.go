package main

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/gocolly/colly/v2"
	"github.com/manifoldco/promptui"
)

type topic struct {
	Content string
	URL     string
}

type entry struct {
	Author   string
	Date     string
	Content  string
	FavCount string
}

var topicListString []string
var topicList []topic

var page int = 1

const baseURL string = "https://eksisozluk.com"

func validateActionInput(input string) error {
	if input == "e" || input == "n" {
		return nil
	}
	if input == "p" {
		if page <= 1 {
			return errors.New("currently you are on first page")
		}
		return nil
	}
	return errors.New("invalid action")
}

func main() {
	topics := promptui.Select{
		Label: "Başlıklar",
		Items: getTopics(),
		Size:  20,
	}
	selectedIndex, _, err := topics.Run()
	if err != nil {
		panic("Prompt failed " + err.Error())
	}
	fmt.Printf("\n\n")

	selectAction := promptui.Prompt{
		Label:    "Action (next: n, previous: p, exit: e)",
		Validate: validateActionInput,
	}
	for {
		entries := getEntries(topicList[selectedIndex].URL, fmt.Sprintf("%v", page))
		for _, entry := range entries {
			prettyPrint(entry)
		}

		selectedAction, err := selectAction.Run()
		if err != nil {
			panic("Prompt failed " + err.Error())
		}

		if selectedAction == "e" {
			break
		}
		switch selectedAction {
		case "p":
			page--
		case "n":
			page++
		}
	}
}

func prettyPrint(entry entry) {
	fmt.Println()

	color.Set(color.BgGreen, color.FgBlack)
	fmt.Printf("%s | %s | %s fav", entry.Author, entry.Date, entry.FavCount)
	color.Unset()
	fmt.Println()

	color.Set(color.BgBlack, color.Bold, color.FgWhite)
	fmt.Println(entry.Content)
	color.Unset()

	fmt.Println()
}

func getTopics() []string {
	c := colly.NewCollector()

	c.OnHTML("ul.topic-list > li > a", func(e *colly.HTMLElement) {
		newTopic := topic{}
		newTopic.URL = e.Attr("href")
		newTopic.Content = e.Text

		topicList = append(topicList, newTopic)
		topicListString = append(topicListString, e.Text)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Print("Visiting ", r.URL, "\n\n")
	})

	c.Visit(baseURL + "/basliklar/gundem")
	return topicListString
}

func getEntries(url, page string) []entry {
	var entries []entry
	c := colly.NewCollector()

	c.OnHTML("#entry-item-list > li", func(e *colly.HTMLElement) {
		newEntry := entry{}
		newEntry.Author = e.Attr("data-author")
		newEntry.Content = e.ChildText("div.content")
		newEntry.Date = e.ChildText("footer > div.info > a.entry-date.permalink")
		newEntry.FavCount = e.Attr("data-favorite-count")

		entries = append(entries, newEntry)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Print("Visiting ", r.URL, "\n\n")
	})

	c.Visit(baseURL + url + fmt.Sprintf("&p=%v", page))
	return entries
}
