package netFramework

import "golang.org/x/sys/windows/registry"

const (
	V45MinValue  = 378389
	V451MinValue = 378675
	V452MinValue = 379893
	V46MinValue  = 393295
	V461MinValue = 394254
	V462MinValue = 394802
	V47MinValue  = 460798
	V471MinValue = 461308
	V472MinValue = 461808
	V48MinValue  = 528040
	V481MinValue = 533320
)

type V2Info struct {
	Install uint64
	Version string
}

type V3Info struct {
	Install uint64
	Version string
}

type V35Info struct {
	Install     uint64
	InstallPath string
	Version     string
}

type V4Info struct {
	Install     uint64
	InstallPath string
	Release     uint64
	Version     string
}

func GetV2Info() (V2Info, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\NET Framework Setup\\NDP\\v2.0.50727", registry.READ)
	info := V2Info{}
	if err != nil {
		return info, err
	}
	install, _, err := key.GetIntegerValue("Install")
	if err != nil {
		return info, err
	}
	version, _, err := key.GetStringValue("Version")
	if err != nil {
		return info, err
	}
	info.Install = install
	info.Version = version
	return info, nil
}

func GetV3Info() (V3Info, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\NET Framework Setup\\NDP\\v3.0\\Setup", registry.READ)
	info := V3Info{}
	if err != nil {
		return info, err
	}
	install, _, err := key.GetIntegerValue("InstallSuccess")
	if err != nil {
		return info, err
	}
	version, _, err := key.GetStringValue("Version")
	if err != nil {
		return info, err
	}
	info.Install = install
	info.Version = version
	return info, nil
}

func GetV35Info() (V35Info, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\NET Framework Setup\\NDP\\v3.5", registry.READ)
	info := V35Info{}
	if err != nil {
		return info, err
	}
	install, _, err := key.GetIntegerValue("Install")
	if err != nil {
		return info, err
	}
	installPath, _, err := key.GetStringValue("InstallPath")
	if err != nil {
		return info, err
	}
	version, _, err := key.GetStringValue("Version")
	if err != nil {
		return info, err
	}
	info.Install = install
	info.InstallPath = installPath
	info.Version = version
	return info, nil
}

func GetV4Info() (V4Info, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\NET Framework Setup\\NDP\\v4\\Full", registry.READ)
	info := V4Info{}
	if err != nil {
		return info, err
	}
	install, _, err := key.GetIntegerValue("Install")
	if err != nil {
		return info, err
	}
	installPath, _, err := key.GetStringValue("InstallPath")
	if err != nil {
		return info, err
	}
	release, _, err := key.GetIntegerValue("Release")
	if err != nil {
		return info, err
	}
	version, _, err := key.GetStringValue("Version")
	if err != nil {
		return info, err
	}
	info.Install = install
	info.InstallPath = installPath
	info.Release = release
	info.Version = version
	return info, nil
}

func CheckFor45PlusVersion(releaseKey int) string {
	if releaseKey >= 533320 {
		return "4.8.1 or later"
	}
	if releaseKey >= 528040 {
		return "4.8"
	}
	if releaseKey >= 461808 {
		return "4.7.2"
	}
	if releaseKey >= 461308 {
		return "4.7.1"
	}
	if releaseKey >= 460798 {
		return "4.7"
	}
	if releaseKey >= 394802 {
		return "4.6.2"
	}
	if releaseKey >= 394254 {
		return "4.6.1"
	}
	if releaseKey >= 393295 {
		return "4.6"
	}
	if releaseKey >= 379893 {
		return "4.5.2"
	}
	if releaseKey >= 378675 {
		return "4.5.1"
	}
	if releaseKey >= 378389 {
		return "4.5"
	}
	// This code should never execute. A non-null release key should mean
	// 安装4.5或更高版本。
	return "No 4.5 or later version detected"
}
