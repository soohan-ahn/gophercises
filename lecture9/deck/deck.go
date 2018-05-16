package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
)

type Card struct {
	Suit    string
	Numeric string
}

var SuitMap = map[string]int{
	"Spade":   1,
	"Diamond": 2,
	"Club":    3,
	"Heart":   4,
}

var NumericMap = map[string]int{
	"A": 1,
	"J": 11,
	"Q": 12,
	"K": 13,
}

type Deck []Card

func (d Deck) Len() int      { return len(d) }
func (d Deck) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d Deck) Less(i, j int) bool {
	iSuit := SuitMap[d[i].Suit]
	jSuit := SuitMap[d[j].Suit]
	if iSuit != jSuit {
		return iSuit < jSuit
	}

	iNumeric, err := strconv.Atoi(d[i].Numeric)
	if err != nil {
		iNumeric = NumericMap[d[i].Numeric]
	}
	jNumeric, err := strconv.Atoi(d[j].Numeric)
	if err != nil {
		jNumeric = NumericMap[d[j].Numeric]
	}
	return iNumeric < jNumeric
}

func (d Deck) Shuffle() {
	n := len(d)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		d.Swap(i, j)
	}
}

func (d Deck) AllCards() {
	for i, _ := range d {
		fmt.Printf("%s %s\n", d[i].Suit, d[i].Numeric)
	}
}

func New() Deck {
	var cards Deck
	for suit, _ := range SuitMap {
		for n, _ := range NumericMap {
			card := Card{
				Suit:    suit,
				Numeric: n,
			}
			cards = append(cards, card)
		}
		for i := 2; i <= 10; i++ {
			n := strconv.Itoa(i)
			card := Card{
				Suit:    suit,
				Numeric: n,
			}
			cards = append(cards, card)
		}
	}

	sort.Sort(cards)
	return cards
}
