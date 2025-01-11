package certificate_repository

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"github.com/google/uuid"
	"time"
)

type repository struct {
	database *database.Database
	ctx      context.Context
}

func New(database *database.Database) certificate.Repository {
	return repository{
		database: database,
		ctx:      context.Background(),
	}
}

func (r repository) FindByID(id uuid.UUID) (*certificate.Certificate, error) {
	var model certificateModel

	err := r.database.Select().
		Model(&model).
		Where("id = ?", id).
		Scan(r.ctx)

	if err != nil {
		return nil, err
	}

	if result, err := toDomain(&model); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (r repository) ExistsByID(id uuid.UUID) (bool, error) {
	count, err := r.database.Select().
		Model((*certificateModel)(nil)).
		Where("id = ?", id).
		Count(r.ctx)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r repository) DeleteByID(id uuid.UUID) error {
	tx, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.NewDelete().
		Model((*certificateModel)(nil)).
		Where("id = ?", id).
		Exec(r.ctx)

	return err
}

func (r repository) Save(certificate *certificate.Certificate) error {
	tx, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	model, err := toModel(certificate)
	if err != nil {
		return err
	}

	_, err = tx.NewMerge().Model(model).Exec(r.ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r repository) FindPage(pageSize, pageNumber int, searchTerms *string) (*pagination.Page[certificate.Certificate], error) {
	var certificates []certificateModel

	query := r.database.Select().Model(&certificates)
	if searchTerms != nil {
		query = query.Where("domain_names ILIKE ?", "%"+*searchTerms+"%")
	}

	count, err := query.Count(r.ctx)

	if err != nil {
		return nil, err
	}

	query = r.database.Select().Model(&certificates)
	if searchTerms != nil {
		query = query.Where("domain_names ILIKE ?", "%"+*searchTerms+"%")
	}

	err = query.
		Limit(pageSize).
		Offset(pageSize * pageNumber).
		Order("domain_names").
		Scan(r.ctx)

	if err != nil {
		return nil, err
	}

	var result []certificate.Certificate
	for _, model := range certificates {
		if domain, err := toDomain(&model); err != nil {
			return nil, err
		} else {
			result = append(result, *domain)
		}
	}

	return pagination.New(pageNumber, pageSize, count, &result), nil
}

func (r repository) FindAllDueToRenew() (*[]certificate.Certificate, error) {
	var certificates []certificateModel

	err := r.database.Select().
		Model(&certificates).
		Where("renew_after IS NOT NULL AND renew_after <= ?", time.Now()).
		Scan(r.ctx)

	if err != nil {
		return nil, err
	}

	var result []certificate.Certificate
	for _, model := range certificates {
		if domain, err := toDomain(&model); err != nil {
			return nil, err
		} else {
			result = append(result, *domain)
		}
	}

	return &result, nil
}
