package types

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type EmptyArgError struct {
	funcName string
	position int
}

func (e EmptyArgError) Error() string {
	return e.funcName + " arg" + fmt.Sprintf("%v", e.position) + " is empty"
}

const (
	FileFuncName = "file"
	EnvFuncName  = "env"
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

func (f EnvFunc) Exec() String {
	val, exists := os.LookupEnv(f.Arg1)
	if !exists {
		return String(f.Arg2)
	}

	return String(val)
}

type FileFunc struct {
	Path string
}

func NewFileFuncFromArgs(args []string) (FileFunc, error) {
	if len(args) != 1 {
		return FileFunc{}, errors.New("invalid number of arguments for file")
	}

	return NewFileFunc(args[0])
}

func NewFileFunc(path string) (FileFunc, error) {
	if strings.TrimSpace(path) == "" {
		return FileFunc{}, EmptyArgError{funcName: "file", position: 1}
	}

	return FileFunc{Path: path}, nil
}

func (f FileFunc) Exec() File {
	return File{Path: f.Path}
}
