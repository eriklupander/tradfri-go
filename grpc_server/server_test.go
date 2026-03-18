package grpc_server

import (
	"context"
	"testing"

	pb "github.com/eriklupander/tradfri-go/grpc_server/golang"
	"github.com/eriklupander/tradfri-go/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mockClient implements TradfriClient without any DTLS or gateway dependency.
type mockClient struct {
	device model.Device
	group  model.Group
	groups []model.Group
	result model.Result
	err    error
}

func (m *mockClient) GetDevice(_ int) (model.Device, error)                       { return m.device, m.err }
func (m *mockClient) GetGroup(_ int) (model.Group, error)                         { return m.group, m.err }
func (m *mockClient) ListGroups() ([]model.Group, error)                          { return m.groups, m.err }
func (m *mockClient) PutDeviceColor(_ int, _, _ int) (model.Result, error)        { return m.result, m.err }
func (m *mockClient) PutDeviceColorRGB(_ int, _ string) (model.Result, error)     { return m.result, m.err }
func (m *mockClient) PutDeviceDimming(_ int, _ int) (model.Result, error)         { return m.result, m.err }
func (m *mockClient) PutDevicePower(_ int, _ int) (model.Result, error)           { return m.result, m.err }
func (m *mockClient) PutDevicePositioning(_ int, _ float32) (model.Result, error) { return m.result, m.err }

func newTestServer(mc *mockClient) *server {
	return &server{tradfriClient: mc}
}

// ── ListGroups ────────────────────────────────────────────────────────────────

func TestListGroups(t *testing.T) {
	mc := &mockClient{
		groups: []model.Group{
			{Name: "Living room", DeviceId: 1},
			{Name: "Bedroom", DeviceId: 2},
		},
	}
	s := newTestServer(mc)
	resp, err := s.ListGroups(context.Background(), &pb.ListGroupsRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(resp.Groups))
	}
}

// ── GetGroup ──────────────────────────────────────────────────────────────────

func TestGetGroup(t *testing.T) {
	mc := &mockClient{group: model.Group{Name: "Hall", DeviceId: 10}}
	s := newTestServer(mc)
	resp, err := s.GetGroup(context.Background(), &pb.GetGroupRequest{Id: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Group.Id != 10 {
		t.Fatalf("expected group id 10, got %d", resp.Group.Id)
	}
}

func TestGetGroup_MissingId(t *testing.T) {
	s := newTestServer(&mockClient{})
	_, err := s.GetGroup(context.Background(), &pb.GetGroupRequest{Id: 0})
	assertCode(t, err, codes.InvalidArgument)
}

// ── ListDevices ───────────────────────────────────────────────────────────────

func TestListDevices(t *testing.T) {
	mc := &mockClient{
		group: model.Group{
			DeviceId: 5,
			Content: struct {
				DeviceList struct {
					DeviceIds []int `json:"9003"`
				} `json:"15002"`
			}{DeviceList: struct {
				DeviceIds []int `json:"9003"`
			}{DeviceIds: []int{101, 102}}},
		},
		device: model.Device{
			Name:     "Bulb",
			DeviceId: 101,
			LightControl: []struct {
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
			}{{Power: 1, Dimmer: 128}},
		},
	}
	s := newTestServer(mc)
	resp, err := s.ListDevices(context.Background(), &pb.ListDevicesRequest{GroupId: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Devices) != 2 {
		t.Fatalf("expected 2 devices, got %d", len(resp.Devices))
	}
}

func TestListDevices_MissingGroupId(t *testing.T) {
	s := newTestServer(&mockClient{})
	_, err := s.ListDevices(context.Background(), &pb.ListDevicesRequest{GroupId: 0})
	assertCode(t, err, codes.InvalidArgument)
}

// ── ListDeviceIDs ─────────────────────────────────────────────────────────────

func TestListDeviceIDs(t *testing.T) {
	mc := &mockClient{
		group: model.Group{
			Content: struct {
				DeviceList struct {
					DeviceIds []int `json:"9003"`
				} `json:"15002"`
			}{DeviceList: struct {
				DeviceIds []int `json:"9003"`
			}{DeviceIds: []int{10, 20, 30}}},
		},
	}
	s := newTestServer(mc)
	resp, err := s.ListDeviceIDs(context.Background(), &pb.ListDeviceIDsRequest{GroupId: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Ids) != 3 {
		t.Fatalf("expected 3 ids, got %d", len(resp.Ids))
	}
}

func TestListDeviceIDs_MissingGroupId(t *testing.T) {
	s := newTestServer(&mockClient{})
	_, err := s.ListDeviceIDs(context.Background(), &pb.ListDeviceIDsRequest{GroupId: 0})
	assertCode(t, err, codes.InvalidArgument)
}

// ── GetDevice ─────────────────────────────────────────────────────────────────

func TestGetDevice(t *testing.T) {
	mc := &mockClient{
		device: model.Device{Name: "Blind 1", DeviceId: 7, BlindControl: []struct {
			Position float32 `json:"5536"`
			DeviceId int     `json:"9003"`
		}{{Position: 75.0}}},
	}
	s := newTestServer(mc)
	resp, err := s.GetDevice(context.Background(), &pb.GetDeviceRequest{Id: 7})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Device.Position != 75.0 {
		t.Fatalf("expected position 75.0, got %f", resp.Device.Position)
	}
}

func TestGetDevice_MissingId(t *testing.T) {
	s := newTestServer(&mockClient{})
	_, err := s.GetDevice(context.Background(), &pb.GetDeviceRequest{Id: 0})
	assertCode(t, err, codes.InvalidArgument)
}

// ── ChangeDeviceColor ─────────────────────────────────────────────────────────

func TestChangeDeviceColor_RGB(t *testing.T) {
	s := newTestServer(&mockClient{})
	_, err := s.ChangeDeviceColor(context.Background(), &pb.ChangeDeviceColorRequest{Id: 7, Rgb: "ff0000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestChangeDeviceColor_XY(t *testing.T) {
	s := newTestServer(&mockClient{})
	_, err := s.ChangeDeviceColor(context.Background(), &pb.ChangeDeviceColorRequest{Id: 7, Xcolor: 30000, Ycolor: 20000})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestChangeDeviceColor_MissingId(t *testing.T) {
	s := newTestServer(&mockClient{})
	_, err := s.ChangeDeviceColor(context.Background(), &pb.ChangeDeviceColorRequest{Id: 0})
	assertCode(t, err, codes.InvalidArgument)
}

// ── ChangeDeviceDimming ───────────────────────────────────────────────────────

func TestChangeDeviceDimming(t *testing.T) {
	s := newTestServer(&mockClient{})
	_, err := s.ChangeDeviceDimming(context.Background(), &pb.ChangeDeviceDimmingRequest{Id: 7, Value: 200})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestChangeDeviceDimming_MissingId(t *testing.T) {
	s := newTestServer(&mockClient{})
	_, err := s.ChangeDeviceDimming(context.Background(), &pb.ChangeDeviceDimmingRequest{Id: 0})
	assertCode(t, err, codes.InvalidArgument)
}

// ── TurnDeviceOn / TurnDeviceOff ──────────────────────────────────────────────

func TestTurnDeviceOn(t *testing.T) {
	s := newTestServer(&mockClient{})
	_, err := s.TurnDeviceOn(context.Background(), &pb.TurnDeviceOnRequest{Id: 7})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTurnDeviceOn_MissingId(t *testing.T) {
	s := newTestServer(&mockClient{})
	_, err := s.TurnDeviceOn(context.Background(), &pb.TurnDeviceOnRequest{Id: 0})
	assertCode(t, err, codes.InvalidArgument)
}

func TestTurnDeviceOff(t *testing.T) {
	s := newTestServer(&mockClient{})
	_, err := s.TurnDeviceOff(context.Background(), &pb.TurnDeviceOffRequest{Id: 7})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTurnDeviceOff_MissingId(t *testing.T) {
	s := newTestServer(&mockClient{})
	_, err := s.TurnDeviceOff(context.Background(), &pb.TurnDeviceOffRequest{Id: 0})
	assertCode(t, err, codes.InvalidArgument)
}

// ── ChangeDevicePositioning ───────────────────────────────────────────────────

func TestChangeDevicePositioning(t *testing.T) {
	s := newTestServer(&mockClient{})
	_, err := s.ChangeDevicePositioning(context.Background(), &pb.ChangeDevicePositioningRequest{Id: 7, Value: 50.0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestChangeDevicePositioning_MissingId(t *testing.T) {
	s := newTestServer(&mockClient{})
	_, err := s.ChangeDevicePositioning(context.Background(), &pb.ChangeDevicePositioningRequest{Id: 0})
	assertCode(t, err, codes.InvalidArgument)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func assertCode(t *testing.T, err error, expected codes.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %v, got nil", expected)
	}
	if s, ok := status.FromError(err); !ok || s.Code() != expected {
		t.Fatalf("expected gRPC code %v, got %v", expected, err)
	}
}
