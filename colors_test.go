package logtail

import (
	"image/color"
	"testing"
)

func TestColor_HexCode(t *testing.T) {
	type fields struct {
		hexCode string
		RGBA    color.RGBA
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "color yellow",
			fields: fields{
				hexCode: "FFFF00",
				RGBA: color.RGBA{
					R: 255,
					G: 255,
					B: 0,
				},
			},
			want: Yellow.hexCode,
		},
		{
			name: "color green",
			fields: fields{
				hexCode: "00FF00",
				RGBA: color.RGBA{
					R: 0,
					G: 255,
					B: 0,
				},
			},
			want: Green.hexCode,
		},
		{
			name: "color red",
			fields: fields{
				hexCode: "FF0000",
				RGBA: color.RGBA{
					R: 255,
					G: 0,
					B: 0,
				},
			},
			want: Red.hexCode,
		},
		{
			name: "color brown",
			fields: fields{
				hexCode: "663300",
				RGBA: color.RGBA{
					R: 102,
					G: 51,
					B: 0,
				},
			},
			want: Brown.hexCode,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Color{
				hexCode: tt.fields.hexCode,
				rgba:    tt.fields.RGBA,
			}
			if got := c.HexCode(); got != tt.want {
				t.Errorf("HexCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
