package service

import (
	"net/http"
	"net/http/cookiejar"
	"sync"
)

type TikTok struct {
	HttpClient *http.Client
	Mu         sync.Mutex
	Wg         sync.WaitGroup
}

func NewTiktok() TikTok {
	jar, _ := cookiejar.New(nil)
	return TikTok{
		HttpClient: &http.Client{Jar: jar},
		Mu:         sync.Mutex{},
		Wg:         sync.WaitGroup{},
	}
}
