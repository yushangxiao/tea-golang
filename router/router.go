package router

import (
	"tea/handlers"
	"tea/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.JWTAuthMiddleware())
	r.Use(middleware.LiveMiddleware())
	r.Use(middleware.AdminMiddleware())
	api := r.Group("/api")
	{
		api.GET("/ping", handlers.PingHandler)
		api.POST("/signup", handlers.CreateUserHandler)
		api.POST("/login", handlers.LoginHandler)
		api.GET("/userinfo", handlers.UserInfoHandler)
		api.GET(("/products"), handlers.GetProductsHandler)
	}
	admin := r.Group("/admin/api")
	{
		admin.GET("/alluser", handlers.GetAllUsersHandler)
		admin.GET("/deleteuser/:name", handlers.DeleteUserHandler)
		admin.GET("/deleteproduct/:name", handlers.DeleteProductHandler)
		admin.GET("/changeSale/:name", handlers.ChangeSaleHandler)
		admin.POST("/product/create", handlers.CreateProductHandler)
		admin.POST("/admin/del", handlers.DeleteAdminHandler)
	}

	r.NoRoute(handlers.WebHandler)
	return r
}
