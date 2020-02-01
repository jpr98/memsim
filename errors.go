package main

import (
	"errors"
	"fmt"
)

func cmdArgsError(command string, argsNum int) error {
	message := fmt.Sprintf("command %s should have %d arguments", command, argsNum)
	return errors.New(message)
}
