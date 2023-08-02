package function

import (
	"github.com/beego/beego/v2/core/config"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

func GetMacStr() string {
	ipMacs := GetMacIps()
	keys := make([]string, 0)
	for _, val := range ipMacs {
		if !InArray(val, keys) {
			keys = append(keys, val)
		}
	}
	sort.Strings(keys)
	return Implode(",", keys)
}

func GetIpStr() string {
	ipMacs := GetMacIps()
	keys := make([]string, 0)
	for val, _ := range ipMacs {
		if !InArray(val, keys) {
			keys = append(keys, val)
		}
	}
	sort.Strings(keys)
	return Implode(",", keys)
}

func InArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack type muset be slice, array or map")
	}

	return false
}

func GetMacIps() (macAddrs map[string]string) {
	macAddrs = map[string]string{}
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return macAddrs
	}

	for _, netInterface := range netInterfaces {
		addrs, err := netInterface.Addrs()
		if err != nil {
			continue
		}
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}

		// 获取IP地址，子网掩码
		for _, addr := range addrs {
			if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
				if ip.IP.To4() != nil {
					ipstr := ip.IP.String()
					//过滤本地IP
					if ipstr == "127.0.0.1" {
						continue
					}
					if len(ipstr) > 7 {
						ipleft := ipstr[0:7]
						if ipleft == "169.254" {
							continue
						}
						macAddrs[ip.IP.String()] = macAddr
					}
				}
			}
		}
	}
	return macAddrs
}

func GetCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func GetConf(fileName, key, valueType string) (interface{}, error) {
	fileName = "atsjkhelper.conf"
	cfg, err := config.NewConfig("ini", fileName)
	if err != nil {
		return nil, err
	}
	if valueType == "strings" {
		val, err := cfg.Strings(key)
		if err != nil {
			return nil, err
		}

		var v []string
		for _, eve := range val {
			v = strings.Split(eve, ",")
		}

		return v, nil
	} else {
		val, err := cfg.String(key)
		if err != nil {
			return nil, err
		}

		return val, nil
	}
}
