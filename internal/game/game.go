package game

import (
	"fmt"
	"sort"
	"strings"

	"github.com/codemicro/cs-toptrumps/internal/cards"
	"github.com/codemicro/cs-toptrumps/internal/helpers"
	"github.com/codemicro/cs-toptrumps/internal/input"

	au "github.com/logrusorgru/aurora"
)

const NumPlayers = 2

type Game struct {
	decks [NumPlayers][]cards.Card
	selectedCardIndexes [NumPlayers]int
	selectedCards       [NumPlayers]*cards.Card
	selectedProperties  [NumPlayers]int
}

func New(decks [NumPlayers][]cards.Card) *Game {
	g := new(Game)
	g.decks = decks

	return g
}

func (g *Game) Run() {

	priorityPlayer := 0

	for { // Infinite loop time!

		g.selectedCards = [NumPlayers]*cards.Card{}

		for playerIndex := 0; playerIndex < NumPlayers; playerIndex += 1 {

			var cardStrings []string
			for _, v := range g.decks[playerIndex] {
				cardStrings = append(cardStrings, v.Name)
			}

			selectedCard, _ := input.Options(fmt.Sprintf("Player %d - pick a card!", playerIndex+1), cardStrings)

			g.selectedCardIndexes[playerIndex] = selectedCard
			g.selectedCards[playerIndex] = &g.decks[playerIndex][selectedCard]

			fmt.Println()

		}

		options := g.selectedCards[priorityPlayer].GetReadableNames()

		_, selectedOption := input.Options(fmt.Sprintf("Okay, player %d - select a property to challenge your opponent with!", priorityPlayer+1), options)

		fmt.Println()

		for playerIndex := 0; playerIndex < NumPlayers; playerIndex += 1 {

			selectedCard := g.selectedCards[playerIndex]
			propVal := selectedCard.GetValueByReadable(selectedOption)

			fmt.Printf("Player %d's %s has a %s of %d\n", playerIndex+1, au.Cyan(selectedCard.Name), au.Magenta(strings.ToLower(selectedOption)), au.Yellow(propVal))

			g.selectedProperties[playerIndex] = propVal
		}

		var winningPlayerNumber int

		{
			var selectedPropCopy []int
			for i := 0; i < NumPlayers; i += 1 {
				selectedPropCopy = append(selectedPropCopy, 0)
			}

			copy(selectedPropCopy, g.selectedProperties[:])
			sort.Ints(selectedPropCopy)

			lastIndex := len(selectedPropCopy) - 1

			if selectedPropCopy[lastIndex] == selectedPropCopy[lastIndex-1] {
				winningPlayerNumber = -1
			} else {
				for i, v := range g.selectedProperties {
					if selectedPropCopy[lastIndex] == v {
						winningPlayerNumber = i
						break
					}
				}
			}
		}

		fmt.Println()

		if winningPlayerNumber == -1 {
			fmt.Println("There was a draw!")
		} else {
			fmt.Printf("Player %d wins!\n", winningPlayerNumber+1)

			// Give the winning player all the losing players' cards

			for i := 0; i < NumPlayers; i += 1 {
				if i != winningPlayerNumber {
					g.decks[winningPlayerNumber] = append(g.decks[winningPlayerNumber], *g.selectedCards[i])
					g.decks[i] = append(g.decks[i][:g.selectedCardIndexes[i]], g.decks[i][g.selectedCardIndexes[i]+1:]...)
				}
			}

			// Check to see if the other player still has any cards left

			for i, deck := range g.decks {
				if len(deck) == 0 {

					var playerWithMostCards int
					{
						var highestCardCount int
						for x, deck := range g.decks {
							if len(deck) > highestCardCount {
								highestCardCount = len(deck)
								playerWithMostCards = x
							}
						}
					}

					fmt.Printf("\n     -----\nPlayer %d has run out of cards! The winner is player %d, as they have the most cards.\n", i+1, playerWithMostCards+1)

					return
				}
			}

		}

		fmt.Printf("Press <ENTER> to continue")
		fmt.Scanf("dfkjghkh")

		helpers.ClearConsole()

		priorityPlayer += 1
		if priorityPlayer == NumPlayers {
			priorityPlayer = 0
		}

	}
}
