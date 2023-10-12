package geoip

import "regexp"

var ipV4Regex = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
var ipV6Regex = regexp.MustCompile(`^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`)

func IsIPv4(ip string) bool {
	return ipV4Regex.MatchString(ip)
}

func IsIPv6(ip string) bool {
	return ipV6Regex.MatchString(ip)
}
