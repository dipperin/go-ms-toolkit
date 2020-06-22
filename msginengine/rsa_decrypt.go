package msginengine

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/dipperin/go-ms-toolkit/log"
	"go.uber.org/zap"
	"runtime"
)

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA6mrG8RLv5TjyCuq5F9hbCZiHktqJHp1EMIsC92p6SEYyPzSj
arF24dO7nxxPxwqt2Ch/8QCRgBubU+gf8D6D+I9XlpHASfv5yNwkRDgEyV/830JV
oQIR1XmJYL/nH7M4Kdj9Flcjbl0rfSyH7ij+tuzcnFPSkDUJdK8E9VaIdfol5e84
sKXEgt9sS3MD0kYlP9RPKPnxgdnN6XEiJpjy/JxLz/KPQYWe+O7Hn2cg4vaPEl6p
p05uxwh7hyAvimOTLM8ts47fFBJJd0OZ+nYXg7c5u1NbJgoGG28U8ivlnfPV/E6b
eU3aFkENPYE8mHrgytAX48WV43ZyowL+0EWO9QIDAQABAoIBAQCfIZclhe52U/7z
bD30Mvox/GpkRZf5wVbOAUAlRxH1yDlJ8OjSf+AtEzf1nhmGC/jRmUSpDPK43YTH
I/eydi3OaThTTWQUlUoOkWrJKKIPNesKgBRy9V235gZdOEikm1wQBG5iYQr7W6Iv
GjC4evnWodptAPYa0PY3UKx0A6clNjsJJzfErcoXfobvZ9HkTunU/ecJioA1quY1
VPYZaio+7a2rZJpWJzUNc1xXhwfWO/U0Q9733hkvegUvV4kIQ0hdiR8tcxJoNeyQ
A5V3nkWD10wE3whLk76BrFHeWHHlvtD8tKsY1efMwRdjX9W650o0DcAKlx9jPVCn
GPIaVNsNAoGBAO017OGnffEWGFvrWuxo6tg/ILTpAKT6AyMVcF8XT9s1R0ppG4+V
Hos4b+duDfL7dbXWSpDgUW3ylnuhLEDOu2SkvV2au6piQpEAHV4Dcxt1avv3wDkZ
NaPo3Z2PrIG8Ch/MXKTOzWENQVP1eLnojAtraTya4PwT+dt1pJZruT4bAoGBAPz8
NIjbQyphl/OcQBB7k16E6BS55aq86MDpWx+5Q5UQorl0xmezKSYgjesjUUDumXZw
YlANz50SOjD2ZmqlVCelDn0Tz8hr3G7XI/46AQEecf1lav8Mo7a8fg3giVf1Z3dX
olN8EK/7+NLYQYpC5VmS+6SoaYp8IfHbFG+c2/gvAoGBAOXigdw6byJq0Faco3RD
RX0myLKqsLYxmUKx70IonHqLiriBXnVrBbvUiRaXIKufqPb9Yyw+SIwuMkpD61gv
QYqK4P5LC55XNb3Ch7Np0m8E/xaLht2PId3kuomNCJh+PK2OZyZNSNrKdspANayt
wrL6eHWEf4+saFOYEla4tUCLAoGAVfRu+QzefjjYgvnUvoTuJlAr9lbPFLrIrjrb
LiUpNC2JzE9D1hoGYiDvdzPxuIkY9SfsD10R7EE0KkydsaBovmuTR6YClAspwzMv
1IHqs3Gfs0PRHcUynrTec2KV55/GvsK0sX7WwKRw/1pgslRWTp/lBiX+bgINGnqg
L+fwyRsCgYA0wOhHE629UwC/RknjUIktZlcNR1M8067r5No2YREeGsfjQ4MNHzxT
lvmbnICu7qANQ599Q5IrTc6HmITJJ3lQFHHQ+L196IUGwzZpQmGO3a5slnPWWBOU
rhRlyfDLSTqZQwEGoWSzjfwF9ntpuf6JWnwD7+pLNv7HwhFZzmAlfw==
-----END RSA PRIVATE KEY-----`

/*
-----BEGIN  WUMAN  RSA PUBLIC KEY -----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA6mrG8RLv5TjyCuq5F9hb
CZiHktqJHp1EMIsC92p6SEYyPzSjarF24dO7nxxPxwqt2Ch/8QCRgBubU+gf8D6D
+I9XlpHASfv5yNwkRDgEyV/830JVoQIR1XmJYL/nH7M4Kdj9Flcjbl0rfSyH7ij+
tuzcnFPSkDUJdK8E9VaIdfol5e84sKXEgt9sS3MD0kYlP9RPKPnxgdnN6XEiJpjy
/JxLz/KPQYWe+O7Hn2cg4vaPEl6pp05uxwh7hyAvimOTLM8ts47fFBJJd0OZ+nYX
g7c5u1NbJgoGG28U8ivlnfPV/E6beU3aFkENPYE8mHrgytAX48WV43ZyowL+0EWO
9QIDAQAB
-----END  WUMAN  RSA PUBLIC KEY -----
*/

var key = []byte(privateKey)

func RsaDecrypt(msg string) ([]byte, error) {

	cryptText, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return []byte{}, err
	}

	block, _ := pem.Decode(key)

	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				//log.Println("runtime err:",err,"Check that the key is correct")
				log.QyLogger.Error("Rsa Decrypt runtime err", zap.Error(err.(runtime.Error)))
			default:
				log.QyLogger.Error("Rsa Decrypt  err", zap.Error(err.(error)))
			}
		}
	}()
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return []byte{}, err
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cryptText)
	if err != nil {
		return []byte{}, err
	}
	return plainText, nil
}
