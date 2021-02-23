package routes

import (
	handleUser "github.com/AnthuanGarcia/Integradora/src/handlers/user"
	"github.com/gin-gonic/gin"
)

// Routes - Estructura para generar los Endpoints de la API
type Routes struct{}

// StartGin - Inicio del servidor
func (c Routes) StartGin() {
	r := gin.Default()

	api := r.Group("/inte")
	{
		api.POST("/login", handleUser.HandleCreateUser)
		api.POST("/sigin", handleUser.HandleSignInUser)
		api.GET("/user/:id", handleUser.HandleGetUserInfo)
		api.POST("/newCommand", handleUser.HandleNewCommand)
		api.POST("/user/:id/newDevice/:type", handleUser.HandleNewDevice)
	}

	r.Run("0.0.0.0:3000")
}
