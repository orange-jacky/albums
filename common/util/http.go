package util

import (
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	HTTP_DIALER_TIMEOUT = 5 // unit:s
)

func DialerWithTimeout() func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, time.Duration(HTTP_DIALER_TIMEOUT)*time.Second)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(10 * time.Second))
		return conn, nil
	}
}

func NewClientWithTimeout() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: DialerWithTimeout(),
		},
	}
}

func HttpGet(url string) (string, error) {
	/*
		client := NewClientWithTimeout()
		resp, err := client.Get(url)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		return "", err
	*/
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
