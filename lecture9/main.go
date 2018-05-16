package main

import (
	"fmt"
	"github.com/gophercises/lecture9/deck"
)

func main() {
	cards := deck.New()
	fmt.Println("Before Shuffle")
	cards.AllCards()
	cards.Shuffle()
	fmt.Println("After Shuffle")
	cards.AllCards()
}
