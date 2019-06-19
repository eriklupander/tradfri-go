package model

import (
	"time"

	pb "github.com/eriklupander/tradfri-go/grpc_server/golang"
)

func ToDeviceResponse(device Device) BulbResponse {
	if device.LightControl != nil && len(device.LightControl) > 0 {

		dr := BulbResponse{
			DeviceMetadata: DeviceMetadata{
				Name:   device.Name,
				Id:     device.DeviceId,
				Type:   device.Metadata.TypeName,
				Vendor: device.Metadata.Vendor},
			Power:      device.LightControl[0].Power == 1,
			CIE_1931_X: device.LightControl[0].CIE_1931_X,
			CIE_1931_Y: device.LightControl[0].CIE_1931_Y,
			RGB:        device.LightControl[0].RGBHex,
			Dimmer:     device.LightControl[0].Dimmer,
		}
		return dr
	}
	return BulbResponse{}
}

func ToDeviceResponseProto(device Device) *pb.Device {
	if device.LightControl != nil && len(device.LightControl) > 0 {
		return &pb.Device{
			Metadata: &pb.DeviceMetadata{
				Name:   device.Name,
				Id:     int32(device.DeviceId),
				Type:   device.Metadata.TypeName,
				Vendor: device.Metadata.Vendor,
			},
			Power:  device.LightControl[0].Power == 1,
			Xcolor: int32(device.LightControl[0].CIE_1931_X),
			Ycolor: int32(device.LightControl[0].CIE_1931_Y),
			Rgb:    device.LightControl[0].RGBHex,
			Dimmer: int32(device.LightControl[0].Dimmer),
		}
	}
	return &pb.Device{}
}

func ToGroupResponse(group Group) GroupResponse {
	gr := GroupResponse{
		Id:         group.DeviceId,
		Power:      group.Power,
		Created:    time.Unix(int64(group.Num9002), 0).Format(time.RFC3339),
		DeviceList: group.Content.DeviceList.DeviceIds,
	}
	return gr
}

func ToGroupResponseProto(group Group) *pb.Group {
	ids := make([]int32, 0, len(group.Content.DeviceList.DeviceIds))
	for _, v := range group.Content.DeviceList.DeviceIds {
		ids = append(ids, int32(v))
	}
	return &pb.Group{
		Id:      int32(group.DeviceId),
		Power:   int32(group.Power),
		Created: time.Unix(int64(group.Num9002), 0).Format(time.RFC3339),
		Devices: ids,
	}
}
