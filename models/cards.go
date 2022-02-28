package models

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrCardCodeValueInvalid = errors.New("card code value is invalid")
	ErrCardCodeSuitInvalid  = errors.New("card code suit is invalid")
)

type Suit string

const (
	SuitClubs    = Suit("CLUBS")
	SuitDiamonds = Suit("DIAMONDS")
	SuitHearts   = Suit("HEARTS")
	SuitSpades   = Suit("SPADES")

	ValueAce   = Value("ACE")
	ValueJack  = Value("JACK")
	ValueQueen = Value("QUEEN")
	ValueKing  = Value("KING")
)

var suitsCodeMap = map[rune]Suit{
	'C': SuitClubs,
	'D': SuitDiamonds,
	'H': SuitHearts,
	'S': SuitSpades,
}

var alphaValuesCodeMap = map[rune]Value{
	'A': ValueAce,
	'J': ValueJack,
	'Q': ValueQueen,
	'K': ValueKing,
}

var suits = []Suit{
	SuitSpades,
	SuitDiamonds,
	SuitClubs,
	SuitHearts,
}

var values = []Value{
	ValueAce,
	Value("2"),
	Value("3"),
	Value("4"),
	Value("5"),
	Value("6"),
	Value("7"),
	Value("8"),
	Value("9"),
	Value("10"),
	ValueJack,
	ValueQueen,
	ValueKing,
}

// Code returns first character of suit as code
func (s Suit) Code() string {
	return string(s[0])
}

type Value string

// Code returns first character of value if value is not numeric
// returns value if numeric
func (v Value) Code() string {
	if _, err := strconv.Atoi(string(v)); err != nil {
		return string([]rune(v)[0])
	}
	return string(v)
}

// Card representation of Card
type Card struct {
	Value Value  `json:"value"`
	Suit  Suit   `json:"suit"`
	Code  string `json:"code"`
}

// NewCard returns new card instance
func NewCard(value Value, suit Suit) *Card {
	return &Card{
		Value: value,
		Suit:  suit,
		Code:  value.Code() + suit.Code(),
	}
}

// CardStorage is used to interact with cards storage
type CardStorage interface {
	// ByCode is used to retrieve single card by code
	ByCode(code string) (*Card, error)
	// All is used to retrieve all cards
	All() ([]*Card, error)
}

type CardService interface {
	CardStorage
	ByCodesStr(codesStr string) ([]*Card, error)
}

// NewCardService returns new CardService instance
func NewCardService() CardService {
	scs := staticCardStorage{}
	cv := cardValidator{&scs}
	return &cardService{&cv}
}

type cardService struct {
	CardStorage
}

// ByCodesStr is used to get codes by comma seperated codes string
func (cs *cardService) ByCodesStr(codesStr string) ([]*Card, error) {
	codes := strings.Split(strings.TrimSpace(codesStr), ",")
	return cs.ByCodes(codes)
}

// ByCodes is used to get cards by code list
// Returns error from storage
func (cs *cardService) ByCodes(codes []string) ([]*Card, error) {
	cards := make([]*Card, len(codes))
	for i := range codes {
		card, err := cs.CardStorage.ByCode(codes[i])
		if err != nil {
			return nil, err
		}
		cards[i] = card
	}
	return cards, nil
}

type cardValFunc func(card *Card) error

func runCardValFuncs(card *Card, fns ...cardValFunc) error {
	for _, fn := range fns {
		if err := fn(card); err != nil {
			return err
		}
	}
	return nil
}

type cardValidator struct {
	CardStorage
}

func (cv *cardValidator) normalizeCode(card *Card) error {
	card.Code = strings.ToUpper(strings.TrimSpace(card.Code))
	return nil
}

func (cv *cardValidator) checkCodeValue(card *Card) error {
	if card.Code != "" {
		r := []rune(card.Code)
		if _, ok := alphaValuesCodeMap[r[0]]; !ok {
			numVal, err := strconv.Atoi(string(r[:len(r)-1]))
			if err != nil {
				return ErrCardCodeValueInvalid
			}
			if numVal > 10 || numVal < 1 {
				return ErrCardCodeValueInvalid
			}
		}
	}
	return nil
}

func (cv *cardValidator) checkCodeSuit(card *Card) error {
	if card.Code != "" {
		r := []rune(card.Code)
		if _, ok := suitsCodeMap[r[len(r)-1]]; !ok {
			return ErrCardCodeSuitInvalid
		}
	}
	return nil
}

// ByCode is used to get Card by code
// Normalizes code before search
// Returns ErrCardCodeValueInvalid error if value part of the code is not a valid card value
// Returns ErrCardCodeSuitInvalid error if suit part of the code is not a valid card suit
func (cv *cardValidator) ByCode(code string) (*Card, error) {
	card := Card{Code: code}
	if err := runCardValFuncs(&card,
		cv.normalizeCode,
		cv.checkCodeValue,
		cv.checkCodeSuit,
	); err != nil {
		return nil, err
	}
	return cv.CardStorage.ByCode(card.Code)
}

type staticCardStorage struct{}

// ByCode is used to get Card by code
func (scs *staticCardStorage) ByCode(code string) (*Card, error) {
	r := []rune(code)
	suit := suitsCodeMap[r[len(r)-1]]
	value, ok := alphaValuesCodeMap[r[0]]
	if !ok {
		value = Value(r[:len(r)-1])
	}

	return NewCard(value, suit), nil
}

// All is used get all cards
func (scs *staticCardStorage) All() ([]*Card, error) {
	var cards []*Card
	for _, suit := range suits {
		for _, value := range values {
			cards = append(cards, NewCard(value, suit))
		}
	}
	return cards, nil
}
