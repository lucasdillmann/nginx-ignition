package client

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"dillmann.com.br/nginx-ignition/core/certificate/client"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

const byClientCertificateIdFilter = "client_certificate_id = ?"

type repository struct {
	database *database.Database
}

func New(database *database.Database) client.Repository {
	return &repository{
		database: database,
	}
}

func (r *repository) IsInUseByID(ctx context.Context, id uuid.UUID) (bool, error) {
	count, err := r.database.Select().Table("host").Where(byClientCertificateIdFilter, id).Count(ctx)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	count, err = r.database.Select().Table("host_route").Where(byClientCertificateIdFilter, id).Count(ctx)
	return count > 0, err
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*client.Certificate, error) {
	var model certificateModel

	err := r.database.Select().Model(&model).Where(constants.ByIdFilter, id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var clients []*clientModel

	err = r.database.Select().Model(&clients).Where(byClientCertificateIdFilter, id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return toDomain(&model, clients)
}

func (r *repository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	_, err = transaction.NewDelete().
		Model((*clientModel)(nil)).
		Where(byClientCertificateIdFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().
		Model((*certificateModel)(nil)).
		Where(constants.ByIdFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) Save(ctx context.Context, certificate *client.Certificate) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	certModel, clientModels, err := toModel(certificate)
	if err != nil {
		return err
	}

	exists, err := transaction.
		NewSelect().
		Model((*certificateModel)(nil)).
		Where(constants.ByIdFilter, certificate.ID).
		Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		err = r.performUpdate(ctx, certModel, clientModels, transaction)
	} else {
		err = r.performInsert(ctx, certModel, clientModels, transaction)
	}

	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) performInsert(ctx context.Context, model *certificateModel, clients []*clientModel, transaction bun.Tx) error {
	_, err := transaction.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return err
	}

	return r.saveClients(ctx, clients, transaction)
}

func (r *repository) performUpdate(
	ctx context.Context,
	model *certificateModel,
	clients []*clientModel,
	transaction bun.Tx,
) error {
	_, err := transaction.NewUpdate().Model(model).Where(constants.ByIdFilter, model.ID).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = transaction.
		NewDelete().
		Table("client_certificate_item").
		Where(byClientCertificateIdFilter, model.ID).
		Exec(ctx)
	if err != nil {
		return err
	}

	return r.saveClients(ctx, clients, transaction)
}

func (r *repository) saveClients(ctx context.Context, clients []*clientModel, transaction bun.Tx) error {
	for _, item := range clients {
		_, err := transaction.NewInsert().Model(item).Exec(ctx)
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
) (*pagination.Page[*client.Certificate], error) {
	var models []certificateModel

	query := r.database.Select().Model(&models)
	if searchTerms != nil && *searchTerms != "" {
		query = query.Where("name ILIKE ?", "%"+*searchTerms+"%")
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

	var result []*client.Certificate
	for _, model := range models {
		var clients []*clientModel
		err = r.database.Select().
			Model(&clients).
			Where(byClientCertificateIdFilter, model.ID).
			Scan(ctx)
		if err != nil {
			return nil, err
		}

		domain, err := toDomain(&model, clients)
		if err != nil {
			return nil, err
		}

		result = append(result, domain)
	}

	return pagination.New(pageNumber, pageSize, count, result), nil
}
