package main

import (
	"log"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "zhaor.com/mapper-grpc-server/pkg/apis/dmi-mapper/v1"
)

const sockPath string = "/tmp/mapper-1.sock"

func UnixConnect(context.Context, string) (net.Conn, error) {
	unixAddress, err := net.ResolveUnixAddr("unix", sockPath)
	conn, err := net.DialUnix("unix", nil, unixAddress)
	return conn, err
}

func main() {
	conn, err := grpc.Dial("/tmp/mapper-1.sock", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithContextDialer(UnixConnect))
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}

	defer conn.Close()

	c := pb.NewMapperClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	deviceName := "device-1"
	twins := []*pb.Twin{&pb.Twin{
		PropertyName: "temperature",
		Desired:      &pb.TwinProperty{
			Value:    "30",
			Metadata: nil,
		},
		Reported:     &pb.TwinProperty{
			Value:    "30",
			Metadata: nil,
		},
	}}
	config := pb.DeviceConfig{
		Model:  &pb.DeviceModel{
			Name: "device-model-1",
			Spec: &pb.DeviceModelSpec{
				Properties: nil,
				Commands:   nil,
			},
		},
		Device: &pb.Device{
			Name:   deviceName,
			Spec:   &pb.DeviceSpec{
				DeviceModelRef:   "device-model-1",
				Protocol:         nil,
				PropertyVisitors: nil,
			},
			Status: &pb.DeviceStatus{
				Twins: twins,
				State: "ok",
			},
		},
	}

	_, err = c.CreateDevice(ctx, &pb.CreateDeviceRequest{
		Config: &config,
	})
	if err != nil {
		log.Fatalf("fail to create device with err %v\n", err)
	}

	log.Printf("success to create device %s\n", deviceName)

	_, err = c.RemoveDevice(ctx, &pb.RemoveDeviceRequest{
		DeviceName: deviceName,
	})
	if err != nil {
		log.Fatalf("fail to remove device with err %v\n", err)
	}

	log.Printf("success to remove device %s\n", deviceName)

	_, err = c.UpdateDevice(ctx, &pb.UpdateDeviceRequest{
		DeviceName: deviceName,
		Config:     &config,
	})
	if err != nil {
		log.Fatalf("fail to update device with err %v\n", err)
	}

	log.Printf("success to update device %s\n", deviceName)

	_, err = c.UpdateDeviceStatus(ctx, &pb.UpdateDeviceStatusRequest{
		DeviceName:    deviceName,
		DesiredDevice: &pb.DeviceStatus{
			Twins: twins,
			State: "",
		},
	})
	if err != nil {
		log.Fatalf("fail to update device status with err %v\n", err)
	}

	log.Printf("success to update device status %s\n", deviceName)

	device, err := c.GetDevice(ctx, &pb.GetDeviceRequest{
		DeviceName: deviceName,
	})
	if err != nil {
		log.Fatalf("fail to get device with err %v\n", err)
	}

	log.Printf("success to get device %+v\n", device)
}