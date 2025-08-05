package userroute

import (
	"github.com/BurakYs/GoAPIExample/internal/middleware"
	"github.com/BurakYs/GoAPIExample/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func Register(router fiber.Router, controller *UserController) {
	router.Add(
		[]string{"GET", "HEAD"},
		"/users",
		limiter.New(middleware.NewLimiter(middleware.LimiterWithMax(50))),
		middleware.ValidateQuery[models.GetAllUsersQuery](),
		controller.GetAllUsers,
	)

	router.Add(
		[]string{"GET", "HEAD"},
		"/users/:id",
		limiter.New(middleware.NewLimiter(middleware.LimiterWithMax(50))),
		middleware.ValidateParams[models.GetUserByIDParams](),
		controller.GetUserByID,
	)

	router.Delete(
		"/delete-account",
		middleware.AuthRequired(),
		controller.DeleteAccount,
	)
}
