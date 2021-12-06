// modified stathat.com/c/consistent
// 1. use treemap to replace map (n^2klognk -> nklognk)
// 2. support murmur3 hash

package consistent

import (
	"testing"
)

var c *Circle

func TestCircle_Add(t *testing.T) {
	type args struct {
		realNode string
	}
	c = NewCircle()
	tests := []struct {
		name string
		c    *Circle
		args args
	}{
		{"1", c, args{"node1"}},
		{"2", c, args{"node2"}},
		{"3", c, args{"node3"}},
		{"4", c, args{"node4"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Add(tt.args.realNode)
		})
	}
}

func TestCircle_Del(t *testing.T) {
	type args struct {
		realNode string
	}
	tests := []struct {
		name string
		c    *Circle
		args args
	}{
		{"1", c, args{"node1"}},
		{"2", c, args{"node2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Del(tt.args.realNode)
		})
	}
}

func TestCircle_GetNode(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		c       *Circle
		args    args
		want    string
		wantErr bool
	}{
		{"1", c, args{"123"}, "node4", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Allocate(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Circle.Allocate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Circle.Allocate() = %v, want %v", got, tt.want)
			}
		})
	}
}
