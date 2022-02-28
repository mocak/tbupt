package controllers

import (
	"errors"
	"github.com/mocak/tbupt/models"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNewDecks(t *testing.T) {
	type args struct {
		ds models.DeckService
	}
	tests := []struct {
		name string
		args args
		want *Decks
	}{
		{
			name: "default",
			args: args{ds: models.NewDeckService(models.NewCardService())},
			want: &Decks{ds: models.NewDeckService(models.NewCardService())},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDecks(tt.args.ds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDecks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecks_Create(t *testing.T) {
	type fields struct {
		ds models.DeckService
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       string
		wantStatus int
	}{
		{
			name: "valid",
			fields: fields{
				ds: mockDeckService{deck: &models.Deck{UUID: "testuuid", Remaining: 3}},
			},
			args: args{
				r: httptest.NewRequest("POST", "/draw", strings.NewReader("{\"shuffled\":true}")),
				w: httptest.NewRecorder(),
			},
			want:       "{\"DeckID\":\"testuuid\",\"Shuffled\":true,\"Remaining\":3}",
			wantStatus: http.StatusCreated,
		},
		{
			name: "create fails",
			fields: fields{
				ds: mockDeckService{err: errors.New(""), deck: &models.Deck{}},
			},
			args: args{
				r: httptest.NewRequest("POST", "/deck", strings.NewReader("{\"a\":\"b\"}")),
				w: httptest.NewRecorder(),
			},
			want:       "\"Unexpected Error\"\n",
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Decks{
				ds: tt.fields.ds,
			}
			d.Create(tt.args.w, tt.args.r)
			resp := tt.args.w.Result()
			defer resp.Body.Close()
			byteSlice, _ := io.ReadAll(resp.Body)
			got := string(byteSlice)
			gotStatus := resp.StatusCode
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestResponse() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(gotStatus, tt.wantStatus) {
				t.Errorf("TestResponse() status code = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestDecks_Draw(t *testing.T) {
	cards := []*models.Card{
		{Value: "ACE", Suit: "SPADES", Code: "AS"},
		{Value: "2", Suit: "SPADES", Code: "2S"},
	}
	type fields struct {
		ds models.DeckService
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       string
		wantStatus int
	}{
		{
			name: "valid",
			fields: fields{
				ds: mockDeckService{
					deck:  &models.Deck{UUID: "testuuid", Remaining: 3},
					cards: cards,
				},
			},
			args: args{
				r: httptest.NewRequest("POST", "/deck/testuuid/draw", strings.NewReader("{\"count\":2}")),
				w: httptest.NewRecorder(),
			},
			want:       "[{\"value\":\"ACE\",\"suit\":\"SPADES\",\"code\":\"AS\"},{\"value\":\"2\",\"suit\":\"SPADES\",\"code\":\"2S\"}]",
			wantStatus: http.StatusOK,
		},
		{
			name: "draw fail",
			fields: fields{
				ds: mockDeckService{
					deck: &models.Deck{UUID: "testuuid", Remaining: 3},
					err:  errors.New("error"),
				},
			},
			args: args{
				r: httptest.NewRequest("POST", "/deck/testuuid/draw", strings.NewReader("{\"count\":2}")),
				w: httptest.NewRecorder(),
			},
			want:       "\"Unexpected Error\"\n",
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Decks{
				ds: tt.fields.ds,
			}
			d.Draw(tt.args.w, tt.args.r)
			resp := tt.args.w.Result()
			defer resp.Body.Close()
			byteSlice, _ := io.ReadAll(resp.Body)
			got := string(byteSlice)
			gotStatus := resp.StatusCode
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestResponse() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(gotStatus, tt.wantStatus) {
				t.Errorf("TestResponse() status code = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestDecks_Open(t *testing.T) {
	cards := []*models.Card{
		{Value: "ACE", Suit: "SPADES", Code: "AS"},
		{Value: "2", Suit: "SPADES", Code: "2S"},
	}
	type fields struct {
		ds models.DeckService
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       string
		wantStatus int
	}{
		{
			name: "valid",
			fields: fields{
				ds: mockDeckService{
					deck:  &models.Deck{UUID: "testuuid", Remaining: 3, Cards: cards},
					cards: cards,
				},
			},
			args: args{
				r: httptest.NewRequest("PUT", "/deck/testuuid/open", strings.NewReader("{\"count\":2}")),
				w: httptest.NewRecorder(),
			},
			want:       "{\"deck_id\":\"testuuid\",\"shuffled\":false,\"remaining\":3,\"cards\":[{\"value\":\"ACE\",\"suit\":\"SPADES\",\"code\":\"AS\"},{\"value\":\"2\",\"suit\":\"SPADES\",\"code\":\"2S\"}]}",
			wantStatus: http.StatusOK,
		},
		{
			name: "open fail",
			fields: fields{
				ds: mockDeckService{
					deck: &models.Deck{UUID: "testuuid", Remaining: 3},
					err:  errors.New("error"),
				},
			},
			args: args{
				r: httptest.NewRequest("PUT", "/deck/testuuid/draw", strings.NewReader("{\"count\":2}")),
				w: httptest.NewRecorder(),
			},
			want:       "\"Unexpected Error\"\n",
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Decks{
				ds: tt.fields.ds,
			}
			d.Open(tt.args.w, tt.args.r)
			resp := tt.args.w.Result()
			defer resp.Body.Close()
			byteSlice, _ := io.ReadAll(resp.Body)
			got := string(byteSlice)
			gotStatus := resp.StatusCode
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestResponse() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(gotStatus, tt.wantStatus) {
				t.Errorf("TestResponse() status code = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

type mockDeckService struct {
	deck  *models.Deck
	err   error
	cards []*models.Card
}

func (m mockDeckService) Update(deck *models.Deck) error {
	return m.err
}

func (m mockDeckService) Open(deck *models.Deck) error {
	return m.err
}

func (m mockDeckService) Create(deck *models.Deck) error {
	deck.UUID = m.deck.UUID
	deck.Remaining = m.deck.Remaining

	return m.err
}

func (m mockDeckService) ByUUID(uuid string) (*models.Deck, error) {
	return m.deck, m.err
}

func (m mockDeckService) Draw(deck *models.Deck, count int) ([]*models.Card, error) {
	return m.cards, m.err
}
