package goads

import (
	"errors"

	"github.com/beranek1/goadsinterface"
)

// This struct/class is for providing legacy support for Go applications that used ADSBridge via adsbridgegolib
type AdsLib struct {
	c         *Connection
	dataTypes map[string]ADS_Data_Type_Entry_Complete
	symbols   map[string]ADS_Symbol_Entry_Complete
}

func NewAdsLib(ip string, target string) (lib *AdsLib, err error) {
	var c *Connection
	c, err = NewConnection(ip, target)
	if err == nil {
		err = c.Open()
		if err == nil {
			lib = &AdsLib{c: c}
		}
	}
	return
}

func (l *AdsLib) GetVersion() (version goadsinterface.AdsVersion, err error) {
	version = goadsinterface.AdsVersion{Version: 1, Revision: 1, Build: 1}
	return
}

func (l *AdsLib) GetState() (state goadsinterface.AdsState, err error) {
	var res ADS_Read_State_Response
	res, err = l.c.ReadState()
	if err == nil {
		state = goadsinterface.AdsState{Ads: res.ADS_State, Device: res.Device_State}
	}
	return
}

func (l *AdsLib) GetDeviceInfo() (info goadsinterface.AdsDeviceInfo, err error) {
	var res ADS_Read_Device_Info_Response
	res, err = l.c.ReadDeviceInfo()
	if err == nil {
		info = goadsinterface.AdsDeviceInfo{Name: string(res.Device_Name[:]), Version: goadsinterface.AdsVersion{Version: res.Major_Version, Revision: res.Minor_Version, Build: res.Version_Build}}
	}
	return
}

func (l *AdsLib) updateSymbolUploadInfo() (err error) {
	var info ADS_Symbol_Upload_Info2
	info, err = l.c.GetAdsSymbolUploadInfo2()
	if err == nil {
		l.dataTypes, err = l.c.GetAdsSymbolDataTypeUpload(info)
		if err == nil {
			l.symbols, err = l.c.GetAdsSymbolUpload(info)
		}
	}
	return
}

func (l *AdsLib) GetSymbol(name string) (symbol goadsinterface.AdsSymbol, err error) {
	if l.symbols == nil || len(l.symbols) == 0 {
		err = l.updateSymbolUploadInfo()
	}
	if err == nil {
		if s, ok := l.symbols[name]; ok {
			symbol = goadsinterface.AdsSymbol{Name: s.Name, IndexGroup: s.Entry.Index_Group, IndexOffset: s.Entry.Index_Offset, Size: s.Entry.Size, Type: s.Type, Comment: s.Comment}
		} else {
			err = errors.New("symbol not found")
		}
	}
	return
}

func (l *AdsLib) GetSymbolInfo() (info goadsinterface.AdsSymbolInfo, err error) {
	if l.symbols == nil || len(l.symbols) == 0 {
		err = l.updateSymbolUploadInfo()
	}
	if err == nil {
		info = make(goadsinterface.AdsSymbolInfo)
		for k := range l.symbols {
			s := l.symbols[k]
			info[k] = goadsinterface.AdsSymbol{Name: s.Name, IndexGroup: s.Entry.Index_Group, IndexOffset: s.Entry.Index_Offset, Size: s.Entry.Size, Type: s.Type, Comment: s.Comment}
		}
	}
	return
}

func (l *AdsLib) GetSymbolValue(name string) (data goadsinterface.AdsData, err error) {
	if l.symbols == nil || len(l.symbols) == 0 {
		err = l.updateSymbolUploadInfo()
	}
	if err == nil {
		if s, ok := l.symbols[name]; ok {
			var value any
			value, err = ReadSymbolValue(l.c, GetDataTypeRecursive(l.dataTypes, s.Type), s.Entry.Index_Group, s.Entry.Index_Offset, false)
			if err == nil {
				data = goadsinterface.AdsData{Data: value}
			}
		} else {
			err = errors.New("symbol not found")
		}
	}
	return
}

func (l *AdsLib) GetSymbolList() (list goadsinterface.AdsSymbolList, err error) {
	if l.symbols == nil || len(l.symbols) == 0 {
		err = l.updateSymbolUploadInfo()
	}
	if err == nil {
		list = make(goadsinterface.AdsSymbolList, len(l.symbols))
		i := 0
		for k := range l.symbols {
			list[i] = k
			i++
		}
	}
	return
}

func (l *AdsLib) SetState(state_in goadsinterface.AdsState) (state_out goadsinterface.AdsState, err error) {
	_, err = l.c.WriteControl(state_in.Ads, state_in.Device)
	if err == nil {
		state_out, err = l.GetState()
	}
	return
}

func (l *AdsLib) SetSymbolValue(name string, value goadsinterface.AdsData) (data goadsinterface.AdsData, err error) {
	if l.symbols == nil || len(l.symbols) == 0 {
		err = l.updateSymbolUploadInfo()
	}
	if err == nil {
		if s, ok := l.symbols[name]; ok {
			err = WriteSymbolValue(l.c, GetDataTypeRecursive(l.dataTypes, s.Type), s.Entry.Index_Group, s.Entry.Index_Offset, value.Data, false)
			if err == nil {
				data = value
			}
		} else {
			err = errors.New("symbol not found")
		}
	}
	return
}
