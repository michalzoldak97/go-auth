package data

import (
	"context"
	"errors"
	"time"

	"aidanwoods.dev/go-paseto"
)

type Token struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (t *Token) Create(user User) (string, error) {
	// create token server
	token := paseto.NewToken()
	// set rules
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())

	expiresAt := time.Now().Add(security.TokenExpiration)

	token.SetExpiration(expiresAt)
	// encrypt
	enc := token.V4Encrypt(security.TokenKey, security.TokenSecret)
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

func (t *Token) deactivate(tokenID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), security.DBTimeout)
	defer cancel()

	query := `
		UPDATE auth.tbl_user_token
		SET deactivated_at = $1
		WHERE
			token_id = $2
			AND deactivated_at IS NULL
	`
	_, err := db.Exec(ctx, query,
		time.Now(),
		tokenID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (t *Token) Validate(token string) error {

	if len(token) != security.TokenLen {
		return errors.New("security token malformed")
	}

	query := `
		SELECT
			t.token_id
		FROM auth.tbl_token t
		INNER JOIN auth.tbl_user_token ut				ON t.token_id = ut.token_id
		WHERE
			t.token = $1
			AND ut.deactivated_at IS NULL
	`
	var tokenID string

	err := selectRow(query, token).Scan(&tokenID)
	if err != nil {
		return err
	}

	parser := paseto.NewParserForValidNow()

	_, err = parser.ParseV4Local(security.TokenKey, token, security.TokenSecret)
	if err != nil {
		t.deactivate(tokenID)
		return err
	}

	return nil
}
