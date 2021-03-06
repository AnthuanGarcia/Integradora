package utils

import (
	model "github.com/AnthuanGarcia/Integradora/src/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RemoveElement(arr interface{}, i int) interface{} {

	switch arr.(type) {
	case []uint16:
		return append(arr.([]uint16)[:i], arr.([]uint16)[i+1:]...)
	case []model.Tv:
		return append(arr.([]model.Tv)[:i], arr.([]model.Tv)[i+1:]...)
	case primitive.A:
		return append(arr.(primitive.A)[:i], arr.(primitive.A)[i+1:]...)
	}

	return nil
}
