package task_handler

import (
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"time"

	taskApplication "task_tracker/internal/application/task"
	"task_tracker/internal/domain/auth"
	domaintask "task_tracker/internal/domain/task"
	dto "task_tracker/internal/handler/task/dto"
	"task_tracker/internal/render"
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler interface {
	ShowTasks(ctx *gin.Context)
	ShowTaskModal(ctx *gin.Context)
	ShowTaskEditModal(ctx *gin.Context)
	UpdateTaskModal(ctx *gin.Context)
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
		"title":    "Задачи",
		"tasks":    []taskView{},
		"actor_id": actor.ID.String(),
		"team_id":  actorTeamID(actor),
		"message":  taskMessage(ctx.Query("message")),
		"error":    taskError(ctx.Query("error")),
	})
}

func (h *handler) ShowTaskModal(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.HTML(http.StatusUnauthorized, "task_modal", gin.H{
			"Error": "Нужно войти, чтобы открыть задачу.",
		})
		return
	}

	taskID, err := uuid.Parse(ctx.Param("task_id"))
	if err != nil {
		ctx.HTML(http.StatusUnprocessableEntity, "task_modal", gin.H{
			"Error": "Некорректный идентификатор задачи.",
		})
		return
	}

	task, err := h.service.GetByID(ctx.Request.Context(), actor, taskID)
	if err != nil {
		_, msg := mapError(err)
		ctx.HTML(http.StatusOK, "task_modal", gin.H{
			"Error": msg,
		})
		return
	}

	ctx.HTML(http.StatusOK, "task_modal", taskModalViewFromTask(task))
}

func (h *handler) ShowTaskEditModal(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.HTML(http.StatusUnauthorized, "task_modal_edit", taskEditModalView{
			Error: "Нужно войти, чтобы редактировать задачу.",
		})
		return
	}

	taskID, err := uuid.Parse(ctx.Param("task_id"))
	if err != nil {
		ctx.HTML(http.StatusOK, "task_modal_edit", taskEditModalView{
			Error: "Некорректный идентификатор задачи.",
		})
		return
	}

	task, err := h.service.GetByID(ctx.Request.Context(), actor, taskID)
	if err != nil {
		_, msg := mapError(err)
		ctx.HTML(http.StatusOK, "task_modal_edit", taskEditModalView{
			Error: msg,
		})
		return
	}

	ctx.HTML(http.StatusOK, "task_modal_edit", taskEditModalViewFromTask(task, ""))
}

func (h *handler) UpdateTaskModal(ctx *gin.Context) {
	actor, ok := middleware.GetActor(ctx)
	if !ok {
		ctx.HTML(http.StatusUnauthorized, "task_modal_edit", taskEditModalView{
			Error: "Нужно войти, чтобы сохранить задачу.",
		})
		return
	}

	taskID, err := uuid.Parse(ctx.Param("task_id"))
	if err != nil {
		ctx.HTML(http.StatusOK, "task_modal_edit", taskEditModalView{
			Error: "Некорректный идентификатор задачи.",
		})
		return
	}

	input, editView, err := newModalUpdateTaskInput(ctx, taskID)
	if err != nil {
		editView.Error = "Проверьте данные задачи."
		ctx.HTML(http.StatusOK, "task_modal_edit", editView)
		return
	}

	_, err = h.service.Update(ctx.Request.Context(), actor, taskID, input)
	if err != nil {
		_, msg := mapError(err)
		editView.Error = msg
		ctx.HTML(http.StatusOK, "task_modal_edit", editView)
		return
	}

	ctx.Header("HX-Trigger", "taskUpdated")
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(""))
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

func newModalUpdateTaskInput(ctx *gin.Context, taskID uuid.UUID) (taskApplication.UpdateTaskInput, taskEditModalView, error) {
	name := ctx.PostForm("name")
	description := ctx.PostForm("description")
	status := domaintask.TaskStatus(ctx.PostForm("status"))
	view := taskEditModalView{
		ID:          taskID.String(),
		Name:        name,
		Description: description,
		Status:      string(status),
		DueToInput:  ctx.PostForm("due_to"),
		ReporterID:  ctx.PostForm("reporter_id"),
		AssigneeID:  ctx.PostForm("assignee_id"),
		BoardID:     ctx.PostForm("board_id"),
	}

	input := taskApplication.UpdateTaskInput{
		Name:        &name,
		Description: &description,
		Status:      &status,
	}

	if rawDueTo := ctx.PostForm("due_to"); rawDueTo != "" {
		dueTo, err := parseOptionalDateTime(rawDueTo)
		if err != nil {
			return taskApplication.UpdateTaskInput{}, view, err
		}
		input.DueTo = &dueTo
	} else {
		dueTo := time.Time{}
		input.DueTo = &dueTo
	}
	if rawReporterID := ctx.PostForm("reporter_id"); rawReporterID != "" {
		reporterID, err := uuid.Parse(rawReporterID)
		if err != nil {
			return taskApplication.UpdateTaskInput{}, view, err
		}
		input.ReporterID = &reporterID
	}
	if rawAssigneeID := ctx.PostForm("assignee_id"); rawAssigneeID != "" {
		assigneeID, err := uuid.Parse(rawAssigneeID)
		if err != nil {
			return taskApplication.UpdateTaskInput{}, view, err
		}
		input.AssigneeID = &assigneeID
	}
	if rawBoardID := ctx.PostForm("board_id"); rawBoardID != "" {
		boardID, err := uuid.Parse(rawBoardID)
		if err != nil {
			return taskApplication.UpdateTaskInput{}, view, err
		}
		input.BoardID = &boardID
	}

	return input, view, nil
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

type taskModalView struct {
	Error           string
	ID              string
	Name            string
	DescriptionHTML template.HTML
	HasDescription  bool
	Status          string
	StatusLabel     string
	StatusClass     string
	DueTo           string
	CreatedAt       string
	UpdatedAt       string
	ReporterID      string
	AssigneeID      string
	BoardID         string
	SprintID        string
}

type taskEditModalView struct {
	Error       string
	ID          string
	Name        string
	Description string
	Status      string
	DueToInput  string
	ReporterID  string
	AssigneeID  string
	BoardID     string
}

func taskModalViewFromTask(task *taskApplication.Task) taskModalView {
	view := taskModalView{
		ID:              task.Id.String(),
		Name:            task.Name,
		HasDescription:  task.Description != "",
		DescriptionHTML: render.RenderMarkdown(task.Description),
		Status:          string(task.Status),
		StatusLabel:     taskStatusLabel(task.Status),
		StatusClass:     taskStatusClass(task.Status, isTaskOverdue(task)),
		DueTo:           "Без срока",
		CreatedAt:       task.CreatedAt.Format("02.01.2006, 15:04"),
		UpdatedAt:       task.UpdatedAt.Format("02.01.2006, 15:04"),
		ReporterID:      task.ReporterId.String(),
	}
	if !task.DueTo.IsZero() {
		view.DueTo = task.DueTo.Format("02.01.2006, 15:04")
	}
	if task.AssigneeId != nil {
		view.AssigneeID = task.AssigneeId.String()
	}
	if task.BoardId != nil {
		view.BoardID = task.BoardId.String()
	}
	if task.SprintId != nil {
		view.SprintID = task.SprintId.String()
	}
	return view
}

func taskEditModalViewFromTask(task *taskApplication.Task, errorMessage string) taskEditModalView {
	view := taskEditModalView{
		Error:       errorMessage,
		ID:          task.Id.String(),
		Name:        task.Name,
		Description: task.Description,
		Status:      string(task.Status),
		ReporterID:  task.ReporterId.String(),
	}
	if !task.DueTo.IsZero() {
		view.DueToInput = task.DueTo.Format("2006-01-02T15:04")
	}
	if task.AssigneeId != nil {
		view.AssigneeID = task.AssigneeId.String()
	}
	if task.BoardId != nil {
		view.BoardID = task.BoardId.String()
	}
	return view
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
