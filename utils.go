package config

import (
	"fmt"
	"strings"

	"github.com/hatchify/errors"
)

const (
	// ErrInvalidRoot is returned whe a root is longer than the request path
	ErrInvalidRoot = errors.Error("invalid root, cannot be longer than request path")
)

func getHandlerParts(handlerKey string) (key, handler string, args []string, err error) {
	spl := strings.SplitN(handlerKey, ".", 2)
	if len(spl) != 2 {
		err = fmt.Errorf("expected key and handler, received \"%s\"", handlerKey)
		return
	}

	key = spl[0]
	handler = spl[1]

	spl = strings.Split(handler, "(")
	if len(spl) == 1 {
		return
	}

	handler = spl[0]
	argsStr := spl[1]

	if argsStr[len(argsStr)-1] != ')' {
		err = ErrExpectedEndParen
		return
	}

	argsStr = argsStr[:len(argsStr)-1]
	args = strings.Split(argsStr, ",")
	return
}
