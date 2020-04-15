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
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600) //0666
	if err != nil {
		fmt.Printf("\n[OpenFilePathExists] File not found creating file @ %s: %v\n", path, err)
		parentPath := filepath.Dir(path)
		if err := os.MkdirAll(parentPath, os.ModePerm); err != nil {
			fmt.Printf("\n[OpenFilePathExists] Parent path could not be created @ %v\n", err)
		}
		_, err = os.Create(path)
		if err != nil {
			fmt.Printf("\n[OpenFilePathExists] error creating file path %s: %v\n", path, err)
			return nil, err
		} else {
			//verify the file is created (optional)
			if !DoesFileExist(path) {
				fmt.Printf("\n[OpenFilePathExists] Tried creating file @ path but failed %s: %v\n", path, err)
				return nil, err
			} else {
				// Fresh connection				
				if e := f.Close(); e!=nil {
					fmt.Println("\n[OpenFilePathExists][ERROR]", e)
				}				
				f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
			}
		}
	}
	return f, err
}

func DoesFileExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func GetType(obj interface{}) string {
	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		return /* "*" + */ t.Elem().Name() //* represents pointer but since this goes back to caller removing the distinction as no one cares
	} else {
		return t.Name()
	}
}
