package model

// Device defines (with JSON tags) a IKEA trådfri device of some kind
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

// Group defines (with JSON tags) a IKEA trådfri Group.
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

// RemoteControl defines (with JSON tags) a IKEA remote control.
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

// ControlOutlet defines (with JSON tags) a IKEA control outlet.
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

// DeviceMetadata defines (with JSON tags) common device metadata. Typically embedded in other structs.
type DeviceMetadata struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Vendor string `json:"vendor"`
	Type   string `json:"type"`
}

// PowerPlugResponse is the response from a power plug device GET.
type PowerPlugResponse struct {
	DeviceMetadata DeviceMetadata `json:"deviceMetadata"`
	Power          bool           `json:"power"`
}

// BulbResponse is the response from a light bulb GET.
type BulbResponse struct {
	DeviceMetadata DeviceMetadata `json:"deviceMetadata"`
	Dimmer         int            `json:"dimmer"`
	CIE_1931_X     int            `json:"xcolor"`
	CIE_1931_Y     int            `json:"ycolor"`
	RGB            string         `json:"rgbcolor"`
	Power          bool           `json:"power"`
}

type Result struct {
	Msg string
}

type TokenExchange struct {
	Token          string `json:"9091"`
	TypeIdentifier string `json:"9029"`
}

// REST API structs
type GroupResponse struct {
	Id         int    `json:"id"`
	Power      int    `json:"power"`
	Created    string `json:"created"`
	DeviceList []int  `json:"deviceList"`
}

// RgbColorRequest allows (trying to) set a bulb color using classic hex RGB string.
type RgbColorRequest struct {
	RGBcolor string `json:"rgbcolor"`
}

// DimmingRequest allows setting the dimmer level from 0-255.
type DimmingRequest struct {
	Dimming int `json:"dimming"`
}

// PowerRequest contains a Power state int, 1 == on, 0 == off.
type PowerRequest struct {
	Power int `json:"power"`
}

// StateRequest allows setting both color, dimmer and power setting in a single PUT.
type StateRequest struct {
	RGBcolor string `json:"rgbcolor"`
	Dimmer   int    `json:"dimmer"`
	Power    int    `json:"power"`
}
