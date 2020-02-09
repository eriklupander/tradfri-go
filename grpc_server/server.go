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
func New(tradfriClient *tradfri.Client) pb.TradfriServiceServer {
	return &server{
		tradfriClient: tradfriClient,
	}
}

type server struct {
	tradfriClient *tradfri.Client
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

func (s *server) ChangeDevicePositioning(ctx context.Context, r *pb.ChangeDevicePositioningRequest) (*pb.ChangeDevicePositioningResponse, error) {
	if r.GetId() < 1 {
		return nil, status.Error(codes.InvalidArgument, "id is mandatory")
	}
	if _, err := s.tradfriClient.PutDevicePositioning(fmt.Sprintf("%d", r.GetId()), float32(r.GetValue())); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ChangeDevicePositioningResponse{}, nil
}
