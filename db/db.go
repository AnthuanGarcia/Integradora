package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	model "github.com/AnthuanGarcia/Integradora/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connTimeout        = 5
	connStringTemplate = "mongodb+srv://%s:%s@devices.pqjzn.mongodb.net/%s?retryWrites=true&w=majority"
	dataBase           = "ControlRemoto"
	collection         = "User"
)

// getConnection - Conexion a MongoDB
func getConnection() (*mongo.Client, context.Context, context.CancelFunc) {
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	clusterEndPoint := os.Getenv("MONGO_ENDPOINT")

	connectionURI := fmt.Sprintf(connStringTemplate, username, password, clusterEndPoint)

	ctx, cancel := context.WithTimeout(context.Background(), connTimeout*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Panicf("Fallo al conectar al cluster %v\n", err)
	}

	/*err = client.Ping(ctx, nil)
	if err != nil {
		log.Panicf("Fallo ping al cluster %v\n", err)
	}*/

	log.Printf("Conectado a MongoDB\n")

	return client, ctx, cancel
}

/*Create - Crea un nuevo documento(Usuario) en la collecion retorna tres valores,
un ID, un error o un entero, si el entero es 0 quiere decir que no se ha repetido
ningun documento si es 1 hay un documento repetido*/
func Create(user *model.User) (primitive.ObjectID, int, error) {
	var exists []bson.M
	client, ctx, cancel := getConnection()

	defer cancel()
	defer client.Disconnect(ctx)

	db := client.Database(dataBase)
	collection := db.Collection(collection)
	result, err := collection.Find(ctx, bson.M{"idgoogle": user.IDGoogle})

	if err != nil {
		log.Printf("Error al buscar el documento: %v\n", err)
		return primitive.NilObjectID, 0, err
	}

	if err := result.All(ctx, &exists); err != nil {
		log.Printf("Error al filtrar el documento: %v\n", err)
		return primitive.NilObjectID, 0, err
	}

	if len(exists) >= 1 {
		log.Printf("Documento existente\n")
		return primitive.NilObjectID, 1, nil
	}

	user.ID = primitive.NewObjectID()

	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("No se ha podido agregar el Usuario: %v\n", err)
		return primitive.NilObjectID, 0, err
	}

	oid := res.InsertedID.(primitive.ObjectID)

	return oid, 0, nil
}

// VerifyUser - Verifica la existencia de un usuario en la base de datos
func VerifyUser(IDGoogle string) (primitive.ObjectID, error) {
	var user model.User

	client, ctx, cancel := getConnection()

	defer cancel()
	defer client.Disconnect(ctx)

	db := client.Database(dataBase)
	collection := db.Collection(collection)
	result := collection.FindOne(ctx, bson.M{"idgoogle": IDGoogle})

	if result == nil {
		log.Printf("Error al buscar el documento\n")
		return primitive.NilObjectID, errors.New("Documento no encontrado")
	}

	if err := result.Decode(&user); err != nil {
		log.Printf("Fallo al decodificar %v\n", err)
		return primitive.NilObjectID, nil
	}

	return user.ID, nil
}

// GetUserInfo - Obtiene todos los datos del usuario
func GetUserInfo(ID string) (*model.User, error) {
	var user model.User

	client, ctx, cancel := getConnection()

	defer cancel()
	defer client.Disconnect(ctx)

	db := client.Database(dataBase)
	collection := db.Collection(collection)

	id, err := primitive.ObjectIDFromHex(ID)

	if err != nil {
		return nil, err
	}

	result := collection.FindOne(ctx, bson.M{"_id": id})

	if result == nil {
		log.Printf("Error al buscar el documento\n")
		return nil, errors.New("Documento no encontrado")
	}

	if err := result.Decode(&user); err != nil {
		log.Printf("Fallo al decodificar %v\n", err)
		return nil, err
	}

	return &user, nil
}

// UpdateUserInfo - Actualiza la informacion del Usuario
func UpdateUserInfo(user *model.User) error {
	client, ctx, cancel := getConnection()

	defer cancel()
	defer client.Disconnect(ctx)

	db := client.Database(dataBase)
	collection := db.Collection(collection)

	res, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)

	if err != nil {
		log.Printf("Fallo al Actualizar documento: %v\n", err)
		return err
	}

	log.Printf("Documento %v actualizado\n", res.ModifiedCount)

	return nil
}
