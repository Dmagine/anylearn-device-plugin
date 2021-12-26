package deviceplugin

import (
	"net"
	"os"
	"time"

	"github.com/NVIDIA/go-gpuallocator/gpuallocator"
	"github.com/dmagine/anylearn-device-plugin/pkg/utils"
	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
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
		log.WithError(err).WithField("Resource", p.resourceName).Error("Could not start device plugin.")
		p.cleanup()
		return err
	}
	log.WithFields(log.Fields{
		"Resource": p.resourceName,
		"Socket":   p.socket,
	}).Info("Starting to serve")

	err = p.Register()
	if err != nil {
		log.WithError(err).Error("Could not register device plugin")
		p.Stop()
		return err
	}
	log.WithField("Resource", p.resourceName).Info("Registered device plugin with Kubelet")

	return nil
}

// Stop stops the gRPC server.
func (p *AnylearnDevicePlugin) Stop() error {
	if p == nil || p.server == nil {
		return nil
	}
	log.WithFields(log.Fields{
		"Resource": p.resourceName,
		"Socket":   p.socket,
	}).Info("Stopping to serve")
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
			log.WithField("Resource", p.resourceName).Info("Starting GRPC server")
			err := p.server.Serve(sock)
			if err == nil {
				break
			}

			log.WithError(err).WithField("Resource", p.resourceName).Error("GRPC server crashed with error.")

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
