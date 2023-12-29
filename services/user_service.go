package services

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/fakhrads/golang-test/models"
	"github.com/fakhrads/golang-test/utils"
	"golang.org/x/crypto/bcrypt"
)

func ValidateCreditCard(creditCard models.CreditCard) error {
	if creditCard.Type == "" || creditCard.Number == "" || creditCard.Name == "" || creditCard.Expired == "" || creditCard.CVV == "" {
		return fmt.Errorf("Invalid credit card data")
	}
	return nil
}

func SaveUserToDB(user models.User) (int64, error) {
	db := utils.GetDB()
	jsonData, err := json.Marshal(user.Photos)
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	_, err = db.Exec(`
		INSERT INTO users (name, email, address, password, photos, creditcard_type, creditcard_number, creditcard_name, creditcard_expired, creditcard_cvv)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Name, user.Email, user.Address, hashedPassword, jsonData, user.CreditCard.Type, user.CreditCard.Number, user.CreditCard.Name, user.CreditCard.Expired, user.CreditCard.CVV,
	)
	if err != nil {
		fmt.Println("Error inserting user:", err)
		return 0, err
	}
	lastInsertID, err := getLastInsertID(db)
	if err != nil {
		fmt.Println("Error getting last insert ID:", err)
		return 0, err
	}

	return lastInsertID, nil
}

func GetUsersFromDatabase(q, ob, sb, of, lt string) ([]models.User, error) {
	db := utils.GetDB()
	query := "SELECT user_id, name, email, address, photos, creditcard_type, creditcard_number, creditcard_name, creditcard_expired, creditcard_cvv FROM users"
	if q != "" {
		query += fmt.Sprintf(" WHERE name LIKE '%%%s%%' OR email LIKE '%%%s%%'", q, q)
	}
	if ob != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", ob, sb)
	}
	if lt != "" && of != "" {
		query += fmt.Sprintf(" LIMIT %s OFFSET %s", lt, of)
	}
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var user models.User
		var photosJSON, creditCardType, creditCardNumber, creditCardName, creditCardExpired, creditCardCVV sql.NullString

		if err := rows.Scan(&user.UserID, &user.Name, &user.Email, &user.Address, &photosJSON, &creditCardType, &creditCardNumber, &creditCardName, &creditCardExpired, &creditCardCVV); err != nil {
			return nil, err
		}

		if photosJSON.Valid {
			var photoMap map[int]string
			if err := json.Unmarshal([]byte(photosJSON.String), &photoMap); err != nil {
				return nil, err
			}
			user.Photos = photoMap
		}

		user.CreditCard = models.CreditCard{
			Type:    creditCardType.String,
			Number:  creditCardNumber.String,
			Name:    creditCardName.String,
			Expired: creditCardExpired.String,
			CVV:     creditCardCVV.String,
		}

		users = append(users, user)
	}

	return users, nil
}

func GetUserDetailFromDatabase(userID string) (*models.User, error) {
	db := utils.GetDB()
	row := db.QueryRow("SELECT user_id, name, email, address FROM users WHERE user_id = ?", userID)
	user := &models.User{}
	err := row.Scan(&user.UserID, &user.Name, &user.Email, &user.Address)
	if err != nil {
		return nil, err
	}
	photos, err := GetPhotosForUser(userID, user)
	if err != nil {
		return nil, err
	}
	user.Photos = photos
	creditCard, err := getCreditCardForUser(userID)
	if err != nil {
		return nil, err
	}
	user.CreditCard = *creditCard
	return user, nil
}

func GetPhotosForUser(userID string, user *models.User) (map[int]string, error) {
	db := utils.GetDB()
	var photosJSON string
	err := db.QueryRow("SELECT photos FROM users WHERE user_id = ?", userID).Scan(&photosJSON)
	if err != nil {
		return nil, err
	}
	var photoMap map[int]string
	if err := json.Unmarshal([]byte(photosJSON), &photoMap); err != nil {
		return nil, err
	}
	user.Photos = photoMap
	return photoMap, nil
}

func getCreditCardForUser(userID string) (*models.CreditCard, error) {
	db := utils.GetDB()
	creditCard := &models.CreditCard{}
	err := db.QueryRow("SELECT creditcard_type, RIGHT(creditcard_number, 4), creditcard_name, creditcard_expired, creditcard_cvv FROM users WHERE user_id = ?", userID).
		Scan(&creditCard.Type, &creditCard.Number, &creditCard.Name, &creditCard.Expired, &creditCard.CVV)
	if err != nil {
		return nil, err
	}

	return creditCard, nil
}

func UpdateUserInDatabase(user *models.User) error {
	db := utils.GetDB()
	var hashedPassword string
	if user.Password != "" {
		var err error
		hashedPassword, err = hashPassword(user.Password)
		if err != nil {
			return err
		}
	}

	result, err := db.Exec("UPDATE users SET name=?, address=?, email=?, password=?, creditcard_type=?, creditcard_number=?, creditcard_name=?, creditcard_expired=?, creditcard_cvv=? WHERE user_id=?",
		user.Name, user.Address, user.Email, hashedPassword, user.CreditCard.Type, user.CreditCard.Number, user.CreditCard.Name, user.CreditCard.Expired, user.CreditCard.CVV, user.UserID)

	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Rows were not affected, indicating that the user with the specified user_id was not found
		return fmt.Errorf("Something went wrong. Please try again later.")
	}

	return nil
}

func ValidateRequest(user *models.User) error {

	if user.Name == "" {
		return fmt.Errorf("Please provide name fields.")
	}
	if user.Email == "" {
		return fmt.Errorf("Please provide email fields.")
	}
	if user.Address == "" {
		return fmt.Errorf("Please provide address fields.")
	}
	if user.Password == "" {
		return fmt.Errorf("Please provide password fields.")
	}
	if user.CreditCard.Type == "" {
		return fmt.Errorf("Please provide creditcard_type fields.")
	}
	if user.CreditCard.Number == "" {
		return fmt.Errorf("Please provide creditcard_number fields.")
	}
	if user.CreditCard.Name == "" {
		return fmt.Errorf("Please provide creditcard_name fields.")
	}
	if user.CreditCard.Expired == "" {
		return fmt.Errorf("Please provide creditcard_expired fields.")
	}
	if user.CreditCard.CVV == "" {
		return fmt.Errorf("Please provide creditcard_cvv fields.")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func getLastInsertID(db *sql.DB) (int64, error) {
	// Execute the query to retrieve the last insert ID
	var lastInsertID int64
	err := db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}
