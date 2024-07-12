package types

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jibaru/do/internal/utils"
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
	UuidFuncName = "uuid"
	DateFuncName = "date"
)

type DateGenerator interface {
	Now() time.Time
}

type Func struct {
	Name string
	Args []interface{}
}

func NewFunc(name string, args []interface{}) (Func, error) {
	for i, arg := range args {
		switch arg.(type) {
		case Map:
			return Func{}, errors.New("map not allowed as argument for argument " + fmt.Sprintf("%v", i+1))
		}
	}

	return Func{Name: name, Args: args}, nil
}

func (f Func) HasReferences() bool {
	for _, arg := range f.Args {
		switch arg.(type) {
		case ReferenceToVariable:
			return true
		case Func:
			if arg.(Func).HasReferences() {
				return true
			}
		}
	}

	return false
}

func (f Func) Resolve(
	uuidFactory utils.UuidFactory,
	dateFactory utils.DateFactory,
) (interface{}, error) {
	switch f.Name {
	case FileFuncName:
		return NewFileFuncFromArgs(f.Args)
	case EnvFuncName:
		return NewEnvFuncFromArgs(f.Args)
	case UuidFuncName:
		return NewUuidFuncFromArgs(f.Args, uuidFactory)
	case DateFuncName:
		return NewDateFuncFromArgs(f.Args, dateFactory)
	}

	return nil, errors.New("unknown function")
}

type EnvFunc struct {
	Arg1 String
	Arg2 String
}

func NewEnvFuncFromArgs(args []interface{}) (EnvFunc, error) {
	if len(args) < 1 || len(args) > 2 {
		return EnvFunc{}, errors.New("invalid number of arguments for env")
	}

	arg, ok := args[0].(String)
	if !ok {
		return EnvFunc{}, errors.New("invalid argument type for env")
	}

	if len(args) == 1 {
		return NewEnvFunc(arg, "")
	}

	arg2, ok := args[1].(String)
	if !ok {
		return EnvFunc{}, errors.New("invalid argument type for env")
	}

	return NewEnvFunc(arg, arg2)
}

func NewEnvFunc(arg1, arg2 String) (EnvFunc, error) {
	if strings.TrimSpace(string(arg1)) == "" {
		return EnvFunc{}, EmptyArgError{funcName: EnvFuncName, position: 1}
	}

	return EnvFunc{Arg1: arg1, Arg2: arg2}, nil
}

func (f EnvFunc) Exec() String {
	val, exists := os.LookupEnv(string(f.Arg1))
	if !exists {
		return String(f.Arg2)
	}

	return String(val)
}

type FileFunc struct {
	Path String
}

func NewFileFuncFromArgs(args []interface{}) (FileFunc, error) {
	if len(args) != 1 {
		return FileFunc{}, errors.New("invalid number of arguments for file")
	}

	arg, ok := args[0].(String)
	if !ok {
		return FileFunc{}, errors.New("invalid argument type for file")
	}

	return NewFileFunc(arg)
}

func NewFileFunc(path String) (FileFunc, error) {
	if strings.TrimSpace(string(path)) == "" {
		return FileFunc{}, EmptyArgError{funcName: "file", position: 1}
	}

	return FileFunc{Path: path}, nil
}

func (f FileFunc) Exec() File {
	return File{Path: string(f.Path)}
}

type UuidFunc struct {
	UuidFactory utils.UuidFactory
}

func NewUuidFuncFromArgs(args []interface{}, uuidFactory utils.UuidFactory) (UuidFunc, error) {
	if len(args) != 0 {
		return UuidFunc{}, errors.New("invalid number of arguments for uuid")
	}

	return UuidFunc{
		UuidFactory: uuidFactory,
	}, nil
}

func (f UuidFunc) Exec() String {
	return String(f.UuidFactory.New())
}

type DateFunc struct {
	Format      String
	DateFactory utils.DateFactory
}

func NewDateFuncFromArgs(args []interface{}, dateFactory utils.DateFactory) (DateFunc, error) {
	if len(args) != 1 {
		return DateFunc{}, errors.New("invalid number of arguments for date")
	}

	arg, ok := args[0].(String)
	if !ok {
		return DateFunc{}, errors.New("invalid argument type for date")
	}

	return NewDateFunc(arg, dateFactory)
}

func NewDateFunc(formatName String, dateFactory utils.DateFactory) (DateFunc, error) {
	valid := map[string]string{
		"ISO8601": "2006-01-02T15:04:05Z",
	}

	format, ok := valid[string(formatName)]

	if !ok {
		return DateFunc{}, errors.New("invalid date format")
	}

	return DateFunc{Format: String(format), DateFactory: dateFactory}, nil
}

func (f DateFunc) Exec() String {
	return String(f.DateFactory.Now().UTC().Format(string(f.Format)))
}
