package userroute

import (
	"github.com/BurakYs/GoAPIExample/internal/config"
	"github.com/BurakYs/GoAPIExample/internal/middleware"
	"github.com/BurakYs/GoAPIExample/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func Register(router fiber.Router, controller *UserController) {
	defaultUserLimiter := limiter.New(config.BaseLimiter(config.LimiterWithMax(50)))

	router.Add(
		[]string{"GET", "HEAD"},
		"/users",
		defaultUserLimiter,
		middleware.ValidateQuery[models.GetAllUsersQuery](),
		controller.GetAllUsers,
	)

	router.Add(
		[]string{"GET", "HEAD"},
		"/users/:id",
		defaultUserLimiter,
		middleware.ValidateParams[models.GetUserByIDParams](),
		controller.GetUserByID,
	)

	router.Post(
		"/register",
		limiter.New(config.BaseLimiter(config.LimiterWithMax(1))),
		middleware.ValidateBody[models.RegisterUserBody](),
		controller.Register,
	)

	router.Post(
		"/login",
		limiter.New(config.BaseLimiter(config.LimiterWithMax(5))),
		middleware.ValidateBody[models.LoginUserBody](),
		controller.Login,
	)

	router.Post(
		"/logout",
		defaultUserLimiter,
		middleware.AuthRequired(),
		controller.Logout,
	)

	router.Delete(
		"/delete-account",
		defaultUserLimiter,
		middleware.AuthRequired(),
		controller.DeleteAccount,
	)
}
