package main

import (
	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type Device struct {
	device  *nvml.Device
	id      string
	path    string
	index   string
	gtaken  bool
	betaken bool
	health  string
}

func buildDevice(d *nvml.Device, path string, index string) *Device {
	dev := Device{}
	dev.device = d
	dev.id = d.UUID
	dev.path = path
	dev.index = index
	dev.health = pluginapi.Healthy
	return &dev
}

func buildAPIDevice(d *Device) *pluginapi.Device {
	dev := &pluginapi.Device{}
	dev.ID = d.id
	dev.Health = d.health
	if d.device.CPUAffinity != nil {
		dev.Topology = &pluginapi.TopologyInfo{
			Nodes: []*pluginapi.NUMANode{
				{
					ID: int64(*(d.device.CPUAffinity)),
				},
			},
		}
	}
	return dev
}
