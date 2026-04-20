package user

import (
	"go.uber.org/zap"
)

func logError(err error, logger *zap.SugaredLogger, fields ...any) error {
	logger.Errorw("user operation failure",
		append(fields, "error", err)...,
	)
	return err
}

func logSuccess(logger *zap.SugaredLogger, fields ...any) {
	logger.Infow("user operation succesful", fields...,
	)
}
