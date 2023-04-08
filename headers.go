package goads

type AMSTCP_Header struct {
	Reserved uint16
	Length   uint32
}

type AMS_Address struct {
	NetId [6]byte
	Port  uint16
}

type AMS_Header struct {
	Target      AMS_Address
	Source      AMS_Address
	Command_Id  uint16
	State_Flags uint16
	Data_Length uint32
	Error_Code  uint32
	Invoke_Id   uint32
}

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

type ADS_Symbol_Entry struct {
	Entry_Length   uint32
	Index_Group    uint32
	Index_Offset   uint32
	Size           uint32
	Data_Type      uint32
	Flags          uint32
	Name_Length    uint16
	Type_Length    uint16
	Comment_Length uint16
}

type ADS_Symbol_Entry_Complete struct {
	Entry   ADS_Symbol_Entry
	Name    string
	Type    string
	Comment string
}

type ADS_Array_Info struct {
	Bound uint32
	Size  uint32
}

type ADS_Data_Type_Entry struct {
	Entry_Length uint32
	Version      uint32
	// Hash_Value      uint32
	Offset_Get uint32
	// Type_Hash_Value uint32
	Offset_Set      uint32
	Size            uint32
	Offset          uint32
	Data_Type       uint32
	Flags           uint32
	Name_Length     uint16
	Type_Length     uint16
	Comment_Length  uint16
	Array_Dimension uint16
	Sub_Items       uint16
}

type ADS_Data_Type_Entry_Complete struct {
	Entry      ADS_Data_Type_Entry
	Name       string
	Type       string
	Comment    string
	Sub_Items  map[string]ADS_Data_Type_Entry_Complete
	Array_Info []ADS_Array_Info
}

type ADS_Symbol_Upload_Info struct {
	Symbols     uint32
	Symbol_Size uint32
}

type ADS_Symbol_Upload_Info2 struct {
	Symbols              uint32
	Symbol_Size          uint32
	Data_Types           uint32
	Data_Type_Size       uint32
	Max_Dynamic_Symbols  uint32
	Used_Dynamic_Symbols uint32
}

type ADS_Symbol_Info_By_Name struct {
	Index_Group  uint32
	Index_Offset uint32
	Length       uint32
}
