package listener

import (
	"encoding/json"
	"log"

	model "github.com/AnthuanGarcia/Integradora/src/models"
	puerto "github.com/jacobsa/go-serial/serial"
)

const (
	serialPortName = "COM3"
	bps            = 115200
)

var (
	options = puerto.OpenOptions{
		PortName:        serialPortName,
		BaudRate:        bps,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	ser, _ = puerto.Open(options)
)

func reopen() {
	ser, _ = puerto.Open(options)
}

// Listen - Escucha el puerto seleccionado para recibir bytes
func listen(request []byte) []byte {

	buffer := make([]byte, 64)

	for {
		data, _ := ser.Read(buffer)

		log.Println(string(buffer[:data]))

		if json.Valid(buffer[:data]) {
			return buffer[:data]
		}

		write(request)

	}
}

// Write - Escribe en el puerto seleccionado un conjunto de bytes
func write(request []byte) {

	//log.Println(string(request))
	n, err := ser.Write(request)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(n)

}

// CaptureCommand - Captura los datos especificados
func CaptureCommand(action []byte) (*model.DeviceInfo, error) {

	ser.Close()
	defer reopen()

	port, err := puerto.Open(options)
	if err != nil {
		log.Panic(err)
	}

	buffer := make([]byte, 64)
	infodevice := new(model.DeviceInfo)

	port.Write(action)

	for infodevice.Command == infodevice.Address {

		info, _ := port.Read(buffer)

		if !json.Valid(buffer[:info]) {
			port.Write(action)
			continue
		}

		if err := json.Unmarshal(buffer[:info], infodevice); err != nil {

			log.Println(err)
			return nil, err

		}

		log.Println(infodevice)

	}

	port.Close()

	return infodevice, nil
}

// SendCommand - Envia un comando en especifico al arduino
func SendCommand(request []byte) {

	write(request)

}
