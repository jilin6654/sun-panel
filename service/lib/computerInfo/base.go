package computerInfo

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/v3/cpu"
	"gitlab.com/tingshuo/go-diskstate/diskstate"
)

type Storage struct {
	Name       string
	FileSystem string
	Total      uint64
	Free       uint64
}

type storageInfo struct {
	Name       string
	Size       uint64
	FreeSpace  uint64
	FileSystem string
	Used       uint64
}

// func GetStorageInfo() {
// 	var storageinfo []storageInfo
// 	var loaclStorages []Storage
// 	err := wmi.Query("Select * from Win32_LogicalDisk", &storageinfo)
// 	if err != nil {
// 		return
// 	}

// 	for _, storage := range storageinfo {
// 		info := Storage{
// 			Name:       storage.Name,
// 			FileSystem: storage.FileSystem,
// 			Total:      storage.Size / 1024 / 1024 / 1024,
// 			Free:       storage.FreeSpace / 1024 / 1024 / 1024,
// 		}
// 		if info.Total >= 1 {
// 			fmt.Printf("%s总大小%dG，可用%dG\n", info.Name, info.Total, info.Free)
// 			loaclStorages = append(loaclStorages, info)
// 		}
// 	}
// 	fmt.Printf("localStorages:= %v\n", loaclStorages)
// }

func GetCurrentStorageInfo(path string) storageInfo {
	state := diskstate.DiskUsage(path)
	info := storageInfo{}
	info.Size = uint64(state.All / diskstate.B)
	info.FreeSpace = uint64(state.Free / diskstate.B)
	info.Used = uint64(state.Used / diskstate.B)

	// fmt.Printf("All=%dM, Free=%dM, Available=%dM, Used=%dM, Usage=%d%%",
	// 	state.All/diskstate.B, state.Free/diskstate.MB, state.Available/diskstate.MB, state.Used/diskstate.MB, 100*state.Used/state.All)
	return info
}

type ComputerMonitor struct {
	CPU float64 `json:"cpu"`
	Mem float64 `json:"mem"`
}

// GetCPUPercent 获取CPU使用率
func GetCPUPercent() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Fatalln(err.Error())
		return -1
	}
	return percent[0]
}

// GetMemPercent 获取内存使用率
func GetMemPercent() float64 {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalln(err.Error())
		return -1
	}
	return memInfo.UsedPercent
}

func GetCpuMem() ComputerMonitor {
	var res ComputerMonitor
	res.CPU = GetCPUPercent()
	res.Mem = GetMemPercent()
	return res
}

func RemoveProtocol(ip string) string {
	if strings.HasPrefix(ip, "http://") {
		return strings.TrimPrefix(ip, "http://")
	} else if strings.HasPrefix(ip, "https://") {
		return strings.TrimPrefix(ip, "https://")
	}

	return ip
}

func CheckHostStatus(ip string) int {
	cmd := exec.Command("ping", "-c", "1", "-W", "1", RemoveProtocol(ip))
	err := cmd.Run()
	status := 0
	if err == nil {
		status = 1
	}
	return status
}

func SleepHandler(ip string) bool {
	if ip == "" {
		return false
	}
	sshCommand := fmt.Sprintf("ssh %s@%s systemctl suspend", "root", ip)
	cmd := exec.Command("bash", "-c", sshCommand)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to execute SSH command: %v", err)
		return false
	}

	return true
}

func WakeOnLanHandler(mac string) bool {
	if mac == "" {
		return false
	}
	mac = strings.ToUpper(mac)
	mac = strings.Replace(mac, "-", ":", -1)

	if err := SendMagicPacket(mac); err != nil {
		log.Printf("Failed to send magic packet: %v", err)
		return false
	}

	return true
}

func SendMagicPacket(macAddr string) error {
	hwAddr, err := net.ParseMAC(macAddr)
	if err != nil {
		return err
	}

	var packet []byte
	packet = append(packet, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF)
	for i := 0; i < 16; i++ {
		packet = append(packet, hwAddr...)
	}

	// Broadcast the magic packet to the local network
	conn, err := net.Dial("udp", "255.255.255.255:9")
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.Write(packet); err != nil {
		return err
	}

	return nil
}
