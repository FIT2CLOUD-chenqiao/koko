package proxy

import (
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/jumpserver/koko/pkg/i18n"
)

const (
	UnAuth            = "unable to authenticate"
	ConnectRefusedErr = "connection refused" // 无监听端口或者端口被防火墙阻断
	IoTimeoutErr      = "i/o timeout"
	NoRouteErr        = "No route to host" //

	LoginFailed = "failed login" // telnet 连接失败

	networkUnreachable = "network is unreachable"
)

func ConvertErrorToReadableMsg(e error) string {
	if e == nil {
		return ""
	}
	errMsg := e.Error()
	if strings.Contains(errMsg, UnAuth) || strings.Contains(errMsg, LoginFailed) {
		return i18n.T("Authentication failed")
	}
	if strings.Contains(errMsg, ConnectRefusedErr) {
		return i18n.T("Connection refused")
	}
	if strings.Contains(errMsg, IoTimeoutErr) {
		return i18n.T("i/o timeout")
	}
	if strings.Contains(errMsg, NoRouteErr) {
		return i18n.T("No route to host")
	}
	if strings.Contains(errMsg, networkUnreachable) {
		return i18n.T("network is unreachable")
	}
	return errMsg
}

func ReplaceURLHostAndPort(originUrl *url.URL, ip string, port int) string {
	newHost := net.JoinHostPort(ip, strconv.Itoa(port))
	switch originUrl.Scheme {
	case "https":
		if port == 443 {
			newHost = ip
		}
	default:
		if port == 80 {
			newHost = ip
		}
	}
	newUrl := url.URL{
		Scheme:     originUrl.Scheme,
		Opaque:     originUrl.Opaque,
		User:       originUrl.User,
		Host:       newHost,
		Path:       originUrl.Path,
		RawPath:    originUrl.RawQuery,
		ForceQuery: originUrl.ForceQuery,
		RawQuery:   originUrl.RawQuery,
		Fragment:   originUrl.Fragment,
	}
	return newUrl.String()
}
