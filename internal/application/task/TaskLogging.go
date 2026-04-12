package task_service

import (
	"go.uber.org/zap"
)

func logError(err error, logger *zap.SugaredLogger, fields ...any) error {
	logger.Errorw("failed to change task", //TODO: change message
		append(fields, "error", err)...,
	)
	return err
}

func logSuccess(logger *zap.SugaredLogger, fields ...any) {
	logger.Infow("task was changed", fields...,
	)
}
