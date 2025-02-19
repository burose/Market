package router

import (
	"market/controllers"
	"market/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DElETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/Register", controllers.Register)
	}
	product := r.Group("/api/product")
	product.Use(middlewares.AuthMiddleware())
	{
		product.POST("/create", controllers.Setproduct)
		product.PUT("/update", controllers.Updateproduct)
		product.DELETE("/delete", controllers.Deleteproduct)
		product.GET("/get", controllers.Getproducts)
		product.GET("/get/:id", controllers.GetproductByID)
	}
	cart := r.Group("/api/cart")
	cart.Use(middlewares.AuthMiddleware())
	{
		cart.POST("/create", controllers.Setcart)
		cart.GET("/get", controllers.Getcart)
		cart.PUT("/add/:id", controllers.Addcart)
		cart.DELETE("/clear/:id", controllers.Clearcart)
		cart.DELETE("/delete/:id", controllers.Deletecart)
	}
	order := r.Group("/api/order")
	order.Use(middlewares.AuthMiddleware())
	{
		order.POST("/create", controllers.Createorder)
		order.GET("/get", controllers.Getorder)
		order.PUT("/update/:id", controllers.Updateorder)
		order.DELETE("/cancel/:id", controllers.Cancelorder)
		order.PUT("/pay/:id", controllers.Payorder)
	}

	return r
}
