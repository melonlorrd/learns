package main

import "fmt"

type storyPage struct {
	text     string
	nextPage *storyPage
}

func (page *storyPage) playStory() {
	for page != nil {
		fmt.Println(page.text)
		page = page.nextPage
	}
}

func (page *storyPage) addToEnd(text string) {
	for page.nextPage != nil {
		page = page.nextPage
	}

	page.nextPage = &storyPage{text, nil}
}

func (page *storyPage) addAfter(text string) {
	newPage := &storyPage{text, page.nextPage}
	page.nextPage = newPage
}

func (page *storyPage) deleteNext() {
	page.nextPage = page.nextPage.nextPage
}

func main() {
	story := storyPage{"It was dark and stormy night", nil}

	story.addToEnd("You are alone, and you need to find the sacred helmet before the bad guys do")
	story.addToEnd("You see a troll ahead")

	story.addAfter("Woooo yeah baby")
	story.deleteNext()

	story.playStory()
}
