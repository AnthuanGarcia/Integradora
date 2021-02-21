package models

// Tv - Informacion de Televisiones (Opcion 1)
type Tv struct {
	/*name     string
	protocol uint8
	addr     uint16
	onOff    uint16*/
	OnOff   uint16   `json:"onoff"`
	VolUp   uint16   `json:"volUp"`
	VolDown uint16   `json:"volDown"`
	ChaUp   uint16   `json:"chaUp"`
	ChaDown uint16   `json:"chaDown"`
	Numbers []uint16 `json:"numbers"`
}

// Device - informacion del dispositivo
type Device struct {
	Name     string `json:"name"`
	Protocol uint8  `json:"protocol"`
	Addr     uint16 `json:"addr"`
	Tv       `json:"tv"`
}
