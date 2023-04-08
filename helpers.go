package goads

import (
	"bytes"
	"encoding/binary"
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

func GetDataTypeRecursive(dataTypes map[string]ADS_Data_Type_Entry_Complete, name string) ADS_Data_Type_Entry_Complete {
	oldDataType := dataTypes[name]
	newDataType := ADS_Data_Type_Entry_Complete{Entry: oldDataType.Entry, Name: oldDataType.Name, Type: oldDataType.Type, Comment: oldDataType.Comment, Sub_Items: map[string]ADS_Data_Type_Entry_Complete{}, Array_Info: []ADS_Array_Info{}}
	if oldDataType.Entry.Sub_Items > 0 {
		for k, v := range oldDataType.Sub_Items {
			subDataType := GetDataTypeRecursive(dataTypes, v.Type)
			subDataType.Type = subDataType.Name
			subDataType.Name = k
			subDataType.Entry.Offset = v.Entry.Offset
			subDataType.Comment = v.Comment
			newDataType.Sub_Items[k] = subDataType
		}
	} else if (oldDataType.Type == "" || oldDataType.Entry.Data_Type < ADST_MAXTYPES) && len(oldDataType.Array_Info) == 0 {
	} else if oldDataType.Entry.Data_Type == ADST_BIGTYPE && len(oldDataType.Array_Info) == 0 {
		if oldDataType.Entry.Size == 4 {
			newDataType.Entry.Data_Type = ADST_UINT32
		} else if oldDataType.Entry.Size == 8 {
			newDataType.Entry.Data_Type = ADST_UINT64
		}
	} else if len(oldDataType.Array_Info) > 0 {
		newDataType = GetDataTypeRecursive(dataTypes, oldDataType.Type)
		vec := make([]ADS_Array_Info, oldDataType.Entry.Array_Dimension+newDataType.Entry.Array_Dimension)
		copy(vec, oldDataType.Array_Info)
		for i, v := range newDataType.Array_Info {
			vec[(int(oldDataType.Entry.Array_Dimension) + i)] = v
		}
		newDataType.Array_Info = vec
		newDataType.Entry.Array_Dimension = oldDataType.Entry.Array_Dimension + newDataType.Entry.Array_Dimension
	} else {
		newDataType = GetDataTypeRecursive(dataTypes, oldDataType.Type)
	}
	return newDataType
}

func ReadArrayValue(c *Connection, dataType ADS_Data_Type_Entry_Complete, indexGroup uint32, indexOffset uint32, dim uint16) (value []any, offset uint32) {
	value = make([]any, dataType.Array_Info[dim].Size)
	offset = indexOffset
	for i := 0; i < int(dataType.Array_Info[dim].Size); i++ {
		if (dim + 1) < dataType.Entry.Array_Dimension {
			value[i], offset = ReadArrayValue(c, dataType, indexGroup, offset, dim+1)
		} else {
			value[i], _ = ReadSymbolValue(c, dataType, indexGroup, offset, true)
			offset += dataType.Entry.Size
		}
	}
	return
}

func ReadSymbolValue(c *Connection, dataType ADS_Data_Type_Entry_Complete, indexGroup uint32, indexOffset uint32, aryItem bool) (value any, err error) {
	if (len(dataType.Array_Info) == 0 || aryItem) && dataType.Entry.Sub_Items > 0 {
		dataMap := make(map[string]any, dataType.Entry.Sub_Items)
		offset := indexOffset + dataType.Entry.Offset
		for k, v := range dataType.Sub_Items {
			dataMap[k], err = ReadSymbolValue(c, v, indexGroup, offset, false)
			offset += v.Entry.Size
		}
		value = dataMap
	} else if len(dataType.Array_Info) > 0 && !aryItem {
		value, _ = ReadArrayValue(c, dataType, indexGroup, indexOffset, 0)
	} else {
		if dataType.Entry.Data_Type == ADST_VOID {
			value = nil
			return
		}
		var data []byte
		data, _, err = c.Read(indexGroup, indexOffset, dataType.Entry.Size)
		r := bytes.NewReader(data)
		if err == nil {
			switch dataType.Entry.Data_Type {
			case ADST_BIT:
				var v bool
				if err = binary.Read(r, binary.LittleEndian, &v); err == nil {
					value = v
				}
			case ADST_INT8:
				var v int8
				if err = binary.Read(r, binary.LittleEndian, &v); err == nil {
					value = v
				}
			case ADST_INT16:
				var v int16
				if err = binary.Read(r, binary.LittleEndian, &v); err == nil {
					value = v
				}
			case ADST_INT32:
				var v int32
				if err = binary.Read(r, binary.LittleEndian, &v); err == nil {
					value = v
				}
			case ADST_INT64:
				var v int64
				if err = binary.Read(r, binary.LittleEndian, &v); err == nil {
					value = v
				}
			case ADST_UINT8:
				var v uint8
				if err = binary.Read(r, binary.LittleEndian, &v); err == nil {
					value = v
				}
			case ADST_UINT16:
				var v uint16
				if err = binary.Read(r, binary.LittleEndian, &v); err == nil {
					value = v
				}
			case ADST_UINT32:
				var v uint32
				if err = binary.Read(r, binary.LittleEndian, &v); err == nil {
					value = v
				}
			case ADST_UINT64:
				var v uint64
				if err = binary.Read(r, binary.LittleEndian, &v); err == nil {
					value = v
				}
			case ADST_REAL32:
				var v float32
				if err = binary.Read(r, binary.LittleEndian, &v); err == nil {
					value = v
				}
			case ADST_REAL64:
				var v float64
				if err = binary.Read(r, binary.LittleEndian, &v); err == nil {
					value = v
				}
			case ADST_STRING:
				var v []byte
				if err = binary.Read(r, binary.LittleEndian, &v); err == nil {
					value = string(v)
				}
			default:
				var v []byte
				if err = binary.Read(r, binary.LittleEndian, &v); err == nil {
					value = v
				}
			}
		}
	}
	return
}
