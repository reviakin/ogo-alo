package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type User struct {
	ID      int64
	Email   string
	Age     int
	Country string
}

var users = map[int64]User{}

type (
	CreateUserRequest struct {
		// BEGIN (write your solution here)
		ID      int64  `json:"id" validate:"required,min=0"`
		Email   string `json:"email" validate:"required,email"`
		Age     int    `json:"age" validate:"required,min=18,max=130"`
		Country string `json:"country" validate:"required,allowable_country"`
		// END
	}
)

type UserOperataions interface {
	CreateUser(id int64, email string, age int, country string)
}

type UserStorage struct {
	users map[int64]User
}

func (s *UserStorage) CreateUser(id int64, email string, age int, country string) {
	s.users[id] = User{
		ID:      id,
		Email:   email,
		Age:     age,
		Country: country,
	}
}

type UserHandler struct {
	storage UserOperataions
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var r CreateUserRequest
	if err := c.BodyParser(&r); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	err := validate.Struct(r)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}

	h.storage.CreateUser(r.ID, r.Email, r.Age, r.Country)

	return nil
}

var validate = validator.New()
var allowedCountries = []string{
	"USA",
	"Germany",
	"France",
}

func main() {
	webApp := fiber.New()
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	var vErr = validate.RegisterValidation("allowable_country", func(fl validator.FieldLevel) bool {
		// Проверяем, что текст не содержит запрещенных слов.
		text := fl.Field().String()
		for _, country := range allowedCountries {
			if text == country {
				return true
			}
		}

		return false
	})
	if vErr != nil {
		logrus.Fatal("register validation ", vErr)
	}

	// BEGIN (write your solution here) (write your solution here)
  userHandler := &UserHandler{
    storage: &UserStorage{
      users: map[int64]User{},
    },
  }

	webApp.Post("/users", userHandler.CreateUser)
	// END
	logrus.Fatal(webApp.Listen(":8080"))
}
