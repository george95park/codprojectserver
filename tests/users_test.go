package test

import (
    "encoding/json"
    "testing"
    "net/http"
    "net/http/httptest"
    "codproject/server/middleware"
    "codproject/server/models"
)

func TestGetAllUsers(t *testing.T) {
    req, err := http.NewRequest("GET", "/getallusers", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(middleware.GetAllUsers)
    handler.ServeHTTP(rr, req)
    got := []models.User{}
    json.NewDecoder(rr.Body).Decode(&got)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    expected := []models.User{}
    expected = append(expected, models.User{User_Id: 2, Username: "ladyj"})
    expected = append(expected, models.User{User_Id: 9, Username: "lol"})
    expected = append(expected, models.User{User_Id: 4, Username: "daniel"})
    expected = append(expected, models.User{User_Id: 5, Username: "caleb"})
    expected = append(expected, models.User{User_Id: 12, Username: "jvmcmc"})
    expected = append(expected, models.User{User_Id: 13, Username: "mama"})
    expected = append(expected, models.User{User_Id: 15, Username: "mnmnmnmn"})
    expected = append(expected, models.User{User_Id: 1, Username: "george"})
    expected = append(expected, models.User{User_Id: 6, Username: "hey"})
    expected = append(expected, models.User{User_Id: 7, Username: "philip"})
    expected = append(expected, models.User{User_Id: 16, Username: "aaaa"})
    expected = append(expected, models.User{User_Id: 8, Username: "Pbob"})
    expected = append(expected, models.User{User_Id: 18, Username: "testing"})
    expected = append(expected, models.User{User_Id: 19, Username: "lala"})
    expected = append(expected, models.User{User_Id: 23, Username: "hii"})
    expected = append(expected, models.User{User_Id: 17, Username: "test"})
    for i := 0; i < len(expected); i++ {
        if got[i] != expected[i] {
            t.Errorf("Handler returned unexpected body: got %v want %v", got[i], expected[i])
        }
    }


}
