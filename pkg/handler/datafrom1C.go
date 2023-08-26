package handler

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	inkassback "github.com/Husenjon/InkassBack"
	"github.com/Husenjon/InkassBack/pkg/logger"
)

func (h *Handler) GetDatasFrom1C() {
	var wg sync.WaitGroup
	wg.Add(1)
	branchs := make(chan inkassback.BranchData)
	go getBranch(h, branchs, &wg)
	go func() {
		wg.Wait()
		close(branchs)
	}()
	br := <-branchs
	dates := getDayRange(2023, 8, 20)
	for i := 0; i < len(br.Данные); i++ {
		for j := 0; j < len(dates)-1; j++ {
			var wg sync.WaitGroup
			wg.Add(1)
			contract := make(chan string)
			go getContracts(h, br.Данные[i].Код, dates[j], dates[j+1], contract, &wg, br.Данные[i].Id)
			go func() {
				wg.Wait()
				close(contract)
			}()
			for val := range contract {
				if val == "null" {
					continue
				}
				h.log.Info("1C", logger.Any("contract", val))
			}
		}
	}
	for i := 0; i < len(br.Данные); i++ {
		for j := 0; j < len(dates)-1; j++ {
			var wg sync.WaitGroup
			wg.Add(1)
			revenue := make(chan string)
			go getRevenue(h, br.Данные[i].Код, dates[j], dates[j+1], revenue, &wg)
			go func() {
				wg.Wait()
				close(revenue)
			}()
			for val := range revenue {
				if val == "null" {
					continue
				}
				h.log.Info("1C", logger.Any("revenue", val))
			}
		}
	}
	fmt.Println("tugadi")
}
func getContracts(h *Handler, branchCode, dataFrom, dateTo string, contract chan string, wg *sync.WaitGroup, branchId int) {
	defer wg.Done()
	url := fmt.Sprintf("http://10.0.57.2/inkas/hs/api/getcontracts?branchcode=%s&datefrom=%s&dateto=%s", branchCode, dataFrom, dateTo)
	var cb inkassback.ContractsBody
	err := HttpRequset("GET", url, &cb)
	if err != nil {
		contract <- err.Error()
		return
	}
	d, err := h.services.Data1C.CreateContracts(cb, branchId)
	if err != nil {
		contract <- err.Error()
	}
	jsonData, err := json.Marshal(d)
	if err != nil {
		contract <- err.Error()
		return
	}
	contract <- string(jsonData)
}
func getBranch(h *Handler, branchs chan inkassback.BranchData, wg *sync.WaitGroup) {
	defer wg.Done()
	var bd inkassback.BranchData
	err := HttpRequset("GET", "http://10.0.57.2/inkas/hs/api/getbranches", &bd)
	if err != nil {
		bd.Ошибка = err.Error()
		branchs <- bd
		return
	}
	b, err := h.services.Data1C.CreateBranch(bd)
	if err != nil {
		b.Ошибка = err.Error()
		branchs <- b
		return
	}
	branchs <- b
}
func getRevenue(h *Handler, branchCode, dataFrom, dateTo string, revenue chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	url := fmt.Sprintf("http://10.0.57.2/inkas/hs/api/revenue?branchcode=%s&datefrom=%s&dateto=%s", branchCode, dataFrom, dateTo)
	var re inkassback.RevenueBody
	err := HttpRequset("GET", url, &re)
	if err != nil {
		revenue <- err.Error()
		return
	}
	r, err := h.services.Data1C.CreateRevenue(re)
	if err != nil {
		revenue <- err.Error()
		return
	}
	jsonData, err := json.Marshal(r)
	if err != nil {
		revenue <- err.Error()
		return
	}
	revenue <- string(jsonData)
}
func getDayRange(year, month, day int) []string {
	var dates []string
	startTime := time.Now()
	currentLocation := startTime.Location()
	var t = time.Date(year, time.Month(month), day, 0, 0, 0, 0, currentLocation)
	for i := 0; startTime.Unix() > t.Unix(); i++ {
		t = t.AddDate(0, 0, 1)
		year, month, day := t.Local().Date()
		var d = fmt.Sprintf("%d%02d%02d", year, int(month), day)
		dates = append(dates, d)
	}
	return dates
}
