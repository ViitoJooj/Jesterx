package logger

import "fmt"

func (l *Log) Print() {
	fmt.Printf("\n[%s] \n[%s] \n[%s]", l.Time, l.TraceID, l.Function)
}
