package goads

import (
	"errors"
	"strconv"
	"strings"
)

func StringToNetId(str string) (netid [6]byte, err error) {
	split := strings.Split(str, ".")
	if len(split) > 6 {
		return netid, errors.New("invalid input string")
	}

	for i, a := range split {
		var value uint64
		value, err = strconv.ParseUint(a, 10, 8)
		if err != nil {
			return
		}
		netid[i] = byte(value)
	}
	return
}

func StringToAMSAddress(str string) (address AMS_Address, err error) {
	split := strings.Split(str, ":")
	if len(split) > 2 {
		return address, errors.New("invalid input string")
	}
	address.NetId, err = StringToNetId(split[0])
	if err != nil || len(split) == 1 {
		return
	}
	var port uint64
	port, err = strconv.ParseUint(split[1], 10, 16)
	address.Port = uint16(port)
	return
}
