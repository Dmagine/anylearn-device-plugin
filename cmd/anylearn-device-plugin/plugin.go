package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"time"

	"github.com/NVIDIA/go-gpuallocator/gpuallocator"
	"google.golang.org/grpc"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type AnylearnDevicePlugin struct {
	controller *AnylearnDeviceController

	resourceName string
	socket       string

	server *grpc.Server
	health chan *Device
	taken  chan *Device
	stop   chan interface{}
}

func (m *AnylearnDevicePlugin) initialize() {
	m.server = grpc.NewServer([]grpc.ServerOption{}...)
	m.health = make(chan *Device)
	m.taken = make(chan *Device)
	m.stop = make(chan interface{})
}

func (m *AnylearnDevicePlugin) cleanup() {
	close(m.stop)
	m.server = nil
	m.health = nil
	m.taken = nil
	m.stop = nil
}

func (m *AnylearnDevicePlugin) apiDevices() []*pluginapi.Device {
	var pdevs []*pluginapi.Device
	for _, d := range m.controller.cachedDevices {
		ad := buildAPIDevice(d)
		if m.resourceName == BestEffortGPU && d.gtaken {
			ad.Health = pluginapi.Unhealthy
		}
		pdevs = append(pdevs, ad)
	}
	return pdevs
}

// ListAndWatch lists devices and update that list according to the health status
func (m *AnylearnDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	s.Send(&pluginapi.ListAndWatchResponse{Devices: m.apiDevices()})

	for {
		select {
		case <-m.stop:
			return nil
		case d := <-m.health:
			// FIXME: there is no way to recover from the Unhealthy state.
			d.health = pluginapi.Unhealthy
			log.Printf("'%s' device marked unhealthy: %s", m.resourceName, d.id)
			s.Send(&pluginapi.ListAndWatchResponse{Devices: m.apiDevices()})
		case <-m.taken:
			s.Send(&pluginapi.ListAndWatchResponse{Devices: m.apiDevices()})
		}
	}
}

// Start starts the gRPC server, registers the device plugin with the Kubelet,
// and starts the device healthchecks.
func (m *AnylearnDevicePlugin) Start() error {
	m.initialize()

	err := m.Serve()
	if err != nil {
		log.Printf("Could not start device plugin for '%s': %s", m.resourceName, err)
		m.cleanup()
		return err
	}
	log.Printf("Starting to serve '%s' on %s", m.resourceName, m.socket)

	err = m.Register()
	if err != nil {
		log.Printf("Could not register device plugin: %s", err)
		m.Stop()
		return err
	}
	log.Printf("Registered device plugin for '%s' with Kubelet", m.resourceName)

	return nil
}

// Stop stops the gRPC server.
func (m *AnylearnDevicePlugin) Stop() error {
	if m == nil || m.server == nil {
		return nil
	}
	log.Printf("Stopping to serve '%s' on %s", m.resourceName, m.socket)
	m.server.Stop()
	if err := os.Remove(m.socket); err != nil && !os.IsNotExist(err) {
		return err
	}
	m.cleanup()
	return nil
}

// Serve starts the gRPC server of the device plugin.
func (m *AnylearnDevicePlugin) Serve() error {
	os.Remove(m.socket)
	sock, err := net.Listen("unix", m.socket)
	if err != nil {
		return err
	}

	pluginapi.RegisterDevicePluginServer(m.server, m)

	go func() {
		lastCrashTime := time.Now()
		restartCount := 0
		for {
			log.Printf("Starting GRPC server for '%s'", m.resourceName)
			err := m.server.Serve(sock)
			if err == nil {
				break
			}

			log.Printf("GRPC server for '%s' crashed with error: %v", m.resourceName, err)

			// restart if it has not been too often
			// i.e. if server has crashed more than 5 times and it didn't last more than one hour each time
			if restartCount > 5 {
				// quit
				log.Fatalf("GRPC server for '%s' has repeatedly crashed recently. Quitting", m.resourceName)
			}
			timeSinceLastCrash := time.Since(lastCrashTime).Seconds()
			lastCrashTime = time.Now()
			if timeSinceLastCrash > 3600 {
				// it has been one hour since the last crash.. reset the count
				// to reflect on the frequency
				restartCount = 1
			} else {
				restartCount++
			}
		}
	}()

	// Wait for server to start by launching a blocking connexion
	conn, err := m.dial(m.socket, 5*time.Second)
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}

// Register registers the device plugin for the given resourceName with Kubelet.
func (m *AnylearnDevicePlugin) Register() error {
	conn, err := m.dial(pluginapi.KubeletSocket, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pluginapi.NewRegistrationClient(conn)
	reqt := &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     path.Base(m.socket),
		ResourceName: m.resourceName,
		Options: &pluginapi.DevicePluginOptions{
			GetPreferredAllocationAvailable: (m.allocatePolicy != nil),
		},
	}

	_, err = client.Register(context.Background(), reqt)
	if err != nil {
		return err
	}
	return nil
}

// GetDevicePluginOptions returns the values of the optional settings for this plugin
func (m *AnylearnDevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	options := &pluginapi.DevicePluginOptions{
		GetPreferredAllocationAvailable: (m.allocatePolicy != nil),
	}
	return options, nil
}

// GetPreferredAllocation returns the preferred allocation from the set of devices specified in the request
func (m *AnylearnDevicePlugin) GetPreferredAllocation(ctx context.Context, r *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	response := &pluginapi.PreferredAllocationResponse{}
	for _, req := range r.ContainerRequests {
		available, err := gpuallocator.NewDevicesFrom(req.AvailableDeviceIDs)
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve list of available devices: %v", err)
		}

		required, err := gpuallocator.NewDevicesFrom(req.MustIncludeDeviceIDs)
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve list of required devices: %v", err)
		}

		allocated := m.allocatePolicy.Allocate(available, required, int(req.AllocationSize))

		var deviceIds []string
		for _, device := range allocated {
			deviceIds = append(deviceIds, device.UUID)
		}

		resp := &pluginapi.ContainerPreferredAllocationResponse{
			DeviceIDs: deviceIds,
		}

		response.ContainerResponses = append(response.ContainerResponses, resp)
	}
	return response, nil
}

// Allocate which return list of devices.
func (m *AnylearnDevicePlugin) Allocate(ctx context.Context, reqs *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	responses := pluginapi.AllocateResponse{}
	for _, req := range reqs.ContainerRequests {
		for _, id := range req.DevicesIDs {
			if !m.deviceExists(id) {
				return nil, fmt.Errorf("invalid allocation request for '%s': unknown device: %s", m.resourceName, id)
			}
		}

		response := pluginapi.ContainerAllocateResponse{}

		uuids := req.DevicesIDs
		deviceIDs := m.deviceIDsFromUUIDs(uuids)

		if deviceListStrategyFlag == DeviceListStrategyEnvvar {
			response.Envs = m.apiEnvs(m.deviceListEnvvar, deviceIDs)
		}
		if deviceListStrategyFlag == DeviceListStrategyVolumeMounts {
			response.Envs = m.apiEnvs(m.deviceListEnvvar, []string{deviceListAsVolumeMountsContainerPathRoot})
			response.Mounts = m.apiMounts(deviceIDs)
		}
		if passDeviceSpecsFlag {
			response.Devices = m.apiDeviceSpecs(nvidiaDriverRootFlag, uuids)
		}

		responses.ContainerResponses = append(responses.ContainerResponses, &response)
	}

	return &responses, nil
}

// PreStartContainer is unimplemented for this plugin
func (m *AnylearnDevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}

// dial establishes the gRPC communication with the registered device plugin.
func (m *AnylearnDevicePlugin) dial(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
	c, err := grpc.Dial(unixSocketPath, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(timeout),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)

	if err != nil {
		return nil, err
	}

	return c, nil
}
