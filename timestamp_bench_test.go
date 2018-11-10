package main

import (
	"testing"
	"time"

	types "github.com/gogo/protobuf/types"
)

var t = time.Unix(1, 0)
var ts = types.Timestamp{Seconds: 1, Nanos: 0}

var i = &Int64{1}
var myTimestamp = &MyTimestamp{1, 0}
var buf = make([]byte, 1024*1024)

type serializable interface {
	MarshalTo(data []byte) (int, error)
	Unmarshal(data []byte) error
}

// Benchmark Proto3 Marshal
func BenchmarkProto3(b *testing.B) {
	tests := map[string]serializable{
		"Int64":                                 i,
		"MyTimestamp":                           myTimestamp,
		"Embedded":                              &Embedded{myTimestamp},
		"EmbeddedStdTime":                       &EmbeddedStdTime{&t},
		"EmbeddedStdTimeNonNull":                &EmbeddedStdTimeNonNull{t},
		"EmbeddedGoogleTimestamp":               &EmbeddedGoogleTimestamp{&ts},
		"EmbeddedGoogleTimestampStdTime":        &EmbeddedGoogleTimestamp{&ts},
		"EmbeddedGoogleTimestampNonNull":        &EmbeddedGoogleTimestampNonNull{ts},
		"EmbeddedGoogleTimestampStdTimeNonNull": &EmbeddedGoogleTimestampStdTimeNonNull{t},
	}

	for name, msg := range tests {
		var n int
		var err error
		b.Run("marshal/"+name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				n, err = msg.MarshalTo(buf)
				if err != nil {
					b.Fatal("Marshaling error:", err)
				}
			}
		})

		b.Run("unmarshal/"+name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				err := msg.Unmarshal(buf[:n])
				if err != nil {
					b.Fatal("Marshaling error:", err)
				}
			}
		})
	}
}
