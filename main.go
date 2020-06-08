package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/yokee99/neute/utils"
)

var (
	/*
	 *参数列表
	 */
	count      int
	wg         sync.WaitGroup
	urllist    []string
	fileName   string
	concurrent int
	timeout    int
	h          bool
	ch         chan int
	blockcount int
	finished   int
	failurl    []string
	failurlT   []string
	privateKey string
)

func init() {
	blockcount = 0
	flag.StringVar(&fileName, "c", "", " path  of your URLLIST")
	flag.IntVar(&concurrent, "k", 1, "concurrent")
	flag.IntVar(&timeout, "t", 15, "timeout ")
	flag.StringVar(&privateKey, "K", "", "privateKey")
	flag.Usage = usage

	k := utils.Exist("blocklist")
	d := utils.Exist("video_tmp")
	if k != true {
		file, err := os.Create("blocklist")

		if err != nil {
			fmt.Println(err)
		}
		file.Close()
	}
	if d != true {
		os.Mkdir("./video_tmp", 0777)
	}

}

func main() {
	start := time.Now()

	flag.Parse()
	args := flag.Args()
	if h {
		flag.Usage()
		return
	}

	if fileName != "" {
		fileName := fileName
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		defer file.Close()
		ch = make(chan int, concurrent) /*创建通道（多线程）*/
		fd := bufio.NewReader(file)
		count = 0
		for {
			line, err := fd.ReadString('\n')
			if err != nil {
				if err == io.EOF { //读取结束，会报EOF
					fmt.Println("Read done!")
					break
				}
				break
			}
			line = strings.Replace(line, "\n", "", -1)
			line = strings.Replace(line, "\r", "", -1)
			urllist = append(urllist, line)
			count++

		}
		fmt.Println("Line:", count)

		for i := 0; i < count; i++ {
			wg.Add(1)
			str := "[" + bar((i*10)/count, 10) + "] "
			fmt.Printf("\r%s  %.1f %%  exe: %d finished: %d/%d  block: %d ", str, float32(i)/float32(count)*100, i, finished, count, blockcount)
			ch <- 1
			urlc := urllist[i]
			if privateKey != "" {
				testurl1 := urlc
				urlpath := utils.GetPathInURL(testurl1)
				NowTimestamp := time.Now().Unix()
				EndTimestamp := NowTimestamp + 3600
				sstring := urlpath + "-" + strconv.FormatInt(EndTimestamp, 10) + "-0-0-" + privateKey
				md5str := utils.Md5V(sstring)
				ssurl := urlc + "?auth_key=" + strconv.FormatInt(EndTimestamp, 10) + "-0-0-" + md5str
				go work(ssurl)
			} else {
				go work(urlc)
			}

		}

		failurlT = failurl
		wg.Wait()
		str := "[" + bar((10), 10) + "] "
		fmt.Printf("\r%s  %.1f %%  exe: %d finished: %d/%d  block: %d ", str, float32(count)/float32(count)*100, count, finished, count, blockcount)
		fmt.Printf("\r\nDone!")
		fmt.Println()

	} else {
		if len(args) < 1 {
			fmt.Println("Too few arguments")
			fmt.Println("Usage: neute  [args] URLs...")
			flag.PrintDefaults()

		} else if len(args) == 1 {
			wg.Add(1)

			if privateKey != "" {
				testurl1 := flag.Arg(0)
				urlpath := utils.GetPathInURL(testurl1)
				NowTimestamp := time.Now().Unix()
				EndTimestamp := NowTimestamp + 3600
				sstring := urlpath + "-" + strconv.FormatInt(EndTimestamp, 10) + "-0-0-" + privateKey
				md5str := utils.Md5V(sstring)
				ssurl := flag.Arg(0) + "?auth_key=" + strconv.FormatInt(EndTimestamp, 10) + "-0-0-" + md5str
				singlework(ssurl)
			} else {
				singlework(flag.Arg(0))
			}

		}
	}

	end := time.Now()
	during := end.Sub(start)
	fmt.Println(during)

}

func singlework(urlc string) {

	defer wg.Done()
	num := rand.Int31n(1)
	time.Sleep(time.Duration(num) * time.Second)

	filename, ext, err := utils.GetNameAndExt(urlc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "log message: %s", err)
		return
	}
	pathPre := "./video_tmp/"
	path := pathPre + filename + "." + ext + ".tmp"
	downloadPro(urlc, path)

}
func work(urlc string) {

	defer wg.Done()
	num := rand.Int31n(1)
	time.Sleep(time.Duration(num) * time.Second)

	filename, ext, err := utils.GetNameAndExt(urlc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "log message: #1%s", err)
		return
	}
	pathPre := "./video_tmp/"
	path := pathPre + filename + "." + ext + ".tmp"
	downloadPro(urlc, path)
	<-ch

}

func downloadPro(url string, path string) {
	chx := make(chan string)
	go func() {
		out, err := os.Create(path)

		if err != nil {
			// fmt.Println("ERROR CODE: #1")
			fmt.Fprintf(os.Stderr, "log message: #2%s", err)
			return
		}
		defer out.Close()

		resp, err := http.Get(url)
		if err != nil {
			// fmt.Println("ERROR CODE: #2")
			fmt.Fprintf(os.Stderr, "log message: #3")
			time.Sleep(3 * time.Second)
			resp, err = http.Get(url)

			if err != nil {
				// fmt.Println("ERROR CODE: #11")
				fmt.Fprintf(os.Stderr, "log message:#33 %s", url[0:72])
				failurl = append(failurl, url)
				return
			}
			defer resp.Body.Close()
		}

		resCode := resp.StatusCode
		if resCode == 200 {

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				// fmt.Println("ERROR CODE: #3")
				if err == io.ErrUnexpectedEOF { //读取结束，会报EOF
					fmt.Fprintf(os.Stderr, "log message: #5 %s", url[0:72])
					return
				}
				fmt.Fprintf(os.Stderr, "log message: #5 %s", err)
				return

			}

		} else {

			blockcount++
			err := utils.AppendToFile("blocklist", url+"\n")
			fmt.Fprintf(os.Stderr, "log message: #%d %s", resCode, url[0:72])
			if err != nil {
				// fmt.Println("ERROR CODE: #4")
				fmt.Fprintf(os.Stderr, "log message: #6 %s", err)
				return
			}
		}

		chx <- ":)"

	}()

	select {
	case res := <-chx:
		fmt.Printf("%s", res)
	case <-time.After(time.Second * time.Duration(timeout)):
		// fmt.Println(path + " done!")
		finished++
		// del_ch <- path
		// fmt.Println("del done!")

	}

}

func usage() {
	fmt.Fprintf(os.Stderr, ` version: 1.0.0.1
	Usage: neute  [opts][args] URLs...

Options:
`)
	flag.PrintDefaults()
	fmt.Println(`
Examples:
	neute -c YOURFILEPATH -k 5 -t 30 
	(URL list is "YOURFILEPATH"  concurrent is 5 , download 30s! )`)
}

func bar(count, size int) string {
	str := ""
	for i := 0; i < size; i++ {
		if i < count {
			str += "#"
		} else {
			str += " "
		}
	}
	return str
}

// func del_buff() {
// 	del_path := <-del_ch
// 	err := os.RemoveAll(del_path)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// }
