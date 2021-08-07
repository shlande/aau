package record

import (
	"testing"
	"time"
)

func TestJsonKV_Set(t *testing.T) {
	data := &record{Name: "lafgewga\"", Url: "fdalief", Episode: 10, Time: time.Now()}
	key := "sfe{sae"
	tests := []struct {
		name string
		args string
	}{
		{"empty", "./test/empty.json"},
		{"len0", "./test/len0.json"},
		{"len1", "./test/len1.json"},
		{"len1_pretty", "./test/len1_pretty.json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewJsonKVFromFile(tt.args)
			for i := 0; i < 10; i++ {
				if err := r.Set(key, data); err != nil {
					panic(err)
				}
			}
		})
	}
}
