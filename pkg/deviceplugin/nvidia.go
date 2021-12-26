package deviceplugin

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"github.com/dmagine/anylearn-device-plugin/pkg/utils"

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

func (controller *AnylearnDevicePluginController) Devices() []*Device {
	n, err := nvml.GetDeviceCount()
	utils.FatalWhenError(err)

	var devs []*Device
	for i := uint(0); i < n; i++ {
		d, err := nvml.NewDeviceLite(i)
		utils.FatalWhenError(err)
		devs = append(devs, buildDevice(d, d.Path, fmt.Sprintf("%v", i)))
	}

	return devs
}

func CheckHealth(stop <-chan interface{}, devices []*Device, unhealthy []chan<- *Device) {
	disableHealthChecks := strings.ToLower(os.Getenv(utils.EnvDisableHealthChecks))
	if disableHealthChecks == "all" {
		disableHealthChecks = utils.AllHealthChecks
	}
	if strings.Contains(disableHealthChecks, "xids") {
		return
	}

	// FIXME: formalize the full list and document it.
	// http://docs.nvidia.com/deploy/xid-errors/index.html#topic_4
	// Application errors: the GPU should still be healthy
	applicationErrorXids := []uint64{
		13, // Graphics Engine Exception
		31, // GPU memory page fault
		43, // GPU stopped processing
		45, // Preemptive cleanup, due to previous errors
		68, // Video processor exception
	}

	skippedXids := make(map[uint64]bool)
	for _, id := range applicationErrorXids {
		skippedXids[id] = true
	}

	for _, additionalXid := range getAdditionalXids(disableHealthChecks) {
		skippedXids[additionalXid] = true
	}

	eventSet := nvml.NewEventSet()
	defer nvml.DeleteEventSet(eventSet)

	for _, d := range devices {
		gpu, _, _, err := nvml.ParseMigDeviceUUID(d.id)
		if err != nil {
			gpu = d.id
		}

		err = nvml.RegisterEventForDevice(eventSet, nvml.XidCriticalError, gpu)
		if err != nil && strings.HasSuffix(err.Error(), "Not Supported") {
			log.Printf("Warning: %s is too old to support healthchecking: %s. Marking it unhealthy.", d.id, err)
			for _, ch := range unhealthy {
				ch <- d
			}
			continue
		}
		utils.FatalWhenError(err)
	}

	for {
		select {
		case <-stop:
			return
		default:
		}

		e, err := nvml.WaitForEvent(eventSet, 5000)
		if err != nil && e.Etype != nvml.XidCriticalError {
			continue
		}

		if skippedXids[e.Edata] {
			continue
		}

		if e.UUID == nil || len(*e.UUID) == 0 {
			// All devices are unhealthy
			log.Printf("XidCriticalError: Xid=%d, All devices will go unhealthy.", e.Edata)
			for _, d := range devices {
				for _, ch := range unhealthy {
					ch <- d
				}
			}
			continue
		}

		for _, d := range devices {
			// Please see https://github.com/NVIDIA/gpu-monitoring-tools/blob/148415f505c96052cb3b7fdf443b34ac853139ec/bindings/go/nvml/nvml.h#L1424
			// for the rationale why gi and ci can be set as such when the UUID is a full GPU UUID and not a MIG device UUID.
			gpu, gi, ci, err := nvml.ParseMigDeviceUUID(d.id)
			if err != nil {
				gpu = d.id
				gi = 0xFFFFFFFF
				ci = 0xFFFFFFFF
			}

			if gpu == *e.UUID && gi == *e.GpuInstanceId && ci == *e.ComputeInstanceId {
				log.Printf("XidCriticalError: Xid=%d on Device=%s, the device will go unhealthy.", e.Edata, d.id)
				for _, ch := range unhealthy {
					ch <- d
				}
			}
		}
	}
}

// getAdditionalXids returns a list of additional Xids to skip from the specified string.
// The input is treaded as a comma-separated string and all valid uint64 values are considered as Xid values. Invalid values
// are ignored.
func getAdditionalXids(input string) []uint64 {
	if input == "" {
		return nil
	}

	var additionalXids []uint64
	for _, additionalXid := range strings.Split(input, ",") {
		trimmed := strings.TrimSpace(additionalXid)
		if trimmed == "" {
			continue
		}
		xid, err := strconv.ParseUint(trimmed, 10, 64)
		if err != nil {
			log.Printf("Ignoring malformed Xid value %v: %v", trimmed, err)
			continue
		}
		additionalXids = append(additionalXids, xid)
	}

	return additionalXids
}