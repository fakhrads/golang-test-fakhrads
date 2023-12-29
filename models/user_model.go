package models

type CreditCard struct {
	Type    string `json:"creditcard_type"`
	Number  string `json:"creditcard_number"`
	Name    string `json:"creditcard_name"`
	Expired string `json:"creditcard_expired"`
	CVV     string `json:"creditcard_cvv"`
}

type User struct {
	UserID     int            `json:"user_id"`
	Name       string         `json:"name"`
	Email      string         `json:"email"`
	Address    string         `json:"address"`
	Password   string         `json:"password"`
	Photos     map[int]string `json:"photos"`
	CreditCard CreditCard     `json:"creditcard"`
}
