package models

type Sales_Request_Add struct {
	EmpID      string  `json:"emp_id"`
	EmpName    string  `json:"emp_name"`
	SalesValue float64 `json:"sales_value"`
}

type Sales_Request_Report struct {
	SortBy string `json:"sort_by"` // name or amount
	Order  string `json:"order"`   // asc or desc
	Limit  string `json:"limit"`   // max - 200
}

type SalesRecord struct {
	EmpID      string  `json:"emp_id"`
	EmpName    string  `json:"emp_name"`
	TotalSales float64 `json:"total_sales"`
}

type Sales_Response_Report struct {
	Status     string        `json:"status"`
	Message    string        `json:"message"`
	TotalCount int           `json:"total_count"`
	Data       []SalesRecord `json:"data"`
}
