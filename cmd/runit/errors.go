package main

import "errors"

var ErrScriptPathRequired = errors.New("ERROR.RUNIT.CMD.SCRIPT_PATH_REQUIRED")
var ErrInvalidScriptPath = errors.New("ERROR.RUNIT.CMD.INVALID_SCRIPT_PATH")
var ErrScriptNotFound = errors.New("ERROR.RUNIT.CMD.SCRIPT_NOT_FOUND")
