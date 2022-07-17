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

// Package to easily serialize, deserialize and work with the hierarchical
// hardware topology of a physical machine, as examined by the ActiK8s project.
//
// The package does not provide any synchronization guarantees, since the
// entities it exposes (i.e., a physical machine's hierarchical hardware
// topology), after their deserialization, are expected to only be read in the
// vast majority of the package's use cases.
// Therefore, concurrency should be catered for externally in cases heavy on
// concurrent modifications.
package actitopo
