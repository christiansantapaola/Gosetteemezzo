package main

import (
	"fmt"
	"sort"

	"github.com/qrowsxi/deck"
)

// CardToPoint will calculate the point for a single card
func CardToPoint(card *deck.Card) float64 {
	if card == nil {
		return 0
	}
	if card.Value > 7 {
		return 0.5
	}
	return float64(card.Value)
}

// Game will play a game of "sette e mezzo" and try to draw `stop` card
// if it lose it will exit early and return 0 point and the number of cards drawed.
func Game(deck *deck.Deck, stop int) (float64, int) {
	deck.Shuffle()
	point := 0.0
	for j := 0; j < stop-1; j++ {
		card := deck.Draw()
		point += CardToPoint(card)
		if point > 7.5 {
			deck.Reset()
			return 0.0, j
		}
		if point == 7.5 {
			deck.Reset()
			return -1.0, j
		}
	}
	card := deck.Draw()
	point += CardToPoint(card)
	deck.Reset()
	return point, stop
}

// game will play a game of "sette e mezzo",
// it will shuffle the deck and sent the result over a channel.
func game(ch chan map[float64]int, stop, no int) {
	var result map[float64]int = make(map[float64]int)
	deck := deck.NewDeck()
	for i := 0; i < no; i++ {
		point, _ := Game(deck, stop)
		result[point]++
	}
	ch <- result
}

func main() {
	var no int = 100000000
	//var stop int = 4
	var noThread int = 4
	var ch chan map[float64]int = make(chan map[float64]int)
	var isFirst bool = true
	// Play the game and stop to 1 to 4 card
	for stop := 1; stop < 5; stop++ {
		var result map[float64]int = make(map[float64]int)
		// play the game, divide it in multiple thread for performance
		for i := 0; i < noThread; i++ {
			go game(ch, stop, no/noThread)
		}
		// get the result for the various thread and merge them
		for i := 0; i < noThread; i++ {
			res := <-ch
			for key := range res {
				result[key] += res[key]
			}
		}
		// sort the key (the point) from the result for printing them in order.
		var sortedKey []float64
		for key := range result {
			sortedKey = append(sortedKey, key)
		}
		sort.Float64s(sortedKey)
		// Print csv header
		if isFirst {
			fmt.Println(fmt.Sprintf("%s,%s,%s,%s,%s", "Total", "Cards", "Points", "Cases", "Prob"))
			isFirst = false
		}
		// Print the result in a csv format to stdout
		for _, key := range sortedKey {
			fmt.Println(fmt.Sprintf("%d,%d,%.1f,%d,%.3f",
				no,
				stop,
				key,
				result[key],
				float64(result[key])/float64(no)))
		}
	}
}
