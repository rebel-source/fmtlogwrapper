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
 For Context buffer under one context defined typically by processId + taskId;
 however if there is no process or task, then default to namespace.
*/
func getContextId(namespace string, processId string, taskId string) string {
	ctxId := processId + "." + taskId
	if processId == "" && taskId == "" {
		ctxId = namespace
	}
	return ctxId
}

/*
 @param namespace string - will decide the file name

 @param processId string (optional)- file name depends on this

 @param taskId string (optional) - Allows multiple processes to write to same file;
 however if we plan to use buffered mode, then each can independently maintain its own buffer.
 Use "" for common.

 @param path string (optional). Defaults to ".\log"
*/
func InitAuditLogger(namespace string, processId string, taskId string, path string) *Logger {
	if path == "" {
		path = ".\\log\\"
	}
	pid := processId
	if processId != "" {
		// Address dangling -, if blank
		pid = "-" + processId
	}
	logger := NewLogger(LogSettings{
		FilePath: path + namespace + pid + ".json",
	})
	ContextLoggers()[getContextId(namespace, processId, taskId)] = logger
	return logger
}

/*
 @param namespace string - will decide the file name

 @param processId string (optional)- file name depends on this

 @param taskId string (optional) - Allows multiple processes to write to same file;
 however if we plan to use buffered mode, then each can independently maintain its own buffer.
 Use "" for common.

 @param path string (optional). Defaults to ".\log"
*/
func AuditLogger(namespace string, processId string, taskId string, path string) *Logger {
	logger := ContextLoggers()[getContextId(namespace, processId, taskId)]
	if logger == nil {
		logger = InitAuditLogger(namespace, processId, taskId, path)
	} else {
		logger.Open()
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
