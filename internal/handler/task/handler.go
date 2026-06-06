package task_handler

import (
	"net/http"
	"net/url"
	"sort"
	"time"

	taskApplication "task_tracker/internal/application/task"
	"task_tracker/internal/domain/auth"
	domaintask "task_tracker/internal/domain/task"
	dto "task_tracker/internal/handler/task/dto"
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler interface {
	ShowTasks(ctx *gin.Context)
	CreateFromUI(ctx *gin.Context)
	UpdateFromUI(ctx *gin.Context)
	DeleteFromUI(ctx *gin.Context)
	Create(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	List(ctx *gin.Context)
	GetByBoardID(ctx *gin.Context)
	GetByAssigneeID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type handler struct {
	service taskApplication.TaskService
}

func New(service taskApplication.TaskService) TaskHandler {
	return &handler{
		service: service,
	}
}

func (h *handler) Create(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.TaskRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	task, err := h.service.Create(ctx.Request.Context(), actor, request.ToServiceInput())
	if err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"task": dto.ToTaskResponse(*task),
	})
}

func (h *handler) GetByID(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, ok := parseIDParam(ctx, "id")
	if !ok {
		return
	}

	task, err := h.service.GetByID(ctx.Request.Context(), actor, id)
	if err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"task": dto.ToTaskResponse(*task),
	})
}

func (h *handler) List(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	filters, ok := parseTaskFilters(ctx)
	if !ok {
		return
	}

	tasks, err := h.service.FindMany(ctx.Request.Context(), actor, filters)
	if err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tasks": dto.ToTaskResponses(tasks),
	})
}

func (h *handler) GetByBoardID(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	boardID, ok := parseIDParam(ctx, "board_id")
	if !ok {
		return
	}

	tasks, err := h.service.FindByBoardID(ctx.Request.Context(), actor, boardID)
	if err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tasks": dto.ToTaskResponses(tasks),
	})
}

func (h *handler) GetByAssigneeID(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	assigneeID, ok := parseIDParam(ctx, "assignee_id")
	if !ok {
		return
	}

	tasks, err := h.service.FindByAssigneeID(ctx.Request.Context(), actor, assigneeID)
	if err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"tasks": dto.ToTaskResponses(tasks),
	})
}

func (h *handler) Update(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, ok := parseIDParam(ctx, "id")
	if !ok {
		return
	}

	var request dto.UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	task, err := h.service.Update(ctx.Request.Context(), actor, id, request.ToServiceInput())
	if err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"task": dto.ToTaskResponse(*task),
	})
}

func (h *handler) Delete(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, ok := parseIDParam(ctx, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), actor, id); err != nil {
		status, msg := mapError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"task_deleted": id.String(),
	})
}

func parseTaskFilters(ctx *gin.Context) (taskApplication.TaskFilters, bool) {
	var filters taskApplication.TaskFilters

	if !parseOptionalUUIDQuery(ctx, "board_id", &filters.BoardID) {
		return taskApplication.TaskFilters{}, false
	}
	if !parseOptionalUUIDQuery(ctx, "assignee_id", &filters.AssigneeID) {
		return taskApplication.TaskFilters{}, false
	}
	if !parseOptionalUUIDQuery(ctx, "reporter_id", &filters.ReporterID) {
		return taskApplication.TaskFilters{}, false
	}
	if !parseOptionalUUIDQuery(ctx, "sprint_id", &filters.SprintID) {
		return taskApplication.TaskFilters{}, false
	}
	if rawStatus := ctx.Query("status"); rawStatus != "" {
		status := domaintask.TaskStatus(rawStatus)
		filters.Status = &status
	}

	return filters, true
}

func parseOptionalUUIDQuery(ctx *gin.Context, name string, dest **uuid.UUID) bool {
	raw := ctx.Query(name)
	if raw == "" {
		return true
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid " + name})
		return false
	}

	*dest = &id
	return true
}

func parseIDParam(ctx *gin.Context, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(ctx.Param(name))
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid " + name})
		return uuid.Nil, false
	}
	return id, true
}

func (h *handler) ShowTasks(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.Redirect(http.StatusSeeOther, "/ui/users/create?auth=required")
		return
	}

	ctx.HTML(http.StatusOK, "task_page", gin.H{
		"title":   "Задачи",
		"tasks":   []taskView{},
		"team_id": actorTeamID(actor),
		"message": taskMessage(ctx.Query("message")),
		"error":   taskError(ctx.Query("error")),
	})
}

func (h *handler) CreateFromUI(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.Redirect(http.StatusSeeOther, "/ui/users/create?auth=required")
		return
	}

	input, err := newCreateTaskFormInput(ctx)
	if err != nil {
		redirectTasksUI(ctx, "", "invalid")
		return
	}

	if _, err := h.service.Create(ctx.Request.Context(), actor, input); err != nil {
		redirectTasksUI(ctx, "", "create")
		return
	}

	redirectTasksUI(ctx, "created", "")
}

func (h *handler) UpdateFromUI(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.Redirect(http.StatusSeeOther, "/ui/users/create?auth=required")
		return
	}

	taskID, input, err := newUpdateTaskFormInput(ctx)
	if err != nil {
		redirectTasksUI(ctx, "", "invalid")
		return
	}

	if _, err := h.service.Update(ctx.Request.Context(), actor, taskID, input); err != nil {
		redirectTasksUI(ctx, "", "update")
		return
	}

	redirectTasksUI(ctx, "updated", "")
}

func (h *handler) DeleteFromUI(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.Redirect(http.StatusSeeOther, "/ui/users/create?auth=required")
		return
	}

	taskID, err := uuid.Parse(ctx.PostForm("task_id"))
	if err != nil {
		redirectTasksUI(ctx, "", "invalid")
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), actor, taskID); err != nil {
		redirectTasksUI(ctx, "", "delete")
		return
	}

	redirectTasksUI(ctx, "deleted", "")
}

func newCreateTaskFormInput(ctx *gin.Context) (taskApplication.CreateTaskInput, error) {
	dueTo, err := parseOptionalDateTime(ctx.PostForm("due_to"))
	if err != nil {
		return taskApplication.CreateTaskInput{}, err
	}

	return taskApplication.CreateTaskInput{
		Name:        ctx.PostForm("name"),
		Description: ctx.PostForm("description"),
		DueTo:       dueTo,
	}, nil
}

func newUpdateTaskFormInput(ctx *gin.Context) (uuid.UUID, taskApplication.UpdateTaskInput, error) {
	taskID, err := uuid.Parse(ctx.PostForm("task_id"))
	if err != nil {
		return uuid.Nil, taskApplication.UpdateTaskInput{}, err
	}

	name := ctx.PostForm("name")
	description := ctx.PostForm("description")
	status := domaintask.TaskStatus(ctx.PostForm("status"))

	input := taskApplication.UpdateTaskInput{
		Name:        &name,
		Description: &description,
		Status:      &status,
	}

	if rawDueTo := ctx.PostForm("due_to"); rawDueTo != "" {
		dueTo, err := parseOptionalDateTime(rawDueTo)
		if err != nil {
			return uuid.Nil, taskApplication.UpdateTaskInput{}, err
		}
		input.DueTo = &dueTo
	}

	return taskID, input, nil
}

func parseOptionalDateTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}

	parsed, err := time.ParseInLocation("2006-01-02T15:04", value, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}

func redirectTasksUI(ctx *gin.Context, message string, errorMessage string) {
	values := url.Values{}
	if message != "" {
		values.Set("message", message)
	}
	if errorMessage != "" {
		values.Set("error", errorMessage)
	}
	location := "/ui/tasks"
	if encoded := values.Encode(); encoded != "" {
		location += "?" + encoded
	}
	ctx.Redirect(http.StatusSeeOther, location)
}

func taskMessage(value string) string {
	switch value {
	case "created":
		return "Задача создана"
	case "updated":
		return "Задача обновлена"
	case "deleted":
		return "Задача удалена"
	default:
		return ""
	}
}

func taskError(value string) string {
	switch value {
	case "invalid":
		return "Проверьте данные задачи"
	case "create":
		return "Не удалось создать задачу"
	case "update":
		return "Не удалось обновить задачу"
	case "delete":
		return "Не удалось удалить задачу"
	default:
		return ""
	}
}

type taskView struct {
	ID          string
	Name        string
	Description string
	Status      string
	StatusLabel string
	StatusClass string
	DueTo       string
	DueToInput  string
	DueToClass  string
	CreatedAt   string
	UpdatedAt   string
}

func newTaskViews(tasks []*taskApplication.Task) []taskView {
	views := make([]taskView, 0, len(tasks))
	for _, task := range tasks {
		if task == nil {
			continue
		}

		isOverdue := isTaskOverdue(task)
		view := taskView{
			ID:          task.Id.String(),
			Name:        task.Name,
			Description: task.Description,
			Status:      string(task.Status),
			StatusLabel: taskStatusLabel(task.Status),
			StatusClass: taskStatusClass(task.Status, isOverdue),
			DueTo:       "Без срока",
			DueToClass:  "due-normal",
			CreatedAt:   task.CreatedAt.Format("02.01.2006, 15:04"),
			UpdatedAt:   task.UpdatedAt.Format("02.01.2006, 15:04"),
		}
		if !task.DueTo.IsZero() {
			view.DueTo = task.DueTo.Format("02.01.2006, 15:04")
			view.DueToInput = task.DueTo.Format("2006-01-02T15:04")
		}
		if isOverdue {
			view.DueToClass = "due-overdue"
		}
		views = append(views, view)
	}
	return views
}

func sortTasksForUI(tasks []*taskApplication.Task) {
	sort.SliceStable(tasks, func(i, j int) bool {
		left := tasks[i]
		right := tasks[j]
		if left == nil || right == nil {
			return right != nil
		}

		leftStatus := taskStatusRank(left.Status)
		rightStatus := taskStatusRank(right.Status)
		if leftStatus != rightStatus {
			return leftStatus < rightStatus
		}

		if left.DueTo.IsZero() && right.DueTo.IsZero() {
			return left.CreatedAt.After(right.CreatedAt)
		}
		if left.DueTo.IsZero() {
			return false
		}
		if right.DueTo.IsZero() {
			return true
		}
		return left.DueTo.Before(right.DueTo)
	})
}

func taskStatusRank(status domaintask.TaskStatus) int {
	switch status {
	case domaintask.Todo:
		return 1
	case domaintask.InProgress:
		return 2
	case domaintask.Done:
		return 3
	case domaintask.Closed:
		return 4
	case domaintask.Archieved:
		return 5
	default:
		return 99
	}
}

func taskStatusLabel(status domaintask.TaskStatus) string {
	switch status {
	case domaintask.Todo:
		return "Новая"
	case domaintask.InProgress:
		return "В работе"
	case domaintask.Done:
		return "Готова к тестированию"
	case domaintask.Closed:
		return "Закрыта"
	case domaintask.Archieved:
		return "В архиве"
	default:
		return "Без статуса"
	}
}

func taskStatusClass(status domaintask.TaskStatus, overdue bool) string {
	if overdue {
		return "status-overdue"
	}

	switch status {
	case domaintask.Todo:
		return "status-todo"
	case domaintask.InProgress:
		return "status-progress"
	case domaintask.Done:
		return "status-done"
	case domaintask.Closed:
		return "status-closed"
	case domaintask.Archieved:
		return "status-archived"
	default:
		return "status-empty"
	}
}

func isTaskOverdue(task *taskApplication.Task) bool {
	if task == nil || task.DueTo.IsZero() {
		return false
	}
	if task.Status == domaintask.Done || task.Status == domaintask.Closed || task.Status == domaintask.Archieved {
		return false
	}
	return time.Now().After(task.DueTo)
}

func actorTeamID(actor auth.Actor) string {
	if actor.TeamID == nil {
		return ""
	}
	return actor.TeamID.String()
}
