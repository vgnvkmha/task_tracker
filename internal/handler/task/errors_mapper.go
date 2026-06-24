package task_handler

import (
	"errors"
	"net/http"

	taskApplication "task_tracker/internal/application/task"
)

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, taskApplication.ErrTaskNotFound):
		return http.StatusNotFound, "Задача не найдена"
	case errors.Is(err, taskApplication.ErrBoardRequired):
		return http.StatusUnprocessableEntity, "Выберите доску для задачи"
	case errors.Is(err, taskApplication.ErrBoardNotFound):
		return http.StatusNotFound, "Доска для задачи не найдена"
	case errors.Is(err, taskApplication.ErrActorTeamRequired):
		return http.StatusForbidden, "Нельзя создать задачу: пользователь не состоит в команде"
	case errors.Is(err, taskApplication.ErrBoardTeamMismatch):
		return http.StatusForbidden, "Нельзя создать задачу в доске другой команды"
	case errors.Is(err, taskApplication.ErrPermissionDenied):
		return http.StatusForbidden, "Недостаточно прав для действия с задачей"
	case errors.Is(err, taskApplication.ErrInvalidInput):
		return http.StatusUnprocessableEntity, "Проверьте данные задачи"
	case errors.Is(err, taskApplication.ErrInvalidStatus):
		return http.StatusUnprocessableEntity, "Некорректный статус задачи"
	case errors.Is(err, taskApplication.ErrInvalidTransition):
		return http.StatusUnprocessableEntity, "Недопустимый переход статуса задачи"
	case errors.Is(err, taskApplication.ErrInvalidTaskID):
		return http.StatusUnprocessableEntity, "Некорректный идентификатор задачи"
	default:
		return http.StatusInternalServerError, "Внутренняя ошибка сервера"
	}
}
