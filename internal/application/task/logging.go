package task

import (
	"go.uber.org/zap"
)

func logError(err error, logger *zap.SugaredLogger, fields ...any) error {
	logger.Errorw("task operation failure",
		append(fields, "error", err)...,
	)
	return err
}

func logSuccess(logger *zap.SugaredLogger, fields ...any) {
	logger.Infow("task operation succesful", fields...,
	)
}
