package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "github.com/AnthuanGarcia/Integradora/db"
	arduino "github.com/AnthuanGarcia/Integradora/src/listener"
	model "github.com/AnthuanGarcia/Integradora/src/models"
	scheduler "github.com/AnthuanGarcia/Integradora/src/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/oauth2/v2"
)

var (
	work       = scheduler.NewScheduler()
	httpClient = &http.Client{}
)

//var ultraDevice = new(model.Device)

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
		Devices:   map[string]interface{}{},
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

	c.JSON(http.StatusOK, gin.H{"id": id.Hex()})
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

	c.JSON(http.StatusOK, gin.H{"id": exists.Hex()})
}

// HandleGetUserInfo - Obtiene todos los dispostivos almacenados de un usuario
func HandleGetUserInfo(c *gin.Context) {
	id := strings.ReplaceAll(c.Param("id"), `"`, "")
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

// HandleNewCommand - Captura un nuevo comando
func HandleNewCommand(c *gin.Context) {
	action := new(model.DeviceInfo)

	// Decodificamos el JSON de la peticion, a un struct para manejarlo
	if err := c.BindJSON(action); err != nil {
		log.Printf("Error al deserializar action: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error deserealizar"})
		return
	}

	// El struct generado lo convertimos a un arrglo de bytes
	// para que Arduino pueda leerlo
	devAction, err := json.Marshal(action)

	if err != nil {
		log.Printf("Error al deserializar action(bytes): %v\n", err)
		c.JSON(http.StatusConflict, gin.H{"msg": "Error deserealizar"})
		return
	}

	// Pasamos el arreglo de bytes a Arduino
	deviceData, err := arduino.CaptureCommand(devAction)

	if err != nil {
		log.Printf("Error al capturar datos: %v\n", err)
		c.JSON(http.StatusConflict, gin.H{"msg": "Error deserealizar"})
		return
	}

	// Enviamos una respuesta al cliente de que el proceso fue exitoso
	c.JSON(http.StatusOK, gin.H{"info": deviceData})
}

// HandleSendCommand - Envia un commando en especifico
func HandleSendCommand(c *gin.Context) {
	capture := model.DeviceInfo{}

	if err := c.BindJSON(&capture); err != nil {
		log.Printf("Erro al capturar Json : \n%v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "JSON invalido"})
		return
	}

	deviceInfo, err := json.Marshal(&capture)

	if err != nil {
		log.Printf("Error al deserializar info, en el envio: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error al deserializar, en el envio"})
		return
	}

	arduino.SendCommand(deviceInfo)
}

// HandleNewDevice - Agrega un nuevo dispositvo al usuario
func HandleNewDevice(c *gin.Context) {

	//var device interface{}
	var newDevice interface{}
	var arrDevice interface{}

	//id := strings.ReplaceAll(c.Param("id"), `"`, "")
	typeDev := c.Param("type")

	userData, err := db.GetUserInfo(c.Param("id"))

	if err != nil {
		log.Printf("Error al obtener informacion del usuario: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error al obtener info del usuario"})
		return
	}

	switch typeDev {
	case "Tv": // Tv
		newDevice = model.Tv{}
	case "MediaPlayer":
		newDevice = model.MediaPlayer{}
	}

	if err := c.BindJSON(&newDevice); err != nil {
		log.Printf("Error al deserializar dispositivo: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error deserealizar"})
		return
	}

	_, ok := userData.Devices[typeDev]

	if !ok {

		arrDevice = []interface{}{}
		userData.Devices[typeDev] = arrDevice

		userData.Devices[typeDev] = append(
			userData.Devices[typeDev].([]interface{}),
			newDevice,
		)

	} else {

		userData.Devices[typeDev] = append(
			userData.Devices[typeDev].(primitive.A),
			newDevice,
		)

	}

	if err = db.UpdateUserInfo(userData); err != nil {
		log.Printf("Error al actualizar documento: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error al actualizar doc"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Dispositivo Agregado"})

}

// HandleScheduleDevice - Programa una fecha de encendido para una Television
func HandleScheduleDevice(c *gin.Context) {
	layout := "2006-01-02T15:04:05MST"
	dateOn := model.PowerOnDate{}
	ctx := context.Background()
	numbers := []int{}

	if err := c.BindJSON(&dateOn); err != nil {
		log.Printf("Error al deserializar fecha: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"err": "Error al deserializar"})
		return
	}

	// Si no se quiere programar un canal de encendido
	// El default es 0
	if dateOn.Channel != 0 {
		numbers = scheduler.SplitNumber(int(dateOn.Channel))
	}

	requests := scheduler.CommandsToBytes(&dateOn.Tv, numbers)

	t, err := time.Parse(layout, dateOn.Date)
	if err != nil {
		log.Printf("Fecha Invalida: %v\n", err)
		return
	}

	interval := t.Sub(time.Now())

	work.Add(ctx, arduino.ScheduleCommand, requests, interval)

	c.JSON(http.StatusOK, gin.H{"msg": "Fecha programada con exito"})
}

// HandleFavorite - Cambia segun la lista de canales favoritos
func HandleFavorite(c *gin.Context) {
	favChan := model.Fav{}

	if err := c.BindJSON(&favChan); err != nil {
		log.Printf("Error al deserilizar canal fav: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error al deserializar"})
		return
	}

	num := scheduler.SplitNumber(int(favChan.Channel))
	favComm := scheduler.CommandsToBytesNum(&favChan.Tv, num)

	arduino.SendMultiCommand(favComm)

	c.JSON(http.StatusOK, gin.H{"msg": "Cambio Realizado"})
}

// HandleNewFavorite - Agrega un nuevo canal favorito al usuario
func HandleNewFavorite(c *gin.Context) {

	fav, _ := strconv.Atoi(c.Param("channel"))
	userData, err := db.GetUserInfo(c.Param("id"))

	if err != nil {
		log.Printf("Error al obtener informacion del usuario: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error al obtener info del usuario"})
		return
	}

	userData.Favorites = append(userData.Favorites, uint16(fav))

	if err = db.UpdateUserInfo(userData); err != nil {
		log.Printf("Error al actualizar documento: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error al actualizar doc"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Canal Agregado"})
}
