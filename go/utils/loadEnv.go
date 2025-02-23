package utils

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/FiSeStRo/Ecoland-Backend-Service/authentication"
)

func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		err = SetEnvVariables(parts)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetEnvVariables(parts []string) error {

	var err error = nil
	switch parts[0] {
	case "JWT_SECRET":
		authentication.JwtVariables.Key = []byte(parts[1])
	case "JWT_EXPIRATION_TIME_AT":
		authentication.JwtVariables.ExpirationTime, err = strconv.Atoi(parts[1])
	case "JWT_EXPIRATION_TIME_RT":
		authentication.JwtVariables.ExpirationTimeRT, err = strconv.Atoi(parts[1])
	case "JWT_Issuer":
		authentication.JwtVariables.Issuer = parts[1]
	default:
		err = os.Setenv(parts[0], parts[1])
	}
	return err
}
