package log

import (
	"github.com/stretchr/testify/assert"
	"os/user"
	"testing"
	"time"

	"go.uber.org/zap"
)

type Tst struct {
	A      string    `json:"a"`
	Number uint      `json:"number"`
	Time   time.Time `json:"time"`
}

func TestInitLogger(t *testing.T) {
	QyLogger.Debug("test init debug log", zap.String("a", "asss"))

	tst := &Tst{A: "fff", Number: 44}
	QyLogger.Debug("test debug log struct", zap.Any("tst", tst))

	tst1 := &Tst{A: "fff", Number: 44, Time: time.Now()}
	QyLogger.Debug("test debug log struct", zap.Any("tst", tst1))

	InitLogger(zap.DebugLevel, "", "", true)

	QyLogger.Debug("aaaa")
	LoggerEnd()
}

func TestLoggerEnd(t *testing.T) {
	InitLogger(zap.DebugLevel, "", "", true)

	QyLogger.Debug("aaaa")

	LoggerEnd()
}

func Test_newLogConfig(t *testing.T) {
	assert.NotNil(t, newLogConfig(zap.DebugLevel, "/tmp", "", true))
}

func Test_getOutPaths(t *testing.T) {
	assert.Len(t, getOutPaths("/tmp", ""), 0)

	arr := getOutPaths("/tmp", "11.log")
	assert.Len(t, arr, 1)

	println(arr[0])
}

func Test_pathExists(t *testing.T) {
	usr, err := user.Current()
	assert.NoError(t, err)
	assert.Equal(t, true, pathExists(usr.HomeDir))
	assert.Equal(t, false, pathExists("fdsafds"))
}
