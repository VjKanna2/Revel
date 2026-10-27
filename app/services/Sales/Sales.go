package SalesService

import (
	"SalesManagement/app/models"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	storedRecords = make(map[string]*models.SalesRecord)
	recordsMutex  sync.RWMutex
)

func AddSalesRecord(payload models.Sales_Request_Add) (models.BasicResponse, any) {

	if strings.TrimSpace(payload.EmpID) == "" || strings.TrimSpace(payload.EmpName) == "" {
		return models.BasicResponse{
			Status:  "Fail",
			Message: "Employee ID & Name are Mandatory",
		}, "INVALID_PAYLOAD"
	}

	if payload.SalesValue < 0 {
		return models.BasicResponse{
			Status:  "Fail",
			Message: "Invalid Sales Value",
		}, "INVALID_PAYLOAD"
	}

	recordsMutex.Lock()
	defer recordsMutex.Unlock()

	existing, exists := storedRecords[payload.EmpID]
	if exists {
		if strings.TrimSpace(existing.EmpName) != strings.TrimSpace(payload.EmpName) {
			return models.BasicResponse{
				Status:  "Fail",
				Message: "Employee Name & Id Mismatch with Existing Records",
			}, "NAME_MISMATCH"
		}
		existing.TotalSales += payload.SalesValue
		fmt.Printf("Updated Sales Record For EmpID %s : New Total = %.2f\n", payload.EmpID, existing.TotalSales)
	} else {
		storedRecords[payload.EmpID] = &models.SalesRecord{
			EmpID:      payload.EmpID,
			EmpName:    payload.EmpName,
			TotalSales: payload.SalesValue,
		}
		fmt.Printf("Added New Sales Record For EmpID %s : Total = %.2f\n", payload.EmpID, payload.SalesValue)
	}

	return models.BasicResponse{
		Status:  "Success",
		Message: "Sales Record Added Successfully",
	}, nil
}

func GetSalesReport(params models.Sales_Request_Report) (models.Sales_Response_Report, any) {

	recordsMutex.RLock()
	records := make([]models.SalesRecord, 0, len(storedRecords))
	for _, record := range storedRecords {
		records = append(records, *record)
	}
	recordsMutex.RUnlock()

	if len(records) == 0 {
		return models.Sales_Response_Report{
			Status:     "Success",
			Message:    "No Sales Records Found",
			TotalCount: 0,
			Data:       records,
		}, nil
	}

	sortBy := strings.ToLower(strings.TrimSpace(params.SortBy))
	order := strings.ToLower(strings.TrimSpace(params.Order))

	if sortBy == "" {
		sortBy = "name"
	}
	if order == "" {
		order = "asc"
	}

	isAscending := order == "asc"

	sort.Slice(records, func(a, b int) bool {
		switch sortBy {
		case "value", "amount":
			if isAscending {
				return records[a].TotalSales < records[b].TotalSales
			}
			return records[a].TotalSales > records[b].TotalSales
		default:
			if isAscending {
				return strings.ToLower(records[a].EmpName) < strings.ToLower(records[b].EmpName)
			}
			return strings.ToLower(records[a].EmpName) > strings.ToLower(records[b].EmpName)
		}
	})

	limit := 50
	if trimmedLimit := strings.TrimSpace(params.Limit); trimmedLimit != "" {
		if parsedLimit, err := strconv.Atoi(trimmedLimit); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if limit > 200 {
		limit = 200
	}

	if len(records) > limit {
		records = records[:limit]
	}

	return models.Sales_Response_Report{
		Status:     "Success",
		Message:    "Sales Report Fetched Successfully",
		TotalCount: len(records),
		Data:       records,
	}, nil
}
