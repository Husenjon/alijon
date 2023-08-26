package repository

import (
	inkassback "github.com/Husenjon/InkassBack"
	"github.com/jmoiron/sqlx"
)

type Authoration interface {
	CreateUser(user inkassback.User) (inkassback.User, error)
	GetUser(username, password string) (inkassback.User, error)
	UpdateToken(id int, token string) (inkassback.User, error)
}

type Data1C interface {
	CreateBranch(db inkassback.Branch) (inkassback.Branch, error)
	CreateContracts(cs inkassback.Contracts, branchId int) (inkassback.Contracts, error)
	CreateRevenue(revenue inkassback.Revenue) (inkassback.Revenue, error)
}
type Repository struct {
	Authoration
	Data1C
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Data1C:      NewData1CPostgres(db),
		Authoration: NweAuthPostgres(db),
	}
}
