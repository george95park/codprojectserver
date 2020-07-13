package test

import (
    "testing"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "codproject/server/middleware"
    "codproject/server/models"
)

func TestAuthUser(t *testing.T) {
    // new request
    req, err := http.NewRequest("GET", "/authuser", nil)
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Cookie", "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IlBib2IiLCJ1c2VyX2lkIjo4LCJleHAiOjE2MjU2Mjg2NjN9.RqTPhf5r6E7HUwWdoGidboSkynkFRZ77rKC-HPhkfoI")
    // returns a *ResponseRecorder which is an implementation of Response Writer
    // that records mutations so that it can be used for testing later
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(middleware.AuthUser)
    handler.ServeHTTP(rr, req)
    got := models.User{}
    json.NewDecoder(rr.Body).Decode(&got)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    expected := models.User{
        Username: "Pbob",
        User_Id: 8,
        Logged_In: true,
    }
    if got != expected {
        t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }
}
