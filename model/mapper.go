package model

func ToDeviceResponse(device Device) BulbResponse {
	if device.LightControl != nil && len(device.LightControl) > 0 {

		dr := BulbResponse{
			DeviceMetadata: DeviceMetadata{Name: device.Name,
				Id:     device.DeviceId,
				Type:   device.Metadata.TypeName,
				Vendor: device.Metadata.Vendor,},
			Powered:    device.LightControl[0].Power == 1,
			CIE_1931_X: device.LightControl[0].CIE_1931_X,
			CIE_1931_Y: device.LightControl[0].CIE_1931_Y,
			RGB:        device.LightControl[0].RGBHex,
			Dimmer:     device.LightControl[0].Dimmer,
		}
		return dr
	}
	return BulbResponse{}
}
