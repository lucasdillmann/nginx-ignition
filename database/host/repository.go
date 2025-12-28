package host

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

const (
	byHostIDFilter       = "host_id = ?"
	byAccessListIDFilter = "access_list_id = ?"
)

type repository struct {
	database *database.Database
}

func New(db *database.Database) host.Repository {
	return &repository{
		database: db,
	}
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*host.Host, error) {
	var model hostModel

	err := r.database.Select().
		Model(&model).
		Relation("Bindings").
		Relation("Routes").
		Relation("VPNs").
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

func (r *repository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	_, err = transaction.NewDelete().
		Model((*hostBindingModel)(nil)).
		Where(byHostIDFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().
		Model((*hostRouteModel)(nil)).
		Where(byHostIDFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().
		Model((*hostVpnModel)(nil)).
		Where(byHostIDFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().
		Model((*hostModel)(nil)).
		Where(constants.ByIDFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) Save(ctx context.Context, h *host.Host) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	model, err := toModel(h)
	if err != nil {
		return err
	}

	exists, err := transaction.NewSelect().Model((*hostModel)(nil)).Where(constants.ByIDFilter, h.ID).Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		err = r.performUpdate(ctx, model, transaction)
	} else {
		err = r.performInsert(ctx, model, transaction)
	}

	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) performInsert(ctx context.Context, model *hostModel, transaction bun.Tx) error {
	_, err := transaction.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return err
	}

	return r.saveLinkedModels(ctx, model, transaction)
}

func (r *repository) performUpdate(ctx context.Context, model *hostModel, transaction bun.Tx) error {
	_, err := transaction.NewUpdate().Model(model).Where(constants.ByIDFilter, model.ID).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().Table("host_binding").Where(byHostIDFilter, model.ID).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().Table("host_route").Where(byHostIDFilter, model.ID).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().Table("host_vpn").Where(byHostIDFilter, model.ID).Exec(ctx)
	if err != nil {
		return err
	}

	return r.saveLinkedModels(ctx, model, transaction)
}

func (r *repository) saveLinkedModels(ctx context.Context, model *hostModel, transaction bun.Tx) error {
	for _, binding := range model.Bindings {
		binding.ID = uuid.New()
		binding.HostID = model.ID

		_, err := transaction.NewInsert().Model(&binding).Exec(ctx)
		if err != nil {
			return err
		}
	}

	for _, route := range model.Routes {
		route.ID = uuid.New()
		route.HostID = model.ID

		_, err := transaction.NewInsert().Model(&route).Exec(ctx)
		if err != nil {
			return err
		}
	}

	for _, vpn := range model.VPNs {
		vpn.HostID = model.ID

		_, err := transaction.NewInsert().Model(&vpn).Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) FindPage(
	ctx context.Context,
	pageSize, pageNumber int,
	searchTerms *string,
) (*pagination.Page[host.Host], error) {
	models := make([]hostModel, 0)

	query := r.database.Select().Model(&models)
	if searchTerms != nil && *searchTerms != "" {
		query = query.Where("domain_names::varchar ILIKE ?", "%"+*searchTerms+"%")
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}

	err = query.
		Relation("Bindings").
		Relation("Routes").
		Relation("VPNs").
		Limit(pageSize).
		Offset(pageSize * pageNumber).
		Order("domain_names").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]host.Host, 0)
	for _, model := range models {
		domain, err := toDomain(&model)
		if err != nil {
			return nil, err
		}

		result = append(result, *domain)
	}

	return pagination.New(pageNumber, pageSize, count, result), nil
}

func (r *repository) FindAllEnabled(ctx context.Context) ([]host.Host, error) {
	models := make([]hostModel, 0)

	err := r.database.Select().
		Model(&models).
		Relation("Bindings").
		Relation("Routes").
		Relation("VPNs").
		Where("enabled = ?", true).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]host.Host, 0)
	for _, model := range models {
		domain, err := toDomain(&model)
		if err != nil {
			return nil, err
		}
		result = append(result, *domain)
	}

	return result, nil
}

func (r *repository) FindDefault(ctx context.Context) (*host.Host, error) {
	var model hostModel

	err := r.database.Select().
		Model(&model).
		Relation("Bindings").
		Relation("Routes").
		Relation("VPNs").
		Where("default_server = ?", true).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return toDomain(&model)
}

func (r *repository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.database.Select().
		Model((*hostModel)(nil)).
		Where(constants.ByIDFilter, id).
		Exists(ctx)
}
