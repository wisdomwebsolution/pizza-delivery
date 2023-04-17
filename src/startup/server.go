package startup

import (
	gin "github.com/gin-gonic/gin"
	"github.com/marijakljestan/golang-web-app/src/api"
	repository "github.com/marijakljestan/golang-web-app/src/domain/repository"
	service "github.com/marijakljestan/golang-web-app/src/domain/service"
	"github.com/marijakljestan/golang-web-app/src/infrastructure/persistence"
	"github.com/marijakljestan/golang-web-app/src/middleware"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) Start() {
	pizzaRepository := server.initPizzaRepository()
	pizzaService := server.initPizzaService(pizzaRepository)
	pizzaHandler := api.NewPizzaController(pizzaService)

	orderRepository := server.initOrderRepository()
	orderService := server.initOrderService(orderRepository, pizzaService)
	orderHandler := api.NewOrderController(orderService)

	userRepository := server.initUserRepository()
	userService := server.initUserService(userRepository)
	userHandler := api.NewUserController(userService)

	router := gin.Default()

	pizzaRoutes := router.Group("/pizza")
	{
		pizzaRoutes.GET("", pizzaHandler.GetMenu)
		pizzaRoutes.POST("", middleware.AuthorizeJWT("ADMIN"), pizzaHandler.AddPizzaToMenu)
		pizzaRoutes.DELETE("/:name", middleware.AuthorizeJWT("ADMIN"), pizzaHandler.DeletePizzaFromMenu)
	}

	orderRoutes := router.Group("/order")
	{
		orderRoutes.POST("", orderHandler.CreateOrder)
		orderRoutes.GET("/status/:id", orderHandler.CheckOrderStatus)
		orderRoutes.PUT("/cancel/:id", orderHandler.CancelOrder)
		orderRoutes.PUT("/:id", middleware.AuthorizeJWT("ADMIN"), orderHandler.CancelOrderRegardlessStatus)
	}

	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/register", userHandler.RegisterUser)
		userRoutes.POST("/login", userHandler.Login)
	}

	router.Run("localhost:8080")
}

func (server *Server) initPizzaRepository() repository.PizzaRepository {
	return persistence.NewOrderInMemoryRepository()
}

func (server *Server) initPizzaService(orderRepository repository.PizzaRepository) *service.PizzaService {
	return service.NewPizzaService(orderRepository)
}

func (server *Server) initOrderRepository() repository.OrderRepository {
	return persistence.NewOrderInmemoryRepository()
}

func (server *Server) initOrderService(orderRepository repository.OrderRepository, pizzaService *service.PizzaService) *service.OrderService {
	return service.NewOrderService(orderRepository, pizzaService)
}

func (server *Server) initUserRepository() repository.UserRepository {
	return persistence.NewUserInmemoryRepository()
}

func (server *Server) initUserService(userRepository repository.UserRepository) *service.UserService {
	return service.NewUserService(userRepository)
}
