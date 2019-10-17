package bannersorter

import "time"

// Input defined methods a banner object should implement.
type Input interface {
	StartTime() time.Time
	EndTime() time.Time
}

// wrapper is a concrete implementation of a banner-encompassing object
type wrapper struct { // TODO private
	filling            interface{}
	startTime, endTime time.Time
}

// StartTime returns a start time for a banner
func (w wrapper) StartTime() time.Time {
	return w.startTime
}

// EndTime returns an end time for a banner
func (w wrapper) EndTime() time.Time {
	return w.endTime
}

// unwrap returns underlying object
func (w wrapper) unwrap() interface{} {
	return w.filling
}

// Wrap encloses a nondescript object in an envelop coupled with start and end times.
func Wrap(filling interface{}, startTime, endTime time.Time) Input {
	return wrapper{filling, startTime, endTime}
}
