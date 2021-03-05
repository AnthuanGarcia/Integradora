package utils

import (
	"encoding/json"

	model "github.com/AnthuanGarcia/Integradora/src/models"
)

// SplitNumber - Devuelve un arreglo de digitos de un numero decimal
func SplitNumber(n int) []int {

	if n < 10 {
		return []int{n}
	}

	num := []int{}

	for i := n; i != 0; i /= 10 {
		num = append(num, i%10)
	}

	reverse(num)

	return num
}

func reverse(num []int) {
	i := 0
	j := len(num) - 1

	for i < j {
		num[i], num[j] = num[j], num[i]
		i++
		j--
	}
}

// CommandsToBytes - Convierte varias estructuras a arreglos de Bytes
func CommandsToBytes(tv *model.Tv, num []int) [][]byte {
	var infoComm [][]byte

	infoPower := model.DeviceInfo{
		Capture:  0,
		Protocol: uint16(tv.Protocol),
		Address:  tv.Addr,
		Command:  tv.OnOff,
	}

	byt, _ := json.Marshal(&infoPower)
	infoComm = append(infoComm, byt)

	for _, n := range num {
		infoNumber := model.DeviceInfo{
			Capture:  0,
			Protocol: uint16(tv.Protocol),
			Address:  tv.Addr,
			Command:  tv.Numbers[n],
		}

		b, _ := json.Marshal(&infoNumber)

		infoComm = append(infoComm, b)
	}

	return infoComm
}
