package response

import "aquafarm-management/app/model"

type HTTPResponseAction struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type HTTPResponseDataPonds struct {
	Status int              `json:"status"`
	Data   []model.PondShow `json:"data"`
}

type HTTPResponseDataPond struct {
	Status int            `json:"status"`
	Data   model.PondShow `json:"data"`
}

type HTTPResponseDataFarms struct {
	Status int              `json:"status"`
	Data   []model.FarmShow `json:"data"`
}

type HTTPResponseDataFarm struct {
	Status int            `json:"status"`
	Data   model.FarmShow `json:"data"`
}

type HTTPResponseDataLogs struct {
	Status int          `json:"status"`
	Data   []model.Logs `json:"data"`
}

type HTTPResponseDataStatistics struct {
	Status int                   `json:"status"`
	Data   model.StatisticsGroup `json:"data"`
}
