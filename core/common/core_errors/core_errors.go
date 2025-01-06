package core_errors

import "errors"

func NotImplemented() error {
	return errors.New("not yet implemented")
}
