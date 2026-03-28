package company

import "time"

type DomainEvent struct {
	EventID     string    `json:"event_id"`
	EventType   string    `json:"event_type"`
	OccurredAt  time.Time `json:"occurred_at"`
	SchemaVer   string    `json:"schema_version"`
	Payload     any       `json:"payload"`
}

type CompanyCreatedPayload struct {
	CompanyID string `json:"companyID"`
	Name      string `json:"name"`
	CNPJ      string `json:"cnpj"`
}

type CompanyBranchCreatedPayload struct {
	BranchID  string `json:"branchID"`
	CompanyID string `json:"companyID"`
	Name      string `json:"name"`
}

type CompanyInfoUpdatedPayload struct {
	CompanyID     string            `json:"companyID"`
	UpdatedFields map[string]string `json:"updatedFields"`
}

func NewCompanyCreatedEvent(id CompanyID, name, cnpj string) DomainEvent {
	return DomainEvent{
		EventID:    string(NewCompanyID()),
		EventType:  "CompanyCreated",
		OccurredAt: time.Now().UTC(),
		SchemaVer:  "1.0",
		Payload: CompanyCreatedPayload{
			CompanyID: string(id),
			Name:      name,
			CNPJ:      cnpj,
		},
	}
}

func NewCompanyBranchCreatedEvent(id BranchID, companyID CompanyID, name string) DomainEvent {
	return DomainEvent{
		EventID:    string(NewCompanyID()),
		EventType:  "CompanyBranchCreated",
		OccurredAt: time.Now().UTC(),
		SchemaVer:  "1.0",
		Payload: CompanyBranchCreatedPayload{
			BranchID:  string(id),
			CompanyID: string(companyID),
			Name:      name,
		},
	}
}

func NewCompanyInfoUpdatedEvent(companyID CompanyID, updatedFields map[string]string) DomainEvent {
	return DomainEvent{
		EventID:    string(NewCompanyID()),
		EventType:  "CompanyInfoUpdated",
		OccurredAt: time.Now().UTC(),
		SchemaVer:  "1.0",
		Payload: CompanyInfoUpdatedPayload{
			CompanyID:     string(companyID),
			UpdatedFields: updatedFields,
		},
	}
}
