package msproxy

type Requester interface {
	Post(api string, req interface{}, respData interface{}) error
	//Get(api string, respData interface{}) error
}
