package eventhandler_test

import (
	"math/rand"
	"testing"

	"root.challenge/eventhandler"
	"root.challenge/eventstore"
)

const testEventType eventhandler.EventType = "TestEvent"

type testEventHandler struct {
	id int
}

func newTestEventHandler() *testEventHandler {
	return &testEventHandler{
		id: rand.Int(),
	}
}

func (teh *testEventHandler) Handle(eventhandler.EventArgs, *eventstore.EventStore) error {
	return nil
}

func TestVanillaUsage(t *testing.T) {
	r := eventhandler.GlobalRegistry()

	registeredEventHandler := newTestEventHandler()
	expectedId := registeredEventHandler.id

	if err := r.RegisterEventHandler(testEventType, registeredEventHandler); err != nil {
		t.Fatalf("RegisterEventHandler() expected: no error, got: %v", err)
	}

	retrievedEventHandler, err := r.GetHandlerForEvent(testEventType)
	if err != nil {
		t.Fatalf("GetHandlerForEvent() expected: no error, got: %v", err)
	}

	actualId := retrievedEventHandler.(*testEventHandler).id
	if actualId != expectedId {
		t.Fatalf("GetHandlerForEvent() expected id: %v, got id: %v", expectedId, actualId)
	}
}

func TestRegisteringNilHandler(t *testing.T) {
	r := eventhandler.NewRegistry()

	if err := r.RegisterEventHandler(testEventType, nil); err == nil {
		t.Fatalf("RegisterEventHandler() expected: error, got: no error")
	}
}

func TestDuplicateRegistration(t *testing.T) {
	r := eventhandler.NewRegistry()

	registeredEventHandler := newTestEventHandler()
	expectedId := registeredEventHandler.id

	if err := r.RegisterEventHandler(testEventType, registeredEventHandler); err != nil {
		t.Fatalf("RegisterEventHandler() expected: no error, got: %v", err)
	}

	if err := r.RegisterEventHandler(testEventType, newTestEventHandler()); err == nil {
		t.Fatalf("RegisterEventHandler() expected: error, got: no error")
	}

	retrievedEventHandler, err := r.GetHandlerForEvent(testEventType)
	if err != nil {
		t.Fatalf("GetHandlerForEvent() expected: no error, got: %v", err)
	}

	actualId := retrievedEventHandler.(*testEventHandler).id
	if actualId != expectedId {
		t.Fatalf("GetHandlerForEvent() expected id: %v, got id: %v", expectedId, actualId)
	}
}

func TestRetrievalWithoutRegistration(t *testing.T) {
	r := eventhandler.NewRegistry()

	if _, err := r.GetHandlerForEvent(testEventType); err == nil {
		t.Fatalf("GetHandlerForEvent() expected: error, got: no error")
	}
}
