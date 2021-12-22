package main

import "google.golang.org/grpc"

type AnylearnDevicePlugin struct {
	resourceName     string
	deviceListEnvvar string
	socket           string

	server        *grpc.Server
	cachedDevices []*Device
	health        chan *Device
	stop          chan interface{}
}

func (plugin *AnylearnDevicePlugin) Start() error {
	return nil
}

func (plugin *AnylearnDevicePlugin) Stop() {
}

func (plugin *AnylearnDevicePlugin) Devices() []*Device {
	return nil
}
