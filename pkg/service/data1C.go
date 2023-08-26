package service

import (
	"errors"

	inkassback "github.com/Husenjon/InkassBack"
	"github.com/Husenjon/InkassBack/pkg/repository"
)

type Data1CService struct {
	repo repository.Data1C
}

func NewData1CService(repo repository.Data1C) *Data1CService {
	return &Data1CService{repo: repo}
}
func (s *Data1CService) CreateBranch(branchData inkassback.BranchData) (inkassback.BranchData, error) {
	var bd = inkassback.BranchData{}
	for i := 0; i < len(branchData.Данные); i++ {
		d, err := s.repo.CreateBranch(branchData.Данные[i])
		if err != nil {
			bd.Ошибка = err.Error()
			return bd, err
		}
		bd.Данные = append(bd.Данные, d)
	}
	return bd, nil
}
func (s *Data1CService) CreateContracts(cb inkassback.ContractsBody, branchId int) ([]inkassback.Contracts, error) {
	var cs []inkassback.Contracts
	if cb.Ошибка != "" {
		return cs, errors.New(cb.Ошибка)
	}
	if len(cb.Данные) == 0 {
		return cs, nil
	}

	for i := 0; i < len(cb.Данные); i++ {
		res, err := s.repo.CreateContracts(cb.Данные[i], branchId)
		if err != nil {
			return nil, err
		}
		cs = append(cs, res)
	}
	return cs, nil
}
func (s *Data1CService) CreateRevenue(revenueBody inkassback.RevenueBody) ([]inkassback.Revenue, error) {
	var re []inkassback.Revenue
	if revenueBody.Ошибка != "" {
		return re, errors.New(revenueBody.Ошибка)
	}
	if len(revenueBody.Данные) == 0 {
		return re, nil
	}
	for i := 0; i < len(revenueBody.Данные); i++ {
		res, err := s.repo.CreateRevenue(revenueBody.Данные[i])
		if err != nil {
			return nil, err
		}
		re = append(re, res)
	}
	return re, nil
}
