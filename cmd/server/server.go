package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type probe struct {
	timeout int
}

type request struct {
	Ip   string   `json:"ip"`
	Port []string `json:"port"`
}

type requests []*request

func (r *requests) parseUri() ([]string, error) {
	res := make([]string, 0)

	for index := range *r {
		req := (*r)[index]
		for y := range req.Port {
			port := req.Port[y]
			if strings.Contains(port, "..") {
				ps := strings.Split(port, "..")
				if len(ps) < 2 {
					break
				}
				start, err := strconv.ParseInt(ps[0], 10, 64)
				if err != nil {
					return nil, err
				} else if start > 65535 {
					return nil, err
				}

				end, err := strconv.ParseInt(ps[1], 10, 64)
				if err != nil {
					return nil, err
				} else if end > 65535 {
					return nil, err
				}
				for ; end >= start; end-- {
					res = append(res, fmt.Sprintf("%s:%d", req.Ip, end))
				}
			} else {
				tmpPort, err := strconv.ParseInt(port, 10, 64)
				if err != nil {
					return nil, err
				}
				res = append(res, fmt.Sprintf("%s:%d", req.Ip, tmpPort))
			}
		}
	}
	return res, nil
}

type response struct {
	Ip        string `json:"ip"`
	Port      string `json:"port"`
	IsReached int    `json:"is_reached"`
	ReplyTime int64  `json:"reply_time"`
}

func (p *probe) ping() {

}

func (p *probe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	datas, exist := values["data"]
	if !exist {
		fmt.Fprintf(w, `{"error_code":%d,"msg":%s}`, 551, "请求参数错误")
		return
	}
	if len(datas) == 0 || len(datas) > 1 {
		fmt.Fprintf(w, `{"error_code":%d,"msg":%s}`, 552, "请求参数过多或不匹配")
		return
	}
	reqs := make(requests, 0)
	err := json.Unmarshal([]byte(datas[0]), &reqs)
	if err != nil {
		fmt.Fprintf(w, `{"error_code":%d,"msg":%s,"data":%s}`, 553, "请求参数无法反序列化", datas[0])
		return
	}

	addrs, err := reqs.parseUri()
	if err != nil {
		fmt.Fprintf(w, `{""error_code"":%d,"msg":%s}`, 554, "请求参数转换uri错误")
		return
	}

	resp, err := pingTimeout(time.Duration(p.timeout), addrs...)
	if err != nil {
		fmt.Fprintf(w, `{""error_code"":%d,"msg":%s}`, 555, "PING发生错误")
		return
	}
	bs, err := json.Marshal(&resp)
	if err != nil {
		fmt.Fprintf(w, `{""error_code"":%d,"msg":%s}`, 556, "PING返回response序列化错误")
		return
	}
	fmt.Fprintf(w, `{"data":%s}`, bs)
}

func RunServer(host string, port string, timeout int) error {
	httpServer := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", host, port),
		Handler:        &probe{timeout},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("start http server on %s:%s", host, port)
	return httpServer.ListenAndServe()
}

func pingTimeout(s time.Duration, addrs ...string) ([]response, error) {
	res := make([]response, 0)
	mu := sync.Mutex{}

	additionRes := func(addr string, isReached bool, rt time.Duration) {
		mu.Lock()
		defer mu.Unlock()
		if !strings.Contains(addr, ":") {
			return
		}
		ipaddr := strings.Split(addr, ":")
		if len(addr) < 2 {
			return
		}
		resp := response{
			Ip:        ipaddr[0],
			Port:      ipaddr[1],
			ReplyTime: int64(rt),
		}
		if isReached {
			resp.IsReached = 1
		}
		res = append(res, resp)
	}

	diag := func(addr string) bool {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return false
		}
		defer conn.Close()
		return true
	}

	pingWithContext := func(ctx context.Context, addr string, f *func(addr string, isReached bool, rt time.Duration)) {
		timestart := time.Now().Unix()
		select {
		case <-ctx.Done():
			(*f)(addr, diag(addr), time.Duration(time.Now().Unix()-timestart))
			return
		default:
			(*f)(addr, diag(addr), time.Duration(time.Now().Unix()-timestart))
		}
	}

	wg := sync.WaitGroup{}
	for _, addr := range addrs {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), s)
			defer cancel()
			pingWithContext(ctx, addr, &additionRes)
		}(addr)
	}
	wg.Wait()

	return res, nil
}
