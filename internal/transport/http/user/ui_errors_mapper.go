package user

import (
	"errors"

	userApplication "task_tracker/internal/application/user"
	"task_tracker/internal/common_errors"
	domainUser "task_tracker/internal/domain/user"
	valueobjects "task_tracker/internal/domain/value_objects"
)

func mapUIError(err error) string {
	switch {
	case errors.Is(err, userApplication.ErrInvalidUserID):
		return "Некорректный идентификатор пользователя."

	case errors.Is(err, userApplication.ErrInvalidTeamID):
		return "Некорректный идентификатор команды."

	case errors.Is(err, userApplication.ErrInvalidBirthDate):
		return "Некорректная дата рождения."

	case errors.Is(err, valueobjects.ErrInvalidEmail):
		return "Введите корректный email."

	case errors.Is(err, valueobjects.ErrWeakPassword):
		return "Пароль слишком слабый: минимум 8 символов, строчная и заглавная буква, цифра и специальный символ."

	case errors.Is(err, domainUser.ErrManagerMustHaveTeam):
		return "Для менеджера нужно указать команду."

	case errors.Is(err, userApplication.ErrUserNotFound),
		errors.Is(err, common_errors.ErrNotFound),
		errors.Is(err, domainUser.ErrNotFound):
		return "Пользователь не найден."

	case errors.Is(err, userApplication.ErrPersonalDataNotFound):
		return "Персональные данные не найдены."

	case errors.Is(err, userApplication.ErrTeamNotFound):
		return "Команда не найдена."

	case errors.Is(err, domainUser.ErrEmailAlreadyUsed),
		errors.Is(err, userApplication.ErrUserAlreadyExists),
		errors.Is(err, common_errors.ErrAlreadyExists),
		errors.Is(err, domainUser.ErrAlreadyExists):
		return "Пользователь с такой почтой уже существует."

	case errors.Is(err, userApplication.ErrUserAlreadyDeleted):
		return "Пользователь уже удален."

	case errors.Is(err, userApplication.ErrOnlyManagersCanModify),
		errors.Is(err, common_errors.ErrPermissionDenied),
		errors.Is(err, domainUser.ErrPermissionDenied):
		return "Недостаточно прав для выполнения действия."

	case errors.Is(err, userApplication.ErrInvalidRole),
		errors.Is(err, userApplication.ErrInvalidInput),
		errors.Is(err, userApplication.ErrRoleRequired),
		errors.Is(err, domainUser.ErrInvalidRole):
		return "Проверьте данные пользователя: роль указана неверно."

	case errors.Is(err, domainUser.ErrEmptyData):
		return "Заполните персональные данные пользователя."

	default:
		return "Не удалось создать пользователя. Попробуйте еще раз."
	}
}

func mapUIFormError(err error) string {
	switch err {
	case errInvalidForm:
		return "Не удалось прочитать форму."
	case errInvalidAge:
		return "Возраст должен быть числом от 0 до 255."
	case errInvalidBirthDate:
		return "Дата рождения должна быть в формате ГГГГ-ММ-ДД."
	case errRequiredFieldsMissing:
		return "Заполните обязательные поля: email, пароль, роль, имя и фамилию."
	default:
		return "Проверьте данные формы."
	}
}
