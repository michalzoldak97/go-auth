package data

import (
	"context"
	"time"

	"aidanwoods.dev/go-paseto"
)

type Token struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	TokenHash []byte    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// create private parser
// create implicit

func (t *Token) Create(user User) (string, error) {
	// create token server
	token := paseto.NewToken()
	// set rules
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())

	expiresAt := time.Now().Add(security.TokenExpiration)

	token.SetExpiration(expiresAt)
	// encrypt
	enc := token.V4Encrypt(paseto.NewV4SymmetricKey(), []byte(";ojgniasrgnoagn"))
	// insert to tbl_token and tbl_user_token

	ctx, cancel := context.WithTimeout(context.Background(), security.DBTimeout)
	defer cancel()

	tx, err := db.Begin(ctx)
	if err != nil {
		return "", err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	query := `
		INSERT INTO auth.tbl_token (email, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING token_id
	`
	var newID string

	err = tx.QueryRow(ctx, query,
		user.Email,
		enc,
		expiresAt,
	).Scan(&newID)

	if err != nil {
		return "", err
	}

	query = `
		UPDATE auth.tbl_user_token
		SET deactivated_at = $1
		WHERE
			user_id = $2
			AND deactivated_at IS NULL
	`

	_, err = tx.Exec(ctx, query,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return "", err
	}

	query = `
		INSERT INTO auth.tbl_user_token (user_id, token_id)
		SELECT 
			$1,
			$2
	`

	_, err = tx.Exec(ctx, query,
		user.ID,
		newID,
	)

	if err != nil {
		return "", err
	}

	return enc, nil
}
