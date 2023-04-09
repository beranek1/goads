package goads

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestConnection(t *testing.T) {
	c, err := NewConnection("127.0.0.1", "192.168.178.34.1.1:851")
	if err != nil {
		t.Error(err.Error())
	}
	err = c.Open()
	defer c.Close()
	if err != nil {
		t.Error(err.Error())
	}
	var response []byte
	response, err = c.Request(ADSSRVID_READDEVICEINFO, []byte{})
	if err != nil {
		t.Error(err.Error())
	}
	var deviceinfo ADS_Read_Device_Info_Response
	r := bytes.NewReader(response)
	if err := binary.Read(r, binary.LittleEndian, &deviceinfo); err != nil {
		t.Error(err.Error())
	}
}
