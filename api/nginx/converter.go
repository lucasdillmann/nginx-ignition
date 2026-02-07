package nginx

import "dillmann.com.br/nginx-ignition/core/nginx"

func toTrafficStatsResponseDTO(stats *nginx.Stats) trafficStatsResponseDTO {
	if stats == nil {
		return trafficStatsResponseDTO{}
	}

	return trafficStatsResponseDTO{
		HostName:      stats.HostName,
		Connections:   toConnectionsDTO(stats.Connections),
		ServerZones:   toServerZonesDTO(stats.ServerZones),
		FilterZones:   toFilterZonesDTO(stats.FilterZones),
		UpstreamZones: toUpstreamZonesDTO(stats.UpstreamZones),
	}
}

func toConnectionsDTO(connections nginx.StatsConnections) trafficStatsConnectionsDTO {
	return trafficStatsConnectionsDTO{
		Active:   connections.Active,
		Reading:  connections.Reading,
		Writing:  connections.Writing,
		Waiting:  connections.Waiting,
		Accepted: connections.Accepted,
		Handled:  connections.Handled,
		Requests: connections.Requests,
	}
}

func toServerZonesDTO(
	serverZones map[string]nginx.StatsZoneData,
) map[string]trafficStatsZoneDataDTO {
	if serverZones == nil {
		return nil
	}

	result := make(map[string]trafficStatsZoneDataDTO, len(serverZones))
	for k, v := range serverZones {
		result[k] = toZoneDataDTO(v)
	}
	return result
}

func toFilterZonesDTO(
	filterZones map[string]map[string]nginx.StatsZoneData,
) map[string]map[string]trafficStatsZoneDataDTO {
	if filterZones == nil {
		return nil
	}

	result := make(map[string]map[string]trafficStatsZoneDataDTO, len(filterZones))
	for k, v := range filterZones {
		if v == nil {
			result[k] = nil
			continue
		}
		inner := make(map[string]trafficStatsZoneDataDTO, len(v))
		for ik, iv := range v {
			inner[ik] = toZoneDataDTO(iv)
		}
		result[k] = inner
	}
	return result
}

func toUpstreamZonesDTO(
	upstreamZones map[string][]nginx.StatsUpstreamZoneData,
) map[string][]trafficStatsUpstreamZoneDataDTO {
	if upstreamZones == nil {
		return nil
	}

	result := make(map[string][]trafficStatsUpstreamZoneDataDTO, len(upstreamZones))
	for k, v := range upstreamZones {
		if v == nil {
			result[k] = nil
			continue
		}
		arr := make([]trafficStatsUpstreamZoneDataDTO, len(v))
		for i, item := range v {
			arr[i] = toUpstreamZoneDataDTO(item)
		}
		result[k] = arr
	}
	return result
}

func toZoneDataDTO(data nginx.StatsZoneData) trafficStatsZoneDataDTO {
	return trafficStatsZoneDataDTO{
		RequestCounter:     data.RequestCounter,
		InBytes:            data.InBytes,
		OutBytes:           data.OutBytes,
		Responses:          toResponsesDTO(data.Responses),
		RequestMsec:        data.RequestMsec,
		RequestMsecCounter: data.RequestMsecCounter,
		RequestMsecs:       toTimeSeriesDTO(data.RequestMsecs),
		RequestBuckets:     toBucketsDTO(data.RequestBuckets),
		OverCounts:         toOverCountsDTO(data.OverCounts),
	}
}

func toUpstreamZoneDataDTO(data nginx.StatsUpstreamZoneData) trafficStatsUpstreamZoneDataDTO {
	return trafficStatsUpstreamZoneDataDTO{
		Server:              data.Server,
		RequestCounter:      data.RequestCounter,
		InBytes:             data.InBytes,
		OutBytes:            data.OutBytes,
		Responses:           toUpstreamResponsesDTO(data.Responses),
		RequestMsec:         data.RequestMsec,
		RequestMsecCounter:  data.RequestMsecCounter,
		RequestMsecs:        toTimeSeriesDTO(data.RequestMsecs),
		RequestBuckets:      toBucketsDTO(data.RequestBuckets),
		ResponseMsec:        data.ResponseMsec,
		ResponseMsecCounter: data.ResponseMsecCounter,
		ResponseMsecs:       toTimeSeriesDTO(data.ResponseMsecs),
		ResponseBuckets:     toBucketsDTO(data.ResponseBuckets),
		Weight:              data.Weight,
		MaxFails:            data.MaxFails,
		FailTimeout:         data.FailTimeout,
		Backup:              data.Backup,
		Down:                data.Down,
		OverCounts:          toOverCountsDTO(data.OverCounts),
	}
}

func toResponsesDTO(responses nginx.StatsResponses) trafficStatsResponsesDTO {
	return trafficStatsResponsesDTO{
		Status1xx:   responses.Status1xx,
		Status2xx:   responses.Status2xx,
		Status3xx:   responses.Status3xx,
		Status4xx:   responses.Status4xx,
		Status5xx:   responses.Status5xx,
		Miss:        responses.Miss,
		Bypass:      responses.Bypass,
		Expired:     responses.Expired,
		Stale:       responses.Stale,
		Updating:    responses.Updating,
		Revalidated: responses.Revalidated,
		Hit:         responses.Hit,
		Scarce:      responses.Scarce,
	}
}

func toUpstreamResponsesDTO(
	responses nginx.StatsUpstreamResponses,
) trafficStatsUpstreamResponsesDTO {
	return trafficStatsUpstreamResponsesDTO{
		Status1xx: responses.Status1xx,
		Status2xx: responses.Status2xx,
		Status3xx: responses.Status3xx,
		Status4xx: responses.Status4xx,
		Status5xx: responses.Status5xx,
	}
}

func toTimeSeriesDTO(timeSeries nginx.StatsTimeSeries) trafficStatsTimeSeriesDTO {
	return trafficStatsTimeSeriesDTO{
		Times: timeSeries.Times,
		Msecs: timeSeries.Msecs,
	}
}

func toBucketsDTO(buckets nginx.StatsBuckets) trafficStatsBucketsDTO {
	return trafficStatsBucketsDTO{
		Msecs:    buckets.Msecs,
		Counters: buckets.Counters,
	}
}

func toOverCountsDTO(overCounts nginx.StatsOverCounts) trafficStatsOverCountsDTO {
	return trafficStatsOverCountsDTO{
		RequestCounter:      overCounts.RequestCounter,
		InBytes:             overCounts.InBytes,
		OutBytes:            overCounts.OutBytes,
		Status1xx:           overCounts.Status1xx,
		Status2xx:           overCounts.Status2xx,
		Status3xx:           overCounts.Status3xx,
		Status4xx:           overCounts.Status4xx,
		Status5xx:           overCounts.Status5xx,
		Miss:                overCounts.Miss,
		Bypass:              overCounts.Bypass,
		Expired:             overCounts.Expired,
		Stale:               overCounts.Stale,
		Updating:            overCounts.Updating,
		Revalidated:         overCounts.Revalidated,
		Hit:                 overCounts.Hit,
		Scarce:              overCounts.Scarce,
		RequestMsecCounter:  overCounts.RequestMsecCounter,
		ResponseMsecCounter: overCounts.ResponseMsecCounter,
	}
}
