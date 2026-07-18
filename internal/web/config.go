package web

import "time"

type Config struct {
	AppName    string
	AppVersion string

	CDNDomain string

	ClientCacheDur     time.Duration
	ClientCacheSeconds string
}
