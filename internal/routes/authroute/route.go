package authroute

import (
	"github.com/BurakYs/GoAPIExample/internal/middleware"
	"github.com/BurakYs/GoAPIExample/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func Register(router fiber.Router, controller *AuthController) {
	router.Post(
		"/login",
		limiter.New(middleware.NewLimiter(middleware.LimiterWithMax(5))),
		middleware.ValidateBody[models.LoginUserBody](),
		controller.Login,
	)

	router.Post(
		"/logout",
		limiter.New(middleware.NewLimiter(middleware.LimiterWithMax(5))),
		middleware.AuthRequired(),
		controller.Logout,
	)

	router.Post(
		"/register",
		limiter.New(middleware.NewLimiter(middleware.LimiterWithMax(3))),
		middleware.ValidateBody[models.RegisterUserBody](),
		controller.Register,
	)
}
