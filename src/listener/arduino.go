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
	ser         = new(serial.Port)
	plsContinue = true
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

	plsContinue = true
	log.Println("Abriendo Puerto")
	time.Sleep(1600 * time.Millisecond)

	chanPort <- s
}

// Listen - Escucha el puerto seleccionado para recibir bytes
func listen(chanInfo chan []byte) {

	buffer := make([]byte, 128)

	for plsContinue {
		data, _ := ser.Read(buffer)

		if json.Valid(buffer[:data]) {
			chanInfo <- buffer[:data]
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

// CaptureCommands - Captura los datos especificados
func CaptureCommands(devtype []byte) (*model.DeviceInfo, error) {
	chanPort := make(chan *serial.Port)
	chanInfo := make(chan []byte)

	go preparePort(chanPort)

	succesState := model.GetComm{}
	infodevice := new(model.DeviceInfo)

	ser = <-chanPort

	defer ser.Close()

	write(devtype)
	go listen(chanInfo)

	for len(infodevice.Command) == 0 {

		info := <-chanInfo

		if err := json.Unmarshal(info, &succesState); err != nil {
			log.Println(err)
		} else {
			log.Println(succesState)
		}

		if err := json.Unmarshal(info, infodevice); err != nil {

			log.Println(err)
			return nil, err

		} else {

			log.Println(infodevice)

		}

	}

	plsContinue = false

	return infodevice, nil
}
