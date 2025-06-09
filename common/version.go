package common

import (
	"fmt"
	"strconv"
	"strings"
)

type VersionInfo struct {
	VersionMajor int64 `json:"versionMajor"`
	VersionMinor int64 `json:"versionMinor"`
	VersionPatch int64 `json:"versionPatch"`
	VersionBuild int64 `json:"versionBuild"`
}

func StringToVersionInfo(str string) VersionInfo {
	ver := VersionInfo{0, 0, 0, 0}
	arr := strings.Split(str, ".")
	if len(arr) > 0 {
		ver.VersionMajor, _ = strconv.ParseInt(arr[0], 10, 64)
	}
	if len(arr) > 1 {
		ver.VersionMinor, _ = strconv.ParseInt(arr[1], 10, 64)
	}
	if len(arr) > 2 {
		ver.VersionPatch, _ = strconv.ParseInt(arr[2], 10, 64)
	}
	if len(arr) > 3 {
		ver.VersionBuild, _ = strconv.ParseInt(arr[3], 10, 64)
	}

	return ver
}

// VersionComparison 版本比较 1.大于 -1.小于 0.等于
func VersionComparison(ver1 VersionInfo, ver2 VersionInfo) int {

	if ver1.VersionMajor > ver2.VersionMajor {
		return 1
	} else if ver1.VersionMajor < ver2.VersionMajor {
		return -1
	}

	if ver1.VersionMinor > ver2.VersionMinor {
		return 1
	} else if ver1.VersionMinor < ver2.VersionMinor {
		return -1
	}

	if ver1.VersionPatch > ver2.VersionPatch {
		return 1
	} else if ver1.VersionPatch < ver2.VersionPatch {
		return -1
	}

	if ver1.VersionBuild > ver2.VersionBuild {
		return 1
	} else if ver1.VersionBuild < ver2.VersionBuild {
		return -1
	}

	return 0
}

func VersionToString(ver VersionInfo) string {
	str := fmt.Sprintf("%d.%d.%d.%d", ver.VersionMajor, ver.VersionMinor, ver.VersionPatch, ver.VersionBuild)
	return str
}

// CompareVersion 比较两个版本号 version1 和 version2，如果 version1 > version2 返回 1，如果 version1 < version2 返回 -1， 除此之外返回 0。
func CompareVersion(version1 string, version2 string) int {
	var arr1 = strings.Split(version1, ".")
	var arr2 = strings.Split(version2, ".")
	var length1 = len(arr1)
	var length2 = len(arr2)
	var length = Min(length1, length2)
	v1 := 0
	v2 := 0
	for i := 0; i < length; i++ {

		if i < length1 {
			v1, _ = strconv.Atoi(arr1[i])
		} else {
			v1 = 0
		}

		if i < length2 {
			v2, _ = strconv.Atoi(arr2[i])
		} else {
			v2 = 0
		}

		if v1 != v2 {
			if v1 > v2 {
				return 1
			} else if v1 < v2 {
				return -1
			}
		}

	}

	if v1 > v2 {
		return 1
	} else if v1 < v2 {
		return -1
	}
	return 0

}

func Max(vals ...int) int {
	var max int
	for _, val := range vals {
		if val > max {
			max = val
		}
	}
	return max
}

func Min(vals ...int) int {
	var min int
	for _, val := range vals {
		if min == 0 || val <= min {
			min = val
		}
	}
	return min
}
