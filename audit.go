package fmtlogwrapper

/*
An Audit wrapper over the logger to provide some standard OOB Auditing features.


@author Arjun Dhar
*/

import (
	"encoding/json"
	"fmt"
	"time"

	"crypto/sha1"
	"encoding/base64"
)

func Now() string {
	t := time.Now()
	when := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	return when
}

/*
 Structured log / audit - Consumed by System
 @Default converts to JSON string
*/
func Audit(any interface{}) string {
	// type == api.StatsReport, write to DB

	b, err := json.Marshal(any)
	if err != nil {
		return fmt.Sprintf("{'audit-error':'[audit] %s for %v'}", err.Error(), any)
	}
	return string(b)
}

/*
 @param namespace string - will decide the file name

 @param processID string (optional)- file name depends on this

 @param taskId string (optional) - Allows multiple processes to write to same file;
 however if we plan to use buffered mode, then each can independently maintain its own buffer.
 Use "" for common.

 @param path string (optional). Defaults to ".\log"
*/
func InitAuditLogger(namespace string, processID string, taskId string, path string) *Logger {
	if path == "" {
		path = ".\\log\\"
	}
	logger := NewLogger(LogSettings{
		FilePath: path + namespace + "-" + processID + ".json",
	})
	ContextLoggers()[processID+"."+taskId] = logger
	return logger
}

/*
 @param namespace string - will decide the file name

 @param processID string (optional)- file name depends on this

 @param taskId string (optional) - Allows multiple processes to write to same file;
 however if we plan to use buffered mode, then each can independently maintain its own buffer.
 Use "" for common.

 @param path string (optional). Defaults to ".\log"
*/
func AuditLogger(namespace string, processID string, taskId string, path string) *Logger {
	logger := ContextLoggers()[processID+"."+taskId]
	if logger == nil {
		logger = InitAuditLogger(namespace, processID, taskId, path)
	}
	//Common settings for most audit operations
	logger.WriteToBuffer(true)
	logger.MuteConsole(true)
	logger.Settings.JobRecMeta = namespace
	return logger
}

func AuditToDiskAndClose(audit *Logger) {
	audit.MuteWrite(true) // Dont want the close statement logged
	audit.Close()         //auto-commits buffer
	audit.MuteWrite(false)
}

func GenSHA(any interface{}) string {
	hasher := sha1.New()
	b, _ := json.Marshal(any)
	hasher.Write(b)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}
