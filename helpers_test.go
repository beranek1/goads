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
}
