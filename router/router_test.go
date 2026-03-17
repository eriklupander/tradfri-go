package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eriklupander/tradfri-go/model"
)

// mockClient implements TradfriClient without any DTLS or gateway dependency.
type mockClient struct {
	device model.Device
	group  model.Group
	groups []model.Group
	result model.Result
	err    error
}

func (m *mockClient) GetDevice(_ int) (model.Device, error)                        { return m.device, m.err }
func (m *mockClient) GetGroup(_ int) (model.Group, error)                          { return m.group, m.err }
func (m *mockClient) ListGroups() ([]model.Group, error)                           { return m.groups, m.err }
func (m *mockClient) PutDeviceColor(_ int, _, _ int) (model.Result, error)         { return m.result, m.err }
func (m *mockClient) PutDeviceColorRGB(_ int, _ string) (model.Result, error)      { return m.result, m.err }
func (m *mockClient) PutDeviceDimming(_ int, _ int) (model.Result, error)          { return m.result, m.err }
func (m *mockClient) PutDevicePower(_ int, _ int) (model.Result, error)            { return m.result, m.err }
func (m *mockClient) PutDeviceState(_ int, _, _ int) (model.Result, error)         { return m.result, m.err }
func (m *mockClient) PutDevicePositioning(_ int, _ float32) (model.Result, error)  { return m.result, m.err }

func newTestRouter(mc *mockClient) http.Handler {
	return newRouter(mc)
}

func TestHealth(t *testing.T) {
	r := newTestRouter(&mockClient{})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "OK" {
		t.Fatalf("expected body OK, got %q", rec.Body.String())
	}
}

func TestListGroups(t *testing.T) {
	mc := &mockClient{
		groups: []model.Group{
			{Name: "Living room", DeviceId: 1, Content: struct {
				DeviceList struct {
					DeviceIds []int `json:"9003"`
				} `json:"15002"`
			}{DeviceList: struct {
				DeviceIds []int `json:"9003"`
			}{DeviceIds: []int{10, 11}}}},
		},
	}
	r := newTestRouter(mc)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/groups", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp []model.GroupResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp) != 1 || resp[0].Id != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestGetGroup(t *testing.T) {
	mc := &mockClient{
		group: model.Group{Name: "Bedroom", DeviceId: 42},
	}
	r := newTestRouter(mc)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/groups/42", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp model.GroupResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Id != 42 {
		t.Fatalf("expected group id 42, got %d", resp.Id)
	}
}

func TestGetDevice(t *testing.T) {
	mc := &mockClient{
		device: model.Device{
			Name:     "Bulb 1",
			DeviceId: 7,
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
			}{{Power: 1, Dimmer: 200}},
		},
	}
	r := newTestRouter(mc)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/device/7", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp model.BulbResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !resp.Power || resp.Dimmer != 200 {
		t.Fatalf("unexpected device response: %+v", resp)
	}
}

func TestSetDimming(t *testing.T) {
	mc := &mockClient{result: model.Result{Msg: "2.05 Content"}}
	r := newTestRouter(mc)
	body, _ := json.Marshal(model.DimmingRequest{Dimming: 128})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/api/device/7/dimmer", bytes.NewReader(body)))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestSetPower(t *testing.T) {
	mc := &mockClient{result: model.Result{Msg: "2.05 Content"}}
	r := newTestRouter(mc)
	body, _ := json.Marshal(model.PowerRequest{Power: 1})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/api/device/7/power", bytes.NewReader(body)))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestSetColorRGB(t *testing.T) {
	mc := &mockClient{result: model.Result{Msg: "2.05 Content"}}
	r := newTestRouter(mc)
	body, _ := json.Marshal(model.RgbColorRequest{RGBcolor: "ff0000"})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/api/device/7/rgb", bytes.NewReader(body)))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestSetState(t *testing.T) {
	mc := &mockClient{result: model.Result{Msg: "2.05 Content"}}
	r := newTestRouter(mc)
	body, _ := json.Marshal(model.StateRequest{Power: 1, Dimmer: 100})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/api/device/7", bytes.NewReader(body)))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestSetPositioning(t *testing.T) {
	mc := &mockClient{result: model.Result{Msg: "2.05 Content"}}
	r := newTestRouter(mc)
	body, _ := json.Marshal(model.PositioningRequest{Positioning: 50.0})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/api/device/7/position", bytes.NewReader(body)))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestBadDeviceId(t *testing.T) {
	r := newTestRouter(&mockClient{})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/device/notanumber", nil))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestBadGroupId(t *testing.T) {
	r := newTestRouter(&mockClient{})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/groups/notanumber", nil))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
