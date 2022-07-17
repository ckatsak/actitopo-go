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
	"encoding/json"
	"os"
	"testing"
)

func TestDesererializeTopoFromFile(t *testing.T) {
	const IN_FILE_PATH = "test_artifacts/t4_de.json"
	const OUT_FILE_PATH = "test_artifacts/go_t4_de_topo.json"

	var (
		topo           Topology
		actitree__topo []byte
		err            error
	)

	if actitree__topo, err = os.ReadFile(IN_FILE_PATH); err != nil {
		t.Fatalf("Error reading from file until EOF: %v\n", err)
	}
	if err = json.Unmarshal(actitree__topo, &topo); err != nil {
		t.Fatalf("Error unmarshaling JSON: %v\n", err)
	}
	t.Logf("Unmarshaled topology:\n%v", topo)

	// Remarshal the Topology
	remarshaled, err := json.Marshal(topo)
	if err != nil {
		t.Fatalf("Failed to remarshal the Topology: %v\n", err)
	}
	t.Logf("Remarshaled:\n%s", remarshaled)

	// Write remarshaled tree into a file
	if err = os.WriteFile(OUT_FILE_PATH, remarshaled, 0664); err != nil {
		t.Fatalf("Failed to write remarshaled Topology into file %v: %v\n", OUT_FILE_PATH, err)
	}
}
