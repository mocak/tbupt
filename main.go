package main

import (
	"github.com/mocak/tbupt/controllers"
	"github.com/mocak/tbupt/models"
	"net/http"
)

func main() {
	cardService := models.NewCardService()
	deckService := models.NewDeckService(cardService)
	deckController := controllers.NewDecks(deckService)
	r := controllers.NewServer(deckController)

	http.ListenAndServe(":3000", r)
}
