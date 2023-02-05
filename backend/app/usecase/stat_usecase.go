package usecase

import (
	"aquafarm-management/app/model"
	"aquafarm-management/app/repository"
	"strconv"
)

// StatUsecase menangani proses bisnis item
type StatUsecase struct {
	Repository *repository.LogRepository
}

// NewStatUsecase membuat instance baru StatUsecase
func NewStatUsecase(r *repository.LogRepository) *StatUsecase {
	return &StatUsecase{
		Repository: r,
	}
}

// Fetch mengambil semua item dari repository
func (u *StatUsecase) Fetch() (*model.StatisticsGroup, error) {
	stats, err := u.Repository.FetchStatistics()
	if err != nil {
		return nil, err
	}

	result := model.StatisticsGroup{StatisticsData: make(map[string]model.StatisticsData)}
	for _, v := range stats {
		var data model.StatisticsData
		data.Count, _ = strconv.Atoi(v.Count)
		data.UniqueUserAgent, _ = strconv.Atoi(v.UniqueUserAgent)
		data.Response200, _ = strconv.Atoi(v.Response200)
		data.Response404, _ = strconv.Atoi(v.Response404)
		data.Response500, _ = strconv.Atoi(v.Response500)
		data.ResponseETC, _ = strconv.Atoi(v.ResponseETC)
		result.StatisticsData[v.Request] = data
	}

	return &result, nil
}

// FetchLogs mengambil semua item dari repository
func (u *StatUsecase) FetchLogs() ([]model.Logs, error) {
	return u.Repository.Fetch()
}
