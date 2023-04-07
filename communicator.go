package goads

import (
	"bytes"
	"encoding/binary"
)

type Communicator struct {
	connection *Connection
}

func CreateCommunicator(ip string, target string) (c *Communicator, err error) {
	c.connection, err = CreateConnection(ip, target)
	if err == nil {
		err = c.connection.Open()
	}
	return
}

func (c *Communicator) ReadDeviceInfo() (res ADS_Read_Device_Info_Response, err error) {
	var response []byte
	response, err = c.connection.Request(ADSSRVID_READDEVICEINFO, []byte{})
	if err == nil {
		r := bytes.NewReader(response)
		err = binary.Read(r, binary.LittleEndian, &res)
	}
	return
}

func (c *Communicator) Read(indexGroup uint32, indexOffset uint32, length uint32) (data []byte, res ADS_Read_Response, err error) {
	req := &ADS_Read_Request{Index_Group: indexGroup, Index_Offset: indexOffset, Length: length}
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, req)
	var response []byte
	response, err = c.connection.Request(ADSSRVID_READ, buffer.Bytes())
	if err == nil {
		r := bytes.NewReader(response)
		err = binary.Read(r, binary.LittleEndian, &res)
		if err == nil {
			data = make([]byte, res.Length)
			err = binary.Read(r, binary.LittleEndian, &data)
		}
	}
	return
}

func (c *Communicator) Write(indexGroup uint32, indexOffset uint32, data []byte) (res ADS_Write_Response, err error) {
	req := &ADS_Write_Request{Index_Group: indexGroup, Index_Offset: indexOffset, Length: uint32(len(data))}
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, req)
	binary.Write(buffer, binary.LittleEndian, data)
	var response []byte
	response, err = c.connection.Request(ADSSRVID_WRITE, buffer.Bytes())
	if err == nil {
		r := bytes.NewReader(response)
		err = binary.Read(r, binary.LittleEndian, &res)
	}
	return
}

func (c *Communicator) ReadState() (res ADS_Read_State_Response, err error) {
	var response []byte
	response, err = c.connection.Request(ADSSRVID_READSTATE, []byte{})
	if err == nil {
		r := bytes.NewReader(response)
		err = binary.Read(r, binary.LittleEndian, &res)
	}
	return
}

func (c *Communicator) WriteControl(adsState uint16, deviceState uint16) (res ADS_Write_Control_Response, err error) {
	req := &ADS_Write_Control_Request{ADS_State: adsState, Device_State: deviceState, Length: 0}
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, req)
	var response []byte
	response, err = c.connection.Request(ADSSRVID_WRITECTRL, buffer.Bytes())
	if err == nil {
		r := bytes.NewReader(response)
		err = binary.Read(r, binary.LittleEndian, &res)
	}
	return
}

func (c *Communicator) ReadWrite(indexGroup uint32, indexOffset uint32, length uint32, input []byte) (output []byte, res ADS_Read_Write_Response, err error) {
	req := &ADS_Read_Write_Request{Index_Group: indexGroup, Index_Offset: indexOffset, Read_Length: length, Write_Length: uint32(len(input))}
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, req)
	binary.Write(buffer, binary.LittleEndian, input)
	var response []byte
	response, err = c.connection.Request(ADSSRVID_WRITE, buffer.Bytes())
	if err == nil {
		r := bytes.NewReader(response)
		err = binary.Read(r, binary.LittleEndian, &res)
		if err == nil {
			output = make([]byte, res.Length)
			err = binary.Read(r, binary.LittleEndian, &output)
		}
	}
	return
}
