package signal

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() err = %s", err)
	}
	Handler(w, r)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Fatalf("Handler() status = %d; want %d", resp.StatusCode, 200)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Handler() Content-Type = %q; want %q", contentType, "application/json")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("io.ReadAll() err = %s", err)
	}
	var p Person
	if err = json.Unmarshal(data, &p); err != nil {
		t.Fatalf("json.Unmarshal() err = %s", err)
	}

	if p.Age != 23 {
		t.Errorf("p.Age = %d; want %d", p.Age, 23)
	}
	if p.Name != "Simon" {
		t.Errorf("p.Name = %s; want %s", p.Name, "Simon")
	}
	if p.Occupation != "Dev" {
		t.Errorf("p.Occupation = %s; want %s", p.Occupation, "Dev")
	}
}
