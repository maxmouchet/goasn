package collectors

import (
	"fmt"
	"time"
)

type RISCollector struct {
	name string
}

func (c RISCollector) Name() string {
	return c.name
}

func (c RISCollector) baseURL() string {
	return fmt.Sprintf("http://data.ris.ripe.net/%s", c.name)
}

func (c RISCollector) TableURL(t time.Time) string {
	return fmt.Sprintf(
		"%s/%s/bview.%s.gz",
		c.baseURL(),
		t.UTC().Format("2006.01"),
		t.UTC().Format("20060102.1500"),
	)
}
