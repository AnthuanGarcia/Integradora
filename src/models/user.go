package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User - Informacion que se almacenara del usuario
type User struct {
	ID        primitive.ObjectID     `json:"ID"`
	IDGoogle  string                 `json:"idGoogle"`
	Devices   map[string]interface{} `json:"devices"`
	Favorites []uint16               `json:"favorites"`
}
