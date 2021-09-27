package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mike/inv/data"
)

func GetShows(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response data.Response

	metaSlice := make([]string, 0)
	metaSlice = append(metaSlice, r.Method)
	shows, err := _getShows()
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(err)
	}
	response.Data = shows
	response.Meta = metaSlice
	response.Success = true
	json.NewEncoder(w).Encode(response)

}
func DeleteShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("inside Deleteshow")
	var response data.Response
	metaSlice := make([]string, 0)
	metaSlice = append(metaSlice, r.Method)
	metaSlice = append(metaSlice, r.Method)
	params := mux.Vars(r)
	fmt.Println("inside Deleteshow params: ")
	fmt.Println("params: ", params["uid"])
	res, err := _deleteShowAtUID(params["uid"])
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(err)
	}
	response.Data = res
	response.Meta = metaSlice
	response.Success = true
	fmt.Println("response.data: ", res)
	json.NewEncoder(w).Encode(response)

}
func UpdateShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var show data.Show
	_ = json.NewDecoder(r.Body).Decode(&show)
	//params := mux.Vars(r)
	fmt.Println("show in updateshow: ")
	fmt.Println(show)
	show, err := _updateShow(show)
	//inventory = append(inventory, item)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	json.NewEncoder(w).Encode(show)
}

func _updateShow(show data.Show) (data.Show, error) {
	db := data.CreateConnection()
	fmt.Println("inside _updateShow")
	fmt.Println("value of show: ")
	fmt.Println(show)
	fmt.Println("value of show uid: ")
	fmt.Println(show.UID)
	defer db.Close()

	sqlStatement := `
		UPDATE shows
		SET show_name = $2, production = $3
		WHERE uid = $1`
	res, err := db.Exec(sqlStatement,
		show.UID,
		show.ShowName,
		show.Production,
	)
	if err != nil {
		fmt.Println("failure in sql execution")
		fmt.Println(err.Error())
		return show, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		fmt.Println("failure in sql execution")
		fmt.Println(err.Error())
		return show, err
	}
	fmt.Println("no errors apparent? here is count")
	fmt.Println(count)
	return show, nil
}
func AddShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var show data.Show
	var response data.Response

	metaSlice := make([]string, 0)
	metaSlice = append(metaSlice, r.Method)

	//error handling on struct mapping
	err := json.NewDecoder(r.Body).Decode(&show)
	fmt.Println("show")
	//switching to a client side uid generation until further notice~
	//#TODO remove reference to generate guid string if generation is migrated to clientside
	//show.UID = helpers.GenerateGUIDString()
	fmt.Println("show UID: ")
	fmt.Println(show.UID)
	fmt.Println("show name")
	fmt.Println(show.ShowName)
	if err != nil {
		fmt.Println("inside err triggered by json.NewDecoder(r.Body).Decode(&show)")
		fmt.Println(err)
		json.NewEncoder(w).Encode(err)
	}
	//error handling on insertion of rows
	dbRes, err := _insertShow(show)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	// return response if all checks out
	response.Data = dbRes
	response.Meta = metaSlice
	json.NewEncoder(w).Encode(response)
}

func _insertShow(show data.Show) (string, error) {
	db := data.CreateConnection()
	defer db.Close()
	sqlStatment := `INSERT INTO shows (uid,show_name,Production ) VALUES ($1,$2,$3)`
	_, err := db.Query(sqlStatment, show.UID, show.ShowName, show.Production)
	if err != nil {
		return err.Error(), err
	}
	return "success", nil

}

func _deleteShowAtUID(uid string) (string, error) {
	db := data.CreateConnection()
	var uuid string
	defer db.Close()

	sqlStatement := `DELETE FROM shows WHERE uid = $1 RETURNING uid`

	res := db.QueryRow(sqlStatement, uid).Scan(&uuid)

	if res == nil {
		return uuid, fmt.Errorf("unable to execute DELETE statment for record with id: %v", uid)
	}
	return uuid, nil

}
func _getShows() ([]data.Show, error) {
	db := data.CreateConnection()
	defer db.Close()
	shows := make([]data.Show, 0)
	sqlQuery := `SELECT * FROM shows`
	rows, err := db.Query(sqlQuery)
	if err != nil {
		fmt.Println(err)
		return shows, err
	}
	defer rows.Close()
	for rows.Next() {
		var show data.Show
		err = rows.Scan(&show.UID, &show.ShowName, &show.Production)
		if err != nil {
			fmt.Println(err)
			return shows, err
		}
		shows = append(shows, show)

	}
	return shows, nil
}
