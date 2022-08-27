package cakes

//go:generate mockgen -destination=../../mocks/repository/mock_repository.go -package=mock_repository -source=repository.go

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

const TableName = "cakes"

var (
	QuerySelect = fmt.Sprintf(`SELECT * FROM %s `, TableName)
	QueryInsert = `INSERT INTO ` + TableName + ` 
		(title, description, rating, image) 
		VALUES 
		('%s', '%s', %v, '%s')`
	QueryDelete = fmt.Sprintf(`DELETE FROM %s WHERE ID = ?`, TableName)
)

type repoImplementation struct {
	db *sql.DB
}

type RepoInterface interface {
	List(ctx context.Context, dto ListRequestDto) ([]Cake, int64, error)
	Get(ctx context.Context, id int) (*Cake, error)
	Create(ctx context.Context, dto RequestDto) error
	Update(ctx context.Context, dto UpdateRequestDto) error
	Delete(ctx context.Context, id int) error
}

func NewRepository(db *sql.DB) RepoInterface {
	return repoImplementation{
		db,
	}
}

func (i repoImplementation) List(ctx context.Context, dto ListRequestDto) (result []Cake, total int64, err error) {
	result = []Cake{}
	qWhere := "WHERE true "

	if dto.Title != "" {
		qWhere += fmt.Sprintf(`AND title LIKE '%%%s%%' `, dto.Title)
	}

	if dto.Description != "" {
		qWhere += fmt.Sprintf(`AND description LIKE '%%%s%%' `, dto.Description)
	}

	err = i.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s %s", TableName, qWhere)).Scan(&total)
	if err != nil {
		return
	}

	rows, err := i.db.QueryContext(ctx, QuerySelect+qWhere+"ORDER BY rating DESC, title ASC LIMIT ? OFFSET ?", dto.Limit, dto.Offset)
	if err != nil {
		return
	}

	for rows.Next() {
		var cake Cake
		err = rows.Scan(&cake.ID, &cake.Title, &cake.Description, &cake.Rating, &cake.Image, &cake.CreatedAt, &cake.UpdatedAt)
		if err != nil {
			return
		}
		result = append(result, cake)
	}
	return
}
func (i repoImplementation) Get(ctx context.Context, id int) (*Cake, error) {
	var result Cake
	err := i.db.QueryRowContext(ctx, fmt.Sprintf(QuerySelect+" WHERE id = ?"), id).
		Scan(&result.ID, &result.Title, &result.Description, &result.Rating, &result.Image, &result.CreatedAt, &result.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &result, err
}
func (i repoImplementation) Create(ctx context.Context, dto RequestDto) (err error) {
	_, err = i.db.ExecContext(ctx, fmt.Sprintf(QueryInsert, dto.Title, dto.Description, dto.Rating, dto.Image))
	return
}
func (i repoImplementation) Update(ctx context.Context, dto UpdateRequestDto) (err error) {
	updateQuery := "UPDATE cakes "
	var updated []string
	if dto.Title != "" {
		updated = append(updated, fmt.Sprintf(`title = '%s'`, dto.Title))
	}
	if dto.Description != "" {
		updated = append(updated, fmt.Sprintf(`description = '%s'`, dto.Description))
	}
	if dto.Rating != nil {
		updated = append(updated, fmt.Sprintf(`rating = %v`, *dto.Rating))
	}
	if dto.Image != "" {
		updated = append(updated, fmt.Sprintf(`image = '%s'`, dto.Image))
	}
	if len(updated) > 0 {
		updateQuery += "set updated_at = now(), " + strings.Join(updated, ", ")
	}

	_, err = i.db.ExecContext(ctx, updateQuery+"where id = ?", dto.ID)
	return
}
func (i repoImplementation) Delete(ctx context.Context, id int) error {
	_, err := i.db.ExecContext(ctx, QueryDelete, id)
	if err != nil {
		return err
	}
	return nil
}
