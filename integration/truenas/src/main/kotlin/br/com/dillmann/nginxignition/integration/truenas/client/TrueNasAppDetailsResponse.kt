package br.com.dillmann.nginxignition.integration.truenas.client

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class TrueNasAppDetailsResponse(
    val id: String,
    val name: String,
    @SerialName("active_workloads")
    val activeWorkloads: Workload,
) {
    @Serializable
    data class Workload(
        @SerialName("used_ports")
        val usedPorts: List<WorkloadPort>,
    )

    @Serializable
    data class WorkloadPort(
        val protocol: String,
        @SerialName("container_port")
        val containerPort: Int,
        @SerialName("host_ports")
        val hostPorts: List<HostPort>,
    )

    @Serializable
    data class HostPort(
        @SerialName("host_port")
        val hostPort: Int,
        @SerialName("host_ip")
        val hostIp: String,
    )
}
