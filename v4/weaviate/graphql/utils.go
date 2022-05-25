package graphql

import "bytes"

// marshalStrings marshals []string into a json like []byte
func marshalStrings(xs []string) []byte {
	n := len(xs)
	var buf bytes.Buffer
	buf.WriteByte('[')
	for _, x := range xs {
		buf.WriteByte('"')
		buf.WriteString(x)
		buf.WriteByte('"')
		n--
		if n > 0 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')
	return buf.Bytes()
}
