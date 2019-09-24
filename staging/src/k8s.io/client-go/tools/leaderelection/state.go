/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package leaderelection

import (
	"reflect"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/clock"
	rl "k8s.io/client-go/tools/leaderelection/resourcelock"
)

type leaderElectorState struct {
	mu             sync.RWMutex
	observedRecord rl.LeaderElectionRecord
	observedTime   time.Time

	// clock is wrapper around time to allow for less flaky testing
	clock clock.Clock
}

func (les *leaderElectorState) updateRecord(ler rl.LeaderElectionRecord) {
	les.mu.Lock()
	defer les.mu.Unlock()
	les.observedRecord = ler
	les.observedTime = les.clock.Now()
}

func (les *leaderElectorState) updateRecordIfDifferent(ler rl.LeaderElectionRecord) {
	les.mu.Lock()
	defer les.mu.Unlock()
	if !reflect.DeepEqual(les.observedRecord, ler) {
		les.observedRecord = ler
		les.observedTime = les.clock.Now()
	}
}

func (les *leaderElectorState) currentTime() time.Time {
	les.mu.RLock()
	defer les.mu.RUnlock()
	return les.observedTime
}

func (les *leaderElectorState) elapsedTime() time.Duration {
	les.mu.RLock()
	defer les.mu.RUnlock()
	return les.clock.Since(les.observedTime)
}

func (les *leaderElectorState) currentRecord() rl.LeaderElectionRecord {
	les.mu.RLock()
	defer les.mu.RUnlock()
	return les.observedRecord
}
