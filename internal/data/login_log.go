package data

// attemptfail
// attempt success
type UserLoginAttempt struct {
	Email   string
	Message string
	Success bool
}

func (u *UserLoginAttempt) Create(ula UserLoginAttempt) error {
	query := `
		INSERT INTO log.tbl_auth_user_login_attempt (email, message, success)
		VALUES ($1, $2, $3)
	`

	err := execRow(query,
		ula.Email,
		ula.Message,
		ula.Success,
	)

	return err
}
