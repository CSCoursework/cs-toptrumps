package cards

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

//go:generate go-bindata -pkg cards cards.json

type Card struct {
	Name string

	NumEngines int `readable:"Number of engines"`
	MaxPax     int `readable:"Maximum passenger count"`
	Range      int `readable:"Range"`
	Cost       int `readable:"Cost when new"`
}

func (c Card) GetReadableNames() (names []string) {
	ct := reflect.TypeOf(c)
	for i := 0; i < ct.NumField(); i += 1 {
		field := ct.Field(i)
		tag := field.Tag.Get("readable")
		if tag != "" {
			names = append(names, tag)
		}
	}
	return
}

func (c Card) GetValueByReadable(readable string) int {
	ct := reflect.TypeOf(c)
	for i := 0; i < ct.NumField(); i += 1 {
		tField := ct.Field(i)
		tag := tField.Tag.Get("readable")
		if tag == readable {
			cv := reflect.ValueOf(c)
			vField := cv.Field(i)
			return vField.Interface().(int)
		}
	}
	return 0
}

var (
	AllCards   []Card
	AvailCards []Card
)

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

// Deal will select n cards from the deck of available cards, remove them from that deck, and return them
func Deal(n int) (deck []Card) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < n; i += 1 {
		chosenIndex := r.Intn(len(AvailCards))
		deck = append(deck, AvailCards[chosenIndex])

		AvailCards = append(AvailCards[:chosenIndex], AvailCards[chosenIndex+1:]...) // Remove chosen card from available deck
	}

	return
}

// SplitCards will create two decks of even size and return those. Cards returned are removed from the deck of all available cards
func SplitCards() (deck1 []Card, deck2 []Card) {
	numCards := len(AvailCards)

	if numCards%2 != 0 {
		numCards -= 1
	}

	cardsPerDeck := numCards / 2

	deck1 = Deal(cardsPerDeck)
	deck2 = Deal(cardsPerDeck)

	return
}
