package models

// Tv - Informacion de Televisiones
type Tv struct {
	Device  `json:"device"`
	OnOff   uint16   `json:"onoff"`
	VolUp   uint16   `json:"volUp"`
	VolDown uint16   `json:"volDown"`
	ChaUp   uint16   `json:"chaUp"`
	ChaDown uint16   `json:"chaDown"`
	Numbers []uint16 `json:"numbers"`
}

// MediaPlayer - Informacion de Televisiones
type MediaPlayer struct {
	Device  `json:"device"`
	OnOff   uint16 `json:"onoff"`
	VolUp   uint16 `json:"volUp"`
	VolDown uint16 `json:"volDown"`
	Play    uint16 `json:"play"`
	Stop    uint16 `json:"stop"`
	SkipR   uint16 `json:"skipr"`
	SkipL   uint16 `json:"skipl"`
	Eject   uint16 `json:"eject"`
}

// Device - informacion del dispositivo
type Device struct {
	Name     string `json:"name"`
	Protocol uint8  `json:"protocol"`
	Addr     uint16 `json:"addres"`
	//Tv          `json:"tv"`
	//MediaPlayer `json:"mediaplayer"`
}
