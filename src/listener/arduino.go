package listener

import (
	"encoding/json"
	"log"
	"time"

	model "github.com/AnthuanGarcia/Integradora/src/models"
	"github.com/tarm/serial"
)

const (
	serialPort = "COM3"
	bps        = 115200
)

var (
	ser = new(serial.Port)
	//plsContinue = true
)

// PreparePort - Abre el puerto serial seleccionado
func preparePort(chanPort chan *serial.Port) {

	c := &serial.Config{
		Name:        serialPort,
		Baud:        bps,
		ReadTimeout: time.Millisecond * 5000,
	}

	s, err := serial.OpenPort(c)

	if err != nil {
		log.Fatal(err)
	}

	//plsContinue = true
	log.Println("Abriendo Puerto")
	time.Sleep(1600 * time.Millisecond)

	chanPort <- s
}

// Listen - Escucha el puerto seleccionado para recibir bytes
func listen() []byte {

	buffer := make([]byte, 128)

	for {
		data, _ := ser.Read(buffer)

		log.Println(string(buffer[:data]))

		if json.Valid(buffer[:data]) {
			return buffer[:data]
		} else {
			write([]byte(`{"capture":1}`))
		}
	}

}

// Write - Escribe en el puerto seleccionado un conjunto de bytes
func write(request []byte) {

	n, err := ser.Write(request)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(n)

}

// CaptureCommand - Captura los datos especificados
func CaptureCommand(action []byte) (*model.DeviceInfo, error) {
	chanPort := make(chan *serial.Port)
	//chanInfo := make(chan []byte)

	go preparePort(chanPort)

	//succesState := model.GetComm{}
	infodevice := new(model.DeviceInfo)

	ser = <-chanPort

	defer ser.Close()

	write(action)
	//go listen(chanInfo)

	for infodevice.Command == 0 {

		info := listen()

		if err := json.Unmarshal(info, infodevice); err != nil {

			log.Println(err)
			return nil, err

		} else {

			log.Println(infodevice)

		}

	}

	//plsContinue = false

	return infodevice, nil
}
