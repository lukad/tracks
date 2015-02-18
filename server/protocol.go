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
	Action       action `struc:"int32"`
}

type connectRequest struct {
	TransactionId int32
}

type connectResponse struct {
	Action        action `struc:"int32"`
	TransactionId int32  `struc:"int32"`
	ConnectionId  int64  `struc:"int64"`
}

type announceRequest struct {
	TransactionId int32
	InfoHash      []byte `struc:"[20]byte"`
	PeerId        []byte `struc:"[20]byte"`
	Downloaded    int64
	Left          int64
	Uploaded      int64
	Event         event `struc:"int32"`
	IpAddress     int32
	Key           int32
	NumWant       int32
	Port          int16
}
