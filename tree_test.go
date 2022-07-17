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
	"fmt"
	"log"
	"os"
	"testing"
)

func printStuff(tree *Tree) {
	var err error
	for nodeID, node := range tree.Nodes {
		var (
			leafIDs, immediateDescendantIDs, ancestorIDs []NodeID
			leaves, immediateDescendants, ancestors      []*Element
			parentID                                     NodeID
			parent                                       *Element
		)
		if leafIDs, err = tree.LeafDescendantIDs(NodeID(nodeID)); err != nil {
			log.Fatalf("error retrieving tree.LeafDescendantIDs(%d): %v", nodeID, err)
		}
		if leaves, err = tree.LeafDescendants(NodeID(nodeID)); err != nil {
			log.Fatalf("error retrieving tree.LeafDescendants(%d): %v", nodeID, err)
		}
		if immediateDescendantIDs, err = tree.ImmediateDescendantIDs(NodeID(nodeID)); err != nil {
			log.Fatalf("error retrieving tree.ImmediateDescendantIDs(%d): %v", nodeID, err)
		}
		if immediateDescendants, err = tree.ImmediateDescendants(NodeID(nodeID)); err != nil {
			log.Fatalf("error retrieving tree.ImmediateDescendants(%d): %v", nodeID, err)
		}
		if ancestorIDs, err = tree.AncestorIDs(NodeID(nodeID)); err != nil {
			log.Fatalf("error retrieving tree.AncestorIDs(%d): %v", nodeID, err)
		}
		if ancestors, err = tree.Ancestors(NodeID(nodeID)); err != nil {
			log.Fatalf("error retrieving tree.Ancestors(%d): %v", nodeID, err)
		}
		if parentID, err = tree.ParentID(NodeID(nodeID)); err != nil {
			log.Printf("error retrieving tree.ParentID(%d): %v", nodeID, err)
		}
		if parent, err = tree.Parent(NodeID(nodeID)); err != nil {
			log.Printf("error retrieving tree.Parent(%d): %v", nodeID, err)
		}
		fmt.Printf("- Node %d: %s\n\tParent: (%d) %s\n\tChildren: %v\n\tImmediateDescendantIDs: %v\n\tImmediateDescendants: %v\n\tLeaf IDs: %v\n\tLeaves: %v\n\tAncestorIDs: %v\n\tAncestors: %v\n\n",
			nodeID, node.Data,
			parentID, parent,
			node.Children,
			immediateDescendantIDs, immediateDescendants,
			leafIDs, leaves,
			ancestorIDs, ancestors,
		)
	}
}

func TestDesererializeTreeFromFile(t *testing.T) {
	var (
		tree           *Tree
		actitree__topo []byte
		err            error
	)

	//const IN_FILE_PATH = "test_artifacts/topo__immutree.json"
	const IN_FILE_PATH = "test_artifacts/t4_de.json"
	const OUT_FILE_PATH = "test_artifacts/go_t4_de_tree.json"

	if actitree__topo, err = os.ReadFile(IN_FILE_PATH); err != nil {
		t.Fatalf("Error reading from file until EOF: %v\n", err)
	}
	if err = json.Unmarshal(actitree__topo, &tree); err != nil {
		t.Fatalf("Error unmarshaling JSON: %v\n", err)
	}

	printStuff(tree)

	// Remarshal the tree
	remarshaled, err := json.Marshal(tree)
	if err != nil {
		t.Fatalf("Failed to remarshal tree: %v\n", err)
	}
	t.Logf("Remarshaled:\n%s", remarshaled)

	// Write remarshaled tree into a file
	if err = os.WriteFile(OUT_FILE_PATH, remarshaled, 0664); err != nil {
		t.Fatalf("Failed to write remarshaled tree into file %v: %v\n", OUT_FILE_PATH, err)
	}
}
