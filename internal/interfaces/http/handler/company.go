package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/joaovitorpedrosomonteiro/low-carbon-company-service/internal/application/command"
	"github.com/joaovitorpedrosomonteiro/low-carbon-company-service/internal/application/query"
	"github.com/joaovitorpedrosomonteiro/low-carbon-company-service/internal/domain/company"
)

type CompanyHandler struct {
	createCompany  *command.CreateCompanyHandler
	updateCompany  *command.UpdateCompanyHandler
	addBranch      *command.AddBranchHandler
	renameBranch   *command.RenameBranchHandler
	getCompany     *query.GetCompanyHandler
	listCompanies  *query.ListCompaniesHandler
	listBranches   *query.ListBranchesHandler
	getBranch      *query.GetBranchHandler
}

func NewCompanyHandler(
	createCompany *command.CreateCompanyHandler,
	updateCompany *command.UpdateCompanyHandler,
	addBranch *command.AddBranchHandler,
	renameBranch *command.RenameBranchHandler,
	getCompany *query.GetCompanyHandler,
	listCompanies *query.ListCompaniesHandler,
	listBranches *query.ListBranchesHandler,
	getBranch *query.GetBranchHandler,
) *CompanyHandler {
	return &CompanyHandler{
		createCompany:  createCompany,
		updateCompany:  updateCompany,
		addBranch:      addBranch,
		renameBranch:   renameBranch,
		getCompany:     getCompany,
		listCompanies:  listCompanies,
		listBranches:   listBranches,
		getBranch:      getBranch,
	}
}

type createCompanyRequest struct {
	Name       string `json:"name"`
	CNPJ       string `json:"cnpj"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	id, err := h.createCompany.Handle(r.Context(), command.CreateCompanyInput{
		Name:       req.Name,
		CNPJ:       req.CNPJ,
		Street:     req.Street,
		City:       req.City,
		State:      req.State,
		PostalCode: req.PostalCode,
		Country:    req.Country,
	})
	if err != nil {
		handleDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": string(id)})
}

func (h *CompanyHandler) List(w http.ResponseWriter, r *http.Request) {
	companies, err := h.listCompanies.Handle(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": companies})
}

func (h *CompanyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "Company ID is required")
		return
	}

	c, err := h.getCompany.Handle(r.Context(), company.CompanyID(id))
	if err != nil {
		handleDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, c)
}

type updateCompanyRequest struct {
	Name       *string `json:"name,omitempty"`
	Street     *string `json:"street,omitempty"`
	City       *string `json:"city,omitempty"`
	State      *string `json:"state,omitempty"`
	PostalCode *string `json:"postalCode,omitempty"`
	Country    *string `json:"country,omitempty"`
}

func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "Company ID is required")
		return
	}

	var req updateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if err := h.updateCompany.Handle(r.Context(), company.CompanyID(id), command.UpdateCompanyInput{
		Name:       req.Name,
		Street:     req.Street,
		City:       req.City,
		State:      req.State,
		PostalCode: req.PostalCode,
		Country:    req.Country,
	}); err != nil {
		handleDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type addBranchRequest struct {
	Name string `json:"name"`
}

func (h *CompanyHandler) AddBranch(w http.ResponseWriter, r *http.Request) {
	companyID := r.PathValue("id")
	if companyID == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "Company ID is required")
		return
	}

	var req addBranchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	branchID, err := h.addBranch.Handle(r.Context(), company.CompanyID(companyID), command.AddBranchInput{
		Name: req.Name,
	})
	if err != nil {
		handleDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": string(branchID)})
}

func (h *CompanyHandler) ListBranches(w http.ResponseWriter, r *http.Request) {
	companyID := r.PathValue("id")
	if companyID == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "Company ID is required")
		return
	}

	branches, err := h.listBranches.Handle(r.Context(), company.CompanyID(companyID))
	if err != nil {
		handleDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": branches})
}

func (h *CompanyHandler) GetBranch(w http.ResponseWriter, r *http.Request) {
	branchID := r.PathValue("branchId")
	if branchID == "" {
		writeError(w, http.StatusBadRequest, "MISSING_BRANCH_ID", "Branch ID is required")
		return
	}

	b, err := h.getBranch.Handle(r.Context(), company.BranchID(branchID))
	if err != nil {
		handleDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, b)
}

type renameBranchRequest struct {
	Name string `json:"name"`
}

func (h *CompanyHandler) RenameBranch(w http.ResponseWriter, r *http.Request) {
	branchID := r.PathValue("branchId")
	if branchID == "" {
		writeError(w, http.StatusBadRequest, "MISSING_BRANCH_ID", "Branch ID is required")
		return
	}

	var req renameBranchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if err := h.renameBranch.Handle(r.Context(), company.BranchID(branchID), req.Name); err != nil {
		handleDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, company.ErrCompanyNotFound):
		writeError(w, http.StatusNotFound, "COMPANY_NOT_FOUND", "Company not found")
	case errors.Is(err, company.ErrBranchNotFound):
		writeError(w, http.StatusNotFound, "BRANCH_NOT_FOUND", "Branch not found")
	case errors.Is(err, company.ErrCNPJAlreadyUsed):
		writeError(w, http.StatusConflict, "CNPJ_ALREADY_EXISTS", "CNPJ already registered")
	default:
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{
		Error: errorBody{Code: code, Message: message},
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
