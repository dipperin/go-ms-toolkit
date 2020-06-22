package msginengine

type decrypt struct {
}

func newDecrypt() *decrypt {
	return &decrypt{}
}

func (d *decrypt) PostBodyDecrypt(body string) ([]byte, error) {
	return RsaDecrypt(body)
}
