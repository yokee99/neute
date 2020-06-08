package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"strings"
)

//Exist Exist
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//Md5V md5V
func Md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//GetNameAndExt 获取文件名和后缀
func GetNameAndExt(uri string) (string, string, error) {
	/*
	*获取文件名和后缀
	 */
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", "", err
	}
	s := strings.Split(u.Path, "/")
	filename := strings.Split(s[len(s)-1], ".")
	if len(filename) > 1 {
		return filename[0], filename[1], nil
	}
	return filename[0], ".dowload", nil
}

//GetPathInURL GetPathInURL
func GetPathInURL(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	return u.Path
}

//AppendToFile AppendToFile
func AppendToFile(fileName string, content string) error {

	/*
	*追加至文件
	 */
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	// 以只写的模式，打开文件
	if err != nil {
		fmt.Println(" file create failed. err: " + err.Error())

	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
	return err
}
