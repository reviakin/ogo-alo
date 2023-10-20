package tasks

import (
	"strconv"

	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type (
	GetTaskResponse struct {
		ID       int64  `json:"id"`
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	CreateTaskRequest struct {
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	CreateTaskResponse struct {
		ID int64 `json:"id"`
	}

	UpdateTaskRequest struct {
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	Task struct {
		ID       int64
		Desc     string
		Deadline int64
	}
)

var taskIDCounter int64 = 1
var errorNoTask = errors.New("no task")

type TaskOperations interface {
	CreateTask(desc string, deadline int64) int64
	UpdateTask(id int64, desc string, deadline int64) (Task, error)
	GetTask(id int64) (Task, error)
	DeleteTask(id int64) error
}

type TaskStorage struct {
	tasks map[int64]Task
}

func (s *TaskStorage) CreateTask(desc string, deadline int64) int64 {
	id := taskIDCounter
	taskIDCounter++

	s.tasks[id] = Task{
		ID:       id,
		Desc:     desc,
		Deadline: deadline,
	}

	return id
}
func (s *TaskStorage) UpdateTask(id int64, desc string, deadline int64) (Task, error) {
	task, ok := s.tasks[id]
	if !ok {
		return Task{}, errorNoTask
	} else {
		task.Desc = desc
		task.Deadline = deadline
    s.tasks[id] = task
		return task, nil
	}
}
func (s *TaskStorage) GetTask(id int64) (Task, error) {
	task, ok := s.tasks[id]
	if !ok {
		return Task{}, errorNoTask
	} else {
		return task, nil
	}
}
func (s *TaskStorage) DeleteTask(id int64) error {
	_, ok := s.tasks[id]
	if !ok {
		return errorNoTask
	} else {
		delete(s.tasks, id)
		return nil
	}
}

type TaskHandler struct {
	storage TaskOperations
}

func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var r CreateTaskRequest
	if err := c.BodyParser(&r); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	id := h.storage.CreateTask(r.Desc, r.Deadline)
	return c.JSON(CreateTaskResponse{
		ID: id,
	})
}
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	var r UpdateTaskRequest
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := c.BodyParser(&r); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	_, e := h.storage.UpdateTask(int64(id), r.Desc, r.Deadline)
	if e != nil {
		return c.SendStatus(fiber.StatusNotFound)
	} else {
		return c.SendStatus(fiber.StatusOK)
	}

}
func (h *TaskHandler) GetTask(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	} else {
		t, e := h.storage.GetTask(int64(id))
		if e != nil {
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			return c.JSON(GetTaskResponse{
        ID: t.ID,
        Desc: t.Desc,
        Deadline: t.Deadline,
      })
		}
	}
}
func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	} else {
		e := h.storage.DeleteTask(int64(id))
		if e != nil {
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			return c.SendStatus(fiber.StatusOK)
		}
	}
}

func main() {
	webApp := fiber.New()
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// BEGIN (write your solution here) (write your solution here)
	taskHandler := &TaskHandler{
		storage: &TaskStorage{
			tasks: make(map[int64]Task),
		},
	}
	webApp.Post("/tasks", taskHandler.CreateTask)
	webApp.Patch("/tasks/:id", taskHandler.UpdateTask)
	webApp.Get("/tasks/:id", taskHandler.GetTask)
	webApp.Delete("/tasks/:id", taskHandler.DeleteTask)
	// END

	logrus.Fatal(webApp.Listen(":8080"))
}
