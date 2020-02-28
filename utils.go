package fmtlogwrapper

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

/*
Open a file and ensure if it doesnt exist @ the path, then create the path Dir + file

@return *os.File is success else nill, error
*/
func OpenFilePathExists(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("\n[Logger][getWriter] File not found creating file @ %s: %v\n", path, err)
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
		_, err := os.Create(path)
		if err != nil {
			fmt.Printf("\n[Logger][getWriter] error creating file path %s: %v\n", path, err)
			return nil, err
		}
	}
	//verify the file is created (optional)
	_, e2 := os.Stat(path)
	if os.IsNotExist(e2) {
		fmt.Printf("\n[Logger][getWriter] Tried creating file @ path but failed %s: %v\n", path, err)
		return nil, e2
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
