package server

type Action int32

const (
	ActionConnect Action = iota
	ActionAnnounce
	ActionScrape
	ActionError
)

type Event int32

const (
	EventNone Event = iota
	EventCompleted
	EventStarted
	EventStopped
)

type RequestHeader struct {
	ConnectionId int64
	Action       Action
}

// Connect
type ConnectRequest struct {
	TransactionId int32
}

type ConnectResponse struct {
	Action        Action
	TransactionId int32
	ConnectionId  int64
}

type AnnounceRequest struct {
	TransactionId int32
	InfoHash      [20]byte
	PeerId        [20]byte
	Downloaded    int64
	Left          int64
	Uploaded      int64
	Event         Event
	IpAddress     uint32
	Key           uint32
	NumWant       int32
	Port          uint16
}

type Peer struct {
	Ip   int32
	Port uint16
}

type AnnounceResponse struct {
	Action        Action
	TransactionId int32
	Interval      int32
	Leechers      int32
	Seeders       int32
}

type InfoHash struct {
	InfoHash [20]byte
}

type ScrapeRequest struct {
	TransactionId int32
}

type TorrentInfo struct {
	Seeders   int32
	Completed int32
	Leechers  int32
}

type ScrapeResponse struct {
	Action        Action
	TransactionId int32
}
