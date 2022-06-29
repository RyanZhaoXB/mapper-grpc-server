package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "zhaor.com/mapper-grpc-server/pkg/apis/dmi-mapper/v1"
)

const sockPath string = "/tmp/mapper-1.sock"

type server struct{}

func (s *server) CreateDevice(ctx context.Context, in *pb.CreateDeviceRequest) (*pb.CreateDeviceResponse, error) {
	fmt.Printf("CreateDeviceRequest: %+v\n", in.Config)
	return &pb.CreateDeviceResponse{DeviceName: "device-1"}, nil
}

func (s *server) RemoveDevice(ctx context.Context, in *pb.RemoveDeviceRequest) (*pb.RemoveDeviceResponse, error) {
	fmt.Printf("RemoveDeviceRequest: %+v\n", in.DeviceName)
	return &pb.RemoveDeviceResponse{}, nil
}

func (s *server) UpdateDevice(ctx context.Context, in *pb.UpdateDeviceRequest) (*pb.UpdateDeviceResponse, error) {
	fmt.Printf("UpdateDeviceRequest: %+v\n", in.Config)
	return &pb.UpdateDeviceResponse{}, nil
}

func (s *server) UpdateDeviceStatus(ctx context.Context, in *pb.UpdateDeviceStatusRequest) (*pb.UpdateDeviceStatusResponse, error) {
	fmt.Printf("UpdateDeviceStatusRequest: %+v\n%+v\n", in.DeviceName, in.DesiredDevice)
	return &pb.UpdateDeviceStatusResponse{}, nil
}

func (s *server) GetDevice(ctx context.Context, in *pb.GetDeviceRequest) (*pb.GetDeviceResponse, error) {
	fmt.Printf("GetDeviceRequest: %+v\n", in.DeviceName)
	return &pb.GetDeviceResponse{Status: &pb.DeviceStatus{
		Twins: []*pb.Twin{&pb.Twin{
			PropertyName: "temperature",
			Desired: &pb.TwinProperty{
				Value:    "27",
				Metadata: nil,
			},
			Reported: &pb.TwinProperty{
				Value:    "30",
				Metadata: nil,
			},
		}},
		State: "ok",
	}}, nil
}

func InitSock(sockPath string) error {
	log.Printf("init uds socket: %s", sockPath)
	_, err := os.Stat(sockPath)
	if err == nil {
		err = os.Remove(sockPath)
		if err != nil {
			return err
		}
		return nil
	} else if os.IsNotExist(err) {
		return nil
	} else {
		return fmt.Errorf("fail to stat uds socket path")
	}
}

func main() {
	err := InitSock(sockPath)
	if err != nil {
		log.Fatalf("failed to remove uds socket with err: %v", err)
		return
	}
	lis, err := net.Listen("unix", sockPath)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDeviceMapperServiceServer(s, &server{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
