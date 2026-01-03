package certificate

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

type repository struct {
	database *database.Database
}

func New(db *database.Database) certificate.Repository {
	return &repository{
		database: db,
	}
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*certificate.Certificate, error) {
	var model certificateModel

	err := r.database.Select().
		Model(&model).
		Where(constants.ByIDFilter, id).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	//nolint:revive
	if result, err := toDomain(&model); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (r *repository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.database.Select().
		Model((*certificateModel)(nil)).
		Where(constants.ByIDFilter, id).
		Exists(ctx)
}

func (r *repository) InUseByID(ctx context.Context, id uuid.UUID) (bool, error) {
	exists, err := r.database.
		Select().
		Table("host_binding").
		Where("certificate_id = ?", id).
		Exists(ctx)

	if err != nil || exists {
		return exists, err
	}

	return r.database.
		Select().
		Table("settings_global_binding").
		Where("certificate_id = ?", id).
		Exists(ctx)
}

func (r *repository) GetAutoRenewSettings(
	ctx context.Context,
) (*certificate.AutoRenewSettings, error) {
	var enabled bool
	var intervalUnit string
	var intervalUnitCount int

	err := r.database.
		Select().
		Column("enabled", "interval_unit", "interval_unit_count").
		Table("settings_certificate_auto_renew").
		Limit(1).
		Scan(ctx, &enabled, &intervalUnit, &intervalUnitCount)
	if err != nil {
		return nil, err
	}

	return &certificate.AutoRenewSettings{
		Enabled:           enabled,
		IntervalUnit:      intervalUnit,
		IntervalUnitCount: intervalUnitCount,
	}, nil
}

func (r *repository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	_, err = transaction.NewDelete().
		Model((*certificateModel)(nil)).
		Where(constants.ByIDFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) Save(ctx context.Context, cert *certificate.Certificate) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	model, err := toModel(cert)
	if err != nil {
		return err
	}

	exists, err := transaction.NewSelect().
		Model((*certificateModel)(nil)).
		Where(constants.ByIDFilter, cert.ID).
		Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		_, err = transaction.NewUpdate().
			Model(model).
			Where(constants.ByIDFilter, model.ID).
			Exec(ctx)
	} else {
		_, err = transaction.NewInsert().Model(model).Exec(ctx)
	}

	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) FindPage(
	ctx context.Context,
	pageSize, pageNumber int,
	searchTerms *string,
) (*pagination.Page[certificate.Certificate], error) {
	certificates := make([]certificateModel, 0)

	query := r.database.Select().Model(&certificates)
	if searchTerms != nil {
		query = query.Where("domain_names::varchar ILIKE ?", "%"+*searchTerms+"%")
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}

	query = r.database.Select().Model(&certificates)
	if searchTerms != nil {
		query = query.Where("domain_names::varchar ILIKE ?", "%"+*searchTerms+"%")
	}

	err = query.
		Limit(pageSize).
		Offset(pageSize * pageNumber).
		Order("domain_names").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]certificate.Certificate, 0)
	for _, model := range certificates {
		//nolint:revive
		if domain, err := toDomain(&model); err != nil {
			return nil, err
		} else {
			result = append(result, *domain)
		}
	}

	return pagination.New(pageNumber, pageSize, count, result), nil
}

func (r *repository) FindAllDueToRenew(ctx context.Context) ([]certificate.Certificate, error) {
	certificates := make([]certificateModel, 0)

	err := r.database.Select().
		Model(&certificates).
		Where("renew_after IS NOT NULL AND renew_after <= ?", time.Now()).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]certificate.Certificate, 0)
	for _, model := range certificates {
		//nolint:revive
		if domain, err := toDomain(&model); err != nil {
			return nil, err
		} else {
			result = append(result, *domain)
		}
	}

	return result, nil
}
