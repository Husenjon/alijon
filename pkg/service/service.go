package service

import (
	inkassback "github.com/Husenjon/InkassBack"
	"github.com/Husenjon/InkassBack/pkg/repository"
)

type Authoration interface {
	GenerateHash(password string) string
	CreateUser(user inkassback.User) (inkassback.User, error)
	GetUser(username, password string) (inkassback.User, error)
	ParseToken(token string) (int, error)
}
type Data1C interface {
	CreateBranch(bd inkassback.BranchData) (inkassback.BranchData, error)
	CreateContracts(contracts inkassback.ContractsBody, branchId int) ([]inkassback.Contracts, error)
	CreateRevenue(revenueBody inkassback.RevenueBody) ([]inkassback.Revenue, error)
}
type Service struct {
	Authoration
	Data1C
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Data1C:      NewData1CService(repos.Data1C),
		Authoration: NewAuthService(repos.Authoration),
	}
}
