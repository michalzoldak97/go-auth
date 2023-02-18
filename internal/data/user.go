package data

import (
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Create(user User) (string, error) {

	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), passCost)
	if err != nil {
		return "", err
	}

	query := `
		INSERT INTO auth.tbl_user (email, first_name, last_name, password)
		VALUES ($1, $2, $3, $4)
		RETURNING user_id
	`
	var newID string
	err = selectRow(query,
		user.Email,
		user.FirstName,
		user.LastName,
		hashPass,
	).Scan(&newID)

	if err != nil {
		return "", err
	}

	return newID, err
}

func (u *User) GetByEmail(email string) ([]User, error) {
	query := `
		SELECT 
			u.user_id AS "ID",
			u.email AS "Email",
			u.first_name AS "FirstName",
			u.last_name AS "LastName",
			u.password AS "Password",
			u.created_at AS "CreatedAt",
			u.updated_at AS "UpdatedAt"
		FROM auth.tbl_user u
		WHERE
			u.email = $1
			AND u.deactivated_at IS NULL
	`
	rows, _ := selectRows(query, email)
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[User])
}
