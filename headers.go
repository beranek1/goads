package goads

type AMSTCP_Header struct {
	Reserved [2]byte
	Length   uint32
}

type AMS_Header struct {
	AMSNetId_Target [6]byte
	AMSPort_Target  uint16
	AMSNetId_Source [6]byte
	AMSPort_Source  uint16
	Command_Id      uint16
	State_Flags     [2]byte
	Data_Length     uint32
	Error_Code      uint32
	Invoke_Id       uint32
}

type ADS_Command uint16

const (
	ADS_Command_Invalid ADS_Command = iota
	ADS_Command_Read_Device_Info
	ADS_Command_Read
	ADS_Command_Write
	ADS_Command_Read_State
	ADS_Command_Write_Control
	ADS_Command_Add_Device_Notification
	ADS_Command_Delete_Device_Notification
	ADS_Command_Device_Notification
	ADS_Command_Read_Write
)

type ADS_Read_Device_Info_Response struct {
	Result        uint32
	Major_Version uint8
	Minor_Version uint8
	Version_Build uint16
	Device_Name   [16]byte
}

type ADS_Read_Request struct {
	Index_Group  uint32
	Index_Offset uint32
	Length       uint32
}

type ADS_Read_Response struct {
	Result uint32
	Length uint32
}

type ADS_Write_Request struct {
	Index_Group  uint32
	Index_Offset uint32
	Length       uint32
}

type ADS_Write_Response struct {
	Result uint32
}

type ADS_Read_State_Response struct {
	Result       uint32
	ADS_State    uint16
	Device_State uint16
}

type ADS_Write_Control_Request struct {
	ADS_State    uint16
	Device_State uint16
	Length       uint32
}

type ADS_Write_Control_Response struct {
	Result uint32
}

type ADS_Add_Device_Notification_Request struct {
	Index_Group       uint32
	Index_Offset      uint32
	Length            uint32
	Transmission_Mode uint32
	Max_Delay         uint32
	Cycle_Time        uint32
	Reserved          [16]byte
}

type ADS_Add_Device_Notification_Response struct {
	Result              uint32
	Notification_Handle uint32
}

type ADS_Delete_Device_Notification_Request struct {
	Notification_Handle uint32
}

type ADS_Delete_Device_Notification_Response struct {
	Result uint32
}

type ADS_Read_Write_Request struct {
	Index_Group  uint32
	Index_Offset uint32
	Read_Length  uint32
	Write_Length uint32
}

type ADS_Read_Write_Response struct {
	Result uint32
	Length uint32
}
