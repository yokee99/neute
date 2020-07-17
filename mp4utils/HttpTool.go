package mp4utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

//HttpClient HttpClient
type HttpClient struct {
	Client http.Client
	AddrIp string
}

func NewHttpClientImage(connTimeout time.Duration, readTimeout time.Duration, newhttpclient *HttpClient) {
	client := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, connTimeout)
				if err != nil {
					return nil, err
				}
				newhttpclient.AddrIp = c.RemoteAddr().String()
				c.SetDeadline(time.Now().Add(readTimeout))
				return c, nil
			},
		},
	}
	newhttpclient.Client = client
}

func (this *HttpClient) Get(httpUrl string, postParams map[string]string, referer bool) ([]byte, error) {
	u, err := url.Parse(httpUrl)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	for key, value := range postParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	url := ""
	if postParams != nil {
		url = u.String()
	} else {
		url = httpUrl
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
	if referer {
		req.Header.Set("referer", "https://www.referer.com") //vipkid使用
	}
	resp, reqErr := this.Client.Do(req)
	if reqErr != nil {
		return nil, reqErr
	}
	defer resp.Body.Close()
	var Codeerr error
	switch resp.StatusCode {
	case 200:
		break
	default:
		Codeerr = errors.New(fmt.Sprintf("Server get request error,statuscode: %v", resp.StatusCode))
		break
	}
	if Codeerr != nil {
		return nil, Codeerr
	}

	data, respErr := ioutil.ReadAll(resp.Body)
	if respErr != nil {
		return nil, respErr
	}
	return data, nil
}

//Post Post
func (this *HttpClient) Post(httpUrl string, headers map[string]string, body string) (string, error) {
	req, _ := http.NewRequest("POST", httpUrl, strings.NewReader(body))
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, reqErr := this.Client.Do(req)
	if reqErr != nil {
		return "", reqErr
	}
	defer resp.Body.Close()
	data, respErr := ioutil.ReadAll(resp.Body)
	if respErr != nil {
		return "", respErr
	}
	return string(data), nil
}

//README：在可执行文件当前文件夹下建立images文件夹用于保存截取出来的图片

//参数1：图片地址

// func mp4utilsmain() {

// dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
// fmt.Println("dir:", dir)
// var ffmpegPath string = dir + "/ffmpeg"
// urlpath := os.Args[1]
// savefilename := dir + "/video.mp4"
// st := time.Now().UnixNano()
// //var freq float64 = 60

// //根据视频地址，下载视频到本地
// if true {
// 	client := HttpClient{}
// 	ConnectTimeout := time.Duration(1000) * time.Nanosecond * 1e9
// 	ReadTimeout := time.Duration(1000) * time.Nanosecond * 1e9
// 	NewHttpClientImage(ConnectTimeout, ReadTimeout, &client)
// 	imgbytes, err := client.Get(urlpath, nil, false)
// 	urlpath = savefilename
// 	ioutil.WriteFile(savefilename, imgbytes, 0777)
// 	fmt.Println("savefilename:", savefilename)
// 	fmt.Println("download err:", err)
// }

// if true {
// 	path := "./images/"
// 	if err := os.MkdirAll(path, os.ModePerm); err != nil {
// 		fmt.Println("mkdir failed!!!")
// 	}
//ffmpeg使用命令: ffmpeg -i http://video.pearvideo.com/head/20180301/cont-1288289-11630613.mp4 -r 1 -t 4 -f image2 image-%05d.jpeg
/*
   -t 代表持续时间，单位为秒
   -f 指定保存图片使用的格式，可忽略。
   -r 指定抽取的帧率，即从视频中每秒钟抽取图片的数量。1代表每秒抽取一帧。
   -ss 指定起始时间
   -vframes 指定抽取的帧数
*/
// videoLen, _ := GenerateLength(ffmpegPath, urlpath)
// testFfmpegParams(urlpath, path, ffmpegPath, 60, videoLen)
//invokeFfmpeg(urlpath, path, ffmpegPath, freq)
//getLastFrame(urlpath, path, ffmpegPath)
// }

// }

//通过-ss参数 获取视频中的图片帧
func testFfmpegParams(url string, path string, ffmpegPath string, freq int, videoLen int) string {
	var outputerror string
	for i := 0; i < videoLen; i = i + freq {
		sec := strconv.Itoa(i)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(50000)*time.Millisecond)
		cmd := exec.CommandContext(ctx, ffmpegPath,
			"-loglevel", "error",
			"-y",
			"-ss", sec,
			"-t", "1",
			"-i", url,
			"-vframes", "1",
			path+"/"+sec+".jpg")
		defer cancel()
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			outputerror += fmt.Sprintf("lastframecmderr:%v;", err)
		}
		if stderr.Len() != 0 {
			outputerror += fmt.Sprintf("lastframestderr:%v;", stderr.String())
		}
		if ctx.Err() != nil {
			outputerror += fmt.Sprintf("lastframectxerr:%v;", ctx.Err())
		}
	}
	return outputerror
}

//GenerateLength 获取视频长度
func GenerateLength(ffprobePath string, urlc string) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	cmd := exec.CommandContext(ctx, ffprobePath,
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		"-i", urlc)
	defer cancel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()
	durationstr := stdout.String()
	stderrstr := stderr.String()
	durationstr = strings.Replace(durationstr, "\n", "", -1)
	stderrstr = strings.Replace(stderrstr, "\n", "", -1)
	length, err := strconv.ParseFloat(durationstr, 64)
	fmt.Printf("%s", stderrstr)
	return int(length), err
}

//获取视频中最后一帧的图片
func getLastFrame(url string, path string, ffmpegPath string) string {
	fmt.Println("url:", url)
	fmt.Println("ffmpeg:", ffmpegPath)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(50000)*time.Millisecond)
	cmd := exec.CommandContext(ctx, ffmpegPath,
		"-timeout", "60000000",
		"-loglevel", "error",
		"-y",
		"-ss", "13",
		"-t", "1",
		"-i", url,
		"-vframes", "1",
		path+"/"+"lastfram.jpg")
	defer cancel()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	var outputerror string
	fmt.Println("lastframpath:", path+"/"+"lastfram.jpg")
	err := cmd.Run()
	fmt.Println("zuihouyzihenerr:", err)
	if err != nil {
		outputerror += fmt.Sprintf("lastframecmderr:%v;", err)
	}
	if stderr.Len() != 0 {
		outputerror += fmt.Sprintf("lastframestderr:%v;", stderr.String())
	}
	if ctx.Err() != nil {
		outputerror += fmt.Sprintf("lastframectxerr:%v;", ctx.Err())
	}
	return outputerror
}

//通过帧率获取视频中的图片和testFfmpegParams函数一样功能.
func invokeFfmpeg(urlpath string, path string, ffmpegPath string, Freq float64) {
	fmt.Println("urlpath:", urlpath)
	currentTime := time.Now().UnixNano()
	//ffmpeg -i 'http://ivi.bupt.edu.cn/hls/cctv1hd.m3u8' -r 1 -t 200 -f image2 images/image-%05d.jpeg  //中央电视台视频流中的图片
	//ffmpeg -i 'http://ivi.bupt.edu.cn/hls/cctv1hd.m3u8' -r 10 -vcodec copy video/aaaaa.mp4    //copy中央电视台的视频流中的视频
	//var Freq float64 = 5 //设置多少秒截一帧
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(14400000)*time.Millisecond)
	//cmd := exec.CommandContext(ctx, ffmpegPath,
	//    "-i", urlpath,
	//  "-vcodec", "copy",
	//  path+"/"+"tttt.mp4")
	cmd := exec.CommandContext(ctx, ffmpegPath,
		"-loglevel", "error",
		"-i", urlpath,
		"-f", "image2",
		"-r", strconv.FormatFloat(float64(1/Freq), 'e', -1, 64),
		path+"/"+"%04d.jpg")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	defer cancel()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	var outputerror string
	err := cmd.Run()
	if err != nil {
		outputerror += fmt.Sprintf("cmderr:%v;", err)
	}
	if stderr.Len() != 0 {
		outputerror += fmt.Sprintf("stderr:%v;", stderr.String())
	}
	if ctx.Err() != nil {
		outputerror += fmt.Sprintf("ctxerr:%v;", ctx.Err())
	}
	cost := float64((time.Now().UnixNano() - currentTime) / 1000000)
	fmt.Println("invokeFfmpeg err:", outputerror)
	fmt.Println("invokeFfmpeg videolengthcost:", cost)
}
