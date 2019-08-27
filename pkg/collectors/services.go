package collectors

import (
	"fmt"
	"regexp"
	"time"
)

var collectorPattern = regexp.MustCompile(`^(.+)\.(routeviews|oregon-ix|ripe)\.\w+`)

type Collector interface {
	Name() string
	TableURL(t time.Time) string
}

func NewCollector(fqdn string) (Collector, error) {
	matches := collectorPattern.FindStringSubmatch(fqdn)
	if matches == nil {
		return nil, fmt.Errorf("Cannot parse %s", fqdn)
	}

	name := matches[1]
	service := matches[2]

	switch service {
	case "ripe":
		var collector Collector = RISCollector{name}
		return collector, nil

	case "routeviews", "oregon-ix":
		var collector Collector = RouteViewsCollector{name}
		return collector, nil
	}

	return nil, fmt.Errorf("Unknown service %s", service)
}
