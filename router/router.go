package router

import (
	"go-jwt-challenge/controllers"
	"go-jwt-challenge/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	router := gin.Default()

	userRouter := router.Group("/users")
	{
		userRouter.POST("/register", controllers.UserRegister)

		userRouter.POST("/login", controllers.UserLogin)

		userRouter.GET("", controllers.GetAllUsers)
		userRouter.GET("/:id", controllers.GetUserByID)

	}

	productRouter := router.Group("/products")
	{
		productRouter.Use(middlewares.Authentication())

		productRouter.GET("/", controllers.GetAllProducts)
		productRouter.GET("/:productId", middlewares.ProductAuthorization(), controllers.GetProductById)
		productRouter.POST("/", controllers.CreateProduct)
		productRouter.PUT("/:productId", middlewares.ProductAuthorization(), controllers.UpdateProduct)
		productRouter.DELETE("/:productId", middlewares.ProductAuthorization(), controllers.DeleteProduct)

	}

	return router
}
