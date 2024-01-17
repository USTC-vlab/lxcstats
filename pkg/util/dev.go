package util

func GetDeviceNumbers(devId uint64) (uint64, uint64) {
	major := (devId >> 8) & 0xfff
	minor := (devId & 0xff) | ((devId >> 12) & 0xfff00)
	return major, minor
}
