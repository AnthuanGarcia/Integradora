package user

import (
	"log"
	"net/http"

	db "github.com/AnthuanGarcia/Integradora/db"
	model "github.com/AnthuanGarcia/Integradora/src/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/oauth2/v2"
)

var httpClient = &http.Client{}

type userInfo struct {
	IDToken string `json:"userinfo"`
}

func verifyIDToken(idToken string) (*oauth2.Tokeninfo, error) {
	oauth2Serv, err := oauth2.New(httpClient)
	tokenInfoCall := oauth2Serv.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()

	if err != nil {
		return nil, err
	}

	return tokenInfo, nil
}

// HandleCreateUser - Agrega un nuevo usuario a traves de la Api
func HandleCreateUser(c *gin.Context) {
	user := new(userInfo)

	if err := c.BindJSON(user); err != nil {
		log.Printf("Campo no encontrado %v\n", err)
		return
	}

	userinfo, err := verifyIDToken(user.IDToken)

	if err != nil {
		log.Printf("Error al validar token %v\n", err)
		return
	}

	newUser := model.User{
		ID:        primitive.NilObjectID,
		IDGoogle:  userinfo.UserId,
		Devices:   []model.Device{},
		Favorites: []uint16{},
	}

	id, rep, err := db.Create(&newUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}

	if rep == 1 {
		c.JSON(http.StatusConflict, gin.H{"err": "Usuario Repetido"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// HandleSignInUser - Acceso de un usuario a traves de su cuenta de Google
func HandleSignInUser(c *gin.Context) {
	user := new(userInfo)

	if err := c.BindJSON(user); err != nil {
		log.Printf("Campo no encontrado %v\n", err)
		return
	}

	userinfo, err := verifyIDToken(user.IDToken)

	if err != nil {
		log.Printf("Error al validar token %v\n", err)
		return
	}

	exists, err := db.VerifyUser(userinfo.UserId)

	if err != nil {
		log.Printf("Error al verificar el usuario %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": exists})
}
