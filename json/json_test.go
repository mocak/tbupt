package json

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestResponse(t *testing.T) {
	type TestType struct {
		Test string
	}
	valid := TestType{Test: "test"}
	validJson, _ := json.Marshal(&valid)

	type args struct {
		w *httptest.ResponseRecorder
		v interface{}
		s int
	}
	tests := []struct {
		name       string
		args       args
		want       string
		wantCT     string
		wantStatus int
	}{
		{
			name:       "valid",
			args:       args{httptest.NewRecorder(), &valid, http.StatusCreated},
			want:       string(validJson),
			wantCT:     "application/json",
			wantStatus: http.StatusCreated,
		},
		{
			name:       "invalid type",
			args:       args{httptest.NewRecorder(), func() {}, http.StatusOK},
			want:       "\"Unexpected Error\"",
			wantCT:     "application/json",
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Response(tt.args.w, tt.args.v, tt.args.s)
			resp := tt.args.w.Result()
			byteSlice, _ := io.ReadAll(resp.Body)
			got := strings.TrimSuffix(string(byteSlice), "\n")
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestResponse() = %v, want %v", got, tt.want)
			}
			gotCT := resp.Header.Get("Content-Type")
			if gotCT != tt.wantCT {
				t.Errorf("TestResponse() = %v, want %v", gotCT, tt.wantCT)
			}
			gotStatus := resp.StatusCode
			if gotStatus != tt.wantStatus {
				t.Errorf("TestResponse() = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}

func TestDecodeBody(t *testing.T) {
	type TestType struct {
		Test string
	}
	v := &TestType{Test: "test"}
	validBody, _ := json.Marshal(&v)
	invalidBody, _ := json.Marshal([]string{"test"})

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    interface{}
	}{
		{
			name: "valid",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", bytes.NewReader(validBody)),
				v: &TestType{},
			},
			want:    &TestType{Test: "test"},
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", bytes.NewReader(invalidBody)),
				v: &TestType{},
			},
			want:    &TestType{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DecodeBody(tt.args.w, tt.args.r, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("DecodeBody() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.want, tt.args.v) {
				t.Errorf("TestDecodeBody() = %v, want %v", tt.args.v, tt.want)
			}
		})
	}
}
