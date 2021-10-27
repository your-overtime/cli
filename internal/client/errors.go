package client

import "errors"

func IsConflictErr(err error) bool {

	return err.Error() == "409 Conflict"
}

func IsNoActivityRunningErr(err error) bool {
	return err == errNoActiviyRunning
}

// some errors
var errNoActiviyRunning = errors.New("no activity is running")
