package main

// Homework
// maintain lamp state
// NPC - talk to them, fight
// NPC move around in graph
// items that can pick up or placed down
// accept natural language with input
// checkpoint
// slices for list of choice

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type choices struct {
	cmd         string
	description string
	nextNode    *storyNode
	nextChoice  *choices
}

type storyNode struct {
	text    string
	choices *choices
}

func (node *storyNode) addChoice(cmd string, description string, nextNode *storyNode) {
	choice := &choices{
		cmd:         cmd,
		description: description,
		nextNode:    nextNode,
		nextChoice:  nil,
	}

	if node.choices == nil {
		node.choices = choice
	} else {
		currentChoice := node.choices
		for currentChoice.nextChoice != nil {
			currentChoice = currentChoice.nextChoice
		}
		currentChoice.nextChoice = choice
	}
}

func (node *storyNode) render() {
	fmt.Println(node.text)
	currentChoice := node.choices
	for currentChoice != nil {
		fmt.Println(currentChoice.cmd, ":", currentChoice.description)
		currentChoice = currentChoice.nextChoice
	}
}

func (node *storyNode) executeCmd(cmd string) *storyNode {
	currentChoice := node.choices
	for currentChoice != nil {
		if strings.EqualFold(currentChoice.cmd, cmd) {
			return currentChoice.nextNode
		}
		currentChoice = currentChoice.nextChoice
	}
	fmt.Println("Sorry, I did not understand that.")
	return node
}

var scanner *bufio.Scanner

func (node *storyNode) play() {
	node.render()
	if node.choices != nil {
		scanner.Scan()
		nextNode := node.executeCmd(scanner.Text())
		nextNode.play()
	}
}

func main() {
	scanner = bufio.NewScanner(os.Stdin)

	start := storyNode{
		text: `
		You are in large chamber, deep underground.
		You see three passages leading out. A north passage leads into darkness.
		To the south, a passage appears to head upwards. The eastern passages appears.
		Flat and well travelled.
		`,
	}

	darkRoom := storyNode{text: "It is pitch black dark. You cannot see a thing."}
	darkRoomLit := storyNode{text: "The dark passage is now lit by your lantern. You can continue north or head back south."}
	yeti := storyNode{text: "When stumbling around in darkness, you are eaten by a Yeti."}
	trap := storyNode{text: "You head down the well travelled path, suddeny a trap door opens and you fall into a pit."}
	treasure := storyNode{text: "You arrive at a small chamber, filled with treasure!"}

	start.addChoice("N", "Go North", &darkRoom)
	start.addChoice("S", "Go South", &darkRoomLit)
	start.addChoice("E", "Go East", &trap)

	darkRoom.addChoice("S", "Try to go back South", &yeti)
	darkRoom.addChoice("O", "Turn on lantern", &darkRoomLit)

	darkRoomLit.addChoice("N", "Go North", &treasure)
	darkRoomLit.addChoice("S", "Go South", &start)

	start.play()

	fmt.Printf("\nThe End.\n")
}
