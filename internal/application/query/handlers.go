package query

import (
	"context"
	"time"

	"github.com/joaovitorpedrosomonteiro/low-carbon-company-service/internal/domain/company"
)

type GetCompanyHandler struct {
	companyRepo company.CompanyRepository
}

func NewGetCompanyHandler(companyRepo company.CompanyRepository) *GetCompanyHandler {
	return &GetCompanyHandler{companyRepo: companyRepo}
}

type CompanyResponse struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	CNPJ      string          `json:"cnpj"`
	Address   AddressResponse `json:"address"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
}

type AddressResponse struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

func (h *GetCompanyHandler) Handle(ctx context.Context, id company.CompanyID) (CompanyResponse, error) {
	c, err := h.companyRepo.FindByID(ctx, id)
	if err != nil {
		return CompanyResponse{}, err
	}
	return toCompanyResponse(c), nil
}

type ListCompaniesHandler struct {
	companyRepo company.CompanyRepository
}

func NewListCompaniesHandler(companyRepo company.CompanyRepository) *ListCompaniesHandler {
	return &ListCompaniesHandler{companyRepo: companyRepo}
}

func (h *ListCompaniesHandler) Handle(ctx context.Context) ([]CompanyResponse, error) {
	companies, err := h.companyRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]CompanyResponse, len(companies))
	for i, c := range companies {
		result[i] = toCompanyResponse(c)
	}
	return result, nil
}

type ListBranchesHandler struct {
	branchRepo company.BranchRepository
}

func NewListBranchesHandler(branchRepo company.BranchRepository) *ListBranchesHandler {
	return &ListBranchesHandler{branchRepo: branchRepo}
}

type BranchResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CompanyID string `json:"companyId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (h *ListBranchesHandler) Handle(ctx context.Context, companyID company.CompanyID) ([]BranchResponse, error) {
	branches, err := h.branchRepo.FindByCompanyID(ctx, companyID)
	if err != nil {
		return nil, err
	}
	result := make([]BranchResponse, len(branches))
	for i, b := range branches {
		result[i] = toBranchResponse(b)
	}
	return result, nil
}

type GetBranchHandler struct {
	branchRepo company.BranchRepository
}

func NewGetBranchHandler(branchRepo company.BranchRepository) *GetBranchHandler {
	return &GetBranchHandler{branchRepo: branchRepo}
}

func (h *GetBranchHandler) Handle(ctx context.Context, branchID company.BranchID) (BranchResponse, error) {
	b, err := h.branchRepo.FindByID(ctx, branchID)
	if err != nil {
		return BranchResponse{}, err
	}
	return toBranchResponse(b), nil
}

func toCompanyResponse(c company.Company) CompanyResponse {
	return CompanyResponse{
		ID:   string(c.ID()),
		Name: c.Name(),
		CNPJ: c.CNPJ().String(),
		Address: AddressResponse{
			Street:     c.Address().Street(),
			City:       c.Address().City(),
			State:      c.Address().State(),
			PostalCode: c.Address().PostalCode(),
			Country:    c.Address().Country(),
		},
		CreatedAt: c.CreatedAt().Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt().Format(time.RFC3339),
	}
}

func toBranchResponse(b company.CompanyBranch) BranchResponse {
	return BranchResponse{
		ID:        string(b.ID()),
		Name:      b.Name(),
		CompanyID: string(b.CompanyID()),
		CreatedAt: b.CreatedAt().Format(time.RFC3339),
		UpdatedAt: b.UpdatedAt().Format(time.RFC3339),
	}
}
