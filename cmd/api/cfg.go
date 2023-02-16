package main

import (
	"fmt"
	"os"
	"strings"
)

func getDSN() (string, error) {
	dsnParams := map[string]string{
		"user":                          "POSTGRES_USER",
		"password":                      "POSTGRES_PASSWORD",
		"host":                          "POSTGRES_HOST_DEV",
		"port":                          "POSTGRES_PORT_DEV",
		"dbname":                        "POSTGRES_DB",
		"sslmode":                       "POSTGRES_SSLMODE",
		"pool_max_conns":                "DB_POOL_MAX_CONNS",
		"pool_max_conn_lifetime":        "DB_POOL_MAX_CONN_LIFETIME",
		"pool_max_conn_idle_time":       "DB_POOL_MAX_CONN_IDLE_TIME",
		"pool_max_conn_lifetime_jitter": "DB_POOL_MAX_CONN_LIFETIME_JITTER",
	}

	var sb strings.Builder

	for k, v := range dsnParams {
		param := os.Getenv(v)

		if len(param) < 1 {
			return "", fmt.Errorf("missing environment variable; %v", k)
		}

		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(param)
		sb.WriteString(" ")
	}

	return sb.String(), nil
}
