package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/gophercises/lecture9/deck"
)

func initCards() deck.Deck {
	cards := deck.New()
	cards.Shuffle()
	return cards
}

var numberOfPlayer = 1

func printAllPlayerCards(allPlayerCards []deck.Deck, busted []bool, done []bool) {
	for i := 0; i < numberOfPlayer; i++ {
		fmt.Printf("Player[%d]'s Card: ", i)
		printPlayerCards(allPlayerCards[i], busted[i], done[i])
	}
}

func printPlayerCards(playerCards deck.Deck, busted bool, done bool) {
	for k := 0; k < len(playerCards); k++ {
		fmt.Printf("%s(%s) ", playerCards[k].Numeric, playerCards[k].Suit)
	}
	fmt.Printf("\tSum of cards: %d\t", SumOfCards(playerCards))
	fmt.Printf("Busted: %v, Done: %v\n", busted, done)
}

func SumOfCards(d deck.Deck) int {
	sum := 0
	for i := 0; i < len(d); i++ {
		iNumeric, err := strconv.Atoi(d[i].Numeric)
		if err != nil {
			iNumeric = deck.NumericMap[d[i].Numeric]
		}
		sum += iNumeric
	}

	if sum > 21 {
		for i := 0; i < len(d); i++ {
			if d[i].Numeric == "A" {
				sum -= 10
			}
			if sum <= 21 {
				break
			}
		}
	}
	return sum
}

func getCard(d deck.Deck) (deck.Card, deck.Deck) {
	if len(d) == 0 {
		d = initCards()
	}
	return d.Pop()
}

func main() {
	var cards deck.Deck
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Do you want to play the game? (Y/N) : ")
		text, _ := reader.ReadString('\n')
		if text[0] == 'N' || text[0] == 'n' {
			break
		}

		busted := make([]bool, numberOfPlayer)
		done := make([]bool, numberOfPlayer)
		var playerCards []deck.Deck
		var dealerCards deck.Deck

		for i := 0; i < numberOfPlayer; i++ {
			d := deck.Deck{}
			playerCards = append(playerCards, d)
		}

		for k := 0; k < 2; k++ {
			for i := 0; i < numberOfPlayer; i++ {
				var card deck.Card
				card, cards = getCard(cards)
				playerCards[i] = append(playerCards[i], card)
			}
		}

		for k := 0; k < 2; k++ {
			var card deck.Card
			card, cards = getCard(cards)
			dealerCards = append(dealerCards, card)
		}

		for {
			fmt.Printf("Dealers Card: %s(%s), ?(?)\n", dealerCards[0].Numeric, dealerCards[0].Suit)
			printAllPlayerCards(playerCards, busted, done)

			gameIsDone := true
			allBusted := true
			for i := 0; i < numberOfPlayer; i++ {
				if !busted[i] && !done[i] {
					fmt.Printf("Player[%d] want more card? : ", i)
					text, _ := reader.ReadString('\n')
					if text[0] == 'N' || text[0] == 'n' {
						done[i] = true
					} else {
						var card deck.Card
						card, cards = getCard(cards)
						playerCards[i] = append(playerCards[i], card)
						if SumOfCards(playerCards[i]) > 21 {
							done[i] = true
							busted[i] = true
						}
					}
				}
				gameIsDone = done[i]  // gameIsDone = false if !done[i]
				allBusted = busted[i] // allBusted = false if !busted[i]
			}

			if gameIsDone {
				printAllPlayerCards(playerCards, busted, done)

				if allBusted {
					fmt.Println("All of players busted. Dealer Won!")
					break
				}
				max := -1
				maxScorePlayer := -1
				for i := 0; i < numberOfPlayer; i++ {
					if !busted[i] {
						if s := SumOfCards(playerCards[i]); s > max {
							max = s
							maxScorePlayer = i
						}
					}
				}

				fmt.Printf("Dealer's Card: ")
				printPlayerCards(dealerCards, false, false)
				if dealerSum := SumOfCards(dealerCards); max > dealerSum {
					fmt.Printf("Player %d won!\n", maxScorePlayer)
				} else {
					fmt.Println("Dealer won!")
				}
				break
			}
		}
	}
}
