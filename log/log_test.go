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

	QyLogger.Debug("ddddddd")

	QyLogger.Info("iiiiiiiiii")

	QyLogger.Warn("wwwwwwwwwww")

	QyLogger.Error("eeeeeee")

	LoggerEnd()
}

func TestLoggerEnd(t *testing.T) {
	InitLogger(zap.DebugLevel, "", "", true)

	QyLogger.Debug("console.....dddd")

	QyLogger.Info("console.....iiii")

	QyLogger.Warn("console....wwww")

	QyLogger.Error("console.....eeee")

	LoggerEnd()
}

func Test_newLogCore(t *testing.T) {
	assert.NotNil(t, newLogCore(zap.DebugLevel, "/tmp", "", true, false))
}

func Test_newLogOptions(t *testing.T) {
	assert.True(t, 1 == len(newLogOptions()))
}

func Test_getOutPath(t *testing.T) {
	assert.True(t, getOutPath("/tmp", "") == "")

	assert.True(t, getOutPath("/tmp", "11.log") == "/tmp/11.log")

}

func Test_pathExists(t *testing.T) {
	usr, err := user.Current()
	assert.NoError(t, err)
	assert.Equal(t, true, pathExists(usr.HomeDir))
	assert.Equal(t, false, pathExists("fdsafds"))
}

func TestTrace(t *testing.T) {
	InitLoggerWithCaller(zap.DebugLevel, "", "", true)

	QyLogger.Info("aaaaaaa")

	//fmt.Println(printMyName())

	//QyLogger.Named()
}

func BenchmarkLogCaller(b *testing.B) {
	//  119269	      8697 ns/op
	InitLoggerWithCaller(zap.DebugLevel, "", "", true)
	for i := 0; i < b.N; i++ {
		QyLogger.Info("aaaaaaa")
	}
}

func BenchmarkNoCaller(b *testing.B) {
	// 219262	      5611 ns/op
	InitLogger(zap.DebugLevel, "", "", true)
	for i := 0; i < b.N; i++ {
		QyLogger.Info("aaaaaa")
	}
}
