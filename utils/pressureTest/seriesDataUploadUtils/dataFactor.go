/*
  Package seriesDataUploadUtils
  @Author: Ahsen17
  @Github: https://github.com/Ahsen17
  @Time: 2024/10/16 0:27
  @Description: ...
*/

package seriesDataUploadUtils

import (
	"sync"
	"time"
)

type SeriesData struct {
	Metric    string            `json:"metric"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
	Tags      map[string]string `json:"tags"`
}

type Metrics []SeriesData

type MetricsFactor struct {
	NodeLimit  int
	Node       []sync.RWMutex
	MetricsMap []Metrics
}

var (
	mLock sync.RWMutex
	tLock sync.RWMutex
	iLock sync.RWMutex

	metricCache = make(map[int]string, metricLimit)
	tagsKVCache = make(map[int]string, tagKVLimit)
	identsCache = make(map[int]int, targetLimit)

	metricsFactory = &MetricsFactor{
		NodeLimit:  10,
		Node:       make([]sync.RWMutex, 10),
		MetricsMap: make([]Metrics, 10),
	}
)

// produceTagsKV 生产tagsKV
func produceTagsKV() {
	index := randInt(0, tagKVLimit)
	kOrV := randStrDoubleLimits(tagKVLenMin, tagKVLenMax)

	iLock.RLock()
	_, exist := tagsKVCache[index]
	iLock.RUnlock()

	if !exist {
		tLock.Lock()
		tagsKVCache[index] = kOrV
		tLock.Unlock()
	}
}

// produceIdents 生产idents
func produceIdents() {
	index := randInt(0, targetLimit)

	iLock.RLock()
	_, exist := identsCache[index]
	iLock.RUnlock()

	if !exist {
		iLock.Lock()
		identsCache[index] = index
		iLock.Unlock()
	}
}

// produceMetric 生产metric
func produceMetric() {
	index := randInt(0, tagKVLimit)
	metric := randStrDoubleLimits(metricLenMin, metricLenMax)

	iLock.RLock()
	_, exist := metricCache[index]
	iLock.RUnlock()

	if !exist {
		tLock.Lock()
		metricCache[index] = metric
		tLock.Unlock()
	}
}

// metricAppendTags 为metric添加随机tags
func metricAppendTags(metric *SeriesData) *SeriesData {
	return nil
}

func (ms *Metrics) Append(data *SeriesData) {
	*ms = append(*ms, *data)
}

func (mf *MetricsFactor) produceSeriesData() {
	batchNum := randInt(0, metricProducePerBatch)
	count := 0

	done := make(chan bool)
	for {
		select {
		case <-done:
			return
		default:
			idx := randInt(0, metricLimit)
			mLock.RLock()
			metric, exist := metricCache[idx]
			mLock.RUnlock()
			if !exist {
				continue
			} else {
				nodeIdx := idx % mf.NodeLimit

				// 放入消费池
				mf.Node[nodeIdx].Lock()
				mf.MetricsMap[nodeIdx].Append(
					metricAppendTags(&SeriesData{
						Metric:    metric,
						Value:     randFloat(valMin, valMax),
						Timestamp: time.Now().Unix(),
						Tags:      nil,
					}))
				mf.Node[nodeIdx].Unlock()

				// 当批次生产完成
				count++
				if count == batchNum {
					done <- true
					break
				}
			}
		}
	}
}
