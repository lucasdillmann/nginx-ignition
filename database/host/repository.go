package host

import (
	"context"
	"database/sql"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/database/certificate"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"errors"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type repository struct {
	database *database.Database
	ctx      context.Context
}

func New(database *database.Database) host.Repository {
	return &repository{
		database: database,
		ctx:      context.Background(),
	}
}

func (r *repository) FindByID(id uuid.UUID) (*host.Host, error) {
	var model hostModel

	err := r.database.Select().
		Model(&model).
		Relation("Bindings").
		Relation("Routes").
		Where(constants.ByIdFilter, id).
		Scan(r.ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return toDomain(&model)
}

func (r *repository) DeleteByID(id uuid.UUID) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	_, err = transaction.NewDelete().
		Model((*hostModel)(nil)).
		Where(constants.ByIdFilter, id).
		Exec(r.ctx)

	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) Save(host *host.Host) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	model, err := toModel(host)
	if err != nil {
		return err
	}

	exists, err := transaction.NewSelect().Model((*hostModel)(nil)).Where(constants.ByIdFilter, host.ID).Exists(r.ctx)
	if err != nil {
		return err
	}

	if exists {
		err = r.performUpdate(model, transaction)
	} else {
		_, err = transaction.NewInsert().Model(model).Exec(r.ctx)
	}

	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) performUpdate(model *hostModel, transaction bun.Tx) error {
	_, err := transaction.NewUpdate().Model(model).Where(constants.ByIdFilter, model.ID).Exec(r.ctx)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().Table("host_binding").Where("host_id = ?", model.ID).Exec(r.ctx)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().Table("host_route").Where("host_id = ?", model.ID).Exec(r.ctx)
	if err != nil {
		return err
	}

	for _, binding := range model.Bindings {
		binding.ID = uuid.New()
		binding.HostID = model.ID

		_, err = transaction.NewInsert().Model(binding).Exec(r.ctx)
		if err != nil {
			return err
		}
	}

	for _, route := range model.Routes {
		route.ID = uuid.New()
		route.HostID = model.ID

		_, err = transaction.NewInsert().Model(route).Exec(r.ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) FindPage(pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*host.Host], error) {
	var models []hostModel

	query := r.database.Select().Model(&models)
	if searchTerms != nil && *searchTerms != "" {
		query = query.Where("domain_names::varchar ILIKE ?", "%"+*searchTerms+"%")
	}

	count, err := query.Count(r.ctx)
	if err != nil {
		return nil, err
	}

	err = query.
		Relation("Bindings").
		Relation("Routes").
		Limit(pageSize).
		Offset(pageSize * pageNumber).
		Order("domain_names").
		Scan(r.ctx)

	if err != nil {
		return nil, err
	}

	var result []*host.Host
	for _, model := range models {
		domain, err := toDomain(&model)
		if err != nil {
			return nil, err
		}

		result = append(result, domain)
	}

	return pagination.New(pageNumber, pageSize, count, result), nil
}

func (r *repository) FindAllEnabled() ([]*host.Host, error) {
	var models []hostModel

	err := r.database.Select().
		Model(&models).
		Relation("Bindings").
		Relation("Routes").
		Where("enabled = ?", true).
		Scan(r.ctx)

	if err != nil {
		return nil, err
	}

	var result []*host.Host
	for _, model := range models {
		domain, err := toDomain(&model)
		if err != nil {
			return nil, err
		}
		result = append(result, domain)
	}

	return result, nil
}

func (r *repository) FindDefault() (*host.Host, error) {
	var model hostModel

	err := r.database.Select().
		Model(&model).
		Relation("Bindings").
		Relation("Routes").
		Where("default_server = ?", true).
		Scan(r.ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return toDomain(&model)
}

func (r *repository) ExistsByID(id uuid.UUID) (bool, error) {
	count, err := r.database.Select().
		Model((*hostModel)(nil)).
		Where(constants.ByIdFilter, id).
		Count(r.ctx)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *repository) ExistsCertificateByID(certificateId uuid.UUID) (bool, error) {
	return certificate.New(r.database).ExistsByID(certificateId)
}

func (r *repository) ExistsByCertificateID(certificateId uuid.UUID) (bool, error) {
	count, err := r.database.
		Select().
		Model((*hostBindingModel)(nil)).
		Where("certificate_id = ?", certificateId).
		Count(r.ctx)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *repository) ExistsByAccessListID(accessListId uuid.UUID) (bool, error) {
	count, err := r.database.Select().
		Model((*hostModel)(nil)).
		Where("access_list_id = ?", accessListId).
		Count(r.ctx)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
