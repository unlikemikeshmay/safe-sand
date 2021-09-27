package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/mike/inv/helpers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mike/inv/data"
	uuid "github.com/satori/go.uuid"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

var inventory []data.Tool

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	_deleteItemAtUID(uuid.FromStringOrNil(params["uid"]))

	json.NewEncoder(w).Encode(inventory)
}
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item data.Tool
	_ = json.NewDecoder(r.Body).Decode(&item)
	//params := mux.Vars(r)
	disTool, err := _updateTool(item)
	//inventory = append(inventory, item)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	json.NewEncoder(w).Encode(disTool)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	//@TODO add userid as a parameter to be added to all transactions

	w.Header().Set("Content-Type", "application/json")
	var tool data.Tool
	var response data.Response
	metaSlice := make([]string,0)
	metaSlice = append(metaSlice, r.Method)
	var tmpString string
	for key,element := range r.Header {

		tmpString = fmt.Sprintf("%s: %s",key,element)
		metaSlice = append(metaSlice,tmpString)
		}
	response.Data = r.Body
	response.Meta = metaSlice
	fmt.Println("response struct: ",response)
	err := json.NewDecoder(r.Body).Decode(&tool)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(err)
	}
	item, err := _insertItem(tool)
	if err != nil {
		json.NewEncoder(w).Encode("Error inserting tool in database")
	}
	json.NewEncoder(w).Encode(item)
}
func GetItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	w.Header().Set("X-Total-Count", "100")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	//enableCors(&w)

	items, err := _getTools()
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}

	json.NewEncoder(w).Encode(items)

}

//-------------------handler db functions----------------//
func _updateTool(tool data.Tool) (data.Tool, error) {
	//@TODO add sql statement for updating database at row that matches id
	db := data.CreateConnection()
	defer db.Close()

	var uidTemp = tool.UID
	var toolname = tool.ToolName
	var description = tool.Description
	var showname = tool.ShowName
	var lastusersignout = tool.LastUserSignOut
	var toolThing data.Tool

	//@TODO set a global object with user information to verify and set current user id
	var currentuserid = tool.CurrentUserId
	var signouttime = tool.SignOutTime
	sqlStatement := `
		UPDATE tools
		SET toolname = $2, description = $3, showname = $4, lastusersignout = $5, currentuserid = $6,
		signouttime = $7
		WHERE uid = $1
		RETURNING *;`

	err := db.QueryRow(sqlStatement, uidTemp, toolname,
		description, showname, lastusersignout,
		currentuserid, signouttime).Scan(
		&toolThing.UID,
		&toolThing.ToolName,
		&toolThing.Description,
		&toolThing.LastUserSignOut,
		&toolThing.CurrentUserId,
		&toolThing.SignOutTime,
	)
	if err != nil {
		return toolThing, err
	}
	return toolThing, nil

}
func _getTools() ([]data.Tool, error) {
	db := data.CreateConnection()

	defer db.Close()
	var inventory []data.Tool
	rows, err := db.Query("select * from tools")
	if err != nil {
		return inventory[:0], err
	}

	defer rows.Close()
	for rows.Next() {
		var uidTemp string
		var toolname string
		var description string
		var showname string
		var lastusersignout string
		//var lastuserString string
		//var currentuserString string

		var currentuserid string
		var signouttime string
		var tool data.Tool
		err = rows.Scan(&uidTemp, &toolname, &description, &showname, &lastusersignout, &currentuserid, &signouttime)
		if err != nil {
			return inventory[:0], err
		}
		tool.UID = uidTemp
		tool.ToolName = toolname
		tool.Description = description
		tool.ShowName = showname
		tool.LastUserSignOut = lastusersignout
		tool.CurrentUserId = currentuserid
		tool.SignOutTime = signouttime
		inventory = append(inventory, tool)

	}
	return inventory, nil

	//defer db.Close()
}
func _deleteItemAtUID(uid uuid.UUID) (uuid.UUID, error) {
	db := data.CreateConnection()
	var uuid uuid.UUID
	defer db.Close()

	sqlStatement := `DELETE FROM tools WHERE uid = $1 RETURNING uid`

	res := db.QueryRow(sqlStatement, uid).Scan(&uuid)

	if res == nil {
		return uuid, fmt.Errorf("unable to execute DELETE statment for record with id: %v", uid)
	}
	return uuid, nil

}
func _insertItem(tool data.Tool) (string, error) {
	//create connection to postgres
	db := data.CreateConnection()
	uid := helpers.GenerateGUIDString()
	fmt.Println("uid",uid)
	//close the connection after used
	defer db.Close()
	fmt.Println("tool inside of _insertItem")
	fmt.Println(tool.ToolName)
	fmt.Println(tool.UID)
	//create the insert sql query
	sqlStatement := `INSERT INTO tools (
	uid,
	toolname,
	description,
	showname,
	lastusersignout,
	currentuserid,
	signouttime
	) VALUES ($1,$2,$3,$4,$5,$6,$7)`

	_, err := db.Query(sqlStatement,
		uid,
		tool.ToolName,
		tool.Description,
		tool.ShowName,
		tool.LastUserSignOut,
		tool.CurrentUserId,
		"",
	)
	if err != nil {
		fmt.Println(err)
		return err.Error(), err
	}

	return "Success", nil

}

//------------------------------------------------------//
/* func GetInventory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")
	inventory = inventory[:0]

	inv, err := _getTools()
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	json.NewEncoder(w).Encode(inv)

} */

/* func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "endpoint called : homepage()")
} */

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
