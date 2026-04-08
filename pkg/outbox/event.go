package outbox

type DomainEvent interface {
	EventName() string
}

type EventCollector struct {
	events []DomainEvent
}

func (c *EventCollector) Add(events ...DomainEvent) {
	c.events = append(c.events, events...)
}
