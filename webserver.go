package webserver

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

//Store 内部数据存储结构
var Store struct {
	sync.Mutex
	Data map[string]int
}

//StartServer 启动webserver
func StartServer() error {
	Store.Data = make(map[string]int,0)
	mux := http.NewServeMux()
	mux.HandleFunc("/string",handler)
	return http.ListenAndServeTLS(":10443","./certs/server.pem","./certs/server.key",mux)
}

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	data := []string{}
	
	err = json.Unmarshal(body,&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	exist := make([]bool,len(data))
	Store.Lock()
	for i,v := range data {
		if _,ok := Store.Data[v];ok {
			exist[i] = true
		} else {
			Store.Data[v]++
			exist[i] = false
		}
	}
	Store.Unlock()
	res, err := json.Marshal(exist)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(res)
}

//StringExist 客户端测试函数
func StringExist(req []string) ([]bool, error) { 
    res := []bool{}
	clienCert := "./certs/client.pem"
	byteCert, err := ioutil.ReadFile(clienCert)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(byteCert)
	cert, err := tls.LoadX509KeyPair("./certs/server.pem", "./certs/server.key")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true,
			Certificates: []tls.Certificate{cert},
			ClientCAs:      pool},
		DisableCompression: true,
	}
    client := &http.Client{Transport: tr}
    data, err := json.Marshal(req)
    if err != nil {
        log.Println(err)
        return nil, err
    }
	resp, err := client.Post("https://localhost:10443/string","application/json",bytes.NewReader(data))
	if err != nil {
        log.Println(err)
        return nil, err
	}
    if resp.StatusCode != http.StatusOK {
        log.Println(resp.Status)
        return nil, errors.New(resp.Status)
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Println(err)
        return nil, err
    }
    err = json.Unmarshal(body,&res)
    if err != nil {
        log.Println(err)
        return nil, err
	}
    return res, nil
}