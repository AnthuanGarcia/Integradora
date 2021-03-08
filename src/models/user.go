package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User - Informacion que se almacenara del usuario
type User struct {
	ID        primitive.ObjectID     `json:"ID" bson:"_id"`
	IDGoogle  string                 `json:"idGoogle"`
	Devices   map[string]interface{} `json:"devices"`
	Favorites []uint16               `json:"favorites"`
}

// Fav - Envio de canales favoritos
type Fav struct {
	Tv
	Channel uint16 `json:"channel"`
}

// PowerOnDate - Fecha y canal para el encendido de un dispositivo
type PowerOnDate struct {
	Tv
	Channel uint16 `json:"channel"`
	Date    string `json:"date"`
}
