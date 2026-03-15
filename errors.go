package inquire

import "errors"

// ErrInterrupt is returned when the user presses Ctrl+C.
var ErrInterrupt = errors.New("interrupted")
