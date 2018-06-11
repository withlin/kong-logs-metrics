package model

import "time"

//Logs 日志实体
type Logs struct {
	Took     int    `json:"took"`
	ScrollID string `json:"_scroll_id"`
	Hits     struct {
		Total    int `json:"total"`
		MaxScore int `json:"max_score"`
		Hits     []struct {
			Score     int         `json:"_score"`
			Index     string      `json:"_index"`
			Type      string      `json:"_type"`
			ID        string      `json:"_id"`
			UID       string      `json:"_uid"`
			Routing   string      `json:"_routing"`
			Parent    string      `json:"_parent"`
			Version   interface{} `json:"_version"`
			Sort      interface{} `json:"sort"`
			Highlight interface{} `json:"highlight"`
			Source    struct {
				Request struct {
					Headers struct {
						Host      string `json:"host"`
						Accept    string `json:"accept"`
						UserAgent string `json:"user-agent"`
					} `json:"headers"`
					Size        string `json:"size"`
					Method      string `json:"method"`
					Querystring struct {
					} `json:"querystring"`
					URI string `json:"uri"`
					URL string `json:"url"`
				} `json:"request"`
				// Tries []struct {
				// 	Port            int    `json:"port"`
				// 	BalancerLatency int    `json:"balancer_latency"`
				// 	IP              string `json:"ip"`
				// } `json:"tries"`
				Latencies struct {
					Request int `json:"request"`
					Proxy   int `json:"proxy"`
					Kong    int `json:"kong"`
				} `json:"latencies"`
				UpstreamURI string `json:"upstream_uri"`
				// Message     string    `json:"message"`
				Timestamp time.Time `json:"@timestamp"`
				Port      int       `json:"port"`
				Response  struct {
					Headers struct {
						Date                 string `json:"date"`
						Server               string `json:"server"`
						XKongUpstreamLatency string `json:"x-kong-upstream-latency"`
						XKongProxyLatency    string `json:"x-kong-proxy-latency"`
						ContentType          string `json:"content-type"`
						Connection           string `json:"connection"`
						Via                  string `json:"via"`
					} `json:"headers"`
					Size   string `json:"size"`
					Status int    `json:"status"`
				} `json:"response"`
				Version   string    `json:"@version"`
				Host      string    `json:"host"`
				StartedAt time.Time `json:"started_at"`
				ClientIP  string    `json:"client_ip"`
				API       struct {
					UpstreamSendTimeout    int      `json:"upstream_send_timeout"`
					UpstreamURL            string   `json:"upstream_url"`
					HTTPSOnly              bool     `json:"https_only"`
					Methods                []string `json:"methods"`
					CreatedAt              int64    `json:"created_at"`
					PreserveHost           bool     `json:"preserve_host"`
					HTTPIfTerminated       bool     `json:"http_if_terminated"`
					Retries                int      `json:"retries"`
					Uris                   []string `json:"uris"`
					UpstreamConnectTimeout int      `json:"upstream_connect_timeout"`
					StripURI               bool     `json:"strip_uri"`
					Name                   string   `json:"name"`
					ID                     string   `json:"id"`
					UpstreamReadTimeout    int      `json:"upstream_read_timeout"`
				} `json:"api"`
				Consumer struct {
					CustomID  string `json:"custom_id"`
					CreatedAt int64  `json:"created_at"`
					ID        string `json:"id"`
					Username  string `json:"username"`
				} `json:"consumer"`
			} `json:"_source"`
			Fields         interface{} `json:"fields"`
			Explanation    interface{} `json:"_explanation"`
			MatchedQueries interface{} `json:"matched_queries"`
			InnerHits      interface{} `json:"inner_hits"`
		} `json:"hits"`
	} `json:"hits"`
	Suggest      interface{} `json:"suggest"`
	Aggregations interface{} `json:"aggregations"`
	TimedOut     bool        `json:"timed_out"`
	Shards       struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
}

//Page 分页
type Page struct {
	PageSize   int    `json:"pagesize" binding:"required,numeric"`
	PageNumber int    `json:"pagenumber" binding:"required,numeric"`
	DateValue  string `json:"datevalue" binding:"required"`
}

//ID ID
type ID struct {
	ID        string `json:"id" binding:"required"`
	IndexName string `json:"indexname" binding:"required"`
}

//API API
type API struct {
	Appid      string `json:"appid" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Data       string `json:"datevalue" binding:"required"`
	PageSize   int    `json:"pagesize" binding:"required,numeric"`
	PageNumber int    `json:"pagenumber" binding:"required,numeric"`
}
