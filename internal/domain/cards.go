package domain

type Card struct {
	Id     string `json:"id" database:"id"`
	Number string `json:"number" database:"number"`
	Data   string `json:"data" database:"data"`
	Cvv    string `json:"cvv" database:"cvv"`
	UserId string `json:"user_id" database:"user_id"`
}

type CardReturn struct {
	Id     string `json:"id" database:"id"`
	Number string `json:"number" database:"number"`
}
