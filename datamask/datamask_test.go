package datamask

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_maskImpl_LastName(t *testing.T) {
	impl := &maskImpl{}

	assert.Equal(t, "***AKUMAR", impl.LastName("SIVAKUMAR"))
	assert.Equal(t, "HAH", impl.LastName("HAH"))
	assert.Equal(t, "HA", impl.LastName("HA"))
	assert.Equal(t, "***M", impl.LastName("WRNM"))
}

func Test_maskImpl_FirsName(t *testing.T) {
	impl := &maskImpl{}

	assert.Equal(t, "***AKUMAR", impl.FirsName("SIVAKUMAR"))
	assert.Equal(t, "HAH", impl.FirsName("HAH"))
	assert.Equal(t, "HA", impl.FirsName("HA"))
	assert.Equal(t, "***M", impl.FirsName("WRNM"))
}

func Test_maskImpl_IDPan(t *testing.T) {
	impl := &maskImpl{}

	assert.Equal(t, "BOEP****1O", impl.IDPan("BOEPP2301O"))
	assert.Equal(t, "1234****9O", impl.IDPan("123456789O"))
}

func Test_maskImpl_IDAadhaar(t *testing.T) {
	impl := &maskImpl{}

	assert.Equal(t, "2592****7048", impl.IDAadhaar("259243927048"))
	assert.Equal(t, "1234****9O12", impl.IDAadhaar("123456789O12"))
}

func Test_maskImpl_BankCard(t *testing.T) {
	impl := &maskImpl{}

	assert.Equal(t, "3328****470", impl.BankCard("33283804470"))
	assert.Equal(t, "1234****9O1", impl.BankCard("123456789O1"))
}

func Test_maskImpl_Password(t *testing.T) {
	impl := &maskImpl{}

	assert.Equal(t, "null", impl.Password("33283804470"))
	assert.Equal(t, "null", impl.Password("123456789O1"))
}

func Test_maskImpl_Phone(t *testing.T) {
	impl := &maskImpl{}

	assert.Equal(t, "850****664", impl.Phone("8500346664"))
}

func Test_maskImpl_Mail(t *testing.T) {
	impl := &maskImpl{}
	assert.Equal(t, "rpr*************@gmail.com", impl.Mail("rprabhakaran6664@gmail.com"))
}