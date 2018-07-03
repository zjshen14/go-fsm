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

package main

import (
	"errors"
	"log"

	"github.com/zjshen14/go-fsm"
)

const (
	// Turnstile has two states: Unlocked and Locked.
	unlocked fsm.State = "Unlocked"
	locked   fsm.State = "Locked"
	// Turnstile accepts two types of events: Coin and Push
	coin fsm.EventType = "Coin"
	push fsm.EventType = "Push"
)

type (
	coinEvt struct{ name string }
	pushEvt struct{ name string }
)

func (e *coinEvt) Type() fsm.EventType { return coin }

func (e *pushEvt) Type() fsm.EventType { return push }

// This is a simple example referred in Wikipedia: https://en.wikipedia.org/wiki/Finite-state_machine
func main() {
	// Create the turnstile FSM
	turnstile, err := fsm.NewBuilder().
		// The initial state is Locked
		AddInitialState(locked).
		AddStates(unlocked).
		AddTransition(locked, coin, func(event fsm.Event) (fsm.State, error) {
			cEvt, ok := event.(*coinEvt)
			if !ok {
				return locked, errors.New("invalid event")
			}
			log.Printf("Unlocks the turnstile so that %s can push through.", cEvt.name)
			return unlocked, nil
		}, []fsm.State{unlocked}).
		AddTransition(locked, push, func(event fsm.Event) (fsm.State, error) {
			log.Println("None")
			return locked, nil
		}, []fsm.State{unlocked}).
		AddTransition(unlocked, coin, func(event fsm.Event) (fsm.State, error) {
			log.Println("None")
			return unlocked, nil
		}, []fsm.State{unlocked}).
		AddTransition(unlocked, push, func(event fsm.Event) (fsm.State, error) {
			pEvt, ok := event.(*pushEvt)
			if !ok {
				return locked, errors.New("invalid event")
			}
			log.Printf("When %s has pushed through, locks the turnstile.", pEvt.name)
			return locked, nil
		}, []fsm.State{locked}).
		Build()

	if err != nil {
		log.Fatalf("error when building the FSM: %v", err)
	}

	// Run the turnstile FSM
	turnstile.Handle(&coinEvt{name: "zhijie"})
	turnstile.Handle(&pushEvt{name: "zhijie"})
	turnstile.Handle(&pushEvt{name: "zhijie"})
	turnstile.Handle(&coinEvt{name: "hui"})
	turnstile.Handle(&coinEvt{name: "hui"})
	turnstile.Handle(&pushEvt{name: "hui"})
}
