package main

import (
	"net"

	"github.com/grandcat/zeroconf"
)

// ServiceItem represents the data returned by zeroconf from a service browse
type ServiceItem struct {
	Name     string   `json:"name"`     // Service Name
	Port     int      `json:"port"`     // Service Port
	HostName string   `json:"hostname"` // Host machine DNS name
	Service  string   `json:"type"`     // Service name
	Domain   string   `json:"domain"`   // If blank, assumes "local"
	Text     []string `json:"text"`     // Service info served as a TXT record
	AddrIPv4 []net.IP `json:"ipv4"`     // Host machine IPv4 address
	AddrIPv6 []net.IP `json:"ipv6"`     // Host machine IPv6 address
}

// NewServiceItemFromZeroConf returns a ServiceItem object loaded with the values from the zeroconf service entry record
func NewServiceItemFromZeroConf(e *zeroconf.ServiceEntry) ServiceItem {
	return ServiceItem{
		Name:     e.Instance,
		HostName: e.HostName,
		Port:     e.Port,
		Service:  e.Service,
		Domain:   e.Domain,
		Text:     e.Text,
		AddrIPv4: e.AddrIPv4,
		AddrIPv6: e.AddrIPv6,
	}
}
