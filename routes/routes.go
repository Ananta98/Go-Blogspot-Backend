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
	UserRoute.DELETE("/:id", controllers.DeleteUser)

	CategoryRoute := r.Group("/category")
	CategoryRoute.Use(middlewares.JwtAuthMiddleware())
	CategoryRoute.POST("/", controllers.CreateNewCategory)
	CategoryRoute.PATCH("/:id", controllers.UpdateCategory)
	CategoryRoute.DELETE("/:id", controllers.DeleteCategory)
	CategoryRoute.GET("/", controllers.GetListCategories)

	PostRoute := r.Group("/post")
	PostRoute.Use(middlewares.JwtAuthMiddleware())

	//posts api section
	PostRoute.POST("/", controllers.CreateNewPost)
	PostRoute.GET("/", controllers.GetListBlogs)
	PostRoute.DELETE("/:id", controllers.DeletePost)
	PostRoute.PATCH("/:id", controllers.UpdatePost)
	PostRoute.GET("/:id", controllers.GetDetailPost)

	//comments api section
	PostRoute.POST("/comment", controllers.CreateNewComment)
	PostRoute.PATCH("/:id/comment/:comment_id", controllers.UpdateComment)
	PostRoute.DELETE("/:id/comment/:comment_id", controllers.DeleteComment)
	PostRoute.GET("/:id/comment", controllers.GetListComments)

	//user like post api section
	PostRoute.POST("/:id/like/:status", controllers.LikePostController)
	PostRoute.GET("/:id/user-likes/", controllers.GetListUserLikePost)
	PostRoute.GET("/:id/user-dislikes/", controllers.GetListUserDislikePost)

	//user like comment post api section
	PostRoute.POST("/comment/:id/like/:status", controllers.LikeCommentController)
	PostRoute.GET("/comment/:id/user-likes/", controllers.GetListUserLikeComment)
	PostRoute.GET("/comment/:id/user-dislikes/", controllers.GetListUserDislikeComment)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
