package valueobject

import "errors"

var ErrInvalidAddress = errors.New("invalid address: all fields are required")

type Address struct {
	street     string
	city       string
	state      string
	postalCode string
	country    string
}

func NewAddress(street, city, state, postalCode, country string) (Address, error) {
	if street == "" || city == "" || state == "" || postalCode == "" || country == "" {
		return Address{}, ErrInvalidAddress
	}
	return Address{
		street:     street,
		city:       city,
		state:      state,
		postalCode: postalCode,
		country:    country,
	}, nil
}

func NewAddressFromDB(street, city, state, postalCode, country string) Address {
	return Address{street: street, city: city, state: state, postalCode: postalCode, country: country}
}

func (a Address) Street() string     { return a.street }
func (a Address) City() string       { return a.city }
func (a Address) State() string      { return a.state }
func (a Address) PostalCode() string { return a.postalCode }
func (a Address) Country() string    { return a.country }
