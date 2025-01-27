// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package db_mock

import (
	"github.com/keptn/keptn/shipyard-controller/models"
	"sync"
	"time"
)

// EventQueueRepoMock is a mock implementation of db.EventQueueRepo.
//
// 	func TestSomethingThatUsesEventQueueRepo(t *testing.T) {
//
// 		// make and configure a mocked db.EventQueueRepo
// 		mockedEventQueueRepo := &EventQueueRepoMock{
// 			DeleteQueuedEventFunc: func(eventID string) error {
// 				panic("mock out the DeleteQueuedEvent method")
// 			},
// 			DeleteQueuedEventsFunc: func(scope models.EventScope) error {
// 				panic("mock out the DeleteQueuedEvents method")
// 			},
// 			GetQueuedEventsFunc: func(timestamp time.Time) ([]models.QueueItem, error) {
// 				panic("mock out the GetQueuedEvents method")
// 			},
// 			IsEventInQueueFunc: func(eventID string) (bool, error) {
// 				panic("mock out the IsEventInQueue method")
// 			},
// 			QueueEventFunc: func(item models.QueueItem) error {
// 				panic("mock out the QueueEvent method")
// 			},
// 		}
//
// 		// use mockedEventQueueRepo in code that requires db.EventQueueRepo
// 		// and then make assertions.
//
// 	}
type EventQueueRepoMock struct {
	// DeleteQueuedEventFunc mocks the DeleteQueuedEvent method.
	DeleteQueuedEventFunc func(eventID string) error

	// DeleteQueuedEventsFunc mocks the DeleteQueuedEvents method.
	DeleteQueuedEventsFunc func(scope models.EventScope) error

	// GetQueuedEventsFunc mocks the GetQueuedEvents method.
	GetQueuedEventsFunc func(timestamp time.Time) ([]models.QueueItem, error)

	// IsEventInQueueFunc mocks the IsEventInQueue method.
	IsEventInQueueFunc func(eventID string) (bool, error)

	// QueueEventFunc mocks the QueueEvent method.
	QueueEventFunc func(item models.QueueItem) error

	// calls tracks calls to the methods.
	calls struct {
		// DeleteQueuedEvent holds details about calls to the DeleteQueuedEvent method.
		DeleteQueuedEvent []struct {
			// EventID is the eventID argument value.
			EventID string
		}
		// DeleteQueuedEvents holds details about calls to the DeleteQueuedEvents method.
		DeleteQueuedEvents []struct {
			// Scope is the scope argument value.
			Scope models.EventScope
		}
		// GetQueuedEvents holds details about calls to the GetQueuedEvents method.
		GetQueuedEvents []struct {
			// Timestamp is the timestamp argument value.
			Timestamp time.Time
		}
		// IsEventInQueue holds details about calls to the IsEventInQueue method.
		IsEventInQueue []struct {
			// EventID is the eventID argument value.
			EventID string
		}
		// QueueEvent holds details about calls to the QueueEvent method.
		QueueEvent []struct {
			// Item is the item argument value.
			Item models.QueueItem
		}
	}
	lockDeleteQueuedEvent  sync.RWMutex
	lockDeleteQueuedEvents sync.RWMutex
	lockGetQueuedEvents    sync.RWMutex
	lockIsEventInQueue     sync.RWMutex
	lockQueueEvent         sync.RWMutex
}

// DeleteQueuedEvent calls DeleteQueuedEventFunc.
func (mock *EventQueueRepoMock) DeleteQueuedEvent(eventID string) error {
	if mock.DeleteQueuedEventFunc == nil {
		panic("EventQueueRepoMock.DeleteQueuedEventFunc: method is nil but EventQueueRepo.DeleteQueuedEvent was just called")
	}
	callInfo := struct {
		EventID string
	}{
		EventID: eventID,
	}
	mock.lockDeleteQueuedEvent.Lock()
	mock.calls.DeleteQueuedEvent = append(mock.calls.DeleteQueuedEvent, callInfo)
	mock.lockDeleteQueuedEvent.Unlock()
	return mock.DeleteQueuedEventFunc(eventID)
}

// DeleteQueuedEventCalls gets all the calls that were made to DeleteQueuedEvent.
// Check the length with:
//     len(mockedEventQueueRepo.DeleteQueuedEventCalls())
func (mock *EventQueueRepoMock) DeleteQueuedEventCalls() []struct {
	EventID string
} {
	var calls []struct {
		EventID string
	}
	mock.lockDeleteQueuedEvent.RLock()
	calls = mock.calls.DeleteQueuedEvent
	mock.lockDeleteQueuedEvent.RUnlock()
	return calls
}

// DeleteQueuedEvents calls DeleteQueuedEventsFunc.
func (mock *EventQueueRepoMock) DeleteQueuedEvents(scope models.EventScope) error {
	if mock.DeleteQueuedEventsFunc == nil {
		panic("EventQueueRepoMock.DeleteQueuedEventsFunc: method is nil but EventQueueRepo.DeleteQueuedEvents was just called")
	}
	callInfo := struct {
		Scope models.EventScope
	}{
		Scope: scope,
	}
	mock.lockDeleteQueuedEvents.Lock()
	mock.calls.DeleteQueuedEvents = append(mock.calls.DeleteQueuedEvents, callInfo)
	mock.lockDeleteQueuedEvents.Unlock()
	return mock.DeleteQueuedEventsFunc(scope)
}

// DeleteQueuedEventsCalls gets all the calls that were made to DeleteQueuedEvents.
// Check the length with:
//     len(mockedEventQueueRepo.DeleteQueuedEventsCalls())
func (mock *EventQueueRepoMock) DeleteQueuedEventsCalls() []struct {
	Scope models.EventScope
} {
	var calls []struct {
		Scope models.EventScope
	}
	mock.lockDeleteQueuedEvents.RLock()
	calls = mock.calls.DeleteQueuedEvents
	mock.lockDeleteQueuedEvents.RUnlock()
	return calls
}

// GetQueuedEvents calls GetQueuedEventsFunc.
func (mock *EventQueueRepoMock) GetQueuedEvents(timestamp time.Time) ([]models.QueueItem, error) {
	if mock.GetQueuedEventsFunc == nil {
		panic("EventQueueRepoMock.GetQueuedEventsFunc: method is nil but EventQueueRepo.GetQueuedEvents was just called")
	}
	callInfo := struct {
		Timestamp time.Time
	}{
		Timestamp: timestamp,
	}
	mock.lockGetQueuedEvents.Lock()
	mock.calls.GetQueuedEvents = append(mock.calls.GetQueuedEvents, callInfo)
	mock.lockGetQueuedEvents.Unlock()
	return mock.GetQueuedEventsFunc(timestamp)
}

// GetQueuedEventsCalls gets all the calls that were made to GetQueuedEvents.
// Check the length with:
//     len(mockedEventQueueRepo.GetQueuedEventsCalls())
func (mock *EventQueueRepoMock) GetQueuedEventsCalls() []struct {
	Timestamp time.Time
} {
	var calls []struct {
		Timestamp time.Time
	}
	mock.lockGetQueuedEvents.RLock()
	calls = mock.calls.GetQueuedEvents
	mock.lockGetQueuedEvents.RUnlock()
	return calls
}

// IsEventInQueue calls IsEventInQueueFunc.
func (mock *EventQueueRepoMock) IsEventInQueue(eventID string) (bool, error) {
	if mock.IsEventInQueueFunc == nil {
		panic("EventQueueRepoMock.IsEventInQueueFunc: method is nil but EventQueueRepo.IsEventInQueue was just called")
	}
	callInfo := struct {
		EventID string
	}{
		EventID: eventID,
	}
	mock.lockIsEventInQueue.Lock()
	mock.calls.IsEventInQueue = append(mock.calls.IsEventInQueue, callInfo)
	mock.lockIsEventInQueue.Unlock()
	return mock.IsEventInQueueFunc(eventID)
}

// IsEventInQueueCalls gets all the calls that were made to IsEventInQueue.
// Check the length with:
//     len(mockedEventQueueRepo.IsEventInQueueCalls())
func (mock *EventQueueRepoMock) IsEventInQueueCalls() []struct {
	EventID string
} {
	var calls []struct {
		EventID string
	}
	mock.lockIsEventInQueue.RLock()
	calls = mock.calls.IsEventInQueue
	mock.lockIsEventInQueue.RUnlock()
	return calls
}

// QueueEvent calls QueueEventFunc.
func (mock *EventQueueRepoMock) QueueEvent(item models.QueueItem) error {
	if mock.QueueEventFunc == nil {
		panic("EventQueueRepoMock.QueueEventFunc: method is nil but EventQueueRepo.QueueEvent was just called")
	}
	callInfo := struct {
		Item models.QueueItem
	}{
		Item: item,
	}
	mock.lockQueueEvent.Lock()
	mock.calls.QueueEvent = append(mock.calls.QueueEvent, callInfo)
	mock.lockQueueEvent.Unlock()
	return mock.QueueEventFunc(item)
}

// QueueEventCalls gets all the calls that were made to QueueEvent.
// Check the length with:
//     len(mockedEventQueueRepo.QueueEventCalls())
func (mock *EventQueueRepoMock) QueueEventCalls() []struct {
	Item models.QueueItem
} {
	var calls []struct {
		Item models.QueueItem
	}
	mock.lockQueueEvent.RLock()
	calls = mock.calls.QueueEvent
	mock.lockQueueEvent.RUnlock()
	return calls
}
