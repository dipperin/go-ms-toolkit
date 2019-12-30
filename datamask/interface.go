package datamask

type Mask interface {
	LastName(data string) string
	FirsName(data string) string
	IDPan(data string) string
	IDAadhaar(data string) string
	BankCard(data string) string
	Password(data string) string
	Mail(data string) string
	Phone(data string) string
}
