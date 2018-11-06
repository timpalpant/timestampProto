package main

import (
	"testing"
	"time"

	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
)

var t = time.Unix(1, 0)
var ts = types.Timestamp{Seconds: 1, Nanos: 0}

var i = &Int64{1}
var myTimestamp = &MyTimestamp{1, 0}

// Benchmark Proto3 Marshal
func BenchmarkProto3(b *testing.B) {
	tests := map[string]proto.Message{
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
		var marshaledMsg []byte
		b.Run("marshal/"+name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tmp, err := proto.Marshal(msg)
				if err != nil {
					b.Fatal("Marshaling error:", err)
				}
				marshaledMsg = tmp
			}
		})

		b.Run("unmarshal/"+name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				err := proto.Unmarshal(marshaledMsg, msg)
				if err != nil {
					b.Fatal("Marshaling error:", err)
				}
			}
		})
	}
}
