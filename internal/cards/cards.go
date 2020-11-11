package cards

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

// Card information is originally defined in a JSON file - however, in order to make this simpler for the user (and so they
// don't have to worry about having specific files in the right place), we can bundle the contents of the `cards.json` file into
// our program at compile time. As a result, our game can be distributed as a single binary, instead of multiple files that have
// to be in specific locations relative to each other.

// There is currently no built in way to do this in Go, so we can use a tall called `go-bindata` to generate some Go code that
// contains our file and some helper functions to access it. Because this is a command line program, we can define the command
// that needs to be run here in order to generate our code file using a `go:generate` directive.

// In order to run this command, you can use the `go generate <pkgname>` command - for example, to generate the code for this
// package, you could run `go generate github.com/codemicro/cs-toptrumps/internal/cards`.

//go:generate go-bindata -pkg cards cards.json

type Card struct {
	Name string

	NumEngines int `readable:"Number of engines"` // <- this weird string thing is called a struct tag
	MaxPax     int `readable:"Maximum passenger count"`
	Range      int `readable:"Range"`
	Cost       int `readable:"Cost when new"`
}

// GetReadableNames iterates all attributes of a given card `c`, and generates a slice of any struct tag that has the readable
// field set. This is slightly confusing to do, but works well.
func (c Card) GetReadableNames() (names []string) {
	ct := reflect.TypeOf(c)
	for i := 0; i < ct.NumField(); i += 1 {
		field := ct.Field(i)
		tag := field.Tag.Get("readable")
		if tag != "" {
			names = append(names, tag)
		}
	}
	return // equivalent to `return names`
}

// GetValueByReadable takes the readable name of an attribute and returns the value of that attribute, if it exists. This
// assumes that only integer values have `readable` tags attached to them. If an attribute that is not an integer is read from
// by this function, it will panic and crash the application.
func (c Card) GetValueByReadable(readable string) int {
	ct := reflect.TypeOf(c)
	for i := 0; i < ct.NumField(); i += 1 {
		tField := ct.Field(i)
		tag := tField.Tag.Get("readable")
		if tag == readable {
			cv := reflect.ValueOf(c)
			vField := cv.Field(i)
			return vField.Interface().(int) // The panic would occur here, when the generic interface{} type is asserted into an
			// integer so it can be returned
		}
	}
	return 0
}

var (
	AllCards   []Card
	AvailCards []Card // Like AllCards, but is modified when cards are removed from the deck
)

// init runs automagically on package initialisation
func init() {
	// Load all card info from cards.json, which is a bundled file

	fCont := MustAsset("cards.json")

	err := json.Unmarshal(fCont, &AllCards)
	if err != nil {
		fmt.Println("Unable to load cards.json. Is the format correct?")
		panic(err)
	}

	AvailCards = AllCards
}

// Deal will select n cards from the deck of available cards at random, remove them from that deck, and return them in a new
// mini-deck
func Deal(n int) (deck []Card) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < n; i += 1 {
		chosenIndex := r.Intn(len(AvailCards))
		deck = append(deck, AvailCards[chosenIndex])

		AvailCards = append(AvailCards[:chosenIndex], AvailCards[chosenIndex+1:]...) // Remove chosen card from available deck
		// The `...` syntax means thge items of that slice/array are used as arguments to the `append` function.
	}

	return
}

// SplitCards will create two decks of even size and return those. Cards returned are removed from the deck of all available
// cards
func SplitCards(n int) (decks [][]Card) {
	numCards := len(AvailCards)

	// If there are less cards than there are decks to create, that's never going to work.
	// In the context of this program, this is only ever going to be caused by a programming error and not by anything that a
	// user inputs. Because of this, we don't need to go to all the hassle of properly handling an error, and can instead just
	// call panic and quit.
	if numCards < n {
		panic(fmt.Errorf("there are not enough available cards (have: %d) in order to create %d new deck(s)", numCards, n))
	}

	// While the number of cards in the deck doesn't divide evenly by the number of required decks, reduce the size of the pool
	// of cards to select from
	for numCards%n != 0 {
		numCards -= 1
	}

	cardsPerDeck := numCards / n

	for i := 0; i < n; i += 1 {
		decks = append(decks, Deal(cardsPerDeck))
	}

	return
}
