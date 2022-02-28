package models

import (
	"reflect"
	"testing"
)

func TestNewCard(t *testing.T) {
	type args struct {
		value Value
		suit  Suit
	}
	tests := []struct {
		name string
		args args
		want *Card
	}{
		{
			name: "alpha value",
			args: args{value: ValueAce, suit: SuitSpades},
			want: &Card{Value: ValueAce, Suit: SuitSpades, Code: "AS"},
		}, {
			name: "numeric value",
			args: args{value: Value("10"), suit: SuitDiamonds},
			want: &Card{Value: "10", Suit: SuitDiamonds, Code: "10D"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCard(tt.args.value, tt.args.suit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCardService(t *testing.T) {
	tests := []struct {
		name string
		want CardService
	}{
		{
			name: "default service",
			want: &cardService{&cardValidator{&staticCardStorage{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCardService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCardService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSuit_Code(t *testing.T) {
	tests := []struct {
		name string
		s    Suit
		want string
	}{
		{
			name: "defined suit",
			s:    SuitHearts,
			want: "H",
		},
		{
			name: "random suit",
			s:    Suit("Test"),
			want: "T",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Code(); got != tt.want {
				t.Errorf("Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValue_Code(t *testing.T) {
	tests := []struct {
		name string
		v    Value
		want string
	}{
		{
			name: "alpha value",
			v:    ValueJack,
			want: "J",
		},
		{
			name: "short numeric value",
			v:    Value("3"),
			want: "3",
		},
		{
			name: "long numeric value",
			v:    Value("10"),
			want: "10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Code(); got != tt.want {
				t.Errorf("Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cardService_ByCodes(t *testing.T) {
	cs := cardValidator{&staticCardStorage{}}
	type fields struct {
		CardStorage CardStorage
	}
	type args struct {
		codes []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Card
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				CardStorage: &cs,
			},
			args: args{
				codes: []string{"AS", "10H"},
			},
			want: []*Card{
				{Value: ValueAce, Suit: SuitSpades, Code: "AS"},
				{Value: Value("10"), Suit: SuitHearts, Code: "10H"},
			},
			wantErr: false,
		},
		{
			name: "valid trailing space",
			fields: fields{
				CardStorage: &cs,
			},
			args: args{
				codes: []string{" AS ", " 10H "},
			},
			want: []*Card{
				{Value: ValueAce, Suit: SuitSpades, Code: "AS"},
				{Value: Value("10"), Suit: SuitHearts, Code: "10H"},
			},
			wantErr: false,
		},
		{
			name: "invalid value alpha",
			fields: fields{
				CardStorage: &cs,
			},
			args: args{
				codes: []string{"X1S"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid value big int",
			fields: fields{
				CardStorage: &cs,
			},
			args: args{
				codes: []string{"11S"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid value small int",
			fields: fields{
				CardStorage: &cs,
			},
			args: args{
				codes: []string{"0S"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid suit",
			fields: fields{
				CardStorage: &cs,
			},
			args: args{
				codes: []string{"10X"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &cardService{
				CardStorage: tt.fields.CardStorage,
			}
			got, err := cs.ByCodes(tt.args.codes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ByCodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByCodes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cardService_ByCodesStr(t *testing.T) {
	cs := cardValidator{&staticCardStorage{}}
	type fields struct {
		CardStorage CardStorage
	}
	type args struct {
		codesStr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Card
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				CardStorage: &cs,
			},
			args: args{
				codesStr: "AS,10H",
			},
			want: []*Card{
				{Value: ValueAce, Suit: SuitSpades, Code: "AS"},
				{Value: Value("10"), Suit: SuitHearts, Code: "10H"},
			},
			wantErr: false,
		},
		{
			name: "valid trailing space",
			fields: fields{
				CardStorage: &cs,
			},
			args: args{
				codesStr: " AS , 10H ",
			},
			want: []*Card{
				{Value: ValueAce, Suit: SuitSpades, Code: "AS"},
				{Value: Value("10"), Suit: SuitHearts, Code: "10H"},
			},
			wantErr: false,
		},
		{
			name: "invalid value",
			fields: fields{
				CardStorage: &cs,
			},
			args: args{
				codesStr: "X1S",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid suit",
			fields: fields{
				CardStorage: &cs,
			},
			args: args{
				codesStr: "1XS",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &cardService{
				CardStorage: tt.fields.CardStorage,
			}
			got, err := cs.ByCodesStr(tt.args.codesStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ByCodesStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByCodesStr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_staticCardStorage_All(t *testing.T) {
	tests := []struct {
		name    string
		want    []*Card
		wantErr bool
	}{
		{
			name:    "valid",
			want:    allCards,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scs := &staticCardStorage{}
			got, err := scs.All()
			if (err != nil) != tt.wantErr {
				t.Errorf("All() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("All() length got = %v, want %v", len(got), len(tt.want))
			}

			for i := range got {
				if !reflect.DeepEqual(*got[i], *tt.want[i]) {
					t.Errorf("All() elem at index  %d got = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

var allCards = []*Card{
	{Value: "ACE", Suit: "SPADES", Code: "AS"},
	{Value: "2", Suit: "SPADES", Code: "2S"},
	{Value: "3", Suit: "SPADES", Code: "3S"},
	{Value: "4", Suit: "SPADES", Code: "4S"},
	{Value: "5", Suit: "SPADES", Code: "5S"},
	{Value: "6", Suit: "SPADES", Code: "6S"},
	{Value: "7", Suit: "SPADES", Code: "7S"},
	{Value: "8", Suit: "SPADES", Code: "8S"},
	{Value: "9", Suit: "SPADES", Code: "9S"},
	{Value: "10", Suit: "SPADES", Code: "10S"},
	{Value: "JACK", Suit: "SPADES", Code: "JS"},
	{Value: "QUEEN", Suit: "SPADES", Code: "QS"},
	{Value: "KING", Suit: "SPADES", Code: "KS"},
	{Value: "ACE", Suit: "DIAMONDS", Code: "AD"},
	{Value: "2", Suit: "DIAMONDS", Code: "2D"},
	{Value: "3", Suit: "DIAMONDS", Code: "3D"},
	{Value: "4", Suit: "DIAMONDS", Code: "4D"},
	{Value: "5", Suit: "DIAMONDS", Code: "5D"},
	{Value: "6", Suit: "DIAMONDS", Code: "6D"},
	{Value: "7", Suit: "DIAMONDS", Code: "7D"},
	{Value: "8", Suit: "DIAMONDS", Code: "8D"},
	{Value: "9", Suit: "DIAMONDS", Code: "9D"},
	{Value: "10", Suit: "DIAMONDS", Code: "10D"},
	{Value: "JACK", Suit: "DIAMONDS", Code: "JD"},
	{Value: "QUEEN", Suit: "DIAMONDS", Code: "QD"},
	{Value: "KING", Suit: "DIAMONDS", Code: "KD"},
	{Value: "ACE", Suit: "CLUBS", Code: "AC"},
	{Value: "2", Suit: "CLUBS", Code: "2C"},
	{Value: "3", Suit: "CLUBS", Code: "3C"},
	{Value: "4", Suit: "CLUBS", Code: "4C"},
	{Value: "5", Suit: "CLUBS", Code: "5C"},
	{Value: "6", Suit: "CLUBS", Code: "6C"},
	{Value: "7", Suit: "CLUBS", Code: "7C"},
	{Value: "8", Suit: "CLUBS", Code: "8C"},
	{Value: "9", Suit: "CLUBS", Code: "9C"},
	{Value: "10", Suit: "CLUBS", Code: "10C"},
	{Value: "JACK", Suit: "CLUBS", Code: "JC"},
	{Value: "QUEEN", Suit: "CLUBS", Code: "QC"},
	{Value: "KING", Suit: "CLUBS", Code: "KC"},
	{Value: "ACE", Suit: "HEARTS", Code: "AH"},
	{Value: "2", Suit: "HEARTS", Code: "2H"},
	{Value: "3", Suit: "HEARTS", Code: "3H"},
	{Value: "4", Suit: "HEARTS", Code: "4H"},
	{Value: "5", Suit: "HEARTS", Code: "5H"},
	{Value: "6", Suit: "HEARTS", Code: "6H"},
	{Value: "7", Suit: "HEARTS", Code: "7H"},
	{Value: "8", Suit: "HEARTS", Code: "8H"},
	{Value: "9", Suit: "HEARTS", Code: "9H"},
	{Value: "10", Suit: "HEARTS", Code: "10H"},
	{Value: "JACK", Suit: "HEARTS", Code: "JH"},
	{Value: "QUEEN", Suit: "HEARTS", Code: "QH"},
	{Value: "KING", Suit: "HEARTS", Code: "KH"},
}
