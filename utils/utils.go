package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"runtime"
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

//ByteCountIEC ByteCountIEC
func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

// 前景 背景 颜色
// ---------------------------------------
// 30  40  黑色
// 31  41  红色
// 32  42  绿色
// 33  43  黄色
// 34  44  蓝色
// 35  45  紫红色
// 36  46  青蓝色
// 37  47  白色
//
// 代码 意义
// -------------------------
//  0  终端默认设置
//  1  高亮显示
//  4  使用下划线
//  5  闪烁
//  7  反白显示
//  8  不可见

//PrintColorTable PrintColorTable
func PrintColorTable() {
	if runtime.GOOS != "windows" {
		for b := 40; b <= 47; b++ { // 背景色彩 = 40-47
			for f := 30; f <= 37; f++ { // 前景色彩 = 30-37
				for d := range []int{0, 1, 4, 5, 7, 8} { // 显示方式 = 0,1,4,5,7,8
					fmt.Printf(" %c[%d;%d;%dm%s(f=%d,b=%d,d=%d)%c[0m ", 0x1B, d, b, f, "", f, b, d, 0x1B)
				}
				fmt.Println("")
			}
			fmt.Println("")
		}
	}

}

//ErrorString RedString
func ErrorString(s string) string {
	if runtime.GOOS != "windows" {
		return "\033[;31;m" + s + "\033[0m\n"
	}
	return s
}

//SuccessString GreanString
func SuccessString(s string) string {
	if runtime.GOOS != "windows" {
		return "\033[;32;m" + s + "\033[0m\n"
	}
	return s

}

//InfoString BlueString
func InfoString(s string) string {
	if runtime.GOOS != "windows" {
		return "\033[;34;m" + s + "\033[0m\n"
	}
	return s

}
