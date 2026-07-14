package types

// AppLoginReq is the OPS/MPC app-key WebSocket login body.
//
//easyjson:json
type AppLoginReq struct {
	AppID  string `json:"appID"`
	Sign   string `json:"sign"`
	Nonce  string `json:"nonce"`
	Time   int64  `json:"time"`
	Source string `json:"source"`
}

// CliPlan2LoginReq is the MPC Plan2 login body when no app key is configured.
//
//easyjson:json
type CliPlan2LoginReq struct {
	Source string `json:"source"`
}
