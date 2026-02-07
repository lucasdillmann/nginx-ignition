package nginx

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/nginx"
)

func Test_toTrafficStatsResponseDTO(t *testing.T) {
	t.Run("converts complete stats", func(t *testing.T) {
		stats := &nginx.Stats{
			HostName: "test-host",
			Connections: nginx.StatsConnections{
				Active:   10,
				Reading:  2,
				Writing:  3,
				Waiting:  5,
				Accepted: 100,
				Handled:  99,
				Requests: 150,
			},
			ServerZones: map[string]nginx.StatsZoneData{
				"zone1": {
					RequestCounter: 50,
					InBytes:        1024,
					OutBytes:       2048,
				},
			},
			FilterZones: map[string]map[string]nginx.StatsZoneData{
				"filter1": {
					"subzone1": {
						RequestCounter: 25,
						InBytes:        512,
						OutBytes:       1024,
					},
				},
			},
			UpstreamZones: map[string][]nginx.StatsUpstreamZoneData{
				"upstream1": {
					{
						Server:         "192.168.1.1:8080",
						RequestCounter: 75,
						InBytes:        3072,
						OutBytes:       6144,
					},
				},
			},
		}

		result := toTrafficStatsResponseDTO(stats)

		assert.Equal(t, "test-host", result.HostName)
		assert.Equal(t, uint64(10), result.Connections.Active)
		assert.Len(t, result.ServerZones, 1)
		assert.Len(t, result.FilterZones, 1)
		assert.Len(t, result.UpstreamZones, 1)
	})

	t.Run("returns empty DTO when stats is nil", func(t *testing.T) {
		result := toTrafficStatsResponseDTO(nil)

		assert.Empty(t, result.HostName)
		assert.Empty(t, result.ServerZones)
		assert.Empty(t, result.FilterZones)
		assert.Empty(t, result.UpstreamZones)
	})
}

func Test_toConnectionsDTO(t *testing.T) {
	t.Run("converts all connection fields", func(t *testing.T) {
		connections := nginx.StatsConnections{
			Active:   10,
			Reading:  2,
			Writing:  3,
			Waiting:  5,
			Accepted: 100,
			Handled:  99,
			Requests: 150,
		}

		result := toConnectionsDTO(connections)

		assert.Equal(t, uint64(10), result.Active)
		assert.Equal(t, uint64(2), result.Reading)
		assert.Equal(t, uint64(3), result.Writing)
		assert.Equal(t, uint64(5), result.Waiting)
		assert.Equal(t, uint64(100), result.Accepted)
		assert.Equal(t, uint64(99), result.Handled)
		assert.Equal(t, uint64(150), result.Requests)
	})

	t.Run("converts zero values", func(t *testing.T) {
		connections := nginx.StatsConnections{}

		result := toConnectionsDTO(connections)

		assert.Equal(t, uint64(0), result.Active)
		assert.Equal(t, uint64(0), result.Requests)
	})
}

func Test_toServerZonesDTO(t *testing.T) {
	t.Run("converts server zones map", func(t *testing.T) {
		serverZones := map[string]nginx.StatsZoneData{
			"zone1": {
				RequestCounter: 100,
				InBytes:        1024,
				OutBytes:       2048,
			},
			"zone2": {
				RequestCounter: 200,
				InBytes:        2048,
				OutBytes:       4096,
			},
		}

		result := toServerZonesDTO(serverZones)

		assert.Len(t, result, 2)
		assert.Equal(t, uint64(100), result["zone1"].RequestCounter)
		assert.Equal(t, uint64(1024), result["zone1"].InBytes)
		assert.Equal(t, uint64(200), result["zone2"].RequestCounter)
	})

	t.Run("returns nil when input is nil", func(t *testing.T) {
		result := toServerZonesDTO(nil)
		assert.Nil(t, result)
	})

	t.Run("returns empty map for empty input", func(t *testing.T) {
		result := toServerZonesDTO(make(map[string]nginx.StatsZoneData))
		assert.Empty(t, result)
	})
}

func Test_toFilterZonesDTO(t *testing.T) {
	t.Run("converts nested filter zones", func(t *testing.T) {
		filterZones := map[string]map[string]nginx.StatsZoneData{
			"filter1": {
				"subzone1": {
					RequestCounter: 50,
					InBytes:        512,
					OutBytes:       1024,
				},
				"subzone2": {
					RequestCounter: 75,
					InBytes:        768,
					OutBytes:       1536,
				},
			},
			"filter2": {
				"subzone3": {
					RequestCounter: 100,
					InBytes:        1024,
					OutBytes:       2048,
				},
			},
		}

		result := toFilterZonesDTO(filterZones)

		assert.Len(t, result, 2)
		assert.Len(t, result["filter1"], 2)
		assert.Len(t, result["filter2"], 1)
		assert.Equal(t, uint64(50), result["filter1"]["subzone1"].RequestCounter)
		assert.Equal(t, uint64(100), result["filter2"]["subzone3"].RequestCounter)
	})

	t.Run("handles nil nested maps", func(t *testing.T) {
		filterZones := map[string]map[string]nginx.StatsZoneData{
			"filter1": nil,
			"filter2": {
				"subzone1": {
					RequestCounter: 25,
				},
			},
		}

		result := toFilterZonesDTO(filterZones)

		assert.Len(t, result, 2)
		assert.Nil(t, result["filter1"])
		assert.NotNil(t, result["filter2"])
		assert.Len(t, result["filter2"], 1)
	})

	t.Run("returns nil when input is nil", func(t *testing.T) {
		result := toFilterZonesDTO(nil)
		assert.Nil(t, result)
	})
}

func Test_toUpstreamZonesDTO(t *testing.T) {
	t.Run("converts upstream zones with multiple servers", func(t *testing.T) {
		upstreamZones := map[string][]nginx.StatsUpstreamZoneData{
			"upstream1": {
				{
					Server:         "192.168.1.1:8080",
					RequestCounter: 100,
					InBytes:        2048,
					OutBytes:       4096,
					Weight:         5,
					MaxFails:       3,
					FailTimeout:    10,
					Backup:         false,
					Down:           false,
				},
				{
					Server:         "192.168.1.2:8080",
					RequestCounter: 150,
					InBytes:        3072,
					OutBytes:       6144,
					Weight:         10,
					MaxFails:       2,
					FailTimeout:    5,
					Backup:         true,
					Down:           false,
				},
			},
		}

		result := toUpstreamZonesDTO(upstreamZones)

		assert.Len(t, result, 1)
		assert.Len(t, result["upstream1"], 2)
		assert.Equal(t, "192.168.1.1:8080", result["upstream1"][0].Server)
		assert.Equal(t, uint64(100), result["upstream1"][0].RequestCounter)
		assert.Equal(t, 5, result["upstream1"][0].Weight)
		assert.False(t, result["upstream1"][0].Backup)
		assert.Equal(t, "192.168.1.2:8080", result["upstream1"][1].Server)
		assert.True(t, result["upstream1"][1].Backup)
	})

	t.Run("handles nil slices", func(t *testing.T) {
		upstreamZones := map[string][]nginx.StatsUpstreamZoneData{
			"upstream1": nil,
			"upstream2": {
				{
					Server:         "192.168.1.1:8080",
					RequestCounter: 50,
				},
			},
		}

		result := toUpstreamZonesDTO(upstreamZones)

		assert.Len(t, result, 2)
		assert.Nil(t, result["upstream1"])
		assert.NotNil(t, result["upstream2"])
		assert.Len(t, result["upstream2"], 1)
	})

	t.Run("returns nil when input is nil", func(t *testing.T) {
		result := toUpstreamZonesDTO(nil)
		assert.Nil(t, result)
	})
}

func Test_toZoneDataDTO(t *testing.T) {
	t.Run("converts all zone data fields", func(t *testing.T) {
		data := nginx.StatsZoneData{
			RequestCounter:     100,
			InBytes:            1024,
			OutBytes:           2048,
			RequestMsec:        50,
			RequestMsecCounter: 5000,
			Responses: nginx.StatsResponses{
				Status2xx: 80,
				Status4xx: 15,
				Status5xx: 5,
			},
			RequestMsecs: nginx.StatsTimeSeries{
				Times: []int64{1000, 2000},
				Msecs: []int64{10, 20},
			},
			RequestBuckets: nginx.StatsBuckets{
				Msecs:    []int64{100, 200},
				Counters: []int64{5, 10},
			},
			OverCounts: nginx.StatsOverCounts{
				RequestCounter: 50,
				InBytes:        512,
			},
		}

		result := toZoneDataDTO(data)

		assert.Equal(t, uint64(100), result.RequestCounter)
		assert.Equal(t, uint64(1024), result.InBytes)
		assert.Equal(t, uint64(2048), result.OutBytes)
		assert.Equal(t, uint64(50), result.RequestMsec)
		assert.Equal(t, uint64(5000), result.RequestMsecCounter)
		assert.Equal(t, uint64(80), result.Responses.Status2xx)
		assert.Len(t, result.RequestMsecs.Times, 2)
		assert.Len(t, result.RequestBuckets.Msecs, 2)
		assert.Equal(t, uint64(50), result.OverCounts.RequestCounter)
	})
}

func Test_toUpstreamZoneDataDTO(t *testing.T) {
	t.Run("converts all upstream zone data fields", func(t *testing.T) {
		data := nginx.StatsUpstreamZoneData{
			Server:              "192.168.1.1:8080",
			RequestCounter:      100,
			InBytes:             2048,
			OutBytes:            4096,
			RequestMsec:         25,
			RequestMsecCounter:  2500,
			ResponseMsec:        30,
			ResponseMsecCounter: 3000,
			Weight:              5,
			MaxFails:            3,
			FailTimeout:         10,
			Backup:              false,
			Down:                false,
			Responses: nginx.StatsUpstreamResponses{
				Status2xx: 90,
				Status5xx: 10,
			},
			RequestMsecs: nginx.StatsTimeSeries{
				Times: []int64{1000},
				Msecs: []int64{25},
			},
			ResponseMsecs: nginx.StatsTimeSeries{
				Times: []int64{1000},
				Msecs: []int64{30},
			},
			RequestBuckets: nginx.StatsBuckets{
				Msecs:    []int64{100},
				Counters: []int64{5},
			},
			ResponseBuckets: nginx.StatsBuckets{
				Msecs:    []int64{100},
				Counters: []int64{10},
			},
			OverCounts: nginx.StatsOverCounts{
				RequestCounter: 25,
			},
		}

		result := toUpstreamZoneDataDTO(data)

		assert.Equal(t, "192.168.1.1:8080", result.Server)
		assert.Equal(t, uint64(100), result.RequestCounter)
		assert.Equal(t, uint64(2048), result.InBytes)
		assert.Equal(t, uint64(4096), result.OutBytes)
		assert.Equal(t, uint64(25), result.RequestMsec)
		assert.Equal(t, uint64(30), result.ResponseMsec)
		assert.Equal(t, 5, result.Weight)
		assert.Equal(t, 3, result.MaxFails)
		assert.Equal(t, 10, result.FailTimeout)
		assert.False(t, result.Backup)
		assert.False(t, result.Down)
		assert.Equal(t, uint64(90), result.Responses.Status2xx)
		assert.Len(t, result.RequestMsecs.Times, 1)
		assert.Len(t, result.ResponseMsecs.Times, 1)
	})
}

func Test_toResponsesDTO(t *testing.T) {
	t.Run("converts all response status codes", func(t *testing.T) {
		responses := nginx.StatsResponses{
			Status1xx:   1,
			Status2xx:   100,
			Status3xx:   20,
			Status4xx:   15,
			Status5xx:   5,
			Miss:        10,
			Bypass:      5,
			Expired:     3,
			Stale:       2,
			Updating:    1,
			Revalidated: 4,
			Hit:         80,
			Scarce:      1,
		}

		result := toResponsesDTO(responses)

		assert.Equal(t, uint64(1), result.Status1xx)
		assert.Equal(t, uint64(100), result.Status2xx)
		assert.Equal(t, uint64(20), result.Status3xx)
		assert.Equal(t, uint64(15), result.Status4xx)
		assert.Equal(t, uint64(5), result.Status5xx)
		assert.Equal(t, uint64(10), result.Miss)
		assert.Equal(t, uint64(5), result.Bypass)
		assert.Equal(t, uint64(80), result.Hit)
	})
}

func Test_toUpstreamResponsesDTO(t *testing.T) {
	t.Run("converts all upstream response status codes", func(t *testing.T) {
		responses := nginx.StatsUpstreamResponses{
			Status1xx: 2,
			Status2xx: 150,
			Status3xx: 25,
			Status4xx: 10,
			Status5xx: 3,
		}

		result := toUpstreamResponsesDTO(responses)

		assert.Equal(t, uint64(2), result.Status1xx)
		assert.Equal(t, uint64(150), result.Status2xx)
		assert.Equal(t, uint64(25), result.Status3xx)
		assert.Equal(t, uint64(10), result.Status4xx)
		assert.Equal(t, uint64(3), result.Status5xx)
	})
}

func Test_toTimeSeriesDTO(t *testing.T) {
	t.Run("converts time series", func(t *testing.T) {
		timeSeries := nginx.StatsTimeSeries{
			Times: []int64{1000, 2000, 3000},
			Msecs: []int64{10, 20, 30},
		}

		result := toTimeSeriesDTO(timeSeries)

		assert.Equal(t, []int64{1000, 2000, 3000}, result.Times)
		assert.Equal(t, []int64{10, 20, 30}, result.Msecs)
	})

	t.Run("converts empty time series", func(t *testing.T) {
		timeSeries := nginx.StatsTimeSeries{
			Times: []int64{},
			Msecs: []int64{},
		}

		result := toTimeSeriesDTO(timeSeries)

		assert.Empty(t, result.Times)
		assert.Empty(t, result.Msecs)
	})
}

func Test_toBucketsDTO(t *testing.T) {
	t.Run("converts buckets", func(t *testing.T) {
		buckets := nginx.StatsBuckets{
			Msecs:    []int64{100, 200, 300},
			Counters: []int64{5, 10, 15},
		}

		result := toBucketsDTO(buckets)

		assert.Equal(t, []int64{100, 200, 300}, result.Msecs)
		assert.Equal(t, []int64{5, 10, 15}, result.Counters)
	})
}

func Test_toOverCountsDTO(t *testing.T) {
	t.Run("converts all over count fields", func(t *testing.T) {
		overCounts := nginx.StatsOverCounts{
			RequestCounter:      100,
			InBytes:             1024,
			OutBytes:            2048,
			Status1xx:           1,
			Status2xx:           80,
			Status3xx:           10,
			Status4xx:           5,
			Status5xx:           4,
			Miss:                15,
			Bypass:              8,
			Expired:             3,
			Stale:               2,
			Updating:            1,
			Revalidated:         5,
			Hit:                 60,
			Scarce:              1,
			RequestMsecCounter:  5000,
			ResponseMsecCounter: 6000,
		}

		result := toOverCountsDTO(overCounts)

		assert.Equal(t, uint64(100), result.RequestCounter)
		assert.Equal(t, uint64(1024), result.InBytes)
		assert.Equal(t, uint64(2048), result.OutBytes)
		assert.Equal(t, uint64(1), result.Status1xx)
		assert.Equal(t, uint64(80), result.Status2xx)
		assert.Equal(t, uint64(10), result.Status3xx)
		assert.Equal(t, uint64(5), result.Status4xx)
		assert.Equal(t, uint64(4), result.Status5xx)
		assert.Equal(t, uint64(15), result.Miss)
		assert.Equal(t, uint64(8), result.Bypass)
		assert.Equal(t, uint64(3), result.Expired)
		assert.Equal(t, uint64(2), result.Stale)
		assert.Equal(t, uint64(1), result.Updating)
		assert.Equal(t, uint64(5), result.Revalidated)
		assert.Equal(t, uint64(60), result.Hit)
		assert.Equal(t, uint64(1), result.Scarce)
		assert.Equal(t, uint64(5000), result.RequestMsecCounter)
		assert.Equal(t, uint64(6000), result.ResponseMsecCounter)
	})
}
