package board

import "errors"

var (
	// general
	ErrNotFound      = errors.New("board not found")
	ErrAlreadyExists = errors.New("board already exists")

	// validation
	ErrEmptyName           = errors.New("board name must be provided")
	ErrNameTooLong         = errors.New("board name is too long")
	ErrInvalidVisibility   = errors.New("invalid board visibility")
	ErrInvalidCreationTime = errors.New("created time must be in the past")

	// roles and rights
	ErrPermissionDenied = errors.New("permission denied")
	ErrNotBoardMember   = errors.New("user is not a board member")

	// status
	ErrInvalidStatus = errors.New("invalid board status")
	ErrArchivedBoard = errors.New("board is archived")

	// logical
	ErrCannotArchive   = errors.New("board cannot be archived")
	ErrAlreadyArchived = errors.New("board is already archived")
	ErrOwnerRequired   = errors.New("board must have an owner")
)
