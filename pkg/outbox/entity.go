package outbox

type EntityWithDomainEvent interface {
	GetDomainEvents() []DomainEvent
}
