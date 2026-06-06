package board

import "go.uber.org/zap"

func logError(err error, logger *zap.SugaredLogger, fields ...any) error {
	if logger != nil {
		logger.Errorw("board service operation failed", append([]any{"error", err}, fields...)...)
	}
	return err
}

func logSuccess(logger *zap.SugaredLogger, fields ...any) {
	if logger != nil {
		logger.Infow("board service operation succeeded", fields...)
	}
}
