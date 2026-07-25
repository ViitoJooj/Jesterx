package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Log struct {
	Erro     error
	TraceID  string
	Function string
	Time     string
}

func Warn(erro error) *Log {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)

	log := &Log{
		TraceID:  uuid.NewString(),
		Function: fn.Name(),
		Time:     time.Now().Format("02/01/2006 15:04:05"),
	}

	var erroTexto string

	if erro != nil {
		erroTexto = fmt.Sprintf(" [ERROR: %s]", erro.Error())
	}

	line := fmt.Sprintf(
		"[%s] [%s] [%s]%s\n",
		log.Time,
		log.TraceID,
		log.Function,
		erroTexto,
	)

	const maxLines = 200
	filePath := "./logs.txt"

	content, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return log
	}

	var lines []string

	if len(content) > 0 {
		lines = strings.Split(
			strings.TrimSpace(string(content)),
			"\n",
		)
	}

	if len(lines) >= maxLines {
		lines = lines[1:]
	}

	lines = append(lines, strings.TrimSpace(line))

	err = os.WriteFile(
		filePath,
		[]byte(strings.Join(lines, "\n")+"\n"),
		0644,
	)
	if err != nil {
		return log
	}

	return log
}

func (l *Log) Print() {
	fmt.Printf("\n[%s] \n[%s] \n[%s]", l.Time, l.TraceID, l.Function)
}
