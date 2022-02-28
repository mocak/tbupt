package models

import (
	"errors"
	"github.com/google/uuid"
	"math/rand"
)

var (
	ErrNotEnoughCards = errors.New("there is not enough cards in the deck for the operation")
	ErrDeckOpened     = errors.New("not permitted on opened deck")
	ErrNotFound       = errors.New("resource not found")
	ErrUUIDRequired   = errors.New("uuid is required")
	ErrUUIDInvalid    = errors.New("uuid is not valid")
)

// Deck is the representation of deck of cards
type Deck struct {
	UUID      string  `json:"deck_id"`
	Shuffled  bool    `json:"shuffled"`
	Remaining int     `json:"remaining"`
	Cards     []*Card `json:"cards"`
	CardCodes string  `json:"-"`
	Opened    bool    `json:"-"`
}

type DeckStorage interface {
	Create(deck *Deck) error
	ByUUID(uuid string) (*Deck, error)
	Update(deck *Deck) error
}

type DeckService interface {
	DeckStorage
	Draw(deck *Deck, count int) ([]*Card, error)
	Open(deck *Deck) error
}

// NewDeckService returns deckService instance by defaults
func NewDeckService(cs CardService) DeckService {
	dm := deckMemory{decks: map[string]Deck{}}
	dv := newDeckValidator(&dm, cs)
	return &deckService{
		DeckStorage: dv,
	}
}

type deckService struct {
	DeckStorage
}

// Draw is used to release given amount of cards from the top of the given deck
// Returns ErrDeckOpened if deck is opened before
// Returns ErrNotEnoughCards if deck has not enough cards to draw
// Returns error from DeckStorage if fails
func (ds *deckService) Draw(deck *Deck, count int) ([]*Card, error) {

	if deck.Opened {
		return nil, ErrDeckOpened
	}

	if count > deck.Remaining {
		return nil, ErrNotEnoughCards
	}

	cards := deck.Cards[:count]
	deck.Cards = deck.Cards[count:]
	if err := ds.DeckStorage.Update(deck); err != nil {
		return nil, err
	}
	return cards, nil
}

// Open sets deck status to opened
func (ds *deckService) Open(deck *Deck) error {
	deck.Opened = true
	return ds.DeckStorage.Update(deck)
}

type deckValFunc func(*Deck) error

func runDeckValFuncs(deck *Deck, fns ...deckValFunc) error {
	for _, fn := range fns {
		if err := fn(deck); err != nil {
			return err
		}
	}
	return nil
}

type deckValidator struct {
	DeckStorage
	cs CardService
}

func newDeckValidator(ds DeckStorage, cs CardService) *deckValidator {
	return &deckValidator{
		DeckStorage: ds,
		cs:          cs,
	}
}

func (dv *deckValidator) shuffle(deck *Deck) error {
	if deck.Shuffled == true {
		cards := make([]*Card, deck.Remaining)
		perm := rand.Perm(deck.Remaining)
		for idx, permIdx := range perm {
			cards[idx] = deck.Cards[permIdx]
		}
		deck.Cards = cards
	}
	return nil
}

func (dv *deckValidator) setUUIDIfUnset(deck *Deck) error {
	if deck.UUID == "" {
		deck.UUID = uuid.NewString()
	}
	return nil
}

func (dv *deckValidator) setCardsByCodes(deck *Deck) error {
	if deck.CardCodes != "" {
		cards, err := dv.cs.ByCodesStr(deck.CardCodes)
		if err != nil {
			return err
		}
		deck.CardCodes = ""
		deck.Cards = cards
	}
	return nil
}

func (dv *deckValidator) setCardsIfEmpty(deck *Deck) error {
	if deck.Cards == nil {
		cards, err := dv.cs.All()
		if err != nil {
			return err
		}
		deck.Cards = cards
	}
	return nil
}

func (dv *deckValidator) setRemaining(deck *Deck) error {
	deck.Remaining = len(deck.Cards)
	return nil
}

func (dv *deckValidator) isValidUUID(deck *Deck) error {
	_, err := uuid.Parse(deck.UUID)
	if err != nil {
		return ErrUUIDInvalid
	}
	return nil
}

func (dv *deckValidator) requireUUID(deck *Deck) error {
	if deck.UUID == "" {
		return ErrUUIDRequired
	}
	return nil
}

// Create will create the provided deck and fill data
// like the UUID, Remaining, Cards fields.
func (dv deckValidator) Create(deck *Deck) error {
	err := runDeckValFuncs(deck,
		dv.setUUIDIfUnset,
		dv.isValidUUID,
		dv.setCardsByCodes,
		dv.setCardsIfEmpty,
		dv.setRemaining,
		dv.shuffle,
	)
	if err != nil {
		return err
	}

	return dv.DeckStorage.Create(deck)
}

// Update updates matching deck in the storage by provided deck
func (dv deckValidator) Update(deck *Deck) error {
	err := runDeckValFuncs(deck,
		dv.requireUUID,
		dv.setRemaining,
	)
	if err != nil {
		return err
	}

	return dv.DeckStorage.Update(deck)
}

// ByUUID retrieves deck from storage by provided uuid
func (dv *deckValidator) ByUUID(uuid string) (*Deck, error) {
	deck := Deck{UUID: uuid}
	err := runDeckValFuncs(&deck,
		dv.requireUUID,
		dv.isValidUUID,
	)
	if err != nil {
		return nil, err
	}

	return dv.DeckStorage.ByUUID(deck.UUID)
}

type deckMemory struct {
	decks map[string]Deck
}

// Create persists given deck to storage
func (dm *deckMemory) Create(deck *Deck) error {
	dm.decks[deck.UUID] = *deck
	return nil
}

// ByUUID finds and returns deck by give uuid
func (dm *deckMemory) ByUUID(uuid string) (*Deck, error) {
	deck, ok := dm.decks[uuid]
	if !ok {
		return nil, ErrNotFound
	}
	return &deck, nil
}

// Update updates matching deck in the storage by given deck
func (dm *deckMemory) Update(deck *Deck) error {
	dm.decks[deck.UUID] = *deck
	return nil
}
