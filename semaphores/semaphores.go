package main

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
	speed    = 1
	wg       sync.WaitGroup
)

func main() {
	// lowTraffic()
	mediumTraffic()
	// highTraffic()
	initJsSimulation()

	swimcaps()
	// lockers()
}

func lowTraffic() {
	swimmers = 2
	capacity = 1
	// arrivalPeriod = 10000;
	// minSwimDuration = 3000;
	// maxSwimDuration = 10000;
	speed = 2
}

func mediumTraffic() {
	swimmers = 20
	capacity = 10
	speed = 2
}

func highTraffic() {
	swimmers = 120
	capacity = 20
	// arrivalPeriod = 10000;
	// minSwimDuration = 3000;
	// maxSwimDuration = 10000;
	speed = 2
}

func initJsSimulation() {
	js.Global.Set("N", swimmers)
	js.Global.Set("capacity", capacity)
	js.Global.Set("speed", speed)
	js.Global.Get("makeBasketCaps").Invoke(capacity)
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
	js.Global.Get("arrive").Invoke(s)
	sleep(time.Second)
}

func (s Swimmer) swim() {
	sleep(300 * time.Millisecond) // Delay where the cap is still in the basket
	durationMs := 2000 + rand.Intn(6000)
	fmt.Println(s, "will swim for", durationMs/speed)
	js.Global.Get("swim").Invoke(s, durationMs)
	backDurationMs := 3000
	sleep(time.Duration(durationMs+backDurationMs) * time.Millisecond)
}

func swimcaps() {
	type Cap struct{}
	caps := make(chan Cap, capacity)
	for i := 0; i < capacity; i++ {
		caps <- Cap{}
	}

	for i := 0; i < swimmers; i++ {
		s := Swimmer(i)
		go func() {
			s.arrive()
			cap := <-caps
			fmt.Println(s, "took a cap")
			s.swim()
			caps <- cap
			fmt.Println(s, "gave a cap")
			leave()
		}()
	}

	wg.Wait()
	fmt.Println("All swimmers have left!")
}

func lockers() {
	type GymBag struct{}
	lockers := make(chan GymBag, capacity)

	for i := 0; i < swimmers; i++ {
		i := i
		go func() {
			arrive()
			fmt.Println(i, "\tarrives")
			lockers <- GymBag{}
			fmt.Println(i, "\tgets in")
			swim()
			fmt.Println(i, "\tgets out")
			<-lockers
			wg.Done()
		}()
	}

	wg.Wait()
}

func arrive() {
	sleep(time.Duration(rand.Intn(400)) * time.Millisecond)
}

func swim() {
	sleep(time.Duration(rand.Intn(400)) * time.Millisecond)
}

func leave() {
	fmt.Println("Leaving!")
	wg.Done()
}
