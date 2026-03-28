package company

import "context"

type CompanyRepository interface {
	FindByID(ctx context.Context, id CompanyID) (Company, error)
	FindAll(ctx context.Context) ([]Company, error)
	ExistsByCNPJ(ctx context.Context, cnpj string) (bool, error)
	Save(ctx context.Context, company Company) error
	Update(ctx context.Context, company Company) error
}

type BranchRepository interface {
	FindByID(ctx context.Context, id BranchID) (CompanyBranch, error)
	FindByCompanyID(ctx context.Context, companyID CompanyID) ([]CompanyBranch, error)
	Save(ctx context.Context, branch CompanyBranch) error
	Update(ctx context.Context, branch CompanyBranch) error
}
