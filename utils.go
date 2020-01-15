package fmtlogwrapper

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

/*
Open a file and ensure if it doesnt exist @ the path, then create the path Dir + file
*/
func OpenFilePathExists(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("\n[Logger][getWriter] File not found creating file @ %s: %v\n", path, err)
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
		_, err := os.Create(path)
		if err != nil {
		fmt.Printf("\n[Logger][getWriter] error creating log path %s: %v\n", path, err)
			return nil,err
	}
}
	return f, err
}

func GetType(obj interface{}) string {
	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		return /* "*" + */ t.Elem().Name() //* represents pointer but since this goes back to caller removing the distinction as no one cares
	} else {
		return t.Name()
	}
}
