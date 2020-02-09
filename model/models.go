package model

// Device defines (with JSON tags) a IKEA trådfri device of some kind
type Device struct {
	Metadata struct {
		Vendor   string `json:"0"`
		TypeName string `json:"1"`
		Num2     string `json:"2"`
		TypeId   string `json:"3"`
		Num6     int    `json:"6"`
		Battery  int    `json:"9"`
	} `json:"3"`
	BlindControl []struct {
		Position float32 `json:"5536"`
		Num9003  int     `json:"9003"`
	} `json:"15015"`
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

// DeviceMetadata defines (with JSON tags) common device metadata. Typically embedded in other structs.
type DeviceMetadata struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Vendor  string `json:"vendor"`
	Type    string `json:"type"`
	Battery int    `json:"battery"`
}

// BlindResponse is the response from a bild GET.
type BlindResponse struct {
	DeviceMetadata DeviceMetadata `json:"deviceMetadata"`
	Position float32 `json:"position"`
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
	Created    string `json:"created"`
	DeviceList []int  `json:"deviceList"`
}

// PositioningRequest allows setting the position from 0-100.
type PositioningRequest struct {
	Positioning float32 `json:"positioning"`
}
