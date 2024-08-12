package controller

import (
	"context"
	"log"

	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/repository"
	"github.com/metalpoch/olt-blueprint/update/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddDevice(ip, community string, client *mongo.Client) {
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	u := usecase.DeviceUsecase(usecase.NewDeviceUsecase(repository.NewDeviceRepository(client)))
	if err := u.Create(ip, community); err != nil {
		log.Fatalf("error saving %s - %s: %v\n", ip, community, err)
	}
}

func GetDevices(client *mongo.Client) []*model.Device {
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	u := usecase.DeviceUsecase(usecase.NewDeviceUsecase(repository.NewDeviceRepository(client)))
	devices, err := u.FindAll()
	if err != nil {
		log.Fatalln("error searching for devices:", err)
	}
	return devices
}
