package vpn

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/vpn"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

type repository struct {
	database *database.Database
}

func New(db *database.Database) vpn.Repository {
	return &repository{
		database: db,
	}
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*vpn.VPN, error) {
	var model vpnModel

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

	return toDomain(&model)
}

func (r *repository) ExistsByName(ctx context.Context, name string) (*bool, error) {
	count, err := r.database.Select().
		Model((*vpnModel)(nil)).
		Where("name = ?", name).
		Count(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return ptr.Of(false), nil
	}

	if err != nil {
		return nil, err
	}

	return ptr.Of(count > 0), nil
}

func (r *repository) ExistsByID(ctx context.Context, id uuid.UUID) (*bool, error) {
	exists, err := r.database.Select().
		Model((*vpnModel)(nil)).
		Where(constants.ByIDFilter, id).
		Exists(ctx)

	return &exists, err
}

func (r *repository) InUseByID(ctx context.Context, id uuid.UUID) (*bool, error) {
	count, err := r.database.Select().
		Table("host_vpn").
		Where("vpn_id = ?", id).
		Count(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return ptr.Of(false), nil
	}

	if err != nil {
		return nil, err
	}

	return ptr.Of(count > 0), nil
}

func (r *repository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	tx, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer tx.Rollback()

	_, err = tx.NewDelete().Model((*vpnModel)(nil)).Where(constants.ByIDFilter, id).Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) Save(ctx context.Context, values *vpn.VPN) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	model, err := toModel(values)
	if err != nil {
		return err
	}

	exists, err := transaction.NewSelect().Model((*vpnModel)(nil)).Where(constants.ByIDFilter, values.ID).Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		_, err = transaction.NewUpdate().Model(model).Where(constants.ByIDFilter, values.ID).Exec(ctx)
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
	enabledOnly bool,
) (*pagination.Page[vpn.VPN], error) {
	models := make([]vpnModel, 0)

	query := r.database.Select().Model(&models)
	if searchTerms != nil && *searchTerms != "" {
		query = query.Where("name ILIKE ?", "%"+*searchTerms+"%")
	}

	if enabledOnly {
		query = query.Where("enabled = ?", true)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}

	err = query.
		Limit(pageSize).
		Offset(pageSize * pageNumber).
		Order("name").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]vpn.VPN, 0)
	for _, model := range models {
		domain, err := toDomain(&model)
		if err != nil {
			return nil, err
		}

		result = append(result, *domain)
	}

	return pagination.New(pageNumber, pageSize, count, result), nil
}
