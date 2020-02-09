package model

import (
	pb "github.com/eriklupander/tradfri-go/grpc_server/golang"
	"time"
)

func ToDeviceResponse(device Device) BlindResponse {
	if device.BlindControl != nil && len(device.BlindControl) > 0 {

		dr := BlindResponse{
			DeviceMetadata:  DeviceMetadata{
				Name:    device.Name,
				Id:      device.DeviceId,
				Type:    device.Metadata.TypeName,
				Vendor:  device.Metadata.Vendor,
				Battery: device.Metadata.Battery,
			},
			Position: device.BlindControl[0].Position,
		}
		return dr
	}
	return BlindResponse{}
}

func ToDeviceResponseProto(device Device) *pb.Device {
	if device.BlindControl != nil && len(device.BlindControl) > 0 {
		return &pb.Device{
			Metadata: &pb.DeviceMetadata{
				Name:    device.Name,
				Id:      int32(device.DeviceId),
				Type:    device.Metadata.TypeName,
				Vendor:  device.Metadata.Vendor,
				Battery: int(device.Metadata.Battery),
			},
			Position: float32(device.BlindControl[0].Position),
		}
	}
	return &pb.Device{}
}

func ToGroupResponse(group Group) GroupResponse {
	gr := GroupResponse{
		Id:         group.DeviceId,
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
		Created: time.Unix(int64(group.Num9002), 0).Format(time.RFC3339),
		Devices: ids,
	}
}
