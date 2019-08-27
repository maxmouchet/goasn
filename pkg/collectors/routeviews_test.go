package collectors

import (
	"testing"
	"time"
)

func TestRouteViewsCollector_TableURL(t *testing.T) {
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
		{"base", fields{"route-views.amsix"}, args{t1}, "http://archive.routeviews.org/route-views.amsix/bgpdata/2019.08/RIBS/rib.20190801.0800.bz2"},
		{"route-views2", fields{"route-views2"}, args{t1}, "http://archive.routeviews.org/bgpdata/2019.08/RIBS/rib.20190801.0800.bz2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := RouteViewsCollector{
				name: tt.fields.name,
			}
			if got := c.TableURL(tt.args.t); got != tt.want {
				t.Errorf("RouteViewsCollector.TableURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
