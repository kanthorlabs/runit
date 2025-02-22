package main

import (
	"errors"
	"fmt"
)

func ErrOneArgOnly() error {
	return errors.New("ERROR.RUNIT.CMD.ONE_ARG_ONLY")
}

func ErrScriptPathRequired() error {
	return errors.New("ERROR.RUNIT.CMD.SCRIPT_PATH_REQUIRED")
}

func ErrInvalidScriptPath() error {
	return errors.New("ERROR.RUNIT.CMD.INVALID_SCRIPT_PATH")
}

func ErrScriptNotFound() error {
	return errors.New("ERROR.RUNIT.CMD.SCRIPT_NOT_FOUND")
}

func ErrArgsParams(reason string) error {
	return fmt.Errorf("ERROR.RUNIT.CMD.ARGS_PARAMS: Unable to parse [%s]", reason)
}
