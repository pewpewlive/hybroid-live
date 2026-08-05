package alerts

type Collector struct {
	alerts []Alert
}

func NewCollector() Collector {
	return Collector{
		alerts: make([]Alert, 0),
	}
}

func (c *Collector) Report(alert Alert) {
	c.alerts = append(c.alerts, alert)
}

func (c *Collector) GetAlerts() []Alert {
	return c.alerts
}
