package platformsh

type Upstream struct {
	SocketFamily SocketFamily   `json:"socket_family"`
	Protocol     SocketProtocol `json:"socket_protocol"`
}
