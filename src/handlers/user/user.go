package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	db "github.com/AnthuanGarcia/Integradora/db"
	arduino "github.com/AnthuanGarcia/Integradora/src/listener"
	model "github.com/AnthuanGarcia/Integradora/src/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/oauth2/v2"
)

var httpClient = &http.Client{}

type userInfo struct {
	IDToken string `json:"idtoken"`
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
		log.Printf("Error al crear Usuario: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}

	if rep == 1 {
		c.JSON(http.StatusConflict, gin.H{"err": "Usuario Repetido"})
		return
	}

	idstr, err := id.MarshalJSON()

	if err != nil {
		log.Printf("Error al codificar el id a JSON %v\n", err)
		c.JSON(http.StatusConflict, gin.H{"msg": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": string(idstr)})
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

	idstr, err := exists.MarshalJSON()

	if err != nil {
		log.Printf("Error al codificar el id a JSON %v\n", err)
		c.JSON(http.StatusConflict, gin.H{"msg": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": string(idstr)})
}

// HandleGetUserInfo - Obtiene todos los dispostivos almacenados de un usuario
func HandleGetUserInfo(c *gin.Context) {
	id := strings.Replace(c.Param("id"), `"`, "", -1)
	data, err := db.GetUserInfo(id)

	if err != nil {
		log.Printf("Error al cargar dispositivos: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error al cargar dispositivos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"devices":   data.Devices,
		"favorites": data.Favorites,
	})
}

// Feedback - Mensajes de retroalimentacion para el cliente
func Feedback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": 1})
}

// HandleNewDevice - Agrega un nuevo dispositvo al usuario
func HandleNewDevice(c *gin.Context) {

	var device interface{}

	id := strings.Replace(c.Param("id"), `"`, "", -1)
	name := strings.Replace(c.Param("name"), `"`, "", -1)

	reqType := new(model.DeviceType) // Este se tiene que recibir desde el cliente

	if err := c.BindJSON(reqType); err != nil {
		log.Printf("Error al deserializar: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error deserealizar"})
		return
	}

	bytesType, err := json.Marshal(&reqType)

	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	deviceData, err := arduino.CaptureCommands(bytesType)

	if err != nil {
		log.Printf("Error al capturar dispositivo: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error al capturar dispositivo"})
		return
	}

	n := len(deviceData.Command) - 1

	// Cada Numero corresponde a un dispositivo
	switch reqType.DevType {
	case 1: // 1 Para Tv
		device = model.Tv{
			OnOff:   deviceData.Command[n],
			VolUp:   deviceData.Command[n-1],
			VolDown: deviceData.Command[n-2],
			ChaUp:   deviceData.Command[n-3],
			ChaDown: deviceData.Command[n-4],
			Numbers: deviceData.Command[:n-4],
		}
		// 2 para Reproductor multimedia
		// 3 Aires acondicionados
		// etc...
	}

	userData, err := db.GetUserInfo(id)

	if err != nil {
		log.Printf("Error al obtener informacion del usuario: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error al obtener info del usuario"})
		return
	}

	userData.Devices = append(userData.Devices, model.Device{
		Name:     name,
		Protocol: uint8(deviceData.Protocol),
		Addr:     deviceData.Address,
		Tv:       device.(model.Tv),
	})

	if err = db.UpdateUserInfo(userData); err != nil {
		log.Printf("Error al actualizar documento: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error al actualizar doc"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Dispositivo Agregado"})

}
