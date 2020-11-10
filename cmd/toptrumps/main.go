package main

import (
	"fmt"
	"time"

	"github.com/codemicro/cs-toptrumps/internal/cards"
	"github.com/codemicro/cs-toptrumps/internal/game"
	"github.com/codemicro/cs-toptrumps/internal/helpers"
)

func main() {

	helpers.ClearConsole()

	fmt.Println("Top trumps, but it's planes and only has 8 cards")
	fmt.Println()
	time.Sleep(time.Second)

	a, b := cards.SplitCards()

	g := game.New([...][]cards.Card{a, b})
	g.Run()
}
