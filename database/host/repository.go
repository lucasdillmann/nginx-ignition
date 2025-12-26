package host

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/database/certificate"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

const byHostIdFilter = "host_id = ?"

type repository struct {
	database *database.Database
}

func New(database *database.Database) host.Repository {
	return &repository{
		database: database,
	}
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*host.Host, error) {
	var model hostModel

	err := r.database.Select().
		Model(&model).
		Relation("Bindings").
		Relation("Routes").
		Relation("VPNs").
		Where(constants.ByIdFilter, id).
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

	defer transaction.Rollback()

	_, err = transaction.NewDelete().
		Model((*hostBindingModel)(nil)).
		Where(byHostIdFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().
		Model((*hostRouteModel)(nil)).
		Where(byHostIdFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().
		Model((*hostVpnModel)(nil)).
		Where(byHostIdFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().
		Model((*hostModel)(nil)).
		Where(constants.ByIdFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) Save(ctx context.Context, host *host.Host) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	model, err := toModel(host)
	if err != nil {
		return err
	}

	exists, err := transaction.NewSelect().Model((*hostModel)(nil)).Where(constants.ByIdFilter, host.ID).Exists(ctx)
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
	_, err := transaction.NewUpdate().Model(model).Where(constants.ByIdFilter, model.ID).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().Table("host_binding").Where(byHostIdFilter, model.ID).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().Table("host_route").Where(byHostIdFilter, model.ID).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().Table("host_vpn").Where(byHostIdFilter, model.ID).Exec(ctx)
	if err != nil {
		return err
	}

	return r.saveLinkedModels(ctx, model, transaction)
}

func (r *repository) saveLinkedModels(ctx context.Context, model *hostModel, transaction bun.Tx) error {
	for _, binding := range model.Bindings {
		binding.ID = uuid.New()
		binding.HostID = model.ID

		_, err := transaction.NewInsert().Model(binding).Exec(ctx)
		if err != nil {
			return err
		}
	}

	for _, route := range model.Routes {
		route.ID = uuid.New()
		route.HostID = model.ID

		_, err := transaction.NewInsert().Model(route).Exec(ctx)
		if err != nil {
			return err
		}
	}

	for _, vpn := range model.VPNs {
		vpn.HostID = model.ID

		_, err := transaction.NewInsert().Model(vpn).Exec(ctx)
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
	var models []hostModel

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

	var result []host.Host
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
	var models []hostModel

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

	var result []host.Host
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
	count, err := r.database.Select().
		Model((*hostModel)(nil)).
		Where(constants.ByIdFilter, id).
		Count(ctx)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *repository) ExistsCertificateByID(ctx context.Context, certificateId uuid.UUID) (bool, error) {
	return certificate.New(r.database).ExistsByID(ctx, certificateId)
}

func (r *repository) ExistsByCertificateID(ctx context.Context, certificateId uuid.UUID) (bool, error) {
	count, err := r.database.
		Select().
		Model((*hostBindingModel)(nil)).
		Where("certificate_id = ?", certificateId).
		Count(ctx)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *repository) ExistsByAccessListID(ctx context.Context, accessListId uuid.UUID) (bool, error) {
	count, err := r.database.Select().
		Model((*hostModel)(nil)).
		Where("access_list_id = ?", accessListId).
		Count(ctx)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	count, err = r.database.Select().
		Model((*hostRouteModel)(nil)).
		Where("access_list_id = ?", accessListId).
		Count(ctx)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
