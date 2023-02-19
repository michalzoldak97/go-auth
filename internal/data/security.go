package data

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type SecurityConfig struct {
	PassLower              bool
	PassUpper              bool
	PassNum                bool
	PassSpecial            bool
	PassMinLen             int
	PassMaxLen             int
	PassCost               int
	EmailDomainsRestricted bool
	AllowedDomains         []string
	MaxPOSTBytes           int64
	DBTimeout              time.Duration
}

type cfgReceiver struct {
	PassLower              string
	PassUpper              string
	PassNum                string
	PassSpecial            string
	PassMinLen             string
	PassMaxLen             string
	PassCost               string
	EmailDomainsRestricted string
	AllowedDomains         string
	MaxPOSTBytes           string
	DBTimeout              string
}

func parseAllowedDomains(areRequired bool, domainsCsv string) ([]string, error) {

	if !areRequired {
		return []string{}, nil
	}

	return strings.Split(domainsCsv, ","), nil
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
	s.PassCost, err = strconv.Atoi(c.PassCost)
	s.EmailDomainsRestricted, err = strconv.ParseBool(c.EmailDomainsRestricted)
	s.MaxPOSTBytes, err = strconv.ParseInt(c.MaxPOSTBytes, 10, 64)
	dbTimeout, err := strconv.Atoi(c.DBTimeout)

	if err != nil {
		return SecurityConfig{}, err
	}

	s.AllowedDomains, err = parseAllowedDomains(s.EmailDomainsRestricted, c.AllowedDomains)
	s.DBTimeout = time.Second * time.Duration(dbTimeout)

	if err != nil {
		return SecurityConfig{}, err
	}

	return s, nil
}

func (sc *SecurityConfig) GetConfig() (SecurityConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `
		SELECT
			security_config_jsonb ->> 'pass_lower'               AS "PassLower",
			security_config_jsonb ->> 'pass_upper'               AS "PassUpper",
			security_config_jsonb ->> 'pass_num'                 AS "PassNum",
			security_config_jsonb ->> 'pass_special'             AS "PassSpecial",
			security_config_jsonb ->> 'pass_min_len'             AS "PassMinLen",
			security_config_jsonb ->> 'pass_max_len'             AS "PassMaxLen",
			security_config_jsonb ->> 'pass_cost'             	 AS "PassCost",
			security_config_jsonb ->> 'email_domains_restricted' AS "EmailDomainsRestricted",
			security_config_jsonb ->> 'allowed_domains'          AS "AllowedDomains",
			security_config_jsonb ->> 'max_post_bytes'           AS "MaxPOSTBytes",
			security_config_jsonb ->> 'db_timeout'           	 AS "DBTimeout"
		FROM auth.tbl_config
		WHERE
			is_active
	`
	rows, err := db.Query(ctx, query)

	defer rows.Close()

	var cfgs []cfgReceiver

	cfgs, err = pgx.CollectRows(rows, pgx.RowToStructByName[cfgReceiver])
	if err != nil {
		return SecurityConfig{}, err
	}

	if len(cfgs) > 0 {
		res, err := extractConfig(cfgs[0])
		security = res
		return res, err
	}

	return SecurityConfig{}, errors.New("failed to load security config")

}
