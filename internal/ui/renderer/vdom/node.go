package vdom

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"sync"
)

// VNode represents a virtual DOM node
type VNode struct {
	ID       string         // Unique component ID
	Type     string         // Component type
	Content  string         // Rendered content
	Hash     string         // Content hash for fast comparison
	Children []*VNode       // Child nodes
	Meta     map[string]any // Metadata (focus info, etc.)
}

// VDOM manages virtual DOM tree and reconciliation
type VDOM struct {
	root    *VNode
	rootMu  sync.RWMutex
	patches []Patch
}

// Patch represents a change to apply
type Patch struct {
	Type   PatchType
	NodeID string
	OldStr string
	NewStr string
	Index  int
}

type PatchType int

const (
	PatchReplace PatchType = iota
	PatchUpdate
	PatchInsert
	PatchRemove
)

// NewVDOM creates new virtual DOM
func NewVDOM() *VDOM {
	return &VDOM{
		patches: make([]Patch, 0),
	}
}

// BuildNode creates VNode from component render output
func BuildNode(id, nodeType, content string, children []*VNode, meta map[string]any) *VNode {
	hash := computeHash(content)
	return &VNode{
		ID:       id,
		Type:     nodeType,
		Content:  content,
		Hash:     hash,
		Children: children,
		Meta:     meta,
	}
}

// Reconcile compares new tree with old tree and generates patches
func (v *VDOM) Reconcile(newRoot *VNode) []Patch {
	v.rootMu.Lock()
	defer v.rootMu.Unlock()

	v.patches = make([]Patch, 0)

	if v.root == nil {
		// First render - full replace
		v.root = newRoot
		v.patches = append(v.patches, Patch{
			Type:   PatchReplace,
			NodeID: "root",
			NewStr: v.flattenTree(newRoot),
		})
	} else {
		// Diff trees
		v.diffNodes(v.root, newRoot, "root")
		v.root = newRoot
	}

	return v.patches
}

// diffNodes recursively diffs two nodes
func (v *VDOM) diffNodes(oldNode, newNode *VNode, path string) {
	if oldNode == nil && newNode == nil {
		return
	}

	// Node removed
	if oldNode != nil && newNode == nil {
		v.patches = append(v.patches, Patch{
			Type:   PatchRemove,
			NodeID: path,
			OldStr: oldNode.Content,
		})
		return
	}

	// Node added
	if oldNode == nil {
		v.patches = append(v.patches, Patch{
			Type:   PatchInsert,
			NodeID: path,
			NewStr: newNode.Content,
		})
		return
	}

	// Node type changed - replace entire subtree
	if oldNode.Type != newNode.Type {
		v.patches = append(v.patches, Patch{
			Type:   PatchReplace,
			NodeID: path,
			OldStr: v.flattenTree(oldNode),
			NewStr: v.flattenTree(newNode),
		})
		return
	}

	// Content changed - update
	if oldNode.Hash != newNode.Hash {
		v.patches = append(v.patches, Patch{
			Type:   PatchUpdate,
			NodeID: path,
			OldStr: oldNode.Content,
			NewStr: newNode.Content,
		})
	}

	// Diff children
	oldChildren := oldNode.Children
	newChildren := newNode.Children
	maxLen := maxInt(len(oldChildren), len(newChildren))

	for i := 0; i < maxLen; i++ {
		childPath := fmt.Sprintf("%s.child[%d]", path, i)

		var oldChild, newChild *VNode
		if i < len(oldChildren) {
			oldChild = oldChildren[i]
		}
		if i < len(newChildren) {
			newChild = newChildren[i]
		}

		v.diffNodes(oldChild, newChild, childPath)
	}
}

// flattenTree converts node tree to flat string
func (v *VDOM) flattenTree(node *VNode) string {
	if node == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(node.Content)

	for _, child := range node.Children {
		sb.WriteString(v.flattenTree(child))
	}

	return sb.String()
}

// ApplyPatches applies patches to output string
func ApplyPatches(original string, patches []Patch) string {
	result := original

	for _, patch := range patches {
		switch patch.Type {
		case PatchReplace:
			result = patch.NewStr

		case PatchUpdate:
			result = strings.ReplaceAll(result, patch.OldStr, patch.NewStr)

		case PatchInsert:
			result += patch.NewStr

		case PatchRemove:
			result = strings.ReplaceAll(result, patch.OldStr, "")
		}
	}

	return result
}

// GetRoot returns current root (thread-safe)
func (v *VDOM) GetRoot() *VNode {
	v.rootMu.RLock()
	defer v.rootMu.RUnlock()
	return v.root
}

// computeHash creates fast hash of content
func computeHash(content string) string {
	h := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", h[:8]) // First 8 bytes
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
