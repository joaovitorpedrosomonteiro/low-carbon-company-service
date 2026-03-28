package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joaovitorpedrosomonteiro/low-carbon-company-service/internal/domain/company"
	"github.com/joaovitorpedrosomonteiro/low-carbon-company-service/internal/domain/valueobject"
)

type CompanyPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewCompanyPostgresRepository(pool *pgxpool.Pool) *CompanyPostgresRepository {
	return &CompanyPostgresRepository{pool: pool}
}

func (r *CompanyPostgresRepository) FindByID(ctx context.Context, id company.CompanyID) (company.Company, error) {
	var (
		name, cnpjStr                     string
		street, city, state, postal, country string
		createdAt, updatedAt              time.Time
	)

	err := r.pool.QueryRow(ctx,
		`SELECT name, cnpj, street, city, state, postal_code, country, created_at, updated_at
		 FROM companies WHERE id = $1`, string(id),
	).Scan(&name, &cnpjStr, &street, &city, &state, &postal, &country, &createdAt, &updatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return company.Company{}, company.ErrCompanyNotFound
		}
		return company.Company{}, err
	}

	cnpj := valueobject.NewCNPJFromDB(cnpjStr)
	address := valueobject.NewAddressFromDB(street, city, state, postal, country)
	return company.NewCompanyFromDB(id, name, cnpj, address, createdAt, updatedAt), nil
}

func (r *CompanyPostgresRepository) FindAll(ctx context.Context) ([]company.Company, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, cnpj, street, city, state, postal_code, country, created_at, updated_at
		 FROM companies ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []company.Company
	for rows.Next() {
		var (
			idStr, name, cnpjStr                 string
			street, city, state, postal, country string
			createdAt, updatedAt                 time.Time
		)
		if err := rows.Scan(&idStr, &name, &cnpjStr, &street, &city, &state, &postal, &country, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		cnpj := valueobject.NewCNPJFromDB(cnpjStr)
		address := valueobject.NewAddressFromDB(street, city, state, postal, country)
		companies = append(companies, company.NewCompanyFromDB(company.CompanyID(idStr), name, cnpj, address, createdAt, updatedAt))
	}
	return companies, nil
}

func (r *CompanyPostgresRepository) ExistsByCNPJ(ctx context.Context, cnpj string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM companies WHERE cnpj = $1)`, cnpj).Scan(&exists)
	return exists, err
}

func (r *CompanyPostgresRepository) Save(ctx context.Context, c company.Company) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO companies (id, name, cnpj, street, city, state, postal_code, country, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		string(c.ID()), c.Name(), c.CNPJ().String(),
		c.Address().Street(), c.Address().City(), c.Address().State(),
		c.Address().PostalCode(), c.Address().Country(),
		c.CreatedAt(), c.UpdatedAt(),
	)
	return err
}

func (r *CompanyPostgresRepository) Update(ctx context.Context, c company.Company) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE companies SET name = $1, street = $2, city = $3, state = $4, postal_code = $5, country = $6, updated_at = $7
		 WHERE id = $8`,
		c.Name(), c.Address().Street(), c.Address().City(), c.Address().State(),
		c.Address().PostalCode(), c.Address().Country(), c.UpdatedAt(), string(c.ID()),
	)
	return err
}

type BranchPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewBranchPostgresRepository(pool *pgxpool.Pool) *BranchPostgresRepository {
	return &BranchPostgresRepository{pool: pool}
}

func (r *BranchPostgresRepository) FindByID(ctx context.Context, id company.BranchID) (company.CompanyBranch, error) {
	var (
		name, companyIDStr     string
		createdAt, updatedAt   time.Time
	)

	err := r.pool.QueryRow(ctx,
		`SELECT name, company_id, created_at, updated_at FROM company_branches WHERE id = $1`, string(id),
	).Scan(&name, &companyIDStr, &createdAt, &updatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return company.CompanyBranch{}, company.ErrBranchNotFound
		}
		return company.CompanyBranch{}, err
	}

	return company.NewCompanyBranchFromDB(id, name, company.CompanyID(companyIDStr), createdAt, updatedAt), nil
}

func (r *BranchPostgresRepository) FindByCompanyID(ctx context.Context, companyID company.CompanyID) ([]company.CompanyBranch, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, company_id, created_at, updated_at FROM company_branches WHERE company_id = $1 ORDER BY name`,
		string(companyID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branches []company.CompanyBranch
	for rows.Next() {
		var (
			idStr, name, cidStr  string
			createdAt, updatedAt time.Time
		)
		if err := rows.Scan(&idStr, &name, &cidStr, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		branches = append(branches, company.NewCompanyBranchFromDB(company.BranchID(idStr), name, company.CompanyID(cidStr), createdAt, updatedAt))
	}
	return branches, nil
}

func (r *BranchPostgresRepository) Save(ctx context.Context, b company.CompanyBranch) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO company_branches (id, name, company_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		string(b.ID()), b.Name(), string(b.CompanyID()), b.CreatedAt(), b.UpdatedAt(),
	)
	return err
}

func (r *BranchPostgresRepository) Update(ctx context.Context, b company.CompanyBranch) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE company_branches SET name = $1, updated_at = $2 WHERE id = $3`,
		b.Name(), b.UpdatedAt(), string(b.ID()),
	)
	return err
}
