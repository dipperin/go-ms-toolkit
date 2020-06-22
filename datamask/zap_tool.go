package datamask

import (
	"go.uber.org/zap"
	"sync"
)

const (
	Phone     = "phone"
	Mail      = "mail"
	Password  = "password"
	FirstName = "first_name"
	LastName  = "last_name"
	IDPan     = "id_pan"
	BankCard  = "bank_card"
	IDAadhaar = "id_aadhaar"
)

var (
	mask *maskImpl
	once sync.Once
)

func GetMark() Mask {
	if mask == nil {
		once.Do(func() {
			mask = newMaskImpl()
		})
	}

	return mask
}

// todo 这个可以做成一个map
func String(maskType string, data string) zap.Field {
	switch maskType {
	case Phone:
		data = GetMark().Phone(data)
	case Mail:
		data = GetMark().Mail(data)
	case Password:
		data = GetMark().Password(data)
	case FirstName:
		data = GetMark().FirstName(data)
	case LastName:
		data = GetMark().LastName(data)
	case IDPan:
		data = GetMark().IDPan(data)
	case BankCard:
		data = GetMark().BankCard(data)
	case IDAadhaar:
		data = GetMark().IDAadhaar(data)
		//default:

	}

	return zap.String(maskType, data)
}
