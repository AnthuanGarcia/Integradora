package listener

import (
	"encoding/json"
	"log"

	"github.com/tarm/serial"
)

func listenArduino(serialPort string, bps int, canal chan int) {
	c := &serial.Config{Name: serialPort, Baud: bps}
	s, err := serial.OpenPort(c)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Escuchando Arduino...\n")

	buffer := make([]byte, 128)

	for {
		data, err := s.Read(buffer)

		if err != nil {
			log.Fatal()
		}

		if json.Valid(buffer[:data]) {

		}
	}
}
