package queue

import "github.com/golang/protobuf/ptypes/timestamp"

type ConnResult struct {
	Flag bool
	CommComp
	QueueId string
	SeatId string
	UserId string
}

type IOQueue struct {
	CommComp
	QueueId string
	UserId string
}

type SeatOut struct {
	CommComp
	SeatId string
	Reason string
}

type SeatQueueChange struct {
	CommComp
	SeatId string
	QueueIds []string
}

type CommComp struct {
	TimeStamp timestamp.Timestamp
	CompId string
}