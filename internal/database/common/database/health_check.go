package database

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/healthcheck"
)

type healthCheckProvider struct {
	database *Database
}

func registerHealthCheck(d *Database, healthCheck *healthcheck.HealthCheck) {
	healthCheck.Register(&healthCheckProvider{
		database: d,
	})
}

func (d *healthCheckProvider) ID() string {
	return "database"
}

func (d *healthCheckProvider) Check(ctx context.Context) error {
	return d.database.db.PingContext(ctx)
}
