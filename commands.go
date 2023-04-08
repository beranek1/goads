package goads

import (
	"bytes"
	"encoding/binary"
)

func (c *Connection) ReadDeviceInfo() (res ADS_Read_Device_Info_Response, err error) {
	var response []byte
	response, err = c.Request(ADSSRVID_READDEVICEINFO, []byte{})
	if err == nil {
		r := bytes.NewReader(response)
		err = binary.Read(r, binary.LittleEndian, &res)
	}
	return
}

func (c *Connection) Read(indexGroup uint32, indexOffset uint32, length uint32) (data []byte, res ADS_Read_Response, err error) {
	req := &ADS_Read_Request{Index_Group: indexGroup, Index_Offset: indexOffset, Length: length}
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, req)
	var response []byte
	response, err = c.Request(ADSSRVID_READ, buffer.Bytes())
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

func (c *Connection) Write(indexGroup uint32, indexOffset uint32, data []byte) (res ADS_Write_Response, err error) {
	req := &ADS_Write_Request{Index_Group: indexGroup, Index_Offset: indexOffset, Length: uint32(len(data))}
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, req)
	binary.Write(buffer, binary.LittleEndian, data)
	var response []byte
	response, err = c.Request(ADSSRVID_WRITE, buffer.Bytes())
	if err == nil {
		r := bytes.NewReader(response)
		err = binary.Read(r, binary.LittleEndian, &res)
	}
	return
}

func (c *Connection) ReadState() (res ADS_Read_State_Response, err error) {
	var response []byte
	response, err = c.Request(ADSSRVID_READSTATE, []byte{})
	if err == nil {
		r := bytes.NewReader(response)
		err = binary.Read(r, binary.LittleEndian, &res)
	}
	return
}

func (c *Connection) WriteControl(adsState uint16, deviceState uint16) (res ADS_Write_Control_Response, err error) {
	req := &ADS_Write_Control_Request{ADS_State: adsState, Device_State: deviceState, Length: 0}
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, req)
	var response []byte
	response, err = c.Request(ADSSRVID_WRITECTRL, buffer.Bytes())
	if err == nil {
		r := bytes.NewReader(response)
		err = binary.Read(r, binary.LittleEndian, &res)
	}
	return
}

func (c *Connection) ReadWrite(indexGroup uint32, indexOffset uint32, length uint32, input []byte) (output []byte, res ADS_Read_Write_Response, err error) {
	req := &ADS_Read_Write_Request{Index_Group: indexGroup, Index_Offset: indexOffset, Read_Length: length, Write_Length: uint32(len(input))}
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, req)
	binary.Write(buffer, binary.LittleEndian, input)
	var response []byte
	response, err = c.Request(ADSSRVID_WRITE, buffer.Bytes())
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
