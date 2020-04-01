package msginrouter

type IContext interface {
	GetResult() interface{}
	SetUserID(uid uint)
}
