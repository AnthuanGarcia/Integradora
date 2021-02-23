package models

/* GetComm - Estructura para mensajes de feedback de Arduino
type GetComm struct {
	Ok int8 `json:"ok"`
}*/

/* DeviceType - Datos del cliente para la captura de dispositivos
type DeviceType struct {
	DevType int8 `json:"devtype"`
	Numbers int8 `json:"numbers"`
}*/

/* Action - Accion a realizar en el arduino,
   1 para capturar señal
   0 para enviar señal
*/
type Action struct {
	Capture uint8 `json:"capture"`
}

// DeviceInfo - Datos capturados por el Arduino
type DeviceInfo struct {
	Protocol uint16 `json:"protocol"`
	Address  uint16 `json:"addr"`
	Command  uint16 `json:"command"`
}
