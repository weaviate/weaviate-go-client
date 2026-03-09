package api

import "time"

// TimeLayout used in the Weaviate server.
const TimeLayout = time.RFC3339

// timeFromString parses a timestamp formatted in [TimeLayout].
func timeFromString(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(TimeLayout, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// timeFromUnix parses a UNIX timestamp, which, in the Weaviate server,
// uses a millisecond precision.
func timeFromUnix(ts int64) time.Time {
	return time.UnixMilli(ts)
}
