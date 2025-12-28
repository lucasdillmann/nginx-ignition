package client

type AvailableAppDTO struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	ActiveWorkloads WorkloadDTO `json:"active_workloads"`
}

type WorkloadDTO struct {
	UsedPorts []WorkloadPortDTO `json:"used_ports"`
}

type WorkloadPortDTO struct {
	Protocol      string        `json:"protocol"`
	HostPorts     []HostPortDTO `json:"host_ports"`
	ContainerPort int           `json:"container_port"`
}

type HostPortDTO struct {
	HostIP   string `json:"host_ip"`
	HostPort int    `json:"host_port"`
}
