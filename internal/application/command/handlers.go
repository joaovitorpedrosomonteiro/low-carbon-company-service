package command

import (
	"context"

	"github.com/joaovitorpedrosomonteiro/low-carbon-company-service/internal/domain/company"
	"github.com/joaovitorpedrosomonteiro/low-carbon-company-service/internal/domain/valueobject"
)

type EventPublisher interface {
	Publish(ctx context.Context, event company.DomainEvent) error
}

type CreateCompanyHandler struct {
	companyRepo company.CompanyRepository
	publisher   EventPublisher
}

func NewCreateCompanyHandler(companyRepo company.CompanyRepository, publisher EventPublisher) *CreateCompanyHandler {
	return &CreateCompanyHandler{companyRepo: companyRepo, publisher: publisher}
}

type CreateCompanyInput struct {
	Name       string
	CNPJ       string
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

func (h *CreateCompanyHandler) Handle(ctx context.Context, input CreateCompanyInput) (company.CompanyID, error) {
	cnpj, err := valueobject.NewCNPJ(input.CNPJ)
	if err != nil {
		return "", err
	}

	exists, err := h.companyRepo.ExistsByCNPJ(ctx, cnpj.String())
	if err != nil {
		return "", err
	}
	if exists {
		return "", company.ErrCNPJAlreadyUsed
	}

	address, err := valueobject.NewAddress(input.Street, input.City, input.State, input.PostalCode, input.Country)
	if err != nil {
		return "", err
	}

	c := company.NewCompany(input.Name, cnpj, address)
	if err := h.companyRepo.Save(ctx, c); err != nil {
		return "", err
	}

	for _, event := range c.PullEvents() {
		if err := h.publisher.Publish(ctx, event); err != nil {
			return "", err
		}
	}

	return c.ID(), nil
}

type UpdateCompanyHandler struct {
	companyRepo company.CompanyRepository
}

func NewUpdateCompanyHandler(companyRepo company.CompanyRepository) *UpdateCompanyHandler {
	return &UpdateCompanyHandler{companyRepo: companyRepo}
}

type UpdateCompanyInput struct {
	Name       *string
	Street     *string
	City       *string
	State      *string
	PostalCode *string
	Country    *string
}

func (h *UpdateCompanyHandler) Handle(ctx context.Context, id company.CompanyID, input UpdateCompanyInput) error {
	c, err := h.companyRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	var addr *valueobject.Address
	if input.Street != nil || input.City != nil || input.State != nil || input.PostalCode != nil || input.Country != nil {
		street := c.Address().Street()
		if input.Street != nil {
			street = *input.Street
		}
		city := c.Address().City()
		if input.City != nil {
			city = *input.City
		}
		state := c.Address().State()
		if input.State != nil {
			state = *input.State
		}
		postalCode := c.Address().PostalCode()
		if input.PostalCode != nil {
			postalCode = *input.PostalCode
		}
		country := c.Address().Country()
		if input.Country != nil {
			country = *input.Country
		}
		a, err := valueobject.NewAddress(street, city, state, postalCode, country)
		if err != nil {
			return err
		}
		addr = &a
	}

	c.UpdateInfo(input.Name, addr)
	return h.companyRepo.Update(ctx, c)
}

type AddBranchHandler struct {
	companyRepo company.CompanyRepository
	branchRepo  company.BranchRepository
	publisher   EventPublisher
}

func NewAddBranchHandler(companyRepo company.CompanyRepository, branchRepo company.BranchRepository, publisher EventPublisher) *AddBranchHandler {
	return &AddBranchHandler{companyRepo: companyRepo, branchRepo: branchRepo, publisher: publisher}
}

type AddBranchInput struct {
	Name string
}

func (h *AddBranchHandler) Handle(ctx context.Context, companyID company.CompanyID, input AddBranchInput) (company.BranchID, error) {
	if _, err := h.companyRepo.FindByID(ctx, companyID); err != nil {
		return "", err
	}

	branch := company.NewCompanyBranch(input.Name, companyID)
	if err := h.branchRepo.Save(ctx, branch); err != nil {
		return "", err
	}

	for _, event := range branch.PullEvents() {
		if err := h.publisher.Publish(ctx, event); err != nil {
			return "", err
		}
	}

	return branch.ID(), nil
}

type RenameBranchHandler struct {
	branchRepo company.BranchRepository
}

func NewRenameBranchHandler(branchRepo company.BranchRepository) *RenameBranchHandler {
	return &RenameBranchHandler{branchRepo: branchRepo}
}

func (h *RenameBranchHandler) Handle(ctx context.Context, branchID company.BranchID, newName string) error {
	branch, err := h.branchRepo.FindByID(ctx, branchID)
	if err != nil {
		return err
	}

	branch.Rename(newName)
	return h.branchRepo.Update(ctx, branch)
}
