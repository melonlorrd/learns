package main

import (
	"bufio"
	"fmt"
	"os"
)

type storyNode struct {
	text    string
	yesPath *storyNode
	noPath  *storyNode
}

func (node *storyNode) play() {
	fmt.Println(node.text)

	if node.yesPath == nil || node.noPath == nil {
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		answer := scanner.Text()

		switch answer {
		case "yes":
			node.yesPath.play()
			return
		case "no":
			node.noPath.play()
			return
		default:
			fmt.Println("That answer was not an option. Please answer again (yes/no).")
		}
	}
}

func main() {
	root := storyNode{"You are at entrance of a dark cave. Do you want to do into the cave?", nil, nil}
	winning := storyNode{"You have won!", nil, nil}
	losing := storyNode{"You have lost!", nil, nil}

	root.yesPath = &losing
	root.noPath = &winning

	root.play()
}
