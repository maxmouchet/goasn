package collectors

import (
	"fmt"
	"time"
)

type RouteViewsCollector struct {
	name string
}

func (c RouteViewsCollector) Name() string {
	return c.name
}

func (c RouteViewsCollector) baseURL() string {
	switch c.name {
	case "route-views2":
		return "http://archive.routeviews.org/bgpdata"
	default:
		return fmt.Sprintf("http://archive.routeviews.org/%s/bgpdata", c.name)
	}
}

func (c RouteViewsCollector) TableURL(t time.Time) string {
	return fmt.Sprintf(
		"%s/%s/RIBS/rib.%s.bz2",
		c.baseURL(),
		t.UTC().Format("2006.01"),
		t.UTC().Format("20060102.1500"),
	)
}
