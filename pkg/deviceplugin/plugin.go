package deviceplugin

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/NVIDIA/go-gpuallocator/gpuallocator"
	"github.com/dmagine/anylearn-device-plugin/pkg/utils"
	"google.golang.org/grpc"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type AnylearnDevicePlugin struct {
	controller *AnylearnDevicePluginController

	resourceName   string
	socket         string
	allocatePolicy gpuallocator.Policy

	server   *grpc.Server
	healthCh chan *Device
	takenCh  chan *Device
}

func (p *AnylearnDevicePlugin) initialize() {
	p.server = grpc.NewServer([]grpc.ServerOption{}...)
	p.healthCh = make(chan *Device)
	p.takenCh = make(chan *Device)
}

func (p *AnylearnDevicePlugin) cleanup() {
	p.server = nil
	p.healthCh = nil
	p.takenCh = nil
}

func (p *AnylearnDevicePlugin) apiDevices() []*pluginapi.Device {
	var pdevs []*pluginapi.Device
	for _, d := range p.controller.cachedDevices {
		ad := buildAPIDevice(d)
		if p.resourceName == utils.BestEffortGPU && d.gtaken {
			ad.Health = pluginapi.Unhealthy
		}
		pdevs = append(pdevs, ad)
	}
	return pdevs
}

// dial establishes the gRPC communication with the registered device plugin.
func (p *AnylearnDevicePlugin) dial(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
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

// Start starts the gRPC server, registers the device plugin with the Kubelet,
// and starts the device healthchecks.
func (p *AnylearnDevicePlugin) Start() error {
	p.initialize()

	err := p.Serve()
	if err != nil {
		log.Printf("Could not start device plugin for '%s': %s", p.resourceName, err)
		p.cleanup()
		return err
	}
	log.Printf("Starting to serve '%s' on %s", p.resourceName, p.socket)

	err = p.Register()
	if err != nil {
		log.Printf("Could not register device plugin: %s", err)
		p.Stop()
		return err
	}
	log.Printf("Registered device plugin for '%s' with Kubelet", p.resourceName)

	return nil
}

// Stop stops the gRPC server.
func (p *AnylearnDevicePlugin) Stop() error {
	if p == nil || p.server == nil {
		return nil
	}
	log.Printf("Stopping to serve '%s' on %s", p.resourceName, p.socket)
	p.server.Stop()
	if err := os.Remove(p.socket); err != nil && !os.IsNotExist(err) {
		return err
	}
	p.cleanup()
	return nil
}

// Serve starts the gRPC server of the device plugin.
func (p *AnylearnDevicePlugin) Serve() error {
	os.Remove(p.socket)
	sock, err := net.Listen("unix", p.socket)
	if err != nil {
		return err
	}

	pluginapi.RegisterDevicePluginServer(p.server, p)

	go func() {
		lastCrashTime := time.Now()
		restartCount := 0
		for {
			log.Printf("Starting GRPC server for '%s'", p.resourceName)
			err := p.server.Serve(sock)
			if err == nil {
				break
			}

			log.Printf("GRPC server for '%s' crashed with error: %v", p.resourceName, err)

			// restart if it has not been too often
			// i.e. if server has crashed more than 5 times and it didn't last more than one hour each time
			if restartCount > 5 {
				// quit
				log.Fatalf("GRPC server for '%s' has repeatedly crashed recently. Quitting", p.resourceName)
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
	conn, err := p.dial(p.socket, 5*time.Second)
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}
