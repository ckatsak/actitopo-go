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

import "fmt"

// NodeID serves as a unique identifier of an Element in the Tree.
// It is also its index in the Tree.
type NodeID = uint32

// TreeNode contains an Element of the hardware topology along with the NodeIDs
// of its immediate descendants in the Tree.
type TreeNode struct {
	// Data is the hardware topology Element contained in this TreeNode.
	Data *Element `json:"data"`
	// Children is a list of the NodeIDs that correspond to the immediate
	// descendants (i.e., children) Elements of the Element contained in
	// this TreeNode.
	Children []NodeID `json:"desc,omitempty"`
}

// Tree represents the hierarchy of the hardware topology hierarchy in ActiK8s.
type Tree struct {
	// Nodes contains all TreeNode objects that constitute the Tree, and is
	// indexed by Elements' NodeIDs in the Tree.
	Nodes []TreeNode `json:"nodes"`
}

// Size returns the number of Elements currently stored in the Tree.
func (t *Tree) Size() int {
	if nil == t {
		return 0
	}
	return len(t.Nodes)
}

// IsEmpty returns true if there are no Elements currently stored in the Tree.
func (t *Tree) IsEmpty() bool {
	if nil == t {
		return true
	}
	return 0 == len(t.Nodes)
}

// Root returns the Element stored at the root of the hardware topology Tree,
// or a non-nil error value in case of failure.
func (t *Tree) Root() (*Element, error) {
	if nil == t {
		return nil, fmt.Errorf("Tree is nil")
	}
	if len(t.Nodes) == 0 {
		return nil, fmt.Errorf("Tree is empty")
	}
	return t.Nodes[0].Data, nil
}

// Get returns a reference to the Element that is stored in the Tree under
// the provided NodeID, if it exists, or a non-nil error value otherwise.
func (t *Tree) Get(id NodeID) (*Element, error) {
	if nil == t {
		return nil, fmt.Errorf("Tree is nil")
	}
	if int(id) >= len(t.Nodes) {
		return nil, fmt.Errorf("Invalid NodeID %d", id)
	}

	return t.Nodes[id].Data, nil
}

// ImmediateDescendantIDs returns a list of NodeIDs that correspond to the
// immediate descendant (i.e., children) elements of the element stored in the
// Tree under the provided NodeID.
func (t *Tree) ImmediateDescendantIDs(id NodeID) ([]NodeID, error) {
	if nil == t {
		return nil, fmt.Errorf("Tree is nil")
	}
	if int(id) >= len(t.Nodes) {
		return nil, fmt.Errorf("Invalid NodeID %d", id)
	}

	return t.Nodes[id].Children, nil
}

// ImmediateDescendants returns a list of the immediate descendant (i.e.,
// children) elements of the element stored in the Tree under the provided
// NodeID.
func (t *Tree) ImmediateDescendants(id NodeID) (children []*Element, err error) {
	if nil == t {
		return nil, fmt.Errorf("Tree is nil")
	}
	if int(id) >= len(t.Nodes) {
		return nil, fmt.Errorf("Invalid NodeID %d", id)
	}

	children = make([]*Element, 0, len(t.Nodes[id].Children))
	for i := range t.Nodes[id].Children {
		children = append(children, t.Nodes[t.Nodes[id].Children[i]].Data)
	}
	return
}

// LeafDescendantIDs returns a list of NodeIDs that correspond to the elements
// at the leaves of the Tree, which are also descendants of the element stored
// under the provided NodeID.
func (t *Tree) LeafDescendantIDs(id NodeID) (leafIDs []NodeID, err error) {
	if nil == t {
		return nil, fmt.Errorf("Tree is nil")
	}
	if int(id) >= len(t.Nodes) {
		return nil, fmt.Errorf("Invalid NodeID %d", id)
	}

	leafIDs = make([]NodeID, 0)
	stack := make([]*struct {
		id NodeID
		tn *TreeNode
	}, 0, len(t.Nodes)/2)
	stack = append(stack, &struct {
		id NodeID
		tn *TreeNode
	}{
		id: NodeID(id),
		tn: &t.Nodes[id],
	})
	for len(stack) > 0 {
		last := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if len(last.tn.Children) > 0 {
			for j := range last.tn.Children {
				stack = append(stack, &struct {
					id NodeID
					tn *TreeNode
				}{
					id: NodeID(last.tn.Children[j]),
					tn: &t.Nodes[last.tn.Children[j]],
				})
			}
		} else {
			leafIDs = append(leafIDs, last.id)
		}
	}
	return
}

// LeafDescendants returns a list of the elements at the leaves of the Tree,
// which are also descendants of the element stored under the provided NodeID.
func (t *Tree) LeafDescendants(id NodeID) (leaves []*Element, err error) {
	if nil == t {
		return nil, fmt.Errorf("Tree is nil")
	}
	if int(id) >= len(t.Nodes) {
		return nil, fmt.Errorf("Invalid NodeID %d", id)
	}

	leaves = make([]*Element, 0)
	stack := make([]*TreeNode, 0, len(t.Nodes)/2)
	stack = append(stack, &t.Nodes[id])
	for len(stack) > 0 {
		last := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if len(last.Children) > 0 {
			for child_id := range last.Children {
				stack = append(stack, &t.Nodes[last.Children[child_id]])
			}
		} else {
			leaves = append(leaves, last.Data)
		}
	}
	return
}

// ParentID returns the NodeID of the immediate ancestor (i.e., the parent)
// element of the element stored in the Tree under the provided NodeID, or a
// non-nil error value in case of failure.
//
// Querying for the parent of the root Element returns an error too.
func (t *Tree) ParentID(id NodeID) (NodeID, error) {
	if nil == t {
		return 0, fmt.Errorf("Tree is nil")
	}
	if int(id) >= len(t.Nodes) {
		return 0, fmt.Errorf("Invalid NodeID %d", id)
	}
	if id == 0 {
		return 0, fmt.Errorf("Root element does not have a parent")
	}

	for parentID := range t.Nodes {
		for j := range t.Nodes[parentID].Children {
			if id == t.Nodes[parentID].Children[j] {
				return NodeID(parentID), nil
			}
		}
	}
	panic("UNREACHABLE") // XXX(ckatsak)
}

// Parent returns the immediate ancestor (i.e., the parent) element of the
// element stored in the Tree under the provided NodeID, or a non-nil error
// value in case of failure.
//
// Querying for the parent of the root Element returns an error too.
func (t *Tree) Parent(id NodeID) (*Element, error) {
	if nil == t {
		return nil, fmt.Errorf("Tree is nil")
	}
	if int(id) >= len(t.Nodes) {
		return nil, fmt.Errorf("Invalid NodeID %d", id)
	}
	if id == 0 {
		return nil, fmt.Errorf("Root element does not have a parent")
	}

	for parentID := range t.Nodes {
		for j := range t.Nodes[parentID].Children {
			if id == t.Nodes[parentID].Children[j] {
				return t.Nodes[parentID].Data, nil
			}
		}
	}
	panic("UNREACHABLE") // XXX(ckatsak)
}

// AncestorIDs returns a list of NodeIDs that correspond to the ancestor (i.e.,
// parent) elements of the element stored in the Tree under the provided
// NodeID, all the way up to the root element of the Tree.
func (t *Tree) AncestorIDs(id NodeID) (ancestorIDs []NodeID, err error) {
	if nil == t {
		return nil, fmt.Errorf("Tree is nil")
	}
	if int(id) >= len(t.Nodes) {
		return nil, fmt.Errorf("Invalid NodeID %d", id)
	}

	allAncestorIDs := make([]NodeID, len(t.Nodes))
	for parentID := range t.Nodes {
		for j := range t.Nodes[parentID].Children {
			allAncestorIDs[t.Nodes[parentID].Children[j]] = NodeID(parentID)
		}
	}

	ancestorIDs = make([]NodeID, 0, len(t.Nodes)/2)
	for ; id != NodeID(0); id = allAncestorIDs[id] {
		ancestorIDs = append(ancestorIDs, allAncestorIDs[id])
	}
	return
}

// Ancestors returns a list of the ancestor (i.e., parent) elements of the
// element stored in the Tree under the provided NodeID, all the way up to the
// root element of the Tree.
func (t *Tree) Ancestors(id NodeID) (ancestors []*Element, err error) {
	if nil == t {
		return nil, fmt.Errorf("Tree is nil")
	}
	if int(id) >= len(t.Nodes) {
		return nil, fmt.Errorf("Invalid NodeID %d", id)
	}

	allAncestorIDs := make([]NodeID, len(t.Nodes))
	for parentID := range t.Nodes {
		for j := range t.Nodes[parentID].Children {
			allAncestorIDs[t.Nodes[parentID].Children[j]] = NodeID(parentID)
		}
	}

	ancestors = make([]*Element, 0, len(t.Nodes)/2)
	for ; id != NodeID(0); id = allAncestorIDs[id] {
		ancestors = append(ancestors, t.Nodes[allAncestorIDs[id]].Data)
	}
	return
}
