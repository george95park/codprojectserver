package middleware

import (
	"os"
	"fmt"
	"net/http"
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"codproject/server/models"
)

// opens access to database
func ConnectDB() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to DB")
	return db
}

// creates loadout in database
func CreateLoadout(w http.ResponseWriter, r *http.Request) {
	db := ConnectDB()
	defer db.Close()
	load := &models.Loadout{}
	// Unmarshall request body
	if err := json.NewDecoder(r.Body).Decode(load); err != nil {
		panic(err)
	}

	// insert new loadout to database
	if _, err := db.Query("insert into loadouts (username, gunname, attachments, description) values ($1, $2, $3, $4)",
		load.Username,
		load.Gunname,
		pq.Array(load.Attachments),
		load.Description); err != nil {
			panic(err)
		}
	fmt.Println("Created loadout successfully")
}

// gets loadouts from database according to user
func GetLoadouts(w http.ResponseWriter, r *http.Request) {
	db := ConnectDB()
	defer db.Close()
	currUser := &models.User{}

	// Unmarshall request body
	if err := json.NewDecoder(r.Body).Decode(&currUser); err != nil {
		panic(err)
	}

	// Get the rows according to user
	rows,err := db.Query("select id,gunname,attachments,description from loadouts where username=$1",currUser.Username)
	if err != nil {
		panic(err)
	}

	// For each row, create new loadout object and append to the result
	userLoadouts := []models.Loadout{}
	for rows.Next() {
		l := models.Loadout{}
		l.Username = currUser.Username
		err = rows.Scan(&l.Id, &l.Gunname, pq.Array(&l.Attachments), &l.Description)
		if err != nil {
			panic(err)
		}
		userLoadouts = append(userLoadouts, l)
	}
	fmt.Println("Received loadouts.")

	// return response
	json.NewEncoder(w).Encode(userLoadouts)
}

func DeleteLoadout(w http.ResponseWriter, r *http.Request) {
	// open database
	db := ConnectDB()
	defer db.Close()
	load := &models.Loadout{}
	if err := json.NewDecoder(r.Body).Decode(load); err != nil {
		panic(err)
	}

	res,err := db.Exec("delete from loadouts where id=$1", load.Id)
	if err != nil {
		panic(err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println("Total rows affected: %v", rowsAffected)
}

func UpdateLoadout(w http.ResponseWriter, r *http.Request) {
	// open database
	db := ConnectDB()
	defer db.Close()
	load := &models.Loadout{}
	if err := json.NewDecoder(r.Body).Decode(load); err != nil {
		panic(err)
	}
	res, err := db.Exec("update loadouts set username=$2, gunname=$3, attachments=$4, description=$5 where id=$1",
		load.Id,
		load.Username,
		load.Gunname,
		pq.Array(load.Attachments),
		load.Description)
	if err != nil {
		panic(err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println("Total rows affected: %v", rowsAffected)
}

// Sign-up handler
func Signup(w http.ResponseWriter, r *http.Request) {
	// open database
	db := ConnectDB()
	defer db.Close()
	creds := &models.Credentials{}

	// decode request body and check for error
	if err := json.NewDecoder(r.Body).Decode(creds); err != nil {
		fmt.Println(err)
		// return a 400 status
		w.WriteHeader(http.StatusBadRequest)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		panic(err)
	}

	// insert into credentials table and check for error
	_,err = db.Query("insert into credentials (username, password) values ($1, $2)",creds.Username,string(hashedPassword))
	if err != nil {
		fmt.Println(err)
		// return 500 status
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// insert into USERS table and check for error
		if _, err := db.Query("insert into users (username) values ($1)", creds.Username); err != nil {
				fmt.Println(err)
				// return 500 status
				w.WriteHeader(http.StatusInternalServerError)
			}
	}
}
// Login handler
func Login(w http.ResponseWriter, r *http.Request) {
	// open database
	db := ConnectDB()
	defer db.Close()
	creds := &models.Credentials{}

	// decode request body and check for error
	if err := json.NewDecoder(r.Body).Decode(creds); err != nil {
		panic(err)
	}

	// query database for the password with username in request body
	storedCreds := &models.Credentials{}
	passFromDB := db.QueryRow("SELECT password FROM credentials WHERE username = $1", creds.Username).Scan(&storedCreds.Password)
	switch {
		case passFromDB == sql.ErrNoRows:
			fmt.Println("No user with the username: %d\n", creds.Username)
			// return 401 status
			w.WriteHeader(http.StatusUnauthorized)
		case passFromDB != nil:
			fmt.Println("Query error: %v\n", passFromDB)
			// return 500 status
			w.WriteHeader(http.StatusInternalServerError)
		default:
			if err := bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
				fmt.Println("Access Denied: Wrong Password.")
				// return 401 status
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				fmt.Println("Access Granted.")
			}
	}
}


func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to home page")
	fmt.Println("Endpoint hit: HomePage")
}
