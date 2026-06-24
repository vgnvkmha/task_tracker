package board

import (
	"errors"
	"net/http"

	boardapp "task_tracker/internal/application/board"
	"task_tracker/internal/common_errors"
)

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, boardapp.ErrBoardNotFound):
		return http.StatusNotFound, "Доска не найдена"
	case errors.Is(err, boardapp.ErrBoardAlreadyExists):
		return http.StatusConflict, "Доска уже существует"
	case errors.Is(err, boardapp.ErrBoardHasTasks):
		return http.StatusConflict, "Нельзя удалить доску: на ней есть задачи"
	case errors.Is(err, boardapp.ErrListBoardsFailed):
		return http.StatusInternalServerError, "Не удалось загрузить доски"
	case errors.Is(err, boardapp.ErrPermissionDenied):
		return http.StatusForbidden, "Недостаточно прав для действия с доской"
	case errors.Is(err, boardapp.ErrBoardTeamMismatch):
		return http.StatusForbidden, "Можно удалять только доски своей команды"
	case errors.Is(err, boardapp.ErrInvalidBoardID),
		errors.Is(err, common_errors.ErrInvalidID):
		return http.StatusUnprocessableEntity, "Некорректный идентификатор доски"
	case errors.Is(err, boardapp.ErrInvalidStatus):
		return http.StatusUnprocessableEntity, "Некорректный статус доски"
	case errors.Is(err, boardapp.ErrInvalidInput):
		return http.StatusUnprocessableEntity, "Проверьте данные доски"
	default:
		return http.StatusInternalServerError, "Неожиданная ошибка"
	}
}
