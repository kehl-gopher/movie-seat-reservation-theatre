package utility

import (
	"errors"
	"fmt"
	"strconv"
)

func VerifyPort(port string) (int, error) {
	p, err := strconv.Atoi(port)
	if err != nil {
		return 0, errors.New("invalid port expected integer")
	}

	if p < 0 || p > 65535 {
		fmt.Println("Port provided is an invalid postgres port.... falling back to default postgres port")
		p = 5432
	}

	return p, nil
}
