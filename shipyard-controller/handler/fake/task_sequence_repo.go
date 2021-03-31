// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fake

import (
	"github.com/keptn/keptn/shipyard-controller/handler"
	"sync"
)

// ITaskSequenceExecutionStateRepoMock is a mock implementation of handler.ITaskSequenceExecutionStateRepo.
//
// 	func TestSomethingThatUsesITaskSequenceExecutionStateRepo(t *testing.T) {
//
// 		// make and configure a mocked handler.ITaskSequenceExecutionStateRepo
// 		mockedITaskSequenceExecutionStateRepo := &ITaskSequenceExecutionStateRepoMock{
// 			GetTaskSequenceExecutionStateFunc: func(keptnContext string, stage string) *handler.TaskSequenceExecutionState {
// 				panic("mock out the GetTaskSequenceExecutionState method")
// 			},
// 			StoreTaskSequenceExecutionStateFunc: func(state handler.TaskSequenceExecutionState) error {
// 				panic("mock out the StoreTaskSequenceExecutionState method")
// 			},
// 			UpdateTaskSequenceExecutionStateFunc: func(state handler.TaskSequenceExecutionState) error {
// 				panic("mock out the UpdateTaskSequenceExecutionState method")
// 			},
// 		}
//
// 		// use mockedITaskSequenceExecutionStateRepo in code that requires handler.ITaskSequenceExecutionStateRepo
// 		// and then make assertions.
//
// 	}
type ITaskSequenceExecutionStateRepoMock struct {
	// GetTaskSequenceExecutionStateFunc mocks the GetTaskSequenceExecutionState method.
	GetTaskSequenceExecutionStateFunc func(keptnContext string, stage string) *handler.TaskSequenceExecutionState

	// StoreTaskSequenceExecutionStateFunc mocks the StoreTaskSequenceExecutionState method.
	StoreTaskSequenceExecutionStateFunc func(state handler.TaskSequenceExecutionState) error

	// UpdateTaskSequenceExecutionStateFunc mocks the UpdateTaskSequenceExecutionState method.
	UpdateTaskSequenceExecutionStateFunc func(state handler.TaskSequenceExecutionState) error

	// calls tracks calls to the methods.
	calls struct {
		// GetTaskSequenceExecutionState holds details about calls to the GetTaskSequenceExecutionState method.
		GetTaskSequenceExecutionState []struct {
			// KeptnContext is the keptnContext argument value.
			KeptnContext string
			// Stage is the stage argument value.
			Stage string
		}
		// StoreTaskSequenceExecutionState holds details about calls to the StoreTaskSequenceExecutionState method.
		StoreTaskSequenceExecutionState []struct {
			// State is the state argument value.
			State handler.TaskSequenceExecutionState
		}
		// UpdateTaskSequenceExecutionState holds details about calls to the UpdateTaskSequenceExecutionState method.
		UpdateTaskSequenceExecutionState []struct {
			// State is the state argument value.
			State handler.TaskSequenceExecutionState
		}
	}
	lockGetTaskSequenceExecutionState    sync.RWMutex
	lockStoreTaskSequenceExecutionState  sync.RWMutex
	lockUpdateTaskSequenceExecutionState sync.RWMutex
}

// GetTaskSequenceExecutionState calls GetTaskSequenceExecutionStateFunc.
func (mock *ITaskSequenceExecutionStateRepoMock) GetTaskSequenceExecutionState(keptnContext string, stage string) *handler.TaskSequenceExecutionState {
	if mock.GetTaskSequenceExecutionStateFunc == nil {
		panic("ITaskSequenceExecutionStateRepoMock.GetTaskSequenceExecutionStateFunc: method is nil but ITaskSequenceExecutionStateRepo.GetTaskSequenceExecutionState was just called")
	}
	callInfo := struct {
		KeptnContext string
		Stage        string
	}{
		KeptnContext: keptnContext,
		Stage:        stage,
	}
	mock.lockGetTaskSequenceExecutionState.Lock()
	mock.calls.GetTaskSequenceExecutionState = append(mock.calls.GetTaskSequenceExecutionState, callInfo)
	mock.lockGetTaskSequenceExecutionState.Unlock()
	return mock.GetTaskSequenceExecutionStateFunc(keptnContext, stage)
}

// GetTaskSequenceExecutionStateCalls gets all the calls that were made to GetTaskSequenceExecutionState.
// Check the length with:
//     len(mockedITaskSequenceExecutionStateRepo.GetTaskSequenceExecutionStateCalls())
func (mock *ITaskSequenceExecutionStateRepoMock) GetTaskSequenceExecutionStateCalls() []struct {
	KeptnContext string
	Stage        string
} {
	var calls []struct {
		KeptnContext string
		Stage        string
	}
	mock.lockGetTaskSequenceExecutionState.RLock()
	calls = mock.calls.GetTaskSequenceExecutionState
	mock.lockGetTaskSequenceExecutionState.RUnlock()
	return calls
}

// StoreTaskSequenceExecutionState calls StoreTaskSequenceExecutionStateFunc.
func (mock *ITaskSequenceExecutionStateRepoMock) StoreTaskSequenceExecutionState(state handler.TaskSequenceExecutionState) error {
	if mock.StoreTaskSequenceExecutionStateFunc == nil {
		panic("ITaskSequenceExecutionStateRepoMock.StoreTaskSequenceExecutionStateFunc: method is nil but ITaskSequenceExecutionStateRepo.StoreTaskSequenceExecutionState was just called")
	}
	callInfo := struct {
		State handler.TaskSequenceExecutionState
	}{
		State: state,
	}
	mock.lockStoreTaskSequenceExecutionState.Lock()
	mock.calls.StoreTaskSequenceExecutionState = append(mock.calls.StoreTaskSequenceExecutionState, callInfo)
	mock.lockStoreTaskSequenceExecutionState.Unlock()
	return mock.StoreTaskSequenceExecutionStateFunc(state)
}

// StoreTaskSequenceExecutionStateCalls gets all the calls that were made to StoreTaskSequenceExecutionState.
// Check the length with:
//     len(mockedITaskSequenceExecutionStateRepo.StoreTaskSequenceExecutionStateCalls())
func (mock *ITaskSequenceExecutionStateRepoMock) StoreTaskSequenceExecutionStateCalls() []struct {
	State handler.TaskSequenceExecutionState
} {
	var calls []struct {
		State handler.TaskSequenceExecutionState
	}
	mock.lockStoreTaskSequenceExecutionState.RLock()
	calls = mock.calls.StoreTaskSequenceExecutionState
	mock.lockStoreTaskSequenceExecutionState.RUnlock()
	return calls
}

// UpdateTaskSequenceExecutionState calls UpdateTaskSequenceExecutionStateFunc.
func (mock *ITaskSequenceExecutionStateRepoMock) UpdateTaskSequenceExecutionState(state handler.TaskSequenceExecutionState) error {
	if mock.UpdateTaskSequenceExecutionStateFunc == nil {
		panic("ITaskSequenceExecutionStateRepoMock.UpdateTaskSequenceExecutionStateFunc: method is nil but ITaskSequenceExecutionStateRepo.UpdateTaskSequenceExecutionState was just called")
	}
	callInfo := struct {
		State handler.TaskSequenceExecutionState
	}{
		State: state,
	}
	mock.lockUpdateTaskSequenceExecutionState.Lock()
	mock.calls.UpdateTaskSequenceExecutionState = append(mock.calls.UpdateTaskSequenceExecutionState, callInfo)
	mock.lockUpdateTaskSequenceExecutionState.Unlock()
	return mock.UpdateTaskSequenceExecutionStateFunc(state)
}

// UpdateTaskSequenceExecutionStateCalls gets all the calls that were made to UpdateTaskSequenceExecutionState.
// Check the length with:
//     len(mockedITaskSequenceExecutionStateRepo.UpdateTaskSequenceExecutionStateCalls())
func (mock *ITaskSequenceExecutionStateRepoMock) UpdateTaskSequenceExecutionStateCalls() []struct {
	State handler.TaskSequenceExecutionState
} {
	var calls []struct {
		State handler.TaskSequenceExecutionState
	}
	mock.lockUpdateTaskSequenceExecutionState.RLock()
	calls = mock.calls.UpdateTaskSequenceExecutionState
	mock.lockUpdateTaskSequenceExecutionState.RUnlock()
	return calls
}