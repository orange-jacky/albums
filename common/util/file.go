package util

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

//判断文件或者目录是否存在
func IsFileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

func IsDirExist(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func WriteTempFile(path, prefix, content string) error {
	f, err := ioutil.TempFile(path, prefix)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err1 := io.WriteString(f, content)
	if err1 != nil {
		return err1
	}
	return nil
}

func WriteToFile(filename string, content []string) error {
	exist := IsFileExist(filename)
	var fp *os.File
	if exist {
		//os.O_TRUNC 覆盖写
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		fp = f
		if err != nil {
			return err
		}
	} else {
		f, err := os.Create(filename)
		fp = f
		if err != nil {
			return err
		}
	}
	defer fp.Close()
	for _, line := range content {
		if _, err := fp.WriteString(line); err != nil {
			return err
		}
	}
	return nil
}

func ReadFromFile(filename string) ([]string, error) {
	exist := IsFileExist(filename)
	var ret []string
	var fp *os.File
	if exist == false {
		return nil, errors.New("file not exist")
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		} else {
			fp = f
			defer fp.Close()
			scanner := bufio.NewScanner(fp)
			for scanner.Scan() {
				ret = append(ret, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return ret, err
			}
		}
	}
	return ret, nil
}
