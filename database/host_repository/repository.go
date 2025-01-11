package host_repository

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"github.com/google/uuid"
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
		Where("id = ?", id).
		Scan(r.ctx)

	if err != nil {
		return nil, err
	}

	return toDomain(&model)
}

func (r *repository) DeleteByID(id uuid.UUID) error {
	tx, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.NewDelete().
		Model((*hostModel)(nil)).
		Where("id = ?", id).
		Exec(r.ctx)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) Save(host *host.Host) error {
	tx, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	model, err := toModel(host)
	if err != nil {
		return err
	}

	_, err = tx.NewMerge().
		Model(model).
		Exec(r.ctx)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) FindPage(pageSize, pageNumber int, searchTerms *string) (*pagination.Page[host.Host], error) {
	var models []hostModel

	query := r.database.Select().Model(&models)
	if searchTerms != nil && *searchTerms != "" {
		query = query.Where("domain_names ILIKE ?", "%"+*searchTerms+"%")
	}

	count, err := query.Count(r.ctx)
	if err != nil {
		return nil, err
	}

	err = query.
		Limit(pageSize).
		Offset(pageSize * pageNumber).
		Order("domain_names").
		Scan(r.ctx)

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

	return pagination.New(pageNumber, pageSize, count, &result), nil
}

func (r *repository) FindAllEnabled() ([]*host.Host, error) {
	var models []hostModel

	err := r.database.Select().
		Model(&models).
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
		Where("default_server = ?", true).
		Scan(r.ctx)

	if err != nil {
		return nil, err
	}

	return toDomain(&model)
}

func (r *repository) ExistsByID(id uuid.UUID) (bool, error) {
	count, err := r.database.Select().
		Model((*hostModel)(nil)).
		Where("id = ?", id).
		Count(r.ctx)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *repository) ExistsByCertificateID(certificateId uuid.UUID) (bool, error) {
	count, err := r.database.Select().
		Model((*hostModel)(nil)).
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
