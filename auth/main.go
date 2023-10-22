package main

import (
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type (
	SignUpRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	SignInRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	SignInResponse struct {
		JWTToken string `json:"jwt_token"`
	}

	ProfileResponse struct {
		Email string `json:"email"`
	}

	User struct {
		Email    string
		password string
	}
)

var (
	portWebApi     = ":8080"
	users          = map[string]User{}
	keySecret      = []byte("qwerty123456")
	keyContextUser = "user"
)

func main() {
	webApp := fiber.New()

	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	webApp.Post("/signup", func(c *fiber.Ctx) error {
		var r SignUpRequest
		if err := c.BodyParser(&r); err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		if _, exist := users[r.Email]; exist {
			return c.SendStatus(fiber.StatusConflict)
		}

		users[r.Email] = User{Email: r.Email, password: r.Password}

		return c.SendStatus(fiber.StatusOK)
	})

	webApp.Post("/signin", func(c *fiber.Ctx) error {
		var r SignInRequest

		if err := c.BodyParser(&r); err != nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		user, ok := users[r.Email]

		if !ok {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		if user.password != r.Password {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		payload := jwt.MapClaims{
			"sub": user.Email,
			"exp": time.Now().Add(72 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
		t, err := token.SignedString(keySecret)

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(SignInResponse{JWTToken: t})
	})

	authorizedGroup := webApp.Group("")
	authorizedGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: keySecret,
		},
		ContextKey: keyContextUser,
	}))
	authorizedGroup.Get("/profile", func(c *fiber.Ctx) error {
		jwtToken, ok := c.Context().Value(keyContextUser).(*jwt.Token)

		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		payload, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		user, ok := users[payload["sub"].(string)]

		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.JSON(ProfileResponse{Email: user.Email})
	})

	logrus.Fatal(webApp.Listen(portWebApi))

}
