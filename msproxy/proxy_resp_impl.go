package msproxy

type ProxyResp struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	ErrMsg  string      `json:"err_msg"`
}

func (p *ProxyResp) GetErrMsg() string {
	return p.ErrMsg
}

func (p *ProxyResp) GetData() interface{} {
	return p.Data
}

func (p *ProxyResp) GetSuccess() bool {
	return p.Success
}
