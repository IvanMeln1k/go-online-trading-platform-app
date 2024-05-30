package domain

import (
	"errors"
)

type Product struct {
	Id           int     `json:"id" db:"id"`
	Article      string  `json:"article" db:"article"`
	Name         string  `json:"name" db:"name"`
	Price        int     `json:"price" db:"price"`
	Manufacturer string  `json:"manufacturer" db:"manufacturer"`
	SellerId     int     `json:"seller_id" db:"seller_id"`
	Deleted      bool    `json:"deleted" db:"deleted"`
	Rating       float32 `json:"rating" db:"rating"`
}

type Filter struct {
	Article      *string  `json:"article" db:"article"`
	Name         *string  `json:"name" db:"name"`
	MinPrice     *int     `json:"min_price" db:"min_price"`
	MaxPrice     *int     `json:"max_price" db:"max_price"`
	Manufacturer *string  `json:"manufacturer" db:"manufacturer"`
	Rating       *float32 `json:"rating" db:"rating"`
	Limit        *int     `json:"limit" db:"limit"`
	Offset       *int     `json:"offset" db:"offset"`
}

func (flt *Filter) Fill_defaults() error {
	if flt.Limit == nil {
		*flt.Limit = 10
	}

	if *flt.Limit < 25 {
		return errors.New("limit too high")
	}

	if flt.Offset == nil {
		*flt.Offset = 10
	}

	return nil
}
