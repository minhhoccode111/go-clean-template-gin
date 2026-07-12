package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/restapi/v1/request"
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/restapi/v1/response"
	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
)

// @Summary     Create task
// @Description Create a new task for the current user
// @ID          create-task
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       request body     request.CreateTask true "Task data"
// @Success     201     {object} entity.Task
// @Failure     400     {object} response.Error
// @Failure     401     {object} response.Error
// @Failure     500     {object} response.Error
// @Security    BearerAuth
// @Router      /tasks [post]
func (r *V1) createTask(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")

		return
	}

	var body request.CreateTask

	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "restapi - v1 - createTask")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - createTask")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	task, err := r.tk.Create(c.Request.Context(), userID.(string), body.Title, body.Description)
	if err != nil {
		r.l.Error(err, "restapi - v1 - createTask")
		errorResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.JSON(http.StatusCreated, task)
}

// @Summary     List tasks
// @Description List tasks for the current user with optional filtering
// @ID          list-tasks
// @Tags        tasks
// @Produce     json
// @Param       status query    string false "Filter by status" Enums(todo, in_progress, done)
// @Param       limit  query    int    false "Limit"  default(10)
// @Param       offset query    int    false "Offset" default(0)
// @Success     200    {object} response.TaskList
// @Failure     401    {object} response.Error
// @Failure     500    {object} response.Error
// @Security    BearerAuth
// @Router      /tasks [get]
func (r *V1) listTasks(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")

		return
	}

	var status *entity.TaskStatus

	if s := c.Query("status"); s != "" {
		ts := entity.TaskStatus(s)
		if !ts.Valid() {
			errorResponse(c, http.StatusBadRequest, "invalid task status")

			return
		}

		status = &ts
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}

	tasks, total, err := r.tk.List(c.Request.Context(), userID.(string), status, limit, offset)
	if err != nil {
		r.l.Error(err, "restapi - v1 - listTasks")
		errorResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.JSON(http.StatusOK, response.TaskList{
		Tasks: tasks,
		Total: total,
	})
}

// @Summary     Get task
// @Description Get a task by ID
// @ID          get-task
// @Tags        tasks
// @Produce     json
// @Param       id  path     string true "Task ID"
// @Success     200 {object} entity.Task
// @Failure     401 {object} response.Error
// @Failure     403 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    BearerAuth
// @Router      /tasks/{id} [get]
func (r *V1) getTask(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")

		return
	}

	taskID := c.Param("id")

	task, err := r.tk.Get(c.Request.Context(), userID.(string), taskID)
	if err != nil {
		r.l.Error(err, "restapi - v1 - getTask")

		if errors.Is(err, entity.ErrTaskNotFound) {
			errorResponse(c, http.StatusNotFound, "task not found")

			return
		}

		if errors.Is(err, entity.ErrTaskForbidden) {
			errorResponse(c, http.StatusForbidden, "forbidden")

			return
		}

		errorResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.JSON(http.StatusOK, task)
}

// @Summary     Update task
// @Description Update task title and description
// @ID          update-task
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       id      path     string            true "Task ID"
// @Param       request body     request.UpdateTask  true "Updated task data"
// @Success     200     {object} entity.Task
// @Failure     400     {object} response.Error
// @Failure     401     {object} response.Error
// @Failure     403     {object} response.Error
// @Failure     404     {object} response.Error
// @Failure     500     {object} response.Error
// @Security    BearerAuth
// @Router      /tasks/{id} [put]
func (r *V1) updateTask(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")

		return
	}

	taskID := c.Param("id")

	var body request.UpdateTask

	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "restapi - v1 - updateTask")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - updateTask")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	task, err := r.tk.Update(c.Request.Context(), userID.(string), taskID, body.Title, body.Description)
	if err != nil {
		r.l.Error(err, "restapi - v1 - updateTask")

		if errors.Is(err, entity.ErrTaskNotFound) {
			errorResponse(c, http.StatusNotFound, "task not found")

			return
		}

		if errors.Is(err, entity.ErrTaskForbidden) {
			errorResponse(c, http.StatusForbidden, "forbidden")

			return
		}

		errorResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.JSON(http.StatusOK, task)
}

// @Summary     Transition task status
// @Description Change task status (todo -> in_progress -> done, or in_progress -> todo)
// @ID          transition-task
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       id      path     string                true "Task ID"
// @Param       request body     request.TransitionTask  true "New status"
// @Success     200     {object} entity.Task
// @Failure     400     {object} response.Error
// @Failure     401     {object} response.Error
// @Failure     403     {object} response.Error
// @Failure     404     {object} response.Error
// @Failure     500     {object} response.Error
// @Security    BearerAuth
// @Router      /tasks/{id}/status [patch]
func (r *V1) transitionTask(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")

		return
	}

	taskID := c.Param("id")

	var body request.TransitionTask

	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "restapi - v1 - transitionTask")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - transitionTask")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	task, err := r.tk.Transition(c.Request.Context(), userID.(string), taskID, body.Status)
	if err != nil {
		r.l.Error(err, "restapi - v1 - transitionTask")

		if errors.Is(err, entity.ErrTaskNotFound) {
			errorResponse(c, http.StatusNotFound, "task not found")

			return
		}

		if errors.Is(err, entity.ErrTaskForbidden) {
			errorResponse(c, http.StatusForbidden, "forbidden")

			return
		}

		if errors.Is(err, entity.ErrInvalidTransition) {
			errorResponse(c, http.StatusBadRequest, "invalid status transition")

			return
		}

		errorResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.JSON(http.StatusOK, task)
}

// @Summary     Delete task
// @Description Delete a task by ID
// @ID          delete-task
// @Tags        tasks
// @Param       id  path     string true "Task ID"
// @Success     204 "No Content"
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    BearerAuth
// @Router      /tasks/{id} [delete]
func (r *V1) deleteTask(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")

		return
	}

	taskID := c.Param("id")

	err := r.tk.Delete(c.Request.Context(), userID.(string), taskID)
	if err != nil {
		r.l.Error(err, "restapi - v1 - deleteTask")

		if errors.Is(err, entity.ErrTaskNotFound) {
			errorResponse(c, http.StatusNotFound, "task not found")

			return
		}

		errorResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.Status(http.StatusNoContent)
}
