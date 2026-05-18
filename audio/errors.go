package audio

import "errors"

// ErrNotInitialized is returned when an operation is attempted before Init.
var ErrNotInitialized = errors.New("AudioPlayerEngine 未初始化,请先调用 Init()")
