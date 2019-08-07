package darproxy

// Cfg Struct
type Cfg struct {
	Port      string     `json:"interface"`
	Upstreams []Upstream `json:"upstreams"`
}

// Upstream Struct
type Upstream struct {
	Path        string   `json:"path"`
	Method      string   `json:"method"`
	Back        []string `json:"backends"`
	ProxyMethod string   `json:"proxyMethod"`
}
