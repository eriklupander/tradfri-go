package grpc_server

import (
	"context"
	"fmt"
	"strconv"

	pb "github.com/eriklupander/tradfri-go/grpc_server/golang"
	"github.com/eriklupander/tradfri-go/model"
	"github.com/eriklupander/tradfri-go/tradfri"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// New initialize a new tradfri server
func New(tradfriClient *tradfri.TradfriClient) pb.TradfriServiceServer {
	return &server{
		tradfriClient: tradfriClient,
	}
}

type server struct {
	tradfriClient *tradfri.TradfriClient
}

func (s *server) ListGroups(ctx context.Context, r *pb.ListGroupsRequest) (*pb.ListGroupsResponse, error) {
	res := make([]*pb.Group, 0)
	{
		groups, err := s.tradfriClient.ListGroups()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		for _, g := range groups {
			res = append(res, model.ToGroupResponseProto(g))
		}
	}
	return &pb.ListGroupsResponse{
		Groups: res,
	}, nil
}

func (s *server) GetGroup(ctx context.Context, r *pb.GetGroupRequest) (*pb.GetGroupResponse, error) {
	if r.GetId() < 1 {
		return nil, status.Error(codes.InvalidArgument, "group id is mandatory")
	}
	g, err := s.tradfriClient.GetGroup(fmt.Sprintf("%d", r.GetId()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetGroupResponse{
		Group: model.ToGroupResponseProto(g),
	}, nil
}

func (s *server) ListDevices(ctx context.Context, r *pb.ListDevicesRequest) (*pb.ListDevicesResponse, error) {
	if r.GetGroupId() < 1 {
		return nil, status.Error(codes.InvalidArgument, "group id is mandatory")
	}
	g, err := s.tradfriClient.GetGroup(fmt.Sprintf("%d", r.GetGroupId()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	res := make([]*pb.Device, 0)
	for _, id := range g.Content.DeviceList.DeviceIds {
		d, _ := s.tradfriClient.GetDevice(strconv.Itoa(id))
		res = append(res, model.ToDeviceResponseProto(d))
	}
	return &pb.ListDevicesResponse{
		Devices: res,
	}, nil
}

func (s *server) ListDeviceIDs(ctx context.Context, r *pb.ListDeviceIDsRequest) (*pb.ListDeviceIDsResponse, error) {
	if r.GetGroupId() < 1 {
		return nil, status.Error(codes.InvalidArgument, "group id is mandatory")
	}
	g, err := s.tradfriClient.GetGroup(fmt.Sprintf("%d", r.GetGroupId()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	res := make([]int32, 0)
	for _, id := range g.Content.DeviceList.DeviceIds {
		res = append(res, int32(id))
	}
	return &pb.ListDeviceIDsResponse{
		Ids: res,
	}, nil
}

func (s *server) GetDevice(ctx context.Context, r *pb.GetDeviceRequest) (*pb.GetDeviceResponse, error) {
	if r.GetId() < 1 {
		return nil, status.Error(codes.InvalidArgument, "id is mandatory")
	}
	d, err := s.tradfriClient.GetDevice(fmt.Sprintf("%d", r.GetId()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetDeviceResponse{
		Device: model.ToDeviceResponseProto(d),
	}, nil
}

func (s *server) ChangeDeviceColor(ctx context.Context, r *pb.ChangeDeviceColorRequest) (*pb.ChangeDeviceColorResponse, error) {
	if r.GetId() < 1 {
		return nil, status.Error(codes.InvalidArgument, "id is mandatory")
	}
	// rgb
	if r.GetRgb() != "" {
		if _, err := s.tradfriClient.PutDeviceColorRGB(fmt.Sprintf("%d", r.GetId()), r.GetRgb()); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &pb.ChangeDeviceColorResponse{}, nil
	}
	// we assume it is x and y request
	if _, err := s.tradfriClient.PutDeviceColor(fmt.Sprintf("%d", r.GetId()), int(r.GetXcolor()), int(r.GetYcolor())); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ChangeDeviceColorResponse{}, nil
}

func (s *server) ChangeDeviceDimming(ctx context.Context, r *pb.ChangeDeviceDimmingRequest) (*pb.ChangeDeviceDimmingResponse, error) {
	if r.GetId() < 1 {
		return nil, status.Error(codes.InvalidArgument, "id is mandatory")
	}
	if _, err := s.tradfriClient.PutDeviceDimming(fmt.Sprintf("%d", r.GetId()), int(r.GetValue())); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ChangeDeviceDimmingResponse{}, nil
}

func (s *server) TurnDeviceOn(ctx context.Context, r *pb.TurnDeviceOnRequest) (*pb.TurnDeviceOnResponse, error) {
	if r.GetId() < 1 {
		return nil, status.Error(codes.InvalidArgument, "id is mandatory")
	}
	if _, err := s.tradfriClient.PutDevicePower(fmt.Sprintf("%d", r.GetId()), 1); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.TurnDeviceOnResponse{}, nil
}

func (s *server) TurnDeviceOff(ctx context.Context, r *pb.TurnDeviceOffRequest) (*pb.TurnDeviceOffResponse, error) {
	if r.GetId() < 1 {
		return nil, status.Error(codes.InvalidArgument, "id is mandatory")
	}
	if _, err := s.tradfriClient.PutDevicePower(fmt.Sprintf("%d", r.GetId()), 0); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.TurnDeviceOffResponse{}, nil
}
