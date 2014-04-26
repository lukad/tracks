package server

const (
	actionConnect int32 = iota
	actionAnnounce
	actionScrape
	actionError
)

type requestHeader struct {
	ConnectionId int64
	Action       int32
}

type connectRequest struct {
	TransactionId int32
}

type connectResponse struct {
	Action        int32
	TransactionId int32
	ConnectionId  int64
}

type announceRequest struct {
	TransactionId int32
}
