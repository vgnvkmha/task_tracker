package user

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

func logError(err error, logger *zap.SugaredLogger, fields ...any) error {
	logFailure(logger, "user operation failure", err, fields...)
	return err
}

func logSuccess(logger *zap.SugaredLogger, fields ...any) {
	logger.Infof("user operation successful: %s", formatLogFields(append(defaultLogFields(), fields...)...))
}

func logFailure(logger *zap.SugaredLogger, message string, err error, fields ...any) {
	fields = append(defaultLogFields(), fields...)
	if err != nil {
		logger.Errorf("%s\nerror=%v\ncontext=%s", message, err, formatLogFields(fields...))
		return
	}

	logger.Errorf("%s\ncontext=%s", message, formatLogFields(fields...))
}

func defaultLogFields() []any {
	return []any{
		"module", module,
		"layer", layer,
	}
}

func appendLogFields(base []any, fields ...any) []any {
	result := make([]any, 0, len(base)+len(fields))
	result = append(result, base...)
	result = append(result, fields...)
	return result
}

func formatLogFields(fields ...any) string {
	if len(fields) == 0 {
		return ""
	}

	parts := make([]string, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		key := fmt.Sprint(fields[i])
		if i+1 >= len(fields) {
			parts = append(parts, key)
			continue
		}

		value := fields[i+1]
		if value == nil {
			parts = append(parts, key+"=nil")
			continue
		}

		parts = append(parts, fmt.Sprintf("%s=%v", key, value))
	}

	return strings.Join(parts, " ")
}
