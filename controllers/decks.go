package controllers

import (
	"github.com/gorilla/mux"
	"github.com/mocak/tbupt/json"
	"github.com/mocak/tbupt/models"
	"net/http"
)

func NewDecks(ds models.DeckService) *Decks {
	return &Decks{
		ds: ds,
	}
}

type Decks struct {
	ds models.DeckService
}

type createResponse struct {
	DeckID    string
	Shuffled  bool
	Remaining int
}

// Create is used to create deck resource
// replies the request with created deck resource info and HTTP 201 code if succeed
//
// POST /deck
func (d *Decks) Create(w http.ResponseWriter, r *http.Request) {
	deck := models.Deck{}
	err := json.DecodeBody(w, r, &deck)
	if err != nil {
		return
	}

	deck.CardCodes = r.URL.Query().Get("cards")
	if err := d.ds.Create(&deck); err != nil {
		json.Error(w, "Unexpected Error", http.StatusInternalServerError)
		return
	}

	cr := createResponse{
		DeckID:    deck.UUID,
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining,
	}

	json.Response(w, cr, http.StatusCreated)
}

// Open is used the open deck
// Replies the request with deck resource with cards and HTTP 200 if succeed
//
// PUT /draw/:uid/open
func (d *Decks) Open(w http.ResponseWriter, r *http.Request) {
	deck, err := d.deckByUUID(w, r)
	if err != nil {
		return
	}
	if err := d.ds.Open(deck); err != nil {
		json.Error(w, "Unexpected Error", http.StatusInternalServerError)
		return
	}
	json.Response(w, deck, http.StatusOK)
}

type drawRequest struct {
	Count int `json:"count"`
}

// Draw used to draw cards from deck resource
// Replies the request with drawn card resources and HTTP 200
//
// POST /deck/:uid/draw
func (d *Decks) Draw(w http.ResponseWriter, r *http.Request) {

	drawReq := drawRequest{}

	err := json.DecodeBody(w, r, &drawReq)
	if err != nil {
		return
	}

	deck, err := d.deckByUUID(w, r)
	if err != nil {
		return
	}

	cards, err := d.ds.Draw(deck, drawReq.Count)
	if err != nil {
		json.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.Response(w, cards, http.StatusOK)
}

// deckByUUID used to get models.Deck record by URL
// Returns matched models.Deck record if found
// Returns error models.ErrNotFound if record not found
// Returns related if another error occurs
// Sets error response if fails
func (d *Decks) deckByUUID(w http.ResponseWriter, r *http.Request) (*models.Deck, error) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	deck, err := d.ds.ByUUID(uuid)
	if err != nil {
		if err == models.ErrNotFound {
			json.Error(w, "Deck not found", http.StatusNotFound)
		} else {
			json.Error(w, "Unexpected Error", http.StatusInternalServerError)
		}
		return nil, err
	}

	return deck, nil
}
