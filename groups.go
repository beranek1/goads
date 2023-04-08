package goads

import (
	"bytes"
	"encoding/binary"
	"io"
)

func (c *Connection) GetAdsSymbolUploadInfo() (info ADS_Symbol_Upload_Info, err error) {
	var data []byte
	data, _, err = c.Read(ADSIGRP_SYM_UPLOADINFO, 0, 24)
	if err == nil {
		r := bytes.NewReader(data)
		err = binary.Read(r, binary.LittleEndian, &info)
	}
	return
}

func (c *Connection) GetAdsSymbolUploadInfo2() (info ADS_Symbol_Upload_Info2, err error) {
	var data []byte
	data, _, err = c.Read(ADSIGRP_SYM_UPLOADINFO2, 0, 24)
	if err == nil {
		r := bytes.NewReader(data)
		err = binary.Read(r, binary.LittleEndian, &info)
	}
	return
}

func (c *Connection) GetAdsSymbolUpload(info ADS_Symbol_Upload_Info2) (symbols map[string]ADS_Symbol_Entry_Complete, err error) {
	var data []byte
	data, _, err = c.Read(ADSIGRP_SYM_UPLOAD, 0, info.Symbol_Size)
	if err == nil {
		symbols = make(map[string]ADS_Symbol_Entry_Complete)
		r := bytes.NewReader(data)
		for i := 0; i < int(info.Symbols); i++ {
			var symbol ADS_Symbol_Entry
			err = binary.Read(r, binary.LittleEndian, &symbol)
			if err == nil {
				symdata := make([]byte, symbol.Entry_Length-30)
				_, err = io.ReadFull(r, symdata)
				if err == nil {
					r := bytes.NewReader(symdata)
					Name := make([]byte, symbol.Name_Length)
					err = binary.Read(r, binary.LittleEndian, &Name)
					if err == nil {
						var str_end byte
						err = binary.Read(r, binary.LittleEndian, &str_end)
						if err == nil {
							Type := make([]byte, symbol.Type_Length)
							err = binary.Read(r, binary.LittleEndian, &Type)
							if err == nil {
								err = binary.Read(r, binary.LittleEndian, &str_end)
								if err == nil {
									Comment := make([]byte, symbol.Comment_Length)
									err = binary.Read(r, binary.LittleEndian, &Comment)
									if err == nil {
										symbols[string(Name)] = ADS_Symbol_Entry_Complete{Entry: symbol, Name: string(Name), Type: string(Type), Comment: string(Comment)}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

func extractDataTypes(r io.Reader, data_types int) (dataTypes map[string]ADS_Data_Type_Entry_Complete, err error) {
	dataTypes = make(map[string]ADS_Data_Type_Entry_Complete, data_types)
	for i := 0; i < data_types; i++ {
		var dataType ADS_Data_Type_Entry
		err = binary.Read(r, binary.LittleEndian, &dataType)
		if err == nil {
			symdata := make([]byte, dataType.Entry_Length-42)
			_, err = io.ReadFull(r, symdata)
			if err == nil {
				r := bytes.NewReader(symdata)
				Name := make([]byte, dataType.Name_Length)
				err = binary.Read(r, binary.LittleEndian, &Name)
				if err == nil {
					var str_end byte
					err = binary.Read(r, binary.LittleEndian, &str_end)
					if err == nil {
						Type := make([]byte, dataType.Type_Length)
						err = binary.Read(r, binary.LittleEndian, &Type)
						if err == nil {
							err = binary.Read(r, binary.LittleEndian, &str_end)
							if err == nil {
								Comment := make([]byte, dataType.Comment_Length)
								err = binary.Read(r, binary.LittleEndian, &Comment)
								if err == nil {
									err = binary.Read(r, binary.LittleEndian, &str_end)
									if err == nil {
										Array_Info := make([]ADS_Array_Info, dataType.Array_Dimension)
										err = binary.Read(r, binary.LittleEndian, &Array_Info)
										if err == nil {
											subDataTypes := make(map[string]ADS_Data_Type_Entry_Complete, dataType.Sub_Items)
											if dataType.Sub_Items > 0 {
												subDataTypes, err = extractDataTypes(r, int(dataType.Sub_Items))
											}
											dataTypes[string(Name)] = ADS_Data_Type_Entry_Complete{Entry: dataType, Name: string(Name), Type: string(Type), Comment: string(Comment), Array_Info: Array_Info, Sub_Items: subDataTypes}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

func (c *Connection) GetAdsSymbolDataTypeUpload(info ADS_Symbol_Upload_Info2) (dataTypes map[string]ADS_Data_Type_Entry_Complete, err error) {
	var data []byte
	data, _, err = c.Read(ADSIGRP_SYM_DT_UPLOAD, 0, info.Data_Type_Size)
	if err == nil {
		r := bytes.NewReader(data)
		return extractDataTypes(r, int(info.Data_Types))
	}
	return
}
