package model

import "time"

// AggMetrics 日志聚合对象
type AggMetrics struct {
	Took     int    `json:"took"`
	ScrollID string `json:"_scroll_id"`
	Hits     struct {
		Total    int           `json:"total"`
		MaxScore int           `json:"max_score"`
		Hits     []interface{} `json:"hits"`
	} `json:"hits"`
	Suggest      interface{} `json:"suggest"`
	Aggregations struct {
		DataAggs struct {
			Buckets []Bucket `json:"buckets"`
		} `json:"DataAggs"`
	} `json:"aggregations"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
}

//Bucket Bucket 对象
type Bucket struct {
	KeyAsString time.Time `json:"key_as_string"`
	Key         int64     `json:"key"`
	DocCount    int       `json:"doc_count"`
	MaxAgg      struct {
		Value float64 `json:"value"`
	} `json:"maxAgg"`
	MinAgg struct {
		Value float64 `json:"value"`
	} `json:"minAgg"`
	AvgAgg struct {
		Value float64 `json:"value"`
	} `json:"avgAgg"`
}

//AggResult 聚合的结果
type AggResult struct {
	Min             [24]float64 `json:"min" binding:"required"`
	Max             [24]float64 `json:"max" binding:"required"`
	Avg             [24]float64 `json:"avg" binding:"required"`
	Count           [24]int     `json:"count" binding:"required"`
	TotalCount      int         `json:"totalCount" binding:"required"`
	ShareTotalCount int         `json:"shareTotalCount" binding:"required"`
}

//LoadAggChart post请求过来的数据
type LoadAggChart struct {
	LogstashName string `json:"logstastname" binding:"required`
	Name         string `json:"name" binding:"required`
}

//PieMetrics Pie聚合返回的实体
type PieMetrics struct {
	Took     int    `json:"took"`
	ScrollID string `json:"_scroll_id"`
	Hits     struct {
		Total    int           `json:"total"`
		MaxScore int           `json:"max_score"`
		Hits     []interface{} `json:"hits"`
	} `json:"hits"`
	Suggest      interface{} `json:"suggest"`
	Aggregations struct {
		RangeAgg struct {
			Buckets []struct {
				Key      string  `json:"key"`
				To       float64 `json:"to"`
				DocCount int     `json:"doc_count"`
				From     float64 `json:"from,omitempty"`
			} `json:"buckets"`
		} `json:"rangeAgg"`
	} `json:"aggregations"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
}

//PieResult 返回Pie结果
type PieResult struct {
	Value int    `json:"value"`
	Name  string `json:"name"`
}

//URL 查询urlname
type URL struct {
	Took     int    `json:"took"`
	ScrollID string `json:"_scroll_id"`
	Hits     struct {
		Total    int           `json:"total"`
		MaxScore int           `json:"max_score"`
		Hits     []interface{} `json:"hits"`
	} `json:"hits"`
	Suggest      interface{} `json:"suggest"`
	Aggregations struct {
		TermAgg struct {
			DocCountErrorUpperBound int             `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int             `json:"sum_other_doc_count"`
			Buckets                 []ResultBuckets `json:"buckets"`
		} `json:"termAgg"`
	} `json:"aggregations"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
}

//ResultBuckets bucketResult
type ResultBuckets struct {
	Key string `json:"key"`
	// DocCount int    `json:"doc_count"`
}

//DateValue DateValue
type DateValue struct {
	LogstashName string `json:"logstastname" binding:"required"`
}
