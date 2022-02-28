package models

import (
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestNewDeckService(t *testing.T) {
	cs := cardService{}
	dv := deckValidator{DeckStorage: &deckMemory{decks: map[string]Deck{}}, cs: &cs}
	type args struct {
		cs CardService
	}
	tests := []struct {
		name string
		args args
		want DeckService
	}{
		{
			name: "default service",
			args: args{cs: &cs},
			want: &deckService{DeckStorage: &dv},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeckService(tt.args.cs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeckService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testDeckStorage(ds DeckStorage) func(t *testing.T) {
	return func(t *testing.T) {
		uuidStr := uuid.NewString()
		uuidStrUnused := uuid.NewString()
		deck := Deck{UUID: uuidStr}

		t.Run("Create", func(t *testing.T) {
			if err := ds.Create(&deck); err != nil {
				t.Errorf("Create() err = %s, want nil", err)
			}
		})

		t.Run("Update", func(t *testing.T) {
			deck.Remaining = 1
			if err := ds.Update(&deck); err != nil {
				t.Errorf("Update() err = %s, want nil", err)
			}
		})

		t.Run("ByUID", func(t *testing.T) {
			if _, err := ds.ByUUID(uuidStrUnused); err != ErrNotFound {
				t.Errorf("ByID(UUID) err = nil, want ErrNotFound")
			}

			foundDeck, err := ds.ByUUID(uuidStr)
			if err != nil {
				t.Errorf("ByID(UUID) err = %s, want nil", err)
			}

			if !reflect.DeepEqual(&deck, foundDeck) {
				t.Errorf("ByID(UUID) got = %v, want %v", foundDeck, &deck)
			}
		})
	}
}

func TestDeckMemory(t *testing.T) {
	ds := deckMemory{decks: map[string]Deck{}}
	t.Run("TestDeckMemory", testDeckStorage(&ds))
}

func TestDeckService(t *testing.T) {
	ds := NewDeckService(NewCardService())
	t.Run("TestDeckService", testDeckStorage(ds))
}

func TestDeckValidator(t *testing.T) {
	ds := newDeckValidator(&deckMemory{decks: map[string]Deck{}}, NewCardService())
	t.Run("TestDeckValidator", testDeckStorage(ds))
}

func Test_deckService_Draw(t *testing.T) {
	uuidStr := uuid.NewString()
	deck := Deck{UUID: uuidStr, Cards: allCards, Remaining: len(allCards)}
	openUUIDStr := uuid.NewString()
	openDeck := Deck{UUID: openUUIDStr, Cards: nil, Remaining: 0, Opened: true}
	dm := deckMemory{decks: map[string]Deck{}}
	_ = dm.Create(&deck)
	_ = dm.Create(&openDeck)

	type fields struct {
		DeckStorage DeckStorage
	}
	type args struct {
		deck  *Deck
		count int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Card
		wantErr bool
	}{
		{
			name:    "draw more than remaining",
			fields:  fields{&dm},
			args:    args{&deck, 99},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "valid",
			fields:  fields{&dm},
			args:    args{&deck, 5},
			want:    allCards[0:5],
			wantErr: false,
		},
		{
			name:    "draw after open",
			fields:  fields{&dm},
			args:    args{&openDeck, 3},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &deckService{
				DeckStorage: tt.fields.DeckStorage,
			}
			got, err := ds.Draw(tt.args.deck, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Draw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Draw() got = %+v, want %+v", got, tt.want)
			}
		})
	}

	t.Run("check updated", func(t *testing.T) {
		foundDeck, _ := dm.ByUUID(uuidStr)
		if !reflect.DeepEqual(foundDeck, &deck) {
			t.Errorf("Draw() got = %+v, want %+v", foundDeck, &deck)
		}
	})
}

func Test_newDeckValidator(t *testing.T) {
	type args struct {
		ds DeckStorage
		cs CardService
	}
	tests := []struct {
		name string
		args args
		want *deckValidator
	}{
		{
			name: "default",
			args: args{&deckMemory{}, NewCardService()},
			want: &deckValidator{&deckMemory{}, NewCardService()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newDeckValidator(tt.args.ds, tt.args.cs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDeckValidator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deckValidator_ByUUID(t *testing.T) {
	validUUID := uuid.NewString()
	deck := Deck{UUID: validUUID}
	dm := deckMemory{decks: map[string]Deck{validUUID: deck}}
	cs := NewCardService()
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		args    args
		want    *Deck
		wantErr bool
	}{
		{
			name:    "valid",
			args:    args{uuid: validUUID},
			want:    &deck,
			wantErr: false,
		},
		{
			name:    "invalid uuid",
			args:    args{uuid: "abcd"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty uuid",
			args:    args{uuid: ""},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dv := &deckValidator{
				DeckStorage: &dm,
				cs:          cs,
			}
			got, err := dv.ByUUID(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ByUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByUUID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deckValidator_Create(t *testing.T) {
	ds := deckMemory{decks: map[string]Deck{}}
	cs := NewCardService()
	dv := deckValidator{
		DeckStorage: &ds,
		cs:          cs,
	}
	validUUID := uuid.NewString()
	tests := []struct {
		name    string
		deck    *Deck
		want    *Deck
		wantErr bool
	}{
		{
			name:    "invalid uuid",
			deck:    &Deck{UUID: "invalid"},
			want:    &Deck{UUID: "invalid"},
			wantErr: true,
		},
		{
			name:    "valid fill cards",
			deck:    &Deck{UUID: validUUID},
			want:    &Deck{UUID: validUUID, Remaining: 52, Cards: allCards},
			wantErr: false,
		},
		{
			name: "valid card codes",
			deck: &Deck{UUID: validUUID, CardCodes: "AS,10D"},
			want: &Deck{UUID: validUUID, Remaining: 2, Cards: []*Card{
				NewCard(ValueAce, SuitSpades),
				NewCard(Value("10"), SuitDiamonds),
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := dv.Create(tt.deck); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.deck, tt.want) {
				t.Errorf("Create() got = %v, want %v", tt.deck, tt.want)
			}
		})
	}

	t.Run("empty uuid", func(t *testing.T) {
		deck := Deck{}
		_ = dv.Create(&deck)
		if _, err := uuid.Parse(deck.UUID); err != nil {
			t.Errorf("Create() uuid parse error = %v, want nil", err)
		}
	})

	t.Run("shuffle", func(t *testing.T) {
		deck := Deck{Shuffled: true}
		_ = dv.Create(&deck)
		if reflect.DeepEqual(deck.Cards, allCards) {
			t.Errorf("Create() not shuffled")
		}
		got := len(deck.Cards)
		want := len(allCards)
		if got != want {
			t.Errorf("Create() len cars got %v, want %v", got, want)
		}
	})
}

func Test_deckValidator_Update(t *testing.T) {
	ds := deckMemory{decks: map[string]Deck{}}
	cs := NewCardService()
	dv := deckValidator{
		DeckStorage: &ds,
		cs:          cs,
	}
	validUUID := uuid.NewString()

	tests := []struct {
		name    string
		deck    *Deck
		want    *Deck
		wantErr bool
	}{
		{
			name:    "valid",
			deck:    &Deck{UUID: validUUID},
			want:    &Deck{UUID: validUUID},
			wantErr: false,
		},
		{
			name:    "empty uuid",
			deck:    &Deck{},
			want:    &Deck{},
			wantErr: true,
		},
		{
			name:    "remaining",
			deck:    &Deck{UUID: validUUID, Cards: allCards[0:5]},
			want:    &Deck{UUID: validUUID, Cards: allCards[0:5], Remaining: 5},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := dv.Update(tt.deck); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.deck, tt.want) {
				t.Errorf("Update() got = %v, want %v", tt.deck, tt.want)
			}
		})
	}
}
