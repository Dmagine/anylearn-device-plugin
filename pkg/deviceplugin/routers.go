package deviceplugin

import (
	"context"
	"path"
	"time"

	log "github.com/sirupsen/logrus"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

// Register registers the device plugin for the given resourceName with Kubelet.
func (p *AnylearnDevicePlugin) Register() error {
	conn, err := p.dial(pluginapi.KubeletSocket, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pluginapi.NewRegistrationClient(conn)
	reqt := &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     path.Base(p.socket),
		ResourceName: p.resourceName,
		Options: &pluginapi.DevicePluginOptions{
			GetPreferredAllocationAvailable: (p.allocatePolicy != nil),
		},
	}

	_, err = client.Register(context.Background(), reqt)
	if err != nil {
		return err
	}
	return nil
}

// ListAndWatch lists devices and update that list according to the health status
func (p *AnylearnDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	s.Send(&pluginapi.ListAndWatchResponse{Devices: p.apiDevices()})
	for {
		select {
		case <-p.controller.stopCh:
			return nil
		case d := <-p.healthCh:
			// FIXME: there is no way to recover from the Unhealthy state.
			d.health = pluginapi.Unhealthy
			log.WithFields(log.Fields{
				"Resource": p.resourceName,
				"Device":   d.id,
			}).Error("Device marked unhealthy")
			s.Send(&pluginapi.ListAndWatchResponse{Devices: p.apiDevices()})
		case <-p.takenCh:
			s.Send(&pluginapi.ListAndWatchResponse{Devices: p.apiDevices()})
		}
	}
}

// GetDevicePluginOptions returns the values of the optional settings for this plugin
func (p *AnylearnDevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	options := &pluginapi.DevicePluginOptions{
		GetPreferredAllocationAvailable: (p.allocatePolicy != nil),
	}
	return options, nil
}

// GetPreferredAllocation returns the preferred allocation from the set of devices specified in the request
func (p *AnylearnDevicePlugin) GetPreferredAllocation(ctx context.Context, r *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return &pluginapi.PreferredAllocationResponse{}, nil
}

// Allocate which return list of devices.
func (p *AnylearnDevicePlugin) Allocate(ctx context.Context, reqs *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	return &pluginapi.AllocateResponse{}, nil
}

// PreStartContainer is unimplemented for this plugin
func (p *AnylearnDevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}
