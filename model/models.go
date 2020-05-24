package model

// Device defines (with JSON tags) a IKEA trådfri device of some kind
type Device struct {
	Metadata struct {
		Vendor       string `json:"0"`
		TypeName     string `json:"1"`
		SerialNumber string `json:"2"`
		TypeId       string `json:"3"`
		PowerType    int    `json:"6"`
		Battery      int    `json:"9"`
	} `json:"3"`
	LightControl []struct {
		RGBHex           string  `json:"5706"`
		Hue              int     `json:"5707"`
		Saturation       int     `json:"5708"`
		CIE_1931_X       int     `json:"5709"`
		CIE_1931_Y       int     `json:"5710"`
		ColorTemperature int     `json:"5711"`
		TransitionTime   float64 `json:"5712"`
		Power            int     `json:"5850"`
		Dimmer           int     `json:"5851"`
		DeviceId         int     `json:"9003"`
	} `json:"3311"`
	BlindControl []struct {
		Position float32 `json:"5536"`
		DeviceId int     `json:"9003"`
	} `json:"15015"`
	Type           int    `json:"5750"`
	Name           string `json:"9001"`
	CreatedAt      int    `json:"9002"`
	DeviceId       int    `json:"9003"`
	Alive          int    `json:"9019"`
	LastSeen       int    `json:"9020"`
	OtaUpdateState int    `json:"9054"`
}

// Group defines (with JSON tags) a IKEA trådfri Group.
type Group struct {
	Power     int    `json:"5850"`
	Dimmer    int    `json:"5851"`
	Name      string `json:"9001"`
	CreatedAt int    `json:"9002"`
	DeviceId  int    `json:"9003"`
	Content   struct {
		DeviceList struct {
			DeviceIds []int `json:"9003"`
		} `json:"15002"`
	} `json:"9018"`
	SceneId   int `json:"9039"`
	GroupType int `json:"9108"`
}

// RemoteControl defines (with JSON tags) a IKEA remote control.
type RemoteControl struct {
	Metadata struct {
		Manufacturer    string `json:"0"`
		ModelNumber     string `json:"1"`
		SerialNumber    string `json:"2"`
		FirmwareVersion string `json:"3"`
		PowerType       int    `json:"6"`
		Battery         int    `json:"9"`
	} `json:"3"`
	Type           int    `json:"5750"`
	Name           string `json:"9001"`
	CreatedAt      int    `json:"9002"`
	DeviceId       int    `json:"9003"`
	Alive          int    `json:"9019"`
	LastSeen       int    `json:"9020"`
	OtaUpdateState int    `json:"9054"`
	SwitchList     []struct {
		Num9003 int `json:"9003"`
	} `json:"15009"`
}

// ControlOutlet defines (with JSON tags) a IKEA control outlet.
type ControlOutlet struct {
	Metadata struct {
		Manufacturer    string `json:"0"`
		ModelNumber     string `json:"1"`
		SerialNumber    string `json:"2"`
		FirmwareVersion string `json:"3"`
		PowerType       int    `json:"6"`
	} `json:"3"`
	PowerControl []struct {
		Power    int `json:"5850"`
		Dimmer   int `json:"5851"`
		DeviceId int `json:"9003"`
	} `json:"3312"`
	Type           int    `json:"5750"`
	Name           string `json:"9001"`
	CreatedAt      int    `json:"9002"`
	DeviceId       int    `json:"9003"`
	Alive          int    `json:"9019"`
	LastSeen       int    `json:"9020"`
	OtaUpdateState int    `json:"9054"`
	HashUnknown    string `json:"9084"`
}

// DeviceMetadata defines (with JSON tags) common device metadata. Typically embedded in other structs.
type DeviceMetadata struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Vendor  string `json:"vendor"`
	Type    string `json:"type"`
	Battery int    `json:"battery"`
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

// Result is a generic result containing a plain text message
type Result struct {
	Msg string
}

// TokenExchange maps the human-readable Token and TypeIdentifies into their IKEA specific numeric codes.
type TokenExchange struct {
	Token          string `json:"9091"`
	TypeIdentifier string `json:"9029"`
}

// REST API structs

// GroupResponse defines a Group JSON response
type GroupResponse struct {
	Id         int    `json:"id"`
	Power      int    `json:"power"`
	Created    string `json:"created"`
	DeviceList []int  `json:"deviceList"`
}

// BlindResponse is the response from a blind GET.
type BlindResponse struct {
	DeviceMetadata DeviceMetadata `json:"deviceMetadata"`
	Position       float32        `json:"position"`
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

// PositioningRequest allows setting the position from 0-100.
type PositioningRequest struct {
	Positioning float32 `json:"positioning"`
}
