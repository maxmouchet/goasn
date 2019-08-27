package collectors

import (
	"testing"
	"time"
)

func TestRISCollector_TableURL(t *testing.T) {
	type fields struct {
		name string
	}
	type args struct {
		t time.Time
	}

	t1, _ := time.Parse("2006-01-02 15:04", "2019-08-01 08:00")

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"base", fields{"rrc00"}, args{t1}, "http://data.ris.ripe.net/rrc00/2019.08/bview.20190801.0800.gz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := RISCollector{
				name: tt.fields.name,
			}
			if got := c.TableURL(tt.args.t); got != tt.want {
				t.Errorf("RISCollector.TableURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
