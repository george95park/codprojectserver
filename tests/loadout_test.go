package test

import (
    "reflect"
    "encoding/json"
    "bytes"
    "testing"
    "net/http"
    "net/http/httptest"
    "codproject/server/middleware"
    "codproject/server/models"
    "github.com/gorilla/mux"
)

func TestCreateLoadout(t *testing.T) {
    body := []byte(`{"user_id":1, "type":"...", "gun":"///", "attachments":["hi","hello"], "subattachments":["bye"], "description":"testing"}`)
    req, err := http.NewRequest("POST", "/createloadout", bytes.NewBuffer(body))
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(middleware.CreateLoadout)
    handler.ServeHTTP(rr, req)

    got := models.Loadout{}
    json.NewDecoder(rr.Body).Decode(&got)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    expected := models.Loadout{
        Loadout_Id: 0,
        User_Id: 1,
        Type: "...",
        Gun: "///",
        Attachments: []string{"hi","hello"},
        SubAttachments: []string{"bye"},
        Description: "testing",
    }
    if !reflect.DeepEqual(got, expected) {
        t.Errorf("Handler returned unexpected body: got %v want %v", got, expected)
    }
}

func TestGetLoadouts(t *testing.T) {
    req, err := http.NewRequest("GET", "/getloadouts", nil)
    if err != nil {
        t.Fatal(err)
    }
    req = mux.SetURLVars(req, map[string]string{
        "id": "1",
    })
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(middleware.GetLoadouts)
    handler.ServeHTTP(rr, req)
    got := []models.Loadout{}
    json.NewDecoder(rr.Body).Decode(&got)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    load1 := models.Loadout{
        Loadout_Id: 1,
        User_Id: 1,
        Type: "Handgun",
        Gun: ".50 GS",
        Attachments: []string{"Muzzle","Barrel","Laser","Optic","Perks"},
        SubAttachments: []string{"Compensator (Lv. 21)\n", "FORGE TAC Enforcer (Lv. 38)\n", "5mW Laser (Lv. 22)\n", "Cronen LP945 Mini Reflex (Lv. 5)\n", "Sleight of Hand (Lv. 20)\n"},
        Description: "123",
    }
    load2 := models.Loadout{
        Loadout_Id: 12,
        User_Id: 1,
        Type: "Marksman_Rifle",
        Gun: "SKS",
        Attachments: []string{"Muzzle","Barrel","Stock","Ammunition","Perks"},
        SubAttachments: []string{"Muzzle Brake\n", "16\" FSS Para\n", "Sawed-off Stock\n", "10 Round Mags\n", "Focus\n"},
        Description: "hehe",
    }
    load3 := models.Loadout{
        Loadout_Id: 13,
        User_Id: 1,
        Type: "Sniper_Rifle",
        Gun: "AX-50",
        Attachments: []string{"hi","hello"},
        SubAttachments: []string{"bye"},
        Description: "testing",
    }
    expected := []models.Loadout{load1, load2, load3}
    for i := 0; i < 3; i++ {
        if !reflect.DeepEqual(got[i], expected[i]) {
            t.Errorf("Handler returned unexpected body: got %v want %v", got[i], expected[i])
        }
    }
}

func TestDeleteLoadout(t *testing.T) {
    req, err := http.NewRequest("DELETE", "/deleteloadout", nil)
    if err != nil {
        t.Fatal(err)
    }
    req = mux.SetURLVars(req, map[string]string{
        "id": "3",
    })
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(middleware.DeleteLoadout)
    handler.ServeHTTP(rr, req)
    var got string
    json.NewDecoder(rr.Body).Decode(&got)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    if got != "3" {
        t.Errorf("Handler returned unexpected body: got %v want %v", got, 3)
    }
}

func TestUpdateLoadout(t *testing.T) {
    body := []byte(`{"loadout_id":13, "user_id":1, "type":"blahblah", "gun":"testing123", "attachments":["popo"], "subattachments":["papa","ppupuas"], "description":"update"}`)
    req, err := http.NewRequest("PUT", "/updateloadout", bytes.NewBuffer(body))
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(middleware.UpdateLoadout)
    handler.ServeHTTP(rr, req)
    got := models.Loadout{}
    json.NewDecoder(rr.Body).Decode(&got)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    expected := models.Loadout{
        Loadout_Id: 13,
        User_Id: 1,
        Type: "blahblah",
        Gun: "testing123",
        Attachments: []string{"popo"},
        SubAttachments: []string{"papa","ppupuas"},
        Description: "update",
    }
    if !reflect.DeepEqual(got, expected) {
        t.Errorf("Handler returned unexpected body: got %v want %v", got, expected)
    }
}
