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
	decks               [NumPlayers][]cards.Card
	selectedCardIndexes [NumPlayers]int
	selectedCards       [NumPlayers]*cards.Card
	selectedProperties  [NumPlayers]int
}

// New creates a new instance of `Game`, based on the decks provided.
func New(decks [][]cards.Card) *Game {

	g := new(Game)

	// `g.decks` is a fixed-length array, whereas the `decks` argument is a variable-length slice of decks. The easiest way turn
	// this slice into the required array it just to iterate over it, adding elements as we go.
	// If there are less decks than the number of players, this will panic. This is acceptable, as (in the context of this
	// program) this is only ever caused by a programming error, never by something the user has inputted. Were it something the
	// user did, proper error handling would be put in place.
	for i := 0; i < NumPlayers; i += 1 {
		g.decks[i] = decks[i]
	}

	return g
}

// Run runs the game. Shocking, I know.
func (g *Game) Run() {

	// priorityPlayer is used to determine which player gets to choose the card property to compare with.
	// It's iterated and/or reset at the end of each round, and it corresponds to an index in the arrays found in the `Game`
	// struct.
	priorityPlayer := 0

	for { // Infinite loop time!

		// Reset the `g.SelectedCards` array to blank. This is not always the case because it is not cleared at the end of a round.
		g.selectedCards = [NumPlayers]*cards.Card{}

		// ----- Ask each player which card they'd like to use out of all the ones in their deck -----

		for playerIndex := 0; playerIndex < NumPlayers; playerIndex += 1 {

			// Make a slice of strings to show to the user, to represent each card struct.
			var cardStrings []string
			for _, v := range g.decks[playerIndex] {
				cardStrings = append(cardStrings, v.Name)
			}

			selectedCard, _ := input.Options(fmt.Sprintf("Player %d - pick a card!", playerIndex+1), cardStrings)

			g.selectedCardIndexes[playerIndex] = selectedCard // this, while it works, it me being lazy and not wanting to do
			// any searching later on in the round logic
			g.selectedCards[playerIndex] = &g.decks[playerIndex][selectedCard]

			fmt.Println()

		}

		// ----- Prompt the player with priority to select a property -----

		options := g.selectedCards[priorityPlayer].GetReadableNames()

		_, selectedOption := input.Options(fmt.Sprintf("Okay, player %d - select a property to challenge your opponent with!",
			priorityPlayer+1), options)

		fmt.Println()

		// ----- Gather values for each of those properties out of all selected cards -----

		for playerIndex, selectedCard := range g.selectedCards {

			propVal := selectedCard.GetValueByReadable(selectedOption)

			// au.ColourName functions only exist to give the output some funky colours with ANSI colour codes. They take a 
			// string in and return a string wrapped with the corresponding colour codes.

			fmt.Printf("Player %d's %s has a %s of %d\n", playerIndex+1, au.Cyan(selectedCard.Name),
				au.Magenta(strings.ToLower(selectedOption)), au.Yellow(propVal))

			g.selectedProperties[playerIndex] = propVal
		}

		// ----- Determine which player won -----

		var winningPlayerNumber int

		{ // These curly brackets are here to scope a load of random, temporary variables to this section of the logic, as to
			// not pollute the namespace elsewhere in the app.
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
					// Code like this is code that I can write really easily when I have a crystal clear mental model of the
					// data structures involved, and code that I can't do anything with after I initially wrote it because I
					// have no idea what the hell is going on
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
		fmt.Scanf("dfkjghkh") // I don't quite know why it needed a string of random characters here, but it made it work so
		// I'm just going to back away slowly and leave it alone

		helpers.ClearConsole()

		priorityPlayer += 1
		if priorityPlayer == NumPlayers {
			priorityPlayer = 0
		}

	}
}
