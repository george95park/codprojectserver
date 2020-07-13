package test

import (
    "bytes"
    "testing"
    "net/http"
    "net/http/httptest"
    "codproject/server/middleware"
)

func TestSignup(t *testing.T) {
    body := []byte(`{"username":"hii","password":"123"}`)
    req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(middleware.Signup)
    handler.ServeHTTP(rr, req)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
}
func TestLogin(t *testing.T) {
    body := []byte(`{"username":"test","password":"123"}`)
    req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(middleware.Login)
    handler.ServeHTTP(rr, req)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
}

func TestLogout(t *testing.T) {
    req, err := http.NewRequest("POST", "/logout", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(middleware.Logout)
    handler.ServeHTTP(rr, req)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
}
