// Copyright 2018 Zhijie Shen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fsm

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const (
	// states
	s1 State = "s1"
	s2 State = "s2"
	s3 State = "s3"
	s4 State = "s4"
	// event types
	et1 EventType = "et1"
	et2 EventType = "et2"
	et3 EventType = "et3"
	et4 EventType = "et4"
)

type (
	evt1 struct{}
	evt2 struct{}
	evt3 struct{}
	evt4 struct {
		flag bool
	}
)

func (e *evt1) Type() EventType {
	return et1
}

func (e *evt2) Type() EventType {
	return et2
}

func (e *evt3) Type() EventType {
	return et3
}

func (e *evt4) Type() EventType {
	return et4
}

func TestSimpleFSM(t *testing.T) {
	// Building a simple 4 state FSM
	// s1 -- evt1 --> s2 -- evt2 --> s3 -- evt3 --> s4 --> evt4 --
	// ^               ^                                         |
	// |               |                                         |
	// |---------------------------------------------------------|
	fsm, err := NewBuilder().
		AddInitialState(s1).
		AddStates(s2, s3, s4).
		AddTransition(s1, et1, func(evt Event) (State, error) { return s2, nil }, []State{s2}).
		AddTransition(s2, et2, func(evt Event) (State, error) { return s3, nil }, []State{s3}).
		AddTransition(s3, et3, func(evt Event) (State, error) { return s4, nil }, []State{s4}).
		AddTransition(s4, et4, func(evt Event) (State, error) {
			e, ok := evt.(*evt4)
			if !ok {
				return s4, errors.New("invalid event")
			}
			if e.flag {
				return s1, nil
			}
			return s2, nil
		}, []State{s1, s2}).
		Build()

	require.Nil(t, err)
	require.NotNil(t, fsm)
	require.Equal(t, s1, fsm.CurrentState())

	fsm.Handle(&evt1{})
	require.Equal(t, s2, fsm.CurrentState())

	fsm.Handle(&evt2{})
	require.Equal(t, s3, fsm.CurrentState())

	fsm.Handle(&evt3{})
	require.Equal(t, s4, fsm.CurrentState())

	fsm.Handle(&evt4{flag: false})
	require.Equal(t, s2, fsm.CurrentState())

	fsm.Handle(&evt2{})
	require.Equal(t, s3, fsm.CurrentState())

	fsm.Handle(&evt3{})
	require.Equal(t, s4, fsm.CurrentState())

	fsm.Handle(&evt4{flag: true})
	require.Equal(t, s1, fsm.CurrentState())
}

func TestBuilder_BuildDupInitState(t *testing.T) {
	_, err := NewBuilder().
		AddInitialState(s1).
		AddInitialState(s2).
		AddStates(s3, s4).
		AddTransition(s1, et1, func(evt Event) (State, error) { return s2, nil }, []State{s2}).
		AddTransition(s2, et2, func(evt Event) (State, error) { return s3, nil }, []State{s3}).
		AddTransition(s3, et3, func(evt Event) (State, error) { return s4, nil }, []State{s4}).
		AddTransition(s4, et4, func(evt Event) (State, error) { return s1, nil }, []State{s1}).
		Build()
	require.NotNil(t, err)
	require.Equal(t, ErrBuild, errors.Cause(err))
}

func TestBuilder_BuildNoInitState(t *testing.T) {
	_, err := NewBuilder().
		AddStates(s1, s2, s3, s4).
		AddTransition(s1, et1, func(evt Event) (State, error) { return s2, nil }, []State{s2}).
		AddTransition(s2, et2, func(evt Event) (State, error) { return s3, nil }, []State{s3}).
		AddTransition(s3, et3, func(evt Event) (State, error) { return s4, nil }, []State{s4}).
		AddTransition(s4, et4, func(evt Event) (State, error) { return s1, nil }, []State{s1}).
		Build()
	require.NotNil(t, err)
	require.Equal(t, ErrBuild, errors.Cause(err))
}

func TestBuilder_BuildUndefinedSrc(t *testing.T) {
	_, err := NewBuilder().
		AddInitialState(s1).
		AddStates(s2, s4).
		AddTransition(s1, et1, func(evt Event) (State, error) { return s2, nil }, []State{s2}).
		AddTransition(s2, et2, func(evt Event) (State, error) { return s3, nil }, []State{s3}).
		AddTransition(s3, et3, func(evt Event) (State, error) { return s4, nil }, []State{s4}).
		Build()
	require.NotNil(t, err)
	require.Equal(t, ErrBuild, errors.Cause(err))
}

func TestBuilder_BuildUndefinedDst(t *testing.T) {
	_, err := NewBuilder().
		AddInitialState(s1).
		AddStates(s2, s3).
		AddTransition(s1, et1, func(evt Event) (State, error) { return s2, nil }, []State{s2}).
		AddTransition(s2, et2, func(evt Event) (State, error) { return s3, nil }, []State{s4}).
		AddTransition(s3, et3, func(evt Event) (State, error) { return s4, nil }, []State{s4}).
		Build()
	require.NotNil(t, err)
	require.Equal(t, ErrBuild, errors.Cause(err))
}

func TestFSM_InvalidTransition(t *testing.T) {
	fsm, err := NewBuilder().
		AddInitialState(s1).
		AddStates(s2, s3, s4).
		AddTransition(s1, et1, func(evt Event) (State, error) { return s3, nil }, []State{s2}).
		Build()
	require.Nil(t, err)
	require.NotNil(t, fsm)

	err = fsm.Handle(&evt1{})
	require.NotNil(t, err)
	require.Equal(t, ErrInvalidTransition, errors.Cause(err))
}

func TestFSM_CustomizedErr(t *testing.T) {
	e := errors.New("customized error")
	fsm, err := NewBuilder().
		AddInitialState(s1).
		AddStates(s2, s3, s4).
		AddTransition(s1, et1, func(evt Event) (State, error) { return s1, e }, []State{s2}).
		Build()
	require.Nil(t, err)
	require.NotNil(t, fsm)

	err = fsm.Handle(&evt1{})
	require.NotNil(t, err)
	require.Equal(t, e, err)
}
