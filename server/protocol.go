package server

type action int32

const (
	actionConnect action = iota
	actionAnnounce
	actionScrape
	actionError
)

type event int32

const (
	eventNone event = iota
	eventCompleted
	eventStarted
	eventStopped
)

type requestHeader struct {
	ConnectionId int64
	Action       action
}

type connectRequest struct {
	TransactionId int32
}

type connectResponse struct {
	Action        action
	TransactionId int32
	ConnectionId  int64
}

type announceRequest struct {
	TransactionId int32
	InfoHash      [20]byte
	PeerId        [20]byte
	Downloaded    int64
	Left          int64
	Uploaded      int64
	Event         event
	IpAddress     uint32
	Key           uint32
	NumWant       int32
	Port          uint16
}

type peer struct {
	Ip   int32
	Port uint16
}

type announceResponse struct {
	Action        action
	TransactionId int32
	Interval      int32
	Leechers      int32
	Seeders       int32
}

type scrapeRequest struct {
	TransactionId int32
}

type torrentInfo struct {
	Seeders   int32
	Completed int32
	Leechers  int32
}

type scrapeResponse struct {
	Action        action
	TransactionId int32
}
