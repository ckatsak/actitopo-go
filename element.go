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

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
////
////	Element
////
///////////////////////////////////////////////////////////////////////////////

// Element represents a node in the hierarchy of the hardware topology.
//
// Apart from the special case of Machine, which is the root node in the
// hierarchy, an Element can be either a Processing node or a Cache.
type Element struct {
	// Processing is non-nil if the Element represents a computation unit
	// in the hierarchical hardware topology.
	*Processing `json:"processing,omitempty"`
	// Cache is non-nil if the Element represents a caching element in the
	// hierarchical hardware topology.
	*Cache `json:"cache,omitempty"`
}

// IsRoot returns true if the Element is the root node in the hierarchy (i.e.,
// the Machine) and false otherwise.
func (e *Element) IsRoot() bool {
	return nil == e.Processing && nil == e.Cache
}

// IsProcessing returns true if the Element is a Processing node and false
// otherwise.
func (e *Element) IsProcessing() bool {
	return nil == e.Cache && nil != e.Processing
}

// IsCache returns true if the Element is a Cache and false otherwise.
func (e *Element) IsCache() bool {
	return nil == e.Processing && nil != e.Cache
}

// String returns the string representation of the Element.
func (e *Element) String() string {
	switch {
	case e.IsRoot():
		return "Machine"
	case e.IsCache():
		return fmt.Sprintf("%s", e.Cache)
	case e.IsProcessing():
		return fmt.Sprintf("%s", e.Processing)
	default:
		panic("UNREACHABLE") // XXX(ckatsak)
	}
}

// MarshalJSON returns the Element marshalled in JSON, or a non-nil error value
// in case of failure.
func (e *Element) MarshalJSON() ([]byte, error) {
	raw := make(map[string]interface{})
	switch {
	case e.IsRoot():
		return []byte(`"machine"`), nil
	case e.IsCache():
		raw["cache"] = e.Cache
		return json.Marshal(raw)
	case e.IsProcessing():
		raw["processing"] = e.Processing
		return json.Marshal(raw)
	default:
		return nil, fmt.Errorf("Invalid Element")
	}
}

// UnmarshalJSON attempts to unmarshal the Element from the provided byte slice
// and returns a non-nil error if it fails.
func (e *Element) UnmarshalJSON(data []byte) (err error) {
	// If it's the root element (i.e., "machine"), get on with it
	if bytes.HasPrefix(bytes.ToLower(data), []byte(`"machine"`)) {
		e.Processing = nil
		e.Cache = nil
		return nil
	}

	// Unmarshal into a map[..]..
	var raw interface{}
	if err = json.Unmarshal(data, &raw); err != nil {
		return
	}

	// Make sure the root is map[string]..
	root, rootOk := raw.(map[string]interface{})
	if !rootOk {
		return fmt.Errorf("failed to unmarshal Element")
	}

	if content, contentOk := root["processing"]; contentOk {
		// If it is a Processing element:
		e.Cache = nil
		processing, processingOk := content.(map[string]interface{})
		if !processingOk {
			return fmt.Errorf("failed to unmarshal Processing")
		}
		kindStr, kindOk := processing["kind"].(string)
		idF64, idOk := processing["id"].(float64)
		if kindOk && idOk {
			var kind ProcessingKind
			if kind, err = ParseProcessingKind(kindStr); err != nil {
				return fmt.Errorf("failed to unmarshal Processing: failed to unmarshal ProcessingKind: %v", err)
			}
			e.Processing = &Processing{
				Kind: kind,
				ID:   uint32(idF64),
			}
		} else {
			err = fmt.Errorf("failed to unmarshal Processing")
		}
	} else if content, contentOk := root["cache"]; contentOk {
		// If it is a Cache element:
		e.Processing = nil
		cache, cacheOk := content.(map[string]interface{})
		if !cacheOk {
			return fmt.Errorf("failed to unmarshal Cache")
		}
		levelStr, levelOk := cache["lvl"].(string)
		liF64, liOk := cache["li"].(float64)
		attrsVal, attrsOk := cache["attrs"].(map[string]interface{})
		if !attrsOk {
			return fmt.Errorf("failed to unmarshal Cache")
		}
		sizeF64, sizeOk := attrsVal["size"].(float64)
		lineF64, lineOk := attrsVal["line"].(float64)
		waysF64, waysOk := attrsVal["ways"].(float64)
		if levelOk && liOk && sizeOk && lineOk && waysOk {
			var cacheLevel CacheLevel
			if cacheLevel, err = ParseCacheLevel(levelStr); err != nil {
				return fmt.Errorf("failed to unmarshal Cache: failed to unmarshal CacheLevel: %v", err)
			}
			e.Cache = &Cache{
				Level:        cacheLevel,
				LogicalIndex: uint32(liF64),
				Attributes: &CacheAttributes{
					Size:          uint64(sizeF64),
					Linesize:      uint32(lineF64),
					Associativity: int32(waysF64),
				},
			}
		} else {
			err = fmt.Errorf("failed to unmarshal Cache")
		}
	} else {
		err = fmt.Errorf("failed to unmarshal Element")
	}
	return
}

///////////////////////////////////////////////////////////////////////////////
////
////	Processing
////
///////////////////////////////////////////////////////////////////////////////

// Processing represents a computation unit in the hardware topology.
type Processing struct {
	// Kind indicates the type of computation node that this Processing
	// node is.
	Kind ProcessingKind `json:"kind"`
	// ID is the index of the computation unit, assigned by the operating
	// system.
	ID uint32 `json:"id"`
}

// String returns the string representation of the Processing.
func (p *Processing) String() string {
	return fmt.Sprintf("%s(%d)", p.Kind, p.ID)
}

// ProcessingKind enumerates all types of computation units that can be used by
// by this package.
type ProcessingKind byte

const (
	// UnknownProcessingKind is employed to represent any unknown
	// ProcessingKind found in the wild.
	//
	// You should never see this; something is fucked up if you do.
	UnknownProcessingKind ProcessingKind = iota
	// Package represents a physical package (i.e., what goes into a
	// physical socket).
	Package
	// NUMAnode represents a NUMA node (i.e., a set of processors around
	// memory which all processors can directly access via the same
	// physical link).
	NUMANode
	// Core represents a physical core.
	Core
	// Thread represents a logical core (i.e., hardware thread, possibly
	// sharing a physical core with other hardware threads).
	Thread
)

// String returns the string representation of the ProcessingKind.
func (pk ProcessingKind) String() string {
	switch pk {
	case Package:
		return "Package"
	case NUMANode:
		return "NUMANode"
	case Core:
		return "Core"
	case Thread:
		return "Thread"
	default:
		return "UnknownProcessingKind"
	}
}

// ParseProcessingKind returns a ProcessingKind parsed from the provided string
// representation, or a non-nil error value if parsing fails.
func ParseProcessingKind(str string) (ProcessingKind, error) {
	switch strings.ToLower(str) {
	case "package":
		return Package, nil
	case "numa_node", "numanode":
		return NUMANode, nil
	case "core":
		return Core, nil
	case "thread":
		return Thread, nil
	default:
		return UnknownProcessingKind, fmt.Errorf("unknown processing kind: '%s'", str)
	}
}

// MarshalJSON returns the ProcessingKind marshalled in JSON, or a non-nil
// error value in case of failure.
func (pk ProcessingKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.ToLower(pk.String()))
}

// UnmarshalJSON attempts to unmarshal the ProcessingKind from the provided
// byte slice and returns a non-nil error if it fails.
func (pk *ProcessingKind) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, pk)
}

///////////////////////////////////////////////////////////////////////////////
////
////	Cache
////
///////////////////////////////////////////////////////////////////////////////

// Cache represents a data caching element (e.g., L3 cache, etc).
type Cache struct {
	// Level is the level of the cache.
	Level CacheLevel `json:"lvl"`
	// LogicalIndex is a logical index assigned by libhwloc.
	LogicalIndex uint32 `json:"li"`
	// Attributes contains any characteristics of the cache that may have
	// been detected.
	Attributes *CacheAttributes `json:"attrs"`
}

// String returns the string representation of the Cache.
func (c *Cache) String() string {
	return fmt.Sprintf("Cache{ %s(L#%d), attrs: %s }", c.Level, c.LogicalIndex, c.Attributes)
}

// CacheLevel represents the level of the cache (e.g., L1, L2, etc).
type CacheLevel byte

const (
	// UnknownCacheLevel is employed to represent any unknown CacheLevel
	// found in the wild.
	//
	// You should never see this; something is fucked up if you do.
	UnknownCacheLevel CacheLevel = iota
	// L1 cache.
	L1
	// L2 cache.
	L2
	// L3 cache.
	L3
	// L4 cache.
	L4
	// L5 cache.
	L5
)

// String returns the string representation of the CacheLevel.
func (cl CacheLevel) String() string {
	switch cl {
	case L1:
		return "L1"
	case L2:
		return "L2"
	case L3:
		return "L3"
	case L4:
		return "L4"
	case L5:
		return "L5"
	default:
		return fmt.Sprintf("Unknown cache level %d", cl)
	}
}

// ParseCacheLevel returns a CacheLevel parsed from the provided string
// representation, or a non-nil error value if parsing fails.
func ParseCacheLevel(level string) (CacheLevel, error) {
	switch level {
	case "L1":
		return L1, nil
	case "L2":
		return L2, nil
	case "L3":
		return L3, nil
	case "L4":
		return L4, nil
	case "L5":
		return L5, nil
	default:
		return UnknownCacheLevel, fmt.Errorf("Unknown cache level '%s'", level)
	}
}

// MarshalJSON returns the CacheLevel marshalled in JSON, or a non-nil error
// value in case of failure.
func (cl CacheLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(cl.String())
}

// UnmarshalJSON attempts to unmarshal the CacheLevel from the provided byte
// slice and returns a non-nil error if it fails.
func (cl *CacheLevel) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, cl)
}

// CacheAttributes represents various characteristics of the cache that may
// have been detected.
type CacheAttributes struct {
	// Size is the total size of the Cache, in bytes.
	Size uint64 `json:"size"`
	// Linesize is the size of the cache line, in bytes.
	Linesize uint32 `json:"line"`
	// Associativity is the associativity of the cache, in # ways.
	Associativity int32 `json:"ways"`
}

// String returns the string representation of the CacheAttributes.
func (ca *CacheAttributes) String() string {
	return fmt.Sprintf("%dB/%dB/%d-way", ca.Size, ca.Linesize, ca.Associativity)
}
