package domain

type Product struct {
	Id           int    `json:"id" db:"id"`
	Article      string `json:"article" db:"article"`
	Name         string `json:"name" db:"name"`
	Price        int    `json:"price" db:"price"`
	Manufacturer string `json:"manufacturer" db:"manufacturer"`
	SellerId     int `json:"seller_id" db:"seller_id"`
	Deleted      bool   `json:"deleted" db:"deleted"`
}
 