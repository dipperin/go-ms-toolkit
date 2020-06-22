package datamask

import (
	"strings"
)

type maskImpl struct{}

func newMaskImpl() *maskImpl {
	return &maskImpl{}
}

// 判断长度，大于3时，掩码前3位，小于3不脱敏
func (m *maskImpl) LastName(data string) string {
	if len(data) <= 3 {
		return data
	}

	return "***" + data[3:]
}

// 判断长度，大于3时，掩码前3位，小于3不脱敏
func (m *maskImpl) FirstName(data string) string {
	return m.LastName(data)
}

func maskMid4(data string) string {
	if len(data) > 9 {
		return data[0:4] + "****" + data[8:]
	}
	return data
}

// 保留前四，掩码中间4位
func (m *maskImpl) IDPan(data string) string {
	return maskMid4(data)
}

// 保留前四，掩码中间4位
func (m *maskImpl) IDAadhaar(data string) string {
	return maskMid4(data)
}

// 保留前四，掩码中间4位
func (m *maskImpl) BankCard(data string) string {
	return maskMid4(data)
}

// 密码打印为null
func (m *maskImpl) Password(data string) string {
	return "null"
}

// 保留邮箱前缀前三位
func (m *maskImpl) Mail(data string) string {
	if !strings.Contains(data, "@") {
		return data
	}

	// 通过@切割
	s := strings.Split(data, "@")
	if size := len(s[0]); size > 3 {
		tmp := s[0][0:3]

		for i := 0; i < size-3; i++ {
			tmp += "*"
		}

		return tmp + "@" + s[1]
	}

	return data
}

// 保留前三，掩码中间4位
func (m *maskImpl) Phone(data string) string {
	if len(data) > 7 {
		return data[0:3] + "****" + data[7:]
	}
	return data
}
