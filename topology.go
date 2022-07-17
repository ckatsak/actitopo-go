/*
  Copyright 2022 Christos Katsakioris

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

package actitopo

import "encoding/json"

// Topology represents the hierarchical hardware topology of a physical node
// for the purposes of the ActiK8s project.
//
// This type is a simple transparent wrapper on a Tree object, providing a few
// more convenience methods.
type Topology struct {
	*Tree
}

// Packages returns a list of all NodeIDs that correspond to a CPU Package
// processing element in the hierarchical hardware topology.
func (t *Topology) Packages() []NodeID {
	return t.getAllProcessingKind(Package)
}

// NUMANodes returns a list of all NodeIDs that correspond to a NUMA node
// processing element in the hierarchical hardware topology.
func (t *Topology) NUMANodes() []NodeID {
	return t.getAllProcessingKind(NUMANode)
}

// Cores returns a list of all NodeIDs that correspond to a physical core
// processing element in the hierarchical hardware topology.
func (t *Topology) Cores() []NodeID {
	return t.getAllProcessingKind(Core)
}

// Threads returns a list of all NodeIDs that correspond to a hardware thread
// processing element in the hierarchical hardware topology.
func (t *Topology) Threads() []NodeID {
	return t.getAllProcessingKind(Thread)
}

// getAllProcessingKind returns a list of all NodeIDs that correspond to a
// processing element of the provided kind in the hierarchical hardware
// topology.
func (t *Topology) getAllProcessingKind(kind ProcessingKind) []NodeID {
	ret := make([]NodeID, 0)
	for id := range t.Nodes {
		if t.Nodes[id].Data.IsProcessing() && t.Nodes[id].Data.Kind == kind {
			ret = append(ret, NodeID(id))
		}
	}
	return ret
}

// L1Caches returns a list of all NodeIDs that correspond to a L1 cache element
// in the hierarchical hardware topology.
func (t *Topology) L1Caches() []NodeID {
	return t.getAllCacheLevel(L1)
}

// L2Caches returns a list of all NodeIDs that correspond to a L2 cache element
// in the hierarchical hardware topology.
func (t *Topology) L2Caches() []NodeID {
	return t.getAllCacheLevel(L2)
}

// L3Caches returns a list of all NodeIDs that correspond to a L3 cache element
// in the hierarchical hardware topology.
func (t *Topology) L3Caches() []NodeID {
	return t.getAllCacheLevel(L3)
}

// L4Caches returns a list of all NodeIDs that correspond to a L4 cache element
// in the hierarchical hardware topology.
func (t *Topology) L4Caches() []NodeID {
	return t.getAllCacheLevel(L4)
}

// L5Caches returns a list of all NodeIDs that correspond to a L5 cache element
// in the hierarchical hardware topology.
func (t *Topology) L5Caches() []NodeID {
	return t.getAllCacheLevel(L5)
}

// getAllCacheLevel returns a list of all NodeIDs that correspond to a cache
// element of the provided cache level in the hierarchical hardware topology.
func (t *Topology) getAllCacheLevel(level CacheLevel) []NodeID {
	ret := make([]NodeID, 0)
	for id := range t.Nodes {
		if t.Nodes[id].Data.IsCache() && t.Nodes[id].Data.Level == level {
			ret = append(ret, NodeID(id))
		}
	}
	return ret
}

// MarshalJSON returns the Topology marshalled in JSON, or a non-nil error
// value in case of failure.
func (t *Topology) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Tree)
}

// UnmarshalJSON attempts to unmarshal the Topology from the provided byte
// slice and returns a non-nil error if it fails.
func (t *Topology) UnmarshalJSON(data []byte) (err error) {
	return json.Unmarshal(data, &t.Tree)
}
