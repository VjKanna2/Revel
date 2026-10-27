package controllers

import (
	models "SalesManagement/app/models"
	SalesService "SalesManagement/app/services/Sales"
	"encoding/json"
	"fmt"
	"io"

	"github.com/revel/revel"
)

type SalesModule struct {
	*revel.Controller
}

func (c SalesModule) AddSales() revel.Result {

	var Result any
	var payload models.Sales_Request_Add

	body := c.Params.JSON
	if len(body) == 0 && c.Request != nil && c.Request.GetBody() != nil {
		body, _ = io.ReadAll(c.Request.GetBody())
	}

	err := json.Unmarshal(body, &payload)
	if err != nil {
		fmt.Println("Error While Adding Sales :", err)
		c.Response.Status = 400
		Result = models.BasicResponse{
			Status:  "Fail",
			Message: "Invalid payload",
		}
		return c.RenderJSON(Result)
	}

	Result, Error := SalesService.AddSalesRecord(payload)
	if Error != nil {
		c.Response.Status = 500
	}

	return c.RenderJSON(Result)
}

func (c SalesModule) GetReport() revel.Result {

	params := models.Sales_Request_Report{
		SortBy: c.Params.Get("sort_by"),
		Order:  c.Params.Get("order"),
		Limit:  c.Params.Get("limit"),
	}

	Result, Error := SalesService.GetSalesReport(params)
	if Error != nil {
		c.Response.Status = 500
	}

	return c.RenderJSON(Result)
}
