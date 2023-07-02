package goads

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Connection struct {
	ip         string
	target     AMS_Address
	out        chan []byte
	done       chan bool
	in         map[uint32]chan []byte
	invoke_id  uint32
	mutex      *sync.Mutex
	inmutex    sync.RWMutex
	connection net.Conn
	source     AMS_Address
}

func NewConnection(ip string, target string) (c *Connection, err error) {
	var target_addr AMS_Address
	target_addr, err = StringToAMSAddress(target)
	if err != nil {
		return
	}
	c = &Connection{ip: ip, target: target_addr, out: make(chan []byte), done: make(chan bool), in: make(map[uint32]chan []byte), invoke_id: 0, mutex: &sync.Mutex{}}
	return
}

func (c *Connection) Open() (err error) {
	c.connection, err = net.Dial("tcp", fmt.Sprintf("%s:48898", c.ip))
	if err != nil {
		return
	}
	var data struct {
		Header AMSTCP_Header
		Data   uint16
	}
	data.Header = AMSTCP_Header{Reserved: 4096, Length: 2}
	data.Data = 0
	packet_buffer := new(bytes.Buffer)
	binary.Write(packet_buffer, binary.LittleEndian, data)
	_, err = c.connection.Write(packet_buffer.Bytes())
	if err == nil {
		res := make([]byte, 6)
		_, err = c.connection.Read(res)
		if err == nil {
			r := bytes.NewReader(res)
			var header AMSTCP_Header
			if err = binary.Read(r, binary.LittleEndian, &header); err == nil {
				res := make([]byte, header.Length)
				_, err = c.connection.Read(res)
				if err == nil {
					r := bytes.NewReader(res)
					var address AMS_Address
					if err = binary.Read(r, binary.LittleEndian, &address); err == nil {
						c.source = address
						c.send()
						c.receive()
					}
				}
			}
		}
	}
	return
}

func (c *Connection) Close() {
	close(c.done)
}

func (c *Connection) Request(command ADSSRVID, data []byte) (response []byte, err error) {
	if c == nil {
		err = errors.New("no connection established")
		return
	}
	id := c.NextInvokeID()

	respCh := make(chan []byte)
	defer close(respCh)

	c.inmutex.Lock()
	c.in[id] = respCh
	c.inmutex.Unlock()

	defer func() {
		c.inmutex.Lock()
		delete(c.in, id)
		c.inmutex.Unlock()
	}()

	var packet struct {
		AMSTCP AMSTCP_Header
		AMS    AMS_Header
	}
	packet.AMS = AMS_Header{Target: c.target, Source: c.source, Command_Id: uint16(command), State_Flags: 4, Data_Length: uint32(len(data)), Invoke_Id: id}
	packet.AMSTCP = AMSTCP_Header{Length: (32 + packet.AMS.Data_Length)}
	packet_buffer := new(bytes.Buffer)
	binary.Write(packet_buffer, binary.LittleEndian, packet)
	binary.Write(packet_buffer, binary.LittleEndian, data)
	select {
	case c.out <- packet_buffer.Bytes():
	case <-time.After(time.Second * 5):
		return response, errors.New("timeout send")
	case <-c.done:
		return response, errors.New("connection closed")
	}
	select {
	case response = <-respCh:
		return
	case <-time.After(time.Second * 5):
		return response, errors.New("timeout receive")
	case <-c.done:
		return response, errors.New("connection closed")
	}
}

func (c *Connection) NextInvokeID() (id uint32) {
	c.mutex.Lock()
	c.invoke_id++
	id = c.invoke_id
	c.mutex.Unlock()
	return
}

// Adapted from https://github.com/stamp/goADS/blob/master/main.go
func (c *Connection) listen() <-chan []byte {
	buffer := make(chan []byte)
	go func(c *Connection) {
		b := make([]byte, 1024)
		for {
			n, err := c.connection.Read(b)
			if n > 0 {
				res := make([]byte, n)
				copy(res, b[:n])
				buffer <- res
			}
			if err == io.EOF {
				break
			} else if err != nil {
				buffer <- nil
				return
			}
		}
	}(c)
	return buffer
}

func (c *Connection) receive() error {
	if c == nil {
		return errors.New("no connection established")
	}

	// Adapted from https://github.com/stamp/goADS/blob/master/main.go
	go func(c *Connection) {
		var buffer bytes.Buffer
		read := c.listen()
		for {
			select {
			case data := <-read:
				if data == nil {
					return
				}
				buffer.Write(data)
				for buffer.Len() >= 38 {
					var data struct {
						AMSTCP AMSTCP_Header
						AMS    AMS_Header
					}
					header := make([]byte, 38)
					buffer.Read(header)
					r := bytes.NewReader(header)
					if err := binary.Read(r, binary.LittleEndian, &data); err != nil {
						continue
					}
					packet := make([]byte, data.AMS.Data_Length)
					n, _ := buffer.Read(packet)
					if n != int(data.AMS.Data_Length) {
						buffer.Write(header)
						buffer.Write(packet[:n])
						break
					}
					c.inmutex.RLock()
					if respCh, test := c.in[data.AMS.Invoke_Id]; test {
						c.inmutex.RUnlock()
						respCh <- packet
					} else {
						c.inmutex.RUnlock()
					}
				}
			case <-c.done:
				return
			}
		}
	}(c)
	return nil
}

func (c *Connection) send() error {
	if c == nil {
		return errors.New("no connection established")
	}

	go func(c *Connection) {
		for {
			select {
			case data := <-c.out:
				go c.connection.Write(data)
			case <-c.done:
				return
			}
		}
	}(c)
	return nil
}
