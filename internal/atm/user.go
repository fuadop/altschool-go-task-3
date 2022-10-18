package atm

const DefaultPIN = "0000"

type User struct {
	Username string  `json:"username"`
	PINHash  string  `json:"pin_hash"`
	Balance  float64 `json:"balance"`
}
