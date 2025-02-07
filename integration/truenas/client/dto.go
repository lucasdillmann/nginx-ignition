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
	ContainerPort int           `json:"container_port"`
	HostPorts     []HostPortDTO `json:"host_ports"`
}

type HostPortDTO struct {
	HostPort int    `json:"host_port"`
	HostIp   string `json:"host_ip"`
}
