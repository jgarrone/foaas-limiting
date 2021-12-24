package utils

import neturl "net/url"

func BuildURL(protocol, host, path string) string {
	u := &neturl.URL{
		Scheme: protocol,
		Host:   host,
		Path:   path,
	}

	return u.String()
}
