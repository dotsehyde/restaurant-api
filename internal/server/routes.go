package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	// r.Use(gin.Logger())
	//Home Route '/'
	home := r.Group("/")
	home.GET("/", s.HelloWorldHandler)
	home.GET("/health", s.healthHandler)
	//TODO: Add Authentication Middleware
	//Food Route '/food'
	food := r.Group("/food")
	food.GET("/all", s.GetFoods)
	food.GET("/:id", s.GetFood)
	food.POST("/create", s.CreateFood)
	food.PUT("/update/:id", s.UpdateFood)
	food.DELETE("/delete/:id", s.DeleteFood)

	//Invoice Route '/invoice'
	invoice := r.Group("/invoice")
	invoice.GET("/all", s.GetInvoices)
	invoice.GET("/:id", s.GetInvoice)
	invoice.POST("/create", s.CreateInvoice)
	invoice.PUT("/update/:id", s.UpdateInvoice)
	invoice.DELETE("/delete/:id", s.DeleteInvoice)

	//Menu Route '/menu'
	menu := r.Group("/menu")
	menu.GET("/all", s.GetMenus)
	menu.GET("/:id", s.GetMenu)
	menu.POST("/create", s.CreateMenu)
	menu.PUT("/update/:id", s.UpdateMenu)
	menu.DELETE("/delete/:id", s.DeleteMenu)

	//Order Route '/order'
	order := r.Group("/order")
	order.GET("/all", s.GetOrders)
	order.GET("/:id", s.GetOrder)
	order.POST("/create", s.CreateOrder)
	order.PUT("/update/:id", s.UpdateOrder)
	order.DELETE("/delete/:id", s.DeleteOrder)

	//Table Route '/table'
	table := r.Group("/table")
	table.GET("/all", s.GetTables)
	table.GET("/:id", s.GetTable)
	table.POST("/create", s.CreateTable)
	table.PUT("/update/:id", s.UpdateTable)
	table.DELETE("/delete/:id", s.DeleteTable)

	//OrderItem Route '/orderItem'
	orderItem := r.Group("/orderItem")
	r.GET("/oderItems-order/:id", s.GetOrderItemsByOrder)
	orderItem.GET("/all", s.GetOrderItems)
	orderItem.GET("/:id", s.GetOrderItem)
	orderItem.POST("/create", s.CreateOrderItem)
	orderItem.PUT("/update/:id", s.UpdateOrderItem)
	orderItem.DELETE("/delete/:id", s.DeleteOrderItem)

	//User Route '/user'
	user := r.Group("/user")
	user.GET("/all", s.GetUsers)
	user.GET("/:id", s.GetUser)
	user.POST("/login", s.LoginUser)
	user.POST("/signup", s.SignupUser)
	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
