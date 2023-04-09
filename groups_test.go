package goads

import (
	"testing"
)

func TestAdsSymbolUpload(t *testing.T) {
	c, err := NewConnection("127.0.0.1", "192.168.178.34.1.1:851")
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
	_, err = c.GetAdsSymbolUpload(info)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAdsSymbolDataTypeUpload(t *testing.T) {
	c, err := NewConnection("127.0.0.1", "192.168.178.34.1.1:851")
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
	_, err = c.GetAdsSymbolDataTypeUpload(info)
	if err != nil {
		t.Error(err.Error())
	}
}
