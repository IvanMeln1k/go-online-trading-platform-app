package domain

type Card struct {
	Id     int    `json:"id" db:"id"`
	Number string `json:"number" db:"number"`
	Data   string `json:"data,omitempty" db:"data"`
	Cvv    string `json:"cvv,omitempty" db:"cvv"`
	UserId int    `json:"user_id,omitempty" db:"user_id"`
}
