package types

import "strings"

type (
	NICType string
)

// TODO: Implement other types user,socket,brtidge

const (
	TAP    NICType = "tap"
	Bridge NICType = "bridge"
	User   NICType = "user"
	Socket NICType = "socket"
)

type TapOptions struct {
	ID     string `json:"id,omitempty"`
	IFName string `json:"ifname,omitempty"`
}

func (t *TapOptions) ExpandOptions() []string {
	return []string{"tap," + strings.Join(Expand(t), ",")}
}

type BridgeOptions struct {
}

func (b *BridgeOptions) ExpandOptions() []string {
	return []string{}
}

type UserOptions struct {
}

func (u *UserOptions) ExpandOptions() []string {
	return []string{}
}

type SocketOptions struct {
	ID           string `json:"id,omitempty"`
	FD           string `json:"fd,omitempty"`
	MCast        string `json:"mcast,omitempty"`
	UDP          string `json:"udp,omitempty"`
	LocalAddress string `json:"localaddr,omitempty"`
	Listen       string `json:"listen,omitempty"`
	Connect      string `json:"connect,omitempty"`
}

func (s *SocketOptions) ExpandOptions() []string {
	return []string{}
}

type Options interface {
	ExpandOptions() []string
}

type NIC struct {
	Type       NICType `json:"type,omitempty" omit:"true"`
	MacAddress string  `json:"mac,omitempty"`
	Options    Options `json:"options,omitempty"`
}
