package data

type User struct {
	UID          string `json:"id"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Show         string    `json:"show"`
	FirstName    string    `json:"firstname"`
	LastName     string    `json:"lastname"`
	Department   string    `json:"department"`
	PasswordHash string    `json:"passwordhash"`
}
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
