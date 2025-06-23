package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddTodoHandler(t *testing.T) {
	// --- Arrange ---
	body := strings.NewReader(`{"description":"Test task","status":"not started"}`)
	req := httptest.NewRequest(http.MethodPost, "/create", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// --- Act ---
	AddTodoHandler(w, req)

	// --- Assert ---
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201 Created, got %d", resp.StatusCode)
	}
}
