package data

import (
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type SecurityConfig struct {
	PassLower              bool
	PassUpper              bool
	PassNum                bool
	PassSpecial            bool
	PassMinLen             int
	PassMaxLen             int
	EmailDomainsRestricted bool
	AllowedDomains         string
	MaxPOSTBytes           int64
}

type cfgReceiver struct {
	PassLower              string
	PassUpper              string
	PassNum                string
	PassSpecial            string
	PassMinLen             string
	PassMaxLen             string
	EmailDomainsRestricted string
	AllowedDomains         string
	MaxPOSTBytes           string
}

func extractConfig(c cfgReceiver) (SecurityConfig, error) {
	var s SecurityConfig
	var err error

	s.PassLower, err = strconv.ParseBool(c.PassLower)
	s.PassUpper, err = strconv.ParseBool(c.PassUpper)
	s.PassNum, err = strconv.ParseBool(c.PassNum)
	s.PassSpecial, err = strconv.ParseBool(c.PassSpecial)
	s.PassMinLen, err = strconv.Atoi(c.PassMinLen)
	s.PassMaxLen, err = strconv.Atoi(c.PassMaxLen)
	s.EmailDomainsRestricted, err = strconv.ParseBool(c.EmailDomainsRestricted)
	s.AllowedDomains = c.AllowedDomains
	s.MaxPOSTBytes, err = strconv.ParseInt(c.MaxPOSTBytes, 10, 64)

	if err != nil {
		return SecurityConfig{}, err
	}

	return s, nil
}

func (sc *SecurityConfig) GetConfig() (SecurityConfig, error) {
	query := `
		SELECT
			security_config_jsonb ->> 'pass_lower'               AS "PassLower",
			security_config_jsonb ->> 'pass_upper'               AS "PassUpper",
			security_config_jsonb ->> 'pass_num'                 AS "PassNum",
			security_config_jsonb ->> 'pass_special'             AS "PassSpecial",
			security_config_jsonb ->> 'pass_min_len'             AS "PassMinLen",
			security_config_jsonb ->> 'pass_max_len'             AS "PassMaxLen",
			security_config_jsonb ->> 'email_domains_restricted' AS "EmailDomainsRestricted",
			security_config_jsonb ->> 'allowed_domains'          AS "AllowedDomains",
			security_config_jsonb ->> 'max_post_bytes'           AS "MaxPOSTBytes"
		FROM auth.tbl_config
		WHERE
			is_active
	`
	rows, _ := selectRows(query)
	defer rows.Close()

	var cfgs []cfgReceiver

	cfgs, err := pgx.CollectRows(rows, pgx.RowToStructByName[cfgReceiver])
	if err != nil {
		return SecurityConfig{}, err
	}

	if len(cfgs) > 0 {
		return extractConfig(cfgs[0])
	}

	return SecurityConfig{}, errors.New("failed to load security config")

}
