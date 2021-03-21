package eventbus

import "errors"

var (
	ErrInvalidFnCallback = errors.New("Invalid FnCallback")
	ErrInvalidPayload    = errors.New("Invalid payload")
	ErrStillRunning      = errors.New("Current FnCallback is still running")
	ErrFnMustBeFunc      = errors.New("Fn must be a func type")
	ErrSignature         = errors.New("Function dosen't match the signature of the event")
)
