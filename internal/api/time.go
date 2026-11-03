package api

import "time"

// TimeLayout used in the Weaviate server.
const TimeLayout = time.RFC3339

// timeFromString parses a timestamp formatted in [TimeLayout].
func timeFromString(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(TimeLayout, s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// timeFromUnix parses a UNIX timestamp, which, in the Weaviate server,
// uses a millisecond precision.
func timeFromUnix(ts int64) time.Time {
	return time.UnixMilli(ts)
}

// durationFromSeconds converts an RPC took value. The protobuf field is a
// float number of seconds, so the fractional part has to be kept.
func durationFromSeconds(seconds float32) time.Duration {
	return time.Duration(float64(seconds) * float64(time.Second))
}
