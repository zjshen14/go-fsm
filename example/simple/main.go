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
	"fmt"

	"github.com/zjshen14/go-fsm"
)

type evt struct {
	t fsm.EventType
}

func (e *evt) Type() fsm.EventType { return e.t }

func main() {
	fsm, _ := fsm.NewBuilder().
		AddInitialState("s1").
		AddStates("s2", "s3").
		AddTransition("s1", "e1", func(event fsm.Event) (fsm.State, error) { return "s2", nil }, []fsm.State{"s2"}).
		AddTransition("s2", "e2", func(event fsm.Event) (fsm.State, error) { return "s3", nil }, []fsm.State{"s3"}).
		AddTransition("s3", "e3", func(event fsm.Event) (fsm.State, error) { return "s1", nil }, []fsm.State{"s1"}).
		Build()

	fmt.Println(fsm.CurrentState())
	fsm.Handle(&evt{t: "e1"})
	fmt.Println(fsm.CurrentState())
	fsm.Handle(&evt{t: "e2"})
	fmt.Println(fsm.CurrentState())
	fsm.Handle(&evt{t: "e3"})
	fmt.Println(fsm.CurrentState())
}
