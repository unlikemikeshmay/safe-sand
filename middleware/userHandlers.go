package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mike/inv/data"
	"github.com/mike/inv/helpers"
)

var userList []data.User

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application-json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Expose-Headers", "*")
	(*w).Header().Set("X-Total-Count", "100")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	w.Header().Set("X-Total-Count", "100")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	enableCors(&w)
	var creds data.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	fmt.Println(creds)
	// check password and user in database
	user, err := _getUserWithCreds(creds.Email)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}
	if user.Email != "" {
		jwt, err := GenerateJWT(user)
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())

		}
		json.NewEncoder(w).Encode(jwt)
	} else {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode("user doesn't exist:")
	}
	//if user exists set jwt to session data if user exists allowing further access to routes
	fmt.Println(user)

	// return user data

	// if user does not exist return appropriate response

}
func addCookie(w http.ResponseWriter, name, jwt string) http.Cookie {
	expire := time.Now().Add(30 * time.Minute)
	cookie := http.Cookie{
		Name:    name,
		Value:   jwt,
		Expires: expire,
		Path:    "/trackerClient/",
	}
	return cookie

}
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	w.Header().Set("X-Total-Count", "100")
	userList = userList[:0]
	enableCors(&w)
	users, err := _getUsers()
	if err != nil {
		fmt.Println("error after assigment fo user or error from getUsers: ")
		fmt.Println(err)
		json.NewEncoder(w).Encode(err)
	}
	formatted := fmt.Sprintf("users: %s ",users)
	fmt.Println(formatted)
	json.NewEncoder(w).Encode(users)

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	uuid, err := _deleteUserAtUID(params["userid"])
	enableCors(&w)
	//inventory = append(inventory, item)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	json.NewEncoder(w).Encode(uuid)
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var user data.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	enableCors(&w)
	userRecord, err := _updateUser(user)
	//inventory = append(inventory, item)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	json.NewEncoder(w).Encode(userRecord)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	//@TODO add userid as a parameter to be added to all transactions
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var user data.User
	enableCors(&w)
	_ = json.NewDecoder(r.Body).Decode(&user)
	//inventory = append(inventory, tool)
	//@TODO insert tool into tool list
	user, err := _insertUser(user)
	if err != nil {
		json.NewEncoder(w).Encode("Error inserting tool in database")
	}
	json.NewEncoder(w).Encode(user.UID)
}

//-------------------handler db functions----------------//
func _getUserWithCreds(email string) (data.User, error) {
	db := data.CreateConnection()
	var userCont data.User
	defer db.Close()

	sqlStatement := `
		SELECT * FROM users WHERE email = $1
	`
	rows, err := db.Query(sqlStatement, email)
	if err != nil {
		return data.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var phone string
		var email string
		var show string
		var firstname string
		var lastname string
		var department string
		var passwordhash string

		err = rows.Scan(
			&phone, &email,
			&show, &firstname,
			&lastname, &department,
			&id, &passwordhash)
		if err != nil {
			return userCont, err
		}
		userCont.UID = id
		userCont.Phone = phone
		userCont.Email = email
		userCont.Show = show
		userCont.FirstName = firstname
		userCont.LastName = lastname
		userCont.Department = department
		userCont.PasswordHash = passwordhash

	}
	return userCont, nil

	// sanitize user data // check if it is the right values

	//

}
func _updateUser(user data.User) (data.User, error) {
	//@TODO add sql statement for updating database at row that matches id
	db := data.CreateConnection()
	defer db.Close()
	//var user data.User
	//phone,email,show,firstname,lastname,department,userid
	sqlStatement := `
		UPDATE users
		SET phone = $1, email = $2, show = $3, firstname = $4, lastname = $5,
		department = $6,passwordhash = $7
		WHERE userid = $8`
	_, err := db.Exec(sqlStatement,
		user.Phone,
		user.Email,
		user.Show,
		user.FirstName,
		user.LastName,
		user.Department,
		user.UID,
		user.PasswordHash,
	)
	if err != nil {
		return user, err
	}
	return user, nil

}
func _getUsers() ([]data.User, error) {
	db := data.CreateConnection()
	fmt.Println("db: ")
	fmt.Println(db)

	sqlStatement := `select * from users;`

	defer db.Close()
	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println("error if rows dont populate")
		fmt.Println(err)
		return userList[:0], err
	}

	defer rows.Close()


	for rows.Next() {
		var user data.User
		err := rows.Scan(&user.Phone, &user.Email, &user.Show, &user.FirstName, &user.LastName, &user.Department, &user.UID, &user.PasswordHash)
		if err != nil {
			return userList[:0], err
		}

		userList = append(userList, user)
		fmt.Println(userList)


	}
	return userList, nil

	//defer db.Close()
}
func _deleteUserAtUID(uid string) (string, error) {
	db := data.CreateConnection()
	var uuid string
	defer db.Close()

	sqlStatement := `DELETE FROM users WHERE userid = $1 RETURNING userid`

	res := db.QueryRow(sqlStatement, uid).Scan(&uuid)

	if res == nil {
		return uuid, fmt.Errorf("unable to execute DELETE statment for record with id: %v", uid)
	}
	return uuid, nil

}

//TODO: HASH PASSWORD AND SAVE IT
func _insertUser(user data.User) (data.User, error) {
	//create connection to postgres
	db := data.CreateConnection()
	uid := helpers.GenerateGUIDString()
	//close the connection after used
	defer db.Close()
	var newUser data.User
	//create the insert sql query
	sqlStatement := `INSERT INTO users (phone,email,show,firstname,lastname,department,userid,passwordhash) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	res, err := db.Exec(sqlStatement,
		user.Phone,
		user.Email,
		user.Show,
		user.FirstName,
		user.LastName,
		user.Department,
		uid,

		//Todo:: hash this
		helpers.HashAndSalt([]byte(user.PasswordHash)),
	)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record with id: %v", res)
	return newUser, nil

}

//-------------------------
//jwt setup
//set up global string for secret
/* var mysignkey = []byte("xmgjetyveks") */

// handler
/* func Gettokenhandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	w.Header().Set("X-Total-Count", "100")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	//enableCors(&w)
	jwt, err := GenerateJWT()
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	json.NewEncoder(w).Encode(jwt)
} */

//------------------------------------------------------//
func Validate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	w.Header().Set("X-Total-Count", "100")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	json.NewEncoder(w).Encode("asdfjkl")
}
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "endpoint called : homepage()")
}

/*
post request body structure example
{

    "toolname":"computer",
    "description":"used to write this",
    "showname":"test",
    "lastusersignout":"69359037-9599-48e7-b8f2-48393c019135",
    "currentuserid":"69359037-9599-48e7-b8f2-48393c019135",
    "signouttime":"'2001-10-05'"
}
*/
