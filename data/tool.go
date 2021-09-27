package data

//Tool defines tools that will be in the inventory
type Tool struct {
	UID             string `json:"uid"`
	ToolName        string    `json:"toolname"`
	Description     string    `json:"description"`
	ShowName        string    `json:"showname"`
	LastUserSignOut string `json:"lastusersignout"`
	CurrentUserId   string `json:"currentuserid"`
	//make sure to store a string that was parsed from datetime
	SignOutTime string `json:"signouttime"`
}
