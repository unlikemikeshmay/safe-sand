package data

type Data = interface {}
type Meta = interface{}

type Response struct {
	Success bool `json:"success"`
	Data Data `json:"data"`
	Meta Meta `json:"meta"`
}