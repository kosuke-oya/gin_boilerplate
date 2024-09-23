package dto

import "gin_server/models"

type FormatDataRequest struct {
	MachineID string `json:"machine_id"`
	Judgement bool   `json:"judgement"`
	models.Traceability
}

type FormatDataResponse struct {
	Data  models.FormatData
	Error string `json:"error"`
}
