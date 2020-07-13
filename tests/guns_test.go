package test

import (
    "encoding/json"
    "bytes"
    "testing"
    "net/http"
    "net/http/httptest"
    "codproject/server/middleware"
    "codproject/server/models"
)

func TestGetGuns(t *testing.T) {
    body := []byte(`{"type":"Sniper_Rifle"}`)
    req, err := http.NewRequest("POST", "/getguns", bytes.NewBuffer(body))
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(middleware.GetGuns)
    handler.ServeHTTP(rr, req)

    got := []models.Gun{}
    json.NewDecoder(rr.Body).Decode(&got)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    drag := models.Gun{
        Name: "Dragunov",
        Type: "Sniper_Rifle",
        Gun_Id: 37,
    }
    hdr := models.Gun{
        Name: "HDR",
        Type: "Sniper_Rifle",
        Gun_Id: 38,
    }
    ax := models.Gun{
        Name: "AX-50",
        Type: "Sniper_Rifle",
        Gun_Id: 39,
    }
    rytec := models.Gun{
        Name: "Rytec AMR",
        Type: "Sniper_Rifle",
        Gun_Id: 40,
    }
    expected := []models.Gun{drag, hdr, ax, rytec}
    for i := 0; i < 4; i++ {
        if got[i] != expected[i] {
            t.Errorf("Handler returned unexpected body: got %v want %v", got[i], expected[i])
        }
    }
}
