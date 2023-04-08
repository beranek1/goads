package goads

import (
	"testing"
)

func TestReadSymbolValue(t *testing.T) {
	c, err := CreateConnection("127.0.0.1", "192.168.178.34.1.1:851")
	if err != nil {
		t.Error(err.Error())
	}
	err = c.Open()
	defer c.Close()
	if err != nil {
		t.Error(err.Error())
	}
	var info ADS_Symbol_Upload_Info2
	info, err = c.GetAdsSymbolUploadInfo2()
	if err != nil {
		t.Error(err.Error())
	}
	var dataTypes map[string]ADS_Data_Type_Entry_Complete
	dataTypes, err = c.GetAdsSymbolDataTypeUpload(info)
	if err != nil {
		t.Error(err.Error())
	}
	var symbols map[string]ADS_Symbol_Entry_Complete
	symbols, err = c.GetAdsSymbolUpload(info)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = ReadSymbolValue(c, GetDataTypeRecursive(dataTypes, "VERSION"), symbols["Constants.CompilerVersion"].Entry.Index_Group, symbols["Constants.CompilerVersion"].Entry.Index_Offset, false)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = ReadSymbolValue(c, GetDataTypeRecursive(dataTypes, "BOOL"), symbols["MAIN.testbool"].Entry.Index_Group, symbols["MAIN.testbool"].Entry.Index_Offset, false)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = ReadSymbolValue(c, GetDataTypeRecursive(dataTypes, "INT"), symbols["MAIN.testint"].Entry.Index_Group, symbols["MAIN.testint"].Entry.Index_Offset, false)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = ReadSymbolValue(c, GetDataTypeRecursive(dataTypes, "STRING(80)"), symbols["MAIN.teststring"].Entry.Index_Group, symbols["MAIN.teststring"].Entry.Index_Offset, false)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = ReadSymbolValue(c, GetDataTypeRecursive(dataTypes, "DUT"), symbols["MAIN.teststruct"].Entry.Index_Group, symbols["MAIN.teststruct"].Entry.Index_Offset, false)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestWriteSymbolValue(t *testing.T) {
	c, err := CreateConnection("127.0.0.1", "192.168.178.34.1.1:851")
	if err != nil {
		t.Error(err.Error())
	}
	err = c.Open()
	defer c.Close()
	if err != nil {
		t.Error(err.Error())
	}
	var info ADS_Symbol_Upload_Info2
	info, err = c.GetAdsSymbolUploadInfo2()
	if err != nil {
		t.Error(err.Error())
	}
	var dataTypes map[string]ADS_Data_Type_Entry_Complete
	dataTypes, err = c.GetAdsSymbolDataTypeUpload(info)
	if err != nil {
		t.Error(err.Error())
	}
	var symbols map[string]ADS_Symbol_Entry_Complete
	symbols, err = c.GetAdsSymbolUpload(info)
	if err != nil {
		t.Error(err.Error())
	}
	err = WriteSymbolValue(c, GetDataTypeRecursive(dataTypes, "BOOL"), symbols["MAIN.testbool"].Entry.Index_Group, symbols["MAIN.testbool"].Entry.Index_Offset, true, false)
	if err != nil {
		t.Error(err.Error())
	}
	err = WriteSymbolValue(c, GetDataTypeRecursive(dataTypes, "INT"), symbols["MAIN.testint"].Entry.Index_Group, symbols["MAIN.testint"].Entry.Index_Offset, int16(42), false)
	if err != nil {
		t.Error(err.Error())
	}
	err = WriteSymbolValue(c, GetDataTypeRecursive(dataTypes, "STRING(80)"), symbols["MAIN.teststring"].Entry.Index_Group, symbols["MAIN.teststring"].Entry.Index_Offset, "Test", false)
	if err != nil {
		t.Error(err.Error())
	}
	strct := map[string][]int16{}
	strct["point"] = []int16{3, 3}
	err = WriteSymbolValue(c, GetDataTypeRecursive(dataTypes, "DUT"), symbols["MAIN.teststruct"].Entry.Index_Group, symbols["MAIN.teststruct"].Entry.Index_Offset, strct, false)
	if err != nil {
		t.Error(err.Error())
	}
}
