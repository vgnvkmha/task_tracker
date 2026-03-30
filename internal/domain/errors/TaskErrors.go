package taskerrors

import "errors"

var (
	InvalidStatus           = errors.New("Invalid Status. Should be To DO, In Progress, etc.")
	InvalidStatusTransition = errors.New("Invalid Status Transition")
	AdminCanModifyOnly      = errors.New("Only Administartors Can Modify Task At This Stage")
)
