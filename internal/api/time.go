package api

import "time"

// TimeLayout used in the Weaviate server.
const TimeLayout = time.RFC3339

// timeFromString parses a timestamp formatted in [TimeLayout].
func timeFromString(s string) (time.Time, error) {
	return time.Parse(TimeLayout, s)
}

// timeFromUnix parses a UNIX timestamp, which, in the Weaviate server,
// uses a millisecond precision.
func timeFromUnix(ts int64) time.Time {
	return time.UnixMilli(ts)
}
