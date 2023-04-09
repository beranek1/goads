package goads

import (
	"testing"

	"github.com/beranek1/goadsinterface"
)

func TestAdsLib(t *testing.T) {
	lib, err := NewAdsLib("127.0.0.1", "192.168.178.34.1.1:851")
	if err != nil {
		t.Error(err.Error())
	}
	_, err = lib.GetVersion()
	if err != nil {
		t.Error(err.Error())
	}
	var state goadsinterface.AdsState
	state, err = lib.GetState()
	if err != nil {
		t.Error(err.Error())
	}
	_, err = lib.GetDeviceInfo()
	if err != nil {
		t.Error(err.Error())
	}
	_, err = lib.GetSymbol("test")
	if err == nil {
		t.Error("reading unknown symbol didn't result in error")
	}
	if err.Error() != "symbol not found" {
		t.Error(err)
	}
	_, err = lib.GetSymbol("MAIN.testbool")
	if err != nil {
		t.Error(err.Error())
	}
	_, err = lib.GetSymbolInfo()
	if err != nil {
		t.Error(err.Error())
	}
	_, err = lib.GetSymbolValue("test")
	if err == nil {
		t.Error("reading unknown symbol didn't result in error")
	}
	if err.Error() != "symbol not found" {
		t.Error(err)
	}
	var value goadsinterface.AdsData
	value, err = lib.GetSymbolValue("MAIN.testbool")
	if err != nil {
		t.Error(err.Error())
	}
	_, err = lib.GetSymbolList()
	if err != nil {
		t.Error(err.Error())
	}
	_, err = lib.SetState(state)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = lib.SetSymbolValue("test", value)
	if err == nil {
		t.Error("reading unknown symbol didn't result in error")
	}
	if err.Error() != "symbol not found" {
		t.Error(err)
	}
	_, err = lib.SetSymbolValue("MAIN.testbool", value)
	if err != nil {
		t.Error(err.Error())
	}
}
