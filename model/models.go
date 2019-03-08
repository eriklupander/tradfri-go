package model

type Configuration struct {
	ID string `json:"id"`
}

type Device struct {
	Metadata struct {
		Vendor   string `json:"0"`
		TypeName string `json:"1"`
		Num2     string `json:"2"`
		TypeId   string `json:"3"`
		Num6     int    `json:"6"`
	} `json:"3"`
	LightControl []struct {
		RGBHex     string `json:"5706"`
		Num5707    int    `json:"5707"`
		Num5708    int    `json:"5708"`
		CIE_1931_X int    `json:"5709"`
		CIE_1931_Y int    `json:"5710"`
		Power      int    `json:"5850"`
		Dimmer     int    `json:"5851"`
		Num9003    int    `json:"9003"`
	} `json:"3311"`
	Num5750  int    `json:"5750"`
	Name     string `json:"9001"`
	Num9002  int    `json:"9002"`
	DeviceId int    `json:"9003"`
	Num9019  int    `json:"9019"`
	Num9020  int    `json:"9020"`
	Num9054  int    `json:"9054"`
}

type Group struct {
	Power    int    `json:"5850"`
	Num5851  int    `json:"5851"`
	Name     string `json:"9001"`
	Num9002  int    `json:"9002"`
	DeviceId int    `json:"9003"`
	Content  struct {
		DeviceList struct {
			DeviceIds []int `json:"9003"`
		} `json:"15002"`
	} `json:"9018"`
	Num9039 int `json:"9039"`
	Num9108 int `json:"9108"`
}

type RemoteControl struct {
	Metadata struct {
		Num0 string `json:"0"`
		Num1 string `json:"1"`
		Num2 string `json:"2"`
		Num3 string `json:"3"`
		Num6 int    `json:"6"`
		Num9 int    `json:"9"`
	} `json:"3"`
	Num5750  int    `json:"5750"`
	Num9001  string `json:"9001"`
	Num9002  int    `json:"9002"`
	Num9003  int    `json:"9003"`
	Num9019  int    `json:"9019"`
	Num9020  int    `json:"9020"`
	Num9054  int    `json:"9054"`
	Num15009 []struct {
		Num9003 int `json:"9003"`
	} `json:"15009"`
}

type ControlOutlet struct {
	Metadata struct {
		Num0 string `json:"0"`
		Num1 string `json:"1"`
		Num2 string `json:"2"`
		Num3 string `json:"3"`
		Num6 int    `json:"6"`
	} `json:"3"`
	PowerControl []struct {
		Num5850 int `json:"5850"`
		Num5851 int `json:"5851"`
		Num9003 int `json:"9003"`
	} `json:"3312"`
	Num5750 int    `json:"5750"`
	Num9001 string `json:"9001"`
	Num9002 int    `json:"9002"`
	Num9003 int    `json:"9003"`
	Num9019 int    `json:"9019"`
	Num9020 int    `json:"9020"`
	Num9054 int    `json:"9054"`
	Num9084 string `json:"9084"`
}

type DeviceMetadata struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Vendor string `json:"vendor"`
	Type   string `json:"type"`
}

type PowerPlugResponse struct {
	DeviceMetadata DeviceMetadata `json:"deviceMetadata"`
	Powered        bool           `json:"powered"`
}

type BulbResponse struct {
	DeviceMetadata DeviceMetadata `json:"deviceMetadata"`
	Dimmer         int            `json:"dimmer"`
	CIE_1931_X     int            `json:"xcolor"`
	CIE_1931_Y     int            `json:"ycolor"`
	RGB            string         `json:"rgbcolor"`
	Powered        bool           `json:"powered"`
}

type Result struct {
	Msg string
}