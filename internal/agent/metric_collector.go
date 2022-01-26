package agent

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

func CollectRuntimeMetrics() {
	var rtm runtime.MemStats
	log.Println("ticker CollectRuntimeMetrics")
	runtime.ReadMemStats(&rtm)
	PollCount += 1
	rand.Seed(time.Now().Unix())
	Metrics = []Gauge{
		{name: "Alloc", value: float64(rtm.Alloc)},
		{name: "BuckHashSys", value: float64(rtm.BuckHashSys)},
		{name: "Frees", value: float64(rtm.Frees)},
		{name: "GCCPUFraction", value: rtm.GCCPUFraction},
		{name: "GCSys", value: float64(rtm.GCSys)},
		{name: "HeapAlloc", value: float64(rtm.HeapAlloc)},
		{name: "HeapIdle", value: float64(rtm.HeapIdle)},
		{name: "HeapInuse", value: float64(rtm.HeapInuse)},
		{name: "HeapObjects", value: float64(rtm.HeapObjects)},
		{name: "HeapReleased", value: float64(rtm.HeapReleased)},
		{name: "HeapSys", value: float64(rtm.HeapSys)},
		{name: "LastGC", value: float64(rtm.LastGC)},
		{name: "Lookups", value: float64(rtm.Lookups)},
		{name: "MCacheInuse", value: float64(rtm.MCacheInuse)},
		{name: "MCacheSys", value: float64(rtm.MCacheSys)},
		{name: "MSpanInuse", value: float64(rtm.MSpanInuse)},
		{name: "MSpanSys", value: float64(rtm.MSpanSys)},
		{name: "Mallocs", value: float64(rtm.Mallocs)},
		{name: "NextGC", value: float64(rtm.NextGC)},
		{name: "NumForcedGC", value: float64(rtm.NumForcedGC)},
		{name: "NumGC", value: float64(rtm.NumGC)},
		{name: "OtherSys", value: float64(rtm.OtherSys)},
		{name: "PauseTotalNs", value: float64(rtm.PauseTotalNs)},
		{name: "StackInuse", value: float64(rtm.StackInuse)},
		{name: "StackSys", value: float64(rtm.StackSys)},
		{name: "Sys", value: float64(rtm.Sys)},
		{name: "TotalAlloc", value: float64(rtm.TotalAlloc)},
	}
}

//StoreRandomMetrics non-runtime metrics
func StoreRandomMetrics(metrics []Gauge) []Gauge {
	rand.Seed(time.Now().Unix())
	randomMetric := Gauge{name: "RandomValue", value: rand.Float64() * 100}
	metrics = append(metrics, randomMetric)
	return metrics
}

func CollectUtilizationMetrics() {
	m, err := mem.VirtualMemory()
	if err != nil {
		log.Println(err)
	}

	UtilizationData.mu.Lock()
	timeNow := time.Now()
	timeDiff := timeNow.Sub(UtilizationData.CPUutilLastTime)

	UtilizationData.CPUutilLastTime = timeNow
	UtilizationData.TotalMemory = Gauge{
		name:  "TotalMemory",
		value: float64(m.Total),
	}
	UtilizationData.FreeMemory = Gauge{
		name:  "FreeMemory",
		value: float64(m.Free),
	}

	cpus, err := cpu.Times(true)
	if err != nil {
		log.Println(err)
	}
	for i := range cpus {
		newCPUTime := cpus[i].User + cpus[i].System
		cpuUtilization := (newCPUTime - UtilizationData.CPUtime[i]) * 1000 / float64(timeDiff.Milliseconds())
		UtilizationData.CPUutilizations[i] = Gauge{
			name:  "CPUutilization" + strconv.Itoa(i+1),
			value: cpuUtilization,
		}
		UtilizationData.CPUtime[i] = newCPUTime
	}
	UtilizationData.mu.Unlock()
}

func init() {
	cpuStat, err := cpu.Times(true)
	if err != nil {
		log.Println(err)
		return
	}
	numCPU := len(cpuStat)
	UtilizationData.CPUtime = make([]float64, numCPU)
	UtilizationData.CPUutilizations = make([]Gauge, numCPU)
}
