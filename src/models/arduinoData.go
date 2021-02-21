package models

// GetComm - Estructura para mensajes de feedback de Arduino
type GetComm struct {
	Ok int8 `json:"ok"`
}

// DeviceType - Datos del cliente para la captura de dispositivos
type DeviceType struct {
	DevType int8 `json:"devtype"`
	Numbers int8 `json:"numbers"`
}

// DeviceInfo - Datos capturados por el Arduino
type DeviceInfo struct {
	Protocol uint16   `json:"protocol"`
	Address  uint16   `json:"addr"`
	Command  []uint16 `json:"comm"`
}
