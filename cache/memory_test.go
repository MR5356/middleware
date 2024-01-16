package cache

import (
	"reflect"
	"testing"
	"time"
)

var cache *MemoryCache

func TestMain(m *testing.M) {
	cache = NewMemoryCache()
	_ = cache.Set("key", "value")
	_ = cache.SetEx("key-expire-1", "value", time.Second)
	_ = cache.SetEx("key-expire-2", "value", 1)
	m.Run()
}

func TestMemoryCache_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{name: "get", args: args{key: "key"}, want: "value", wantErr: false},
		{name: "get not exist", args: args{key: "key-not-exist"}, want: nil, wantErr: true},
		{name: "get expired", args: args{key: "key-expire-1"}, want: "value", wantErr: false},
		{name: "get expired", args: args{key: "key-expire-2"}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := cache
			got, err := c.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryCache_Set(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "set",
			args: args{
				key:   "key",
				value: "value",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := cache
			if err := c.Set(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
