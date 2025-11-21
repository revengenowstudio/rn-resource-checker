package log

import (
	"log"
)

func Init() {
	// log.Println("init logging instance")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)
}

func Debug(v ...any) { log.Println(append([]any{"[DEBUG]"}, v...)...) }
func Info(v ...any)  { log.Println(append([]any{"[INFO]"}, v...)...) }
func Warn(v ...any)  { log.Println(append([]any{"[WARN]"}, v...)...) }
func Error(v ...any) { log.Println(append([]any{"[ERROR]"}, v...)...) }
func Fatal(v ...any) { log.Fatal(append([]any{"[FATAL]"}, v...)...) }
