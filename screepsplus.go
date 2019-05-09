package main

import (
	"gopkg.in/resty.v1"
)

type screepsplusStatResponse struct {
	Ok        int8
	Timestamp int64 `yaml:"ts"`
	Format    string
}

func pushStats(token string, stats []Stat) (*screepsplusStatResponse, error) {
	data := make(map[string]float64, 0)
	for _, stat := range stats {
		data[stat.Key] = stat.Value
	}
	resp, err := resty.R().
		SetBody(data).
		SetBasicAuth("token", token).
		SetResult(screepsplusStatResponse{}).
		Post("https://screepspl.us/api/stats/submit")
	if err != nil {
		return nil, err
	}
	return resp.Result().(*screepsplusStatResponse), nil
}
