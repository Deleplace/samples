// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/gopherjs/gopherjs/js"
)

var (
	capacity = 20
	swimmers = 20
	speed    = 1.0
	metaphor = "swimcaps"
	launch   = make(chan func())
	wg       sync.WaitGroup
)

func main() {
	js.Global.Set("launchSimulation", launchSimulation)
	simulation := <-launch
	simulation()
}

// This will be exposed in JS, triggerable by the web page
func launchSimulation(Capacity, Swimmers int, Speed float64, Metaphor string) {
	capacity = Capacity
	swimmers = Swimmers
	speed = Speed
	metaphor = Metaphor

	initJsSimulation()

	switch Metaphor {
	case "swimcaps":
		launch <- swimcaps
	case "gymbags":
		launch <- gymbags
	case "nosync":
		launch <- nosync
	default:
		panic(fmt.Sprintf("Unexpected metaphor %q", Metaphor))
	}
}

func initJsSimulation() {
	js.Global.Set("N", swimmers)
	js.Global.Set("capacity", capacity)
	js.Global.Set("speed", speed)
	js.Global.Set("metaphor", metaphor)
	if metaphor == "swimcaps" {
		js.Global.Get("makeBasketCaps").Invoke(capacity)
	}
	if metaphor == "gymbags" {
		js.Global.Get("makeGymbagsShelf").Invoke(capacity)
	}
	wg = sync.WaitGroup{}
	wg.Add(swimmers)
}

func sleep(d time.Duration) {
	time.Sleep(d / time.Duration(speed))
}

type Swimmer int

func (s Swimmer) arrive() {
	arrivalDateMs := rand.Intn(15000)
	sleep(time.Duration(arrivalDateMs) * time.Millisecond)
	fmt.Println(s, "arrives")
	js.Global.Get("arrive").Invoke(s)
	sleep(time.Second)
}

func (s Swimmer) swim() {
	sleep(300 * time.Millisecond) // Delay where the cap is still in the basket
	durationMs := 2000 + rand.Intn(6000)
	fmt.Println(s, "will swim for", float64(durationMs)/speed)
	js.Global.Get("swim").Invoke(s, durationMs)
	backDurationMs := 3000
	sleep(time.Duration(durationMs+backDurationMs) * time.Millisecond)
}

func (s Swimmer) leave() {
	fmt.Println("Leaving!")
	wg.Done()
}

func nosync() {
	for i := 0; i < swimmers; i++ {
		s := Swimmer(i)
		go func() {
			s.arrive()
			s.swim()
			s.leave()
		}()
	}

	wg.Wait()
	fmt.Println("All swimmers have left!")
}

func swimcaps() {
	type SwimCap struct{}
	caps := make(chan SwimCap, capacity)
	for i := 0; i < capacity; i++ {
		caps <- SwimCap{}
	}

	for i := 0; i < swimmers; i++ {
		s := Swimmer(i)
		go func() {
			s.arrive()
			swimcap := <-caps
			fmt.Println(s, "took a cap")
			s.swim()
			caps <- swimcap
			fmt.Println(s, "gave a cap")
			s.leave()
		}()
	}

	wg.Wait()
	fmt.Println("All swimmers have left!")
}

func gymbags() {
	type GymBag struct{}
	shelf := make(chan GymBag, capacity)

	for i := 0; i < swimmers; i++ {
		s := Swimmer(i)
		go func() {
			s.arrive()
			shelf <- GymBag{}
			fmt.Println(s, "gave a gym bag")
			s.swim()
			<-shelf
			fmt.Println(s, "took a gym bag")
			s.leave()
		}()
	}

	wg.Wait()
}
