package types

import (
	"errors"
	"os"
	"strings"
)

type EmptyArgError struct {
	funcName string
	position int
}

func (e EmptyArgError) Error() string {
	return e.funcName + " arg" + string(e.position) + " is empty"
}

const (
	EnvFuncName = "env"
)

type EnvFunc struct {
	Arg1 string
	Arg2 string
}

func NewEnvFuncFromArgs(args []string) (EnvFunc, error) {
	if len(args) < 1 || len(args) > 2 {
		return EnvFunc{}, errors.New("invalid number of arguments for env")
	}

	return NewEnvFunc(args[0], args[1])
}

func NewEnvFunc(arg1, arg2 string) (EnvFunc, error) {
	if strings.TrimSpace(arg1) == "" {
		return EnvFunc{}, EmptyArgError{funcName: EnvFuncName, position: 1}
	}

	return EnvFunc{Arg1: arg1, Arg2: arg2}, nil
}

func (f EnvFunc) Exec() string {
	val, exists := os.LookupEnv(f.Arg1)
	if !exists {
		return f.Arg2
	}

	return val
}
