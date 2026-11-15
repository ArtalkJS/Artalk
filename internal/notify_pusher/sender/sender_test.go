package sender

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLarkNotificationFormats(t *testing.T) {
	const title = "New comment"
	const message = "Reader replied."
	const link = "https://comments.example.org/post"
	requests := make(chan []byte, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil || r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		requests <- data
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	SendLark(server.URL, title, message, false)
	require.Len(t, requests, 1)
	require.JSONEq(t, `{"msg_type":"text","content":{"text":"New comment\n\nReader replied."}}`, string(<-requests))

	SendLark(server.URL, title, message+"\n"+link, true)
	require.Len(t, requests, 1)
	var payload struct {
		Type string `json:"msg_type"`
		Card struct {
			Header struct {
				Title struct {
					Content string `json:"content"`
				} `json:"title"`
			} `json:"header"`
			Elements []struct {
				Text struct {
					Content string `json:"content"`
				} `json:"text"`
				Actions []struct {
					URL string `json:"url"`
				} `json:"actions"`
			} `json:"elements"`
		} `json:"card"`
	}
	require.NoError(t, json.Unmarshal(<-requests, &payload))
	require.Equal(t, "interactive", payload.Type)
	require.Equal(t, title, payload.Card.Header.Title.Content)
	require.Len(t, payload.Card.Elements, 2)
	require.Equal(t, message, payload.Card.Elements[0].Text.Content)
	require.Len(t, payload.Card.Elements[1].Actions, 1)
	require.Equal(t, link, payload.Card.Elements[1].Actions[0].URL)
}

func TestBarkNotificationEscapesPathSegments(t *testing.T) {
	paths := make(chan string, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths <- r.URL.EscapedPath()
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)
	SendBark(server.URL+"/fixture/", "", "A / B?")
	require.Len(t, paths, 1)
	require.Equal(t, "/fixture/Artalk/A+%2F+B%3F", <-paths)
}
