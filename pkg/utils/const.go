package utils

// Annotations
const (
	BestEffortGPUUUIDAnnotationKey = "anylearn.dmagine.github.com/besteffort_gpu_uuid"
)

// ResourceName
const (
	BestEffortGPU = "anylearn.thuml.ai/gpu-besteffort"
	GuaranteeGPU  = "anylearn.thuml.ai/gpu-guarantee"
)

// Env Switches
const (
	EnvDisableHealthChecks = "DP_DISABLE_HEALTHCHECKS"
	AllHealthChecks        = "xids"
	DeviceListEnvvar       = "NVIDIA_VISIBLE_DEVICES"
)

const (
	RetryTimes = 8
)
