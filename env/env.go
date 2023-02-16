package env

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var envVars = map[string]string{}

func addEnvVar(line string) error {
	if !strings.Contains(line, "=") {
		return fmt.Errorf("incorrect format of the statement %v : missing = ", line)
	}

	kv := strings.Split(line, "=")
	envVars[kv[0]] = kv[1]

	return nil
}

func appendEnvVars(fPth *string) error {
	f, err := os.Open(*fPth)
	if err != nil {
		return err
	}
	defer f.Close()

	fScan := bufio.NewScanner(f)
	fScan.Split(bufio.ScanLines)

	for fScan.Scan() {
		err = addEnvVar(fScan.Text())
		if err != nil {
			return err
		}
	}

	return nil
}

func getEnvVars(envPth *string) error {
	envFiles, err := ioutil.ReadDir(*envPth)
	if err != nil {
		return err
	}

	for _, file := range envFiles {
		fName := file.Name()
		if len(fName) > 4 &&
			fName[len(fName)-4:] != ".env" {
			continue
		}

		if _, ok := envVars[fName]; ok {
			continue
		}

		if file.IsDir() {
			pth := fmt.Sprintf("%v/%v/", *envPth, fName)
			getEnvVars(&pth)
			continue
		}

		pth := fmt.Sprintf("%v/%v", *envPth, fName)
		err = appendEnvVars(&pth)
		if err != nil {
			return err
		}
	}

	return nil
}

func UnloadVars() {
	for k := range envVars {
		os.Unsetenv(k)
	}
}

func LoadEnvVars(envPth string) error {
	err := getEnvVars(&envPth)
	if err != nil {
		return err
	}

	for k, v := range envVars {
		if os.Getenv(k) != "" {
			continue
		}

		os.Setenv(k, v)
	}

	return nil
}
