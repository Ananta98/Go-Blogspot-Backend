package routes

import (
	"blogspot-project/controllers"
	"blogspot-project/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	AuthRoute := r.Group("/auth")
	AuthRoute.POST("/login", controllers.LoginUser)
	AuthRoute.POST("/register", controllers.RegisterNewUser)

	updateUserMiddlewareRoute := r.Group("/login")
	updateUserMiddlewareRoute.Use(middlewares.JwtAuthMiddleware())
	updateUserMiddlewareRoute.PATCH("/update-password", controllers.UpdatePassword)
	updateUserMiddlewareRoute.PATCH("/update-current-user", controllers.UpdateCurrentUser)

	UserRoute := r.Group("/user")
	UserRoute.Use(middlewares.JwtAuthMiddleware())
	UserRoute.GET("/", controllers.GetListUsers)
	UserRoute.GET("/profile", controllers.GetCurrentUserProfile)

	CategoryRoute := r.Group("/category")
	CategoryRoute.POST("/", controllers.CreateNewCategory)
	CategoryRoute.PATCH("/:id", controllers.UpdateCategory)
	CategoryRoute.DELETE("/:id", controllers.DeleteCategory)
	CategoryRoute.GET("/", controllers.GetListCategories)

	PostRoute := r.Group("/post")
	PostRoute.POST("/", controllers.CreateNewPost)
	PostRoute.GET("/", controllers.GetListBlogs)
	PostRoute.DELETE("/:id", controllers.DeletePost)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
