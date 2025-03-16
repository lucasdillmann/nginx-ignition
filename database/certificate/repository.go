package certificate

import (
	"context"
	"database/sql"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"errors"
	"github.com/google/uuid"
	"time"
)

type repository struct {
	database *database.Database
}

func New(database *database.Database) certificate.Repository {
	return repository{
		database: database,
	}
}

func (r repository) FindByID(ctx context.Context, id uuid.UUID) (*certificate.Certificate, error) {
	var model certificateModel

	err := r.database.Select().
		Model(&model).
		Where(constants.ByIdFilter, id).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if result, err := toDomain(&model); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (r repository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	count, err := r.database.Select().
		Model((*certificateModel)(nil)).
		Where(constants.ByIdFilter, id).
		Count(ctx)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r repository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	_, err = transaction.NewDelete().
		Model((*certificateModel)(nil)).
		Where(constants.ByIdFilter, id).
		Exec(ctx)

	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r repository) Save(ctx context.Context, certificate *certificate.Certificate) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	model, err := toModel(certificate)
	if err != nil {
		return err
	}

	exists, err := transaction.NewSelect().Model((*certificateModel)(nil)).Where(constants.ByIdFilter, certificate.ID).Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		_, err = transaction.NewUpdate().Model(model).Where(constants.ByIdFilter, model.ID).Exec(ctx)
	} else {
		_, err = transaction.NewInsert().Model(model).Exec(ctx)
	}

	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r repository) FindPage(
	ctx context.Context,
	pageSize, pageNumber int,
	searchTerms *string,
) (*pagination.Page[*certificate.Certificate], error) {
	var certificates []certificateModel

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

	var result []*certificate.Certificate
	for _, model := range certificates {
		if domain, err := toDomain(&model); err != nil {
			return nil, err
		} else {
			result = append(result, domain)
		}
	}

	return pagination.New(pageNumber, pageSize, count, result), nil
}

func (r repository) FindAllDueToRenew(ctx context.Context) ([]*certificate.Certificate, error) {
	var certificates []certificateModel

	err := r.database.Select().
		Model(&certificates).
		Where("renew_after IS NOT NULL AND renew_after <= ?", time.Now()).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	var result []*certificate.Certificate
	for _, model := range certificates {
		if domain, err := toDomain(&model); err != nil {
			return nil, err
		} else {
			result = append(result, domain)
		}
	}

	return result, nil
}
