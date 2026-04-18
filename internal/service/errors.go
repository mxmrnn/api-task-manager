package service

import "errors"

var ErrTitleTooShort = errors.New("title must be at least 3 characters")
var ErrInvalidTaskStatus = errors.New("invalid task status")

var ErrUserFullNameRequired = errors.New("user full name required")
var ErrUserEmailRequired = errors.New("user email required")
