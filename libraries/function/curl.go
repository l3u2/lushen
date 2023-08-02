package function

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type CmdModel struct {
	MessageId string `json:"messageId"`
	Cmd       string `json:"cmd"`
	Timestamp string `json:"timestamp"`
	Code      int    `json:"code"`
	Msg       string `json:"message"`
	Data      string `json:"data"`
}

type RespModel struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	ServerTime string `json:"serverTime"`
}

func Response(url string, cmd string, timestamp string, data string, messageId string, code int, msg string) {
	var model CmdModel
	model.MessageId = messageId
	model.Cmd = cmd
	model.Timestamp = timestamp
	model.Code = code
	model.Msg = msg
	model.Data = data
	b, err := json.Marshal(model)
	if err != nil {
		logs.Error("序列化数据失败，error:", err)
		return
	}

	//向url发送post请求
	res, err := httpPost(url, string(b))
	if err != nil || !res {
		logs.Error("PHP回调应答失败，error：", err)
	} else {
		logs.Info("PHP回调应答成功")
	}

	return
}

func httpPost(url string, data string) (bool, error) {
	if len(url) == 0 {
		return false, errors.New("url is empty")
	}
	defer func() {
		if errtmp := recover(); errtmp != nil {
			errRet := fmt.Errorf("%v", errtmp)
			log.Println("请求服务端失败" + errRet.Error() + NowMicro())
		}
	}()
	content, err := sendPostRetry(url, data, 2, 60)
	if err != nil {
		log.Println("获取服务端返回结果失败" + err.Error())
		return false, err
	}

	log.Println("获取服务端返回结果成功 " + string(content) + NowMicro())

	var model RespModel
	err = json.Unmarshal(content, &model)
	if err != nil {
		log.Println("JSON解析失败" + err.Error())
		return false, err
	}
	if model.Code == 1 {
		return true, nil
	} else {
		log.Println("服务端错误" + model.Message)
		return false, errors.New(model.Message)
	}
}

func sendPostRetry(postUrl string, reqBody string, attempts int, timeout int64) ([]byte, error) {
	var errtmp error
	for index := 0; index < attempts; index++ {
		res, err := sendPost(postUrl, reqBody, timeout)
		errtmp = err
		if err == nil {
			return res, nil
		}
	}
	if errtmp != nil {
		return nil, errors.New("SendRetry err:" + errtmp.Error())
	} else {
		return nil, nil
	}
}

func sendPost(postUrl string, reqBody string, timeout int64) (response []byte, errRet error) {
	response = nil
	errRet = nil
	//可以通过client中transport的Dail函数,在自定义Dail函数里面设置建立连接超时时长和发送接受数据超时
	defer func() {
		if errtmp := recover(); errtmp != nil {
			errRet = fmt.Errorf("%v", errtmp)
			log.Println("sendPost defer error :" + errRet.Error() + NowMicro())
		}
	}()
	log.Println("请求服务端开始：" + NowMicro())
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration((timeout)))
				if err != nil {
					return nil, err
				}
				_ = conn.SetDeadline(time.Now().Add(time.Second * time.Duration((timeout))))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * time.Duration(timeout),
		},
	}
	//提交请求;用指定的方法，网址，可选的主体放回一个新的*Request
	requestDo, err := http.NewRequest("POST", postUrl, strings.NewReader(reqBody))
	if err != nil {
		log.Println("请求服务端失败:" + err.Error() + NowMicro())
		return nil, errors.New("httpPost err:" + err.Error())
	}
	requestDo.Header.Set("Content-Type", "application/json")
	//前面预处理一些参数，状态，Do执行发送；处理返回结果;Do:发送请求,
	log.Println("请求服务端DO：" + NowMicro())
	res, err := client.Do(requestDo)
	if nil != err {
		log.Println("请求服务端失败: " + err.Error() + NowMicro())
		return nil, errors.New("httpPost err:" + err.Error())
	}
	defer res.Body.Close()
	log.Println("请求服务端Read：" + NowMicro())
	data, err := ioutil.ReadAll(res.Body)
	log.Println("请求服务端ReadOver：" + NowMicro())
	if nil != err {
		log.Println("读取数据失败：" + errRet.Error() + NowMicro())
		return nil, errors.New("ReadAll err:" + err.Error())
	}
	return data, nil
}

func CallbackRequest(filename string, targetUrl string, params map[string]string) error {
	logs.Info("callback filename:", filename)
	logs.Info("callback address:", targetUrl)
	logs.Info("callback params:%#v", params)
	bodyBuf := &bytes.Buffer{}
	bodyWrite := multipart.NewWriter(bodyBuf)
	if params != nil {
		for k, v := range params {
			bodyWrite.WriteField(k, v)
		}
	}
	fileWriter, err := bodyWrite.CreateFormFile("uploadfile", filename)
	if err != nil {
		logs.Error("CreateFormFile Error:", err.Error())
		return err
	}
	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		logs.Error("Open Error:", err.Error())
		return err
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		logs.Error("Copy Error:", err.Error())
		return err
	}

	contentType := bodyWrite.FormDataContentType()
	err = bodyWrite.Close()
	if err != nil {
		return err
	}

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		logs.Error("Post Error:", err.Error())
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("ReadAll Error:", err.Error())
		return err
	}

	var model RespModel
	err = json.Unmarshal(respBody, &model)
	if err != nil {
		logs.Error("Unmarshal Error:", err.Error())
		return err
	}
	if model.Code == 1 {
		return nil
	} else {
		logs.Error("Http Response Error:", model.Message)
		return errors.New(model.Message)
	}
}
