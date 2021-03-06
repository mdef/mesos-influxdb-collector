package mesos

import (
	"encoding/json"
	"github.com/kpacha/mesos-influxdb-collector/store"
	"io"
	"io/ioutil"
	"log"
	"time"
)

type SlaveParser struct {
	Node string
}

func (mp SlaveParser) Parse(r io.ReadCloser) ([]store.Point, error) {
	defer r.Close()
	var stats SlaveStats
	body, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println("Error reading from", r)
		return []store.Point{}, err
	}
	if err = json.Unmarshal(body, &stats); err != nil {
		log.Println("Error parsing to SlaveStats")
		return []store.Point{}, err
	}
	stats.Node = mp.Node
	stats.Time = time.Now()
	return mp.getMesosPoints(stats), nil
}

func (mp SlaveParser) getMesosPoints(stats SlaveStats) []store.Point {
	return []store.Point{
		mp.getCpuPoint(stats),
		mp.getDiskPoint(stats),
		mp.getMemPoint(stats),
		mp.getSystemPoint(stats),
		mp.getTasksPoint(stats),
		mp.getExecutorPoint(stats),
		mp.getGlobalPoint(stats),
	}
}

func (mp SlaveParser) getCpuPoint(stats SlaveStats) store.Point {
	return store.Point{
		Measurement: "slave-cpu",
		Tags: map[string]string{
			"node": stats.Node,
		},
		Fields: map[string]interface{}{
			"percent": stats.Slave_cpusPercent,
			"total":   stats.Slave_cpusTotal,
			"used":    stats.Slave_cpusUsed,
		},
		Time: stats.Time,
	}
}

func (mp SlaveParser) getDiskPoint(stats SlaveStats) store.Point {
	return store.Point{
		Measurement: "slave-disk",
		Tags: map[string]string{
			"node": stats.Node,
		},
		Fields: map[string]interface{}{
			"percent": stats.Slave_diskPercent,
			"total":   stats.Slave_diskTotal,
			"used":    stats.Slave_diskUsed,
		},
		Time: stats.Time,
	}
}

func (mp SlaveParser) getMemPoint(stats SlaveStats) store.Point {
	return store.Point{
		Measurement: "slave-mem",
		Tags: map[string]string{
			"node": stats.Node,
		},
		Fields: map[string]interface{}{
			"percent": stats.Slave_memPercent,
			"total":   stats.Slave_memTotal,
			"used":    stats.Slave_memUsed,
		},
		Time: stats.Time,
	}
}

func (mp SlaveParser) getTasksPoint(stats SlaveStats) store.Point {
	return store.Point{
		Measurement: "slave-tasks",
		Tags: map[string]string{
			"node": stats.Node,
		},
		Fields: map[string]interface{}{
			"failed":   stats.Slave_tasksFailed,
			"finished": stats.Slave_tasksFinished,
			"killed":   stats.Slave_tasksKilled,
			"lost":     stats.Slave_tasksLost,
			"running":  stats.Slave_tasksRunning,
			"staging":  stats.Slave_tasksStaging,
			"starting": stats.Slave_tasksStarting,
		},
		Time: stats.Time,
	}
}

func (mp SlaveParser) getSystemPoint(stats SlaveStats) store.Point {
	return store.Point{
		Measurement: "system",
		Tags: map[string]string{
			"node": stats.Node,
		},
		Fields: map[string]interface{}{
			"cpus_total":      stats.System_cpusTotal,
			"load_15min":      stats.System_load15min,
			"load_1min":       stats.System_load1min,
			"load_5min":       stats.System_load5min,
			"mem_free_bytes":  int(stats.System_memFreeBytes),
			"mem_total_bytes": int(stats.System_memTotalBytes),
		},
		Time: stats.Time,
	}
}

func (mp SlaveParser) getExecutorPoint(stats SlaveStats) store.Point {
	return store.Point{
		Measurement: "slave-executor",
		Tags: map[string]string{
			"node": stats.Node,
		},
		Fields: map[string]interface{}{
			"directory_max_allowed_age_secs": stats.Slave_executorDirectoryMaxAllowedAgeSecs,
			"registering":                    stats.Slave_executorsRegistering,
			"running":                        stats.Slave_executorsRunning,
			"terminated":                     stats.Slave_executorsTerminated,
			"terminating":                    stats.Slave_executorsTerminating,
		},
		Time: stats.Time,
	}
}

func (mp SlaveParser) getGlobalPoint(stats SlaveStats) store.Point {
	return store.Point{
		Measurement: "slave-global",
		Tags: map[string]string{
			"node": stats.Node,
		},
		Fields: map[string]interface{}{
			"registered":                 stats.Slave_registered,
			"invalid_framework_messages": stats.Slave_invalidFrameworkMessages,
			"invalid_status_updates":     stats.Slave_invalidStatusUpdates,
			"uptime_secs":                stats.Slave_uptimeSecs,
			"valid_framework_messages":   stats.Slave_validFrameworkMessages,
			"valid_status_updates":       stats.Slave_validStatusUpdates,
			"framewors":                  stats.Slave_frameworksActive,
			"conatiner_launch_errors":    stats.Slave_containerLaunchErrors,
		},
		Time: stats.Time,
	}
}

type SlaveStats struct {
	Containerizer_mesos_containerDestroyErrors float64 `json:"containerizer/mesos/container_destroy_errors"`
	Slave_containerLaunchErrors                float64 `json:"slave/container_launch_errors"`
	Slave_cpusPercent                          float64 `json:"slave/cpus_percent"`
	Slave_cpusRevocablePercent                 float64 `json:"slave/cpus_revocable_percent"`
	Slave_cpusRevocableTotal                   float64 `json:"slave/cpus_revocable_total"`
	Slave_cpusRevocableUsed                    float64 `json:"slave/cpus_revocable_used"`
	Slave_cpusTotal                            float64 `json:"slave/cpus_total"`
	Slave_cpusUsed                             float64 `json:"slave/cpus_used"`
	Slave_diskPercent                          float64 `json:"slave/disk_percent"`
	Slave_diskRevocablePercent                 float64 `json:"slave/disk_revocable_percent"`
	Slave_diskRevocableTotal                   float64 `json:"slave/disk_revocable_total"`
	Slave_diskRevocableUsed                    float64 `json:"slave/disk_revocable_used"`
	Slave_diskTotal                            float64 `json:"slave/disk_total"`
	Slave_diskUsed                             float64 `json:"slave/disk_used"`
	Slave_executorDirectoryMaxAllowedAgeSecs   float64 `json:"slave/executor_directory_max_allowed_age_secs"`
	Slave_executorsRegistering                 float64 `json:"slave/executors_registering"`
	Slave_executorsRunning                     float64 `json:"slave/executors_running"`
	Slave_executorsTerminated                  float64 `json:"slave/executors_terminated"`
	Slave_executorsTerminating                 float64 `json:"slave/executors_terminating"`
	Slave_frameworksActive                     float64 `json:"slave/frameworks_active"`
	Slave_invalidFrameworkMessages             float64 `json:"slave/invalid_framework_messages"`
	Slave_invalidStatusUpdates                 float64 `json:"slave/invalid_status_updates"`
	Slave_memPercent                           float64 `json:"slave/mem_percent"`
	Slave_memRevocablePercent                  float64 `json:"slave/mem_revocable_percent"`
	Slave_memRevocableTotal                    float64 `json:"slave/mem_revocable_total"`
	Slave_memRevocableUsed                     float64 `json:"slave/mem_revocable_used"`
	Slave_memTotal                             float64 `json:"slave/mem_total"`
	Slave_memUsed                              float64 `json:"slave/mem_used"`
	Slave_recoveryErrors                       float64 `json:"slave/recovery_errors"`
	Slave_registered                           float64 `json:"slave/registered"`
	Slave_tasksFailed                          float64 `json:"slave/tasks_failed"`
	Slave_tasksFinished                        float64 `json:"slave/tasks_finished"`
	Slave_tasksKilled                          float64 `json:"slave/tasks_killed"`
	Slave_tasksLost                            float64 `json:"slave/tasks_lost"`
	Slave_tasksRunning                         float64 `json:"slave/tasks_running"`
	Slave_tasksStaging                         float64 `json:"slave/tasks_staging"`
	Slave_tasksStarting                        float64 `json:"slave/tasks_starting"`
	Slave_uptimeSecs                           float64 `json:"slave/uptime_secs"`
	Slave_validFrameworkMessages               float64 `json:"slave/valid_framework_messages"`
	Slave_validStatusUpdates                   float64 `json:"slave/valid_status_updates"`
	System_cpusTotal                           float64 `json:"system/cpus_total"`
	System_load15min                           float64 `json:"system/load_15min"`
	System_load1min                            float64 `json:"system/load_1min"`
	System_load5min                            float64 `json:"system/load_5min"`
	System_memFreeBytes                        float64 `json:"system/mem_free_bytes"`
	System_memTotalBytes                       float64 `json:"system/mem_total_bytes"`
	Time                                       time.Time
	Node                                       string
}
