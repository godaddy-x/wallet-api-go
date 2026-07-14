package types

import (
	"github.com/godaddy-x/freego/node/common"
	"github.com/godaddy-x/freego/ormx/sqlc"
)

// MonitorAlertResult is a monitor rule hit event; pull incrementally by id.
//easyjson:json
type MonitorAlertResult struct {
	ID              int64  `json:"id"`
	AppID           string `json:"appID"`
	Domain          string `json:"domain"`
	RuleID          int64  `json:"ruleID"`
	RuleCode        string `json:"ruleCode"`
	Name            string `json:"name"`
	Category        string `json:"category"`
	Metric          string `json:"metric"`
	Level           int64  `json:"level"`
	WindowType      string `json:"windowType"`
	BucketType      string `json:"bucketType"`
	BucketStart     int64  `json:"bucketStart"`
	Scope           string `json:"scope"`
	Symbol          string `json:"symbol"`
	ContractAddress string `json:"contractAddress"`
	Address         string `json:"address"`
	MainSymbol      string `json:"mainSymbol"`
	MetricValue     string `json:"metricValue"`
	ThresholdValue  string `json:"thresholdValue"`
	Threshold2Value string `json:"threshold2Value"`
	CompareOp       string `json:"compareOp"`
	Message         string `json:"message"`
	Payload         string `json:"payload"`
	AlertLatency    string `json:"alertLatency"`
	Status          int64  `json:"status"`
	CreateAt        int64  `json:"createAt"`
	UpdateAt        int64  `json:"updateAt"`
}

//easyjson:json
type FindMonitorAlertReq struct {
	common.BaseReq
}

//easyjson:json
type FindMonitorAlertRes struct {
	Result []MonitorAlertResult `json:"result"`
	Limit  sqlc.Limit           `json:"limit"`
}
