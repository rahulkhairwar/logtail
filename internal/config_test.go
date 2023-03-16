package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_String(t *testing.T) {
	type fields struct {
		Port       int
		FileToTail string
		ParseKeys  []string
		LogLevels  map[string]level
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Port:       tt.fields.Port,
				FileToTail: tt.fields.FileToTail,
				ParseKeys:  tt.fields.ParseKeys,
				LogLevels:  tt.fields.LogLevels,
			}
			assert.Equalf(t, tt.want, c.String(), "String()")
		})
	}
}

func Test_level_String(t *testing.T) {
	type fields struct {
		Key   string
		Color string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := level{
				Key:   tt.fields.Key,
				Color: tt.fields.Color,
			}
			assert.Equalf(t, tt.want, l.String(), "String()")
		})
	}
}
