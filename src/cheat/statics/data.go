package statics

import (
	"time"


	"bytes"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"

)

//系统信息，分为三部分:宿主信息(操作系统信息),应用信息()
type SystemInfo struct {
	Host         *HostInfo
	//DataBase     *DataBaseInfo
	AppInfo      *AppInfo
	DocumentInfo *DocumentInfo
}

func GetSysmtemInfo() *SystemInfo {
	return &SystemInfo{
		Host:         GetHostInfo(),
		//DataBase:     GetDataBaseInfo(),
		AppInfo:      GetAppInfo(),
		//DocumentInfo: GetDocumenInfo(),
	}
}

type HostInfo struct {
	CPU    int `json:"CPU" bson:"CPU"`
	Memory int `json:"Memory" bson:"Memory"`
	HD     int `json:"HD" bson:"HD"`
}
type DocumentInfo struct {
	Send, Read, Do, Receive int
}

func GetHostInfo() *HostInfo {

	if runtime.GOOS == "linux" {
		return &HostInfo{
			CPU:    calculateCpuUsage(),
			Memory: calMemoryUsage(),
			HD:     calDiskUsage(),
		}
	} else {
		return &HostInfo{}
	}

}

/*
	应用信息，包括应用版本信息，上次启动时间，运行时间，panic次数
*/
type AppInfo struct {
	Version      string            `json:"Version" `
	StartTime    string            `json:"StartTime" `
	UpTime       string            `json:"UpTime"`
	PanicNum     int32             `json:"PanicNum" `
	GoroutineNum int               `json:"GoroutineNum" `
	CPUNum       int               `json:"CPUNum" `
	Mem          *runtime.MemStats `json:"Mem" `
	LoginStatus  string            `json:"LoginStatus" `
	OnLineCount  int               `json:"OnLineCount" ` //在线人数
}

func GetAppInfo() *AppInfo {
	info := &AppInfo{}
	//info.Version = statics.GetVersion()
	//info.StartTime = statics.GetStartTime().Format(timeHelper.SecondFmt)
	//info.UpTime = timeHelper.DurationString(time.Now().Sub(statics.GetStartTime()))
	//info.PanicNum = one.GetPanicNum()
	info.GoroutineNum = runtime.NumGoroutine()
	info.CPUNum = runtime.NumCPU()
	info.Mem = &runtime.MemStats{}
	runtime.ReadMemStats(info.Mem)
	//info.OnLineCount = heartBeat.Count()
	info.LoginStatus = "正常"
	return info

}

/* .
cpu使用率计算，使用/proc/stat的输出做一个差量运算
*/
func calculateCpuUsage() int {
	var cpu1 int64
	var cpu2 int64

	var idle1 int64
	var idle2 int64

	var err error

	if cpu1, idle1, err = getCpuAndIdle(); err != nil {
		return 0
	}
	time.Sleep(time.Millisecond * 100)
	if cpu2, idle2, err = getCpuAndIdle(); err != nil {
		return 0
	}

	for cpu2 == cpu1 {
		time.Sleep(time.Millisecond * 100)
		if cpu2, idle2, err = getCpuAndIdle(); err != nil {
			return 0
		}
	}

	return 100 - int((idle2-idle1)*100/(cpu2-cpu1))
}

/*

 */
func getCpuAndIdle() (int64, int64, error) {
	if result, err := exec.Command("cat", "/proc/stat").CombinedOutput(); err != nil {
		return 0, 0, err
	} else {

		result = bytes.Split(result, []byte("\n"))[0]
		fields := bytes.Split(result, []byte(" "))

		user, _ := strconv.ParseInt(string(fields[1]), 10, 64)
		nice, _ := strconv.ParseInt(string(fields[2]), 10, 64)
		system, _ := strconv.ParseInt(string(fields[3]), 10, 64)
		idle, _ := strconv.ParseInt(string(fields[4]), 10, 64)
		iowait, _ := strconv.ParseInt(string(fields[5]), 10, 64)
		irq, _ := strconv.ParseInt(string(fields[6]), 10, 64)
		softirq, _ := strconv.ParseInt(string(fields[7]), 10, 64)
		stealstolen, _ := strconv.ParseInt(string(fields[8]), 10, 64)
		guest, _ := strconv.ParseInt(string(fields[9]), 10, 64)
		return user + nice + system + idle + iowait + irq + softirq + stealstolen + guest, idle, nil
	}

}
func calMemoryUsage() int {
	if result, err := exec.Command("free", "-m").CombinedOutput(); err != nil {
		return 0
	} else {
		result = bytes.Split(result, []byte("\n"))[2]
		re := regexp.MustCompile("\\D*(\\d+)\\D+(\\d+)")
		fields := re.FindSubmatch(result)
		free, _ := strconv.Atoi(string(fields[2]))
		used, _ := strconv.Atoi(string(fields[1]))
		return used * 100 / (free + used)
	}

}

//返回已用百分比
func calDiskUsage() int {
	if result, err := exec.Command("df", "-h").CombinedOutput(); err != nil {
		return 0
	} else {
		result = bytes.Split(result, []byte("\n"))[1]
		re := regexp.MustCompile("\\D+(\\d+\\D+)(\\d+)%")
		fields := re.FindSubmatch(result)
		usage, _ := strconv.Atoi(string(fields[2]))
		return 100 - usage
	}
}
