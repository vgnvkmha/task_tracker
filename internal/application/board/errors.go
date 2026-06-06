package board

import "errors"

var (
	ErrBoardNotFound      = errors.New("board not found")
	ErrBoardAlreadyExists = errors.New("board already exists")
	ErrCreateBoardFailed  = errors.New("failed to create board")
	ErrUpdateBoardFailed  = errors.New("failed to update board")
	ErrDeleteBoardFailed  = errors.New("failed to delete board")

	ErrInvalidInput     = errors.New("invalid board input")
	ErrInvalidBoardID   = errors.New("invalid board id")
	ErrInvalidStatus    = errors.New("invalid board status")
	ErrPermissionDenied = errors.New("permission denied")
)
