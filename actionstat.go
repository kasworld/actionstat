// Copyright 2015 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// stat package
// total and lap support
package actionstat

import (
	"fmt"
	"sync"
	"time"
)

type element struct {
	Count int64
	Time  time.Time
}

type ActionStat struct {
	Total element
	Laps  []element
	mutex sync.Mutex
}

func intDurToStr(count int64, dur time.Duration) string {
	return fmt.Sprintf("%v/%5.1f/s", count, float64(count)/dur.Seconds())
}

func (a ActionStat) String() string {
	//defer a.UpdateLap()
	dur := time.Now().Sub(a.Total.Time)
	lastlap := a.Laps[len(a.Laps)-1]
	lapcount, lapdur := a.Total.Count-lastlap.Count, time.Now().Sub(lastlap.Time)
	//a.Total.Sub(a.Laps[len(a.Laps)-1])
	return fmt.Sprintf("total:%v lap:%v",
		intDurToStr(a.Total.Count, dur),
		intDurToStr(lapcount, lapdur))
}

func New() *ActionStat {
	r := &ActionStat{
		Total: element{0, time.Now()},
		Laps:  make([]element, 1),
	}
	r.UpdateLap()
	return r
}

func (a *ActionStat) PerSec(min, max float64) float64 {
	dur := time.Now().Sub(a.Total.Time)
	perSec := float64(a.Total.Count) / dur.Seconds()
	if perSec > max {
		perSec = max
	}
	if perSec < min {
		perSec = min
	}
	return perSec
}

func (a *ActionStat) LapPerSec(min, max float64) float64 {
	lastlap := a.Laps[len(a.Laps)-1]
	lapcount, lapdur := a.Total.Count-lastlap.Count, time.Now().Sub(lastlap.Time)
	lapPerSec := float64(lapcount) / lapdur.Seconds()
	if lapPerSec > max {
		lapPerSec = max
	}
	if lapPerSec < min {
		lapPerSec = min
	}
	return lapPerSec
}

func (a *ActionStat) NewLap() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.Laps = append(a.Laps, element{a.Total.Count, time.Now()})
}

func (a *ActionStat) UpdateLap() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.Laps[len(a.Laps)-1] = element{a.Total.Count, time.Now()}
}

func (a *ActionStat) Inc() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.Total.Count++
}
func (a *ActionStat) Add(n int64) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.Total.Count += n
}

func (a *ActionStat) GetCount() int64 {
	return a.Total.Count
}
