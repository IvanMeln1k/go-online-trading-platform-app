package domain

type Card struct {
	Id     int    `json:"id" database:"id"`
	Number string `json:"number" database:"number"`
	Data   string `json:"data,omitempty" database:"data"`
	Cvv    string `json:"cvv,omitempty" database:"cvv"`
	UserId int    `json:"user_id,omitempty" database:"user_id"`
}
