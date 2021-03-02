module github.com/kayrus/tuncfg

go 1.16

// a fork with a FreeBSD default tun name patch
replace golang.zx2c4.com/wireguard v0.0.0-20210225140808-70b7b7158fc9 => github.com/kayrus/wireguard v0.0.0-20210228102730-04afc3c4c795

// a fork with a Windows convertInterfaceIndexToLUID
replace golang.zx2c4.com/wireguard/windows v0.3.8 => github.com/kayrus/wireguard-windows v0.0.0-20210303100507-540e87897140

require (
	github.com/IBM/netaddr v1.4.0
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/vishvananda/netlink v1.1.0
	golang.org/x/sys v0.0.0-20210303074136-134d130e1a04
	golang.zx2c4.com/wireguard v0.0.0-20210225140808-70b7b7158fc9
	golang.zx2c4.com/wireguard/windows v0.3.8
)
