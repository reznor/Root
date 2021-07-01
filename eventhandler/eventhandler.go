package eventhandler

import "root.challenge/eventstore"

// EventType defines the unique keyword with which each event is registered with the system.
type EventType string

// EventArgs defines the arguments for each event that are passed in to its registered handler at runtime.
type EventArgs []string

// Interface defines the runtime operations for handling each event that enters the system.
//
// See this package's README.md for the steps required when adding new implementations of this interface.
//
// ============================================== Maintainer Notes ==============================================
//
// As the handling of each event gets more complex, the list of these operations should be expanded (for example,
// to provide:
//
// a) a one-time system-wide Setup() and Teardown() hook for each handler,
// b) a CheckPreconditions() hook to implement the Template Method design pattern
//    (https://en.wikipedia.org/wiki/Template_method_pattern) that allows each handler to provide quick, light
//    sanity checks that can be invoked by the framework before invoking the potentially more-expensive-to-spin-up
//    `Handle`() method,
// c) observability hooks that are invoked every so often by the framework to collect telemetry (and free each
//    handler from having to worry about the infrastructure setup for said observability),
//
// etc.)
type Interface interface {
	// `Handle`() is invoked at runtime for each new event of the corresponding `EventType` that enters the system.
	//
	// It is worth noting that since `eventstore.EventStore` is passed in to each call to this method (as opposed to
	// being bound to each handler once at startup time), the framework has the flexibility to pass in different
	// instances of `eventstore.EventStore` in successive calls to the same handler (for example, to perform
	// hot-swaps/upgrades/load-shedding/rotation invisibly with zero downtime).
	Handle(EventArgs, *eventstore.EventStore) error
}
