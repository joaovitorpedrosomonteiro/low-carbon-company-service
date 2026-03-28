package company

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/joaovitorpedrosomonteiro/low-carbon-company-service/internal/domain/valueobject"
)

var (
	ErrCompanyNotFound  = errors.New("company not found")
	ErrBranchNotFound   = errors.New("branch not found")
	ErrCNPJAlreadyUsed  = errors.New("CNPJ already registered")
)

type CompanyID string

func NewCompanyID() CompanyID {
	b := make([]byte, 16)
	rand.Read(b)
	return CompanyID(hex.EncodeToString(b))
}

type Company struct {
	id        CompanyID
	name      string
	cnpj      valueobject.CNPJ
	address   valueobject.Address
	createdAt time.Time
	updatedAt time.Time
	events    []DomainEvent
}

func NewCompany(name string, cnpj valueobject.CNPJ, address valueobject.Address) Company {
	now := time.Now().UTC()
	c := Company{
		id:        NewCompanyID(),
		name:      name,
		cnpj:      cnpj,
		address:   address,
		createdAt: now,
		updatedAt: now,
	}
	c.addEvent(NewCompanyCreatedEvent(c.id, c.name, c.cnpj.String()))
	return c
}

func NewCompanyFromDB(id CompanyID, name string, cnpj valueobject.CNPJ, address valueobject.Address, createdAt, updatedAt time.Time) Company {
	return Company{id: id, name: name, cnpj: cnpj, address: address, createdAt: createdAt, updatedAt: updatedAt}
}

func (c *Company) UpdateInfo(name *string, address *valueobject.Address) {
	updatedFields := make(map[string]string)
	if name != nil && *name != c.name {
		c.name = *name
		updatedFields["name"] = *name
	}
	if address != nil {
		c.address = *address
		updatedFields["address"] = "updated"
	}
	if len(updatedFields) > 0 {
		c.updatedAt = time.Now().UTC()
		c.addEvent(NewCompanyInfoUpdatedEvent(c.id, updatedFields))
	}
}

func (c *Company) ID() CompanyID              { return c.id }
func (c *Company) Name() string               { return c.name }
func (c *Company) CNPJ() valueobject.CNPJ     { return c.cnpj }
func (c *Company) Address() valueobject.Address { return c.address }
func (c *Company) CreatedAt() time.Time        { return c.createdAt }
func (c *Company) UpdatedAt() time.Time        { return c.updatedAt }

func (c *Company) PullEvents() []DomainEvent {
	events := c.events
	c.events = nil
	return events
}

func (c *Company) addEvent(event DomainEvent) {
	c.events = append(c.events, event)
}

type BranchID string

func NewBranchID() BranchID {
	b := make([]byte, 16)
	rand.Read(b)
	return BranchID(hex.EncodeToString(b))
}

type CompanyBranch struct {
	id        BranchID
	name      string
	companyID CompanyID
	createdAt time.Time
	updatedAt time.Time
	events    []DomainEvent
}

func NewCompanyBranch(name string, companyID CompanyID) CompanyBranch {
	now := time.Now().UTC()
	b := CompanyBranch{
		id:        NewBranchID(),
		name:      name,
		companyID: companyID,
		createdAt: now,
		updatedAt: now,
	}
	b.addEvent(NewCompanyBranchCreatedEvent(b.id, b.companyID, b.name))
	return b
}

func NewCompanyBranchFromDB(id BranchID, name string, companyID CompanyID, createdAt, updatedAt time.Time) CompanyBranch {
	return CompanyBranch{id: id, name: name, companyID: companyID, createdAt: createdAt, updatedAt: updatedAt}
}

func (b *CompanyBranch) Rename(newName string) {
	if newName != b.name {
		b.name = newName
		b.updatedAt = time.Now().UTC()
	}
}

func (b *CompanyBranch) ID() BranchID    { return b.id }
func (b *CompanyBranch) Name() string    { return b.name }
func (b *CompanyBranch) CompanyID() CompanyID { return b.companyID }
func (b *CompanyBranch) CreatedAt() time.Time { return b.createdAt }
func (b *CompanyBranch) UpdatedAt() time.Time { return b.updatedAt }

func (b *CompanyBranch) PullEvents() []DomainEvent {
	events := b.events
	b.events = nil
	return events
}

func (b *CompanyBranch) addEvent(event DomainEvent) {
	b.events = append(b.events, event)
}
