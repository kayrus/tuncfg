package resolv

import (
	"bufio"
	"bytes"
	"io/fs"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

const (
	ResolvPath = "/etc/resolv.conf"
)

var (
	AppName = "tuncfg"
)

type Handler struct {
	name        string
	dnsServers  []net.IP
	dnsSuffixes []string

	// parsed /etc/resolv.conf
	origDnsServers  []net.IP
	origDnsSuffixes []string
	origDnsOptions  []string

	// macos only
	names []string

	// the above relates only to linux/freebsd

	dnsOptions []string

	// rewrite indicates whether the resolv handler needs to rewrite the
	// /etc/resolv.conf on systems, where it is used, e.g. Linux or FreeBSD
	rewrite bool

	// backup is used to store original /etc/resolv.conf or a path to the
	// original backup resolv.conf to restore the initial config.
	backup         []byte
	backupFilename string

	mode fs.FileMode
}

func New(name string, dnsServers []net.IP, dnsSuffixes []string, rewrite bool) (*Handler, error) {
	return newHandler(name, dnsServers, dnsSuffixes, rewrite)
}

func splitFunc(c rune) bool {
	return c == ' ' || c == '\n' || c == '\r'
}

func parseIP(ip string) net.IP {
	v := net.ParseIP(ip)
	if v == nil {
		return nil
	}
	if v := v.To4(); v != nil {
		return v
	} else if v := v.To16(); v != nil {
		return v
	}
	return nil
}

func (h *Handler) parseResolvConf() error {
	var err error

	h.backup, err = ioutil.ReadFile(ResolvPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if len(h.backup) == 0 {
		return nil
	}

	buf := bufio.NewReader(bytes.NewReader(h.backup))
	for line, isPrefix, err := buf.ReadLine(); err == nil && !isPrefix; line, isPrefix, err = buf.ReadLine() {
		if len(line) > 0 && (line[0] == ';' || line[0] == '#') {
			continue
		}

		f := strings.FieldsFunc(string(line), splitFunc)
		if len(f) < 2 {
			continue
		}
		switch f[0] {
		case "nameserver":
			if len(h.origDnsServers) < 3 {
				v := parseIP(f[1])
				if v == nil {
					continue
				}
				h.origDnsServers = append(h.origDnsServers, v)
			}
		case "search":
			for _, v := range f[1:] {
				if v[0] == ';' || v[0] == '#' {
					// ignore commented data
					break
				}
				h.origDnsSuffixes = append(h.origDnsSuffixes, v)
			}
		case "options":
			for _, v := range f[1:] {
				if v[0] == ';' || v[0] == '#' {
					// ignore commented data
					break
				}
				h.origDnsOptions = append(h.origDnsOptions, v)
			}
		}
	}

	return nil
}

func (h *Handler) GetOriginalDNS() []net.IP {
	return h.origDnsServers
}

func (h *Handler) GetOriginalSuffixes() []string {
	return h.origDnsSuffixes
}

func (h *Handler) GetOriginalOptions() []string {
	return h.origDnsOptions
}

func (h *Handler) SetDNSServers(dnsServers []net.IP) {
	h.dnsServers = dnsServers
}

func (h *Handler) SetSuffixes(dnsSuffixes []string) {
	h.dnsSuffixes = dnsSuffixes
}

func (h *Handler) SetOptions(dnsOptions []string) {
	h.dnsOptions = dnsOptions
}
