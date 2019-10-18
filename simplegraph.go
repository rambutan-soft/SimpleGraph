package simplegraph

import (
	"errors"
	"sync"
)

//Node ... type node defintion
type Node struct {
	Key       string
	Value     interface{}
	EdgesUp   map[*Node]interface{}
	EdgesDown map[*Node]interface{}
	sync.RWMutex
}

//SimpleGraph ... type SimpleGraph defintion
type SimpleGraph struct {
	Nodes map[string]*Node
	sync.RWMutex
}

//New ...Initializes a new graph.
func New() *SimpleGraph {
	return &SimpleGraph{map[string]*Node{}, sync.RWMutex{}}
}

//Add ... add a new node, if existing, overwrite it
func (g *SimpleGraph) Add(key string, data interface{}) {
	g.Lock()
	defer g.Unlock()

	if g.Nodes[key] != nil {
		g.Nodes[key].Value = data
	} else {
		g.Nodes[key] = &Node{key, data, make(map[*Node]interface{}), make(map[*Node]interface{}), sync.RWMutex{}}
	}
}

//Get ... get the node
func (g *SimpleGraph) Get(key string) (v *Node, err error) {
	g.RLock()
	v = g.Nodes[key]
	g.RUnlock()

	if v == nil {
		err = errors.New("invalid key")
	}

	return
}

//Connect ... allow recursive for now, will see
func (g *SimpleGraph) Connect(parentKey string, childKey string, edge interface{}) bool {
	g.RLock()
	defer g.RUnlock()

	val := g.Nodes[parentKey]
	childVal := g.Nodes[childKey]

	if val == nil || childVal == nil {
		return false
	}

	val.Lock()
	childVal.Lock()

	val.EdgesDown[childVal] = edge
	childVal.EdgesUp[val] = edge

	val.Unlock()
	childVal.Unlock()

	return true
}

//Disconnect ... remove connection
func (g *SimpleGraph) Disconnect(parentKey string, childKey string) bool {
	g.RLock()
	defer g.RUnlock()

	val := g.Nodes[parentKey]
	childVal := g.Nodes[childKey]

	if val == nil || childVal == nil {
		return false
	}

	val.Lock()
	childVal.Lock()

	delete(val.EdgesDown, childVal)
	delete(childVal.EdgesUp, val)

	val.Unlock()
	childVal.Unlock()

	return true

}

//Delete ... delete node
func (g *SimpleGraph) Delete(key string) bool {
	g.Lock()
	defer g.Unlock()

	v := g.Nodes[key]
	if v == nil {
		return false
	}

	// remove children relationships
	for e := range v.EdgesDown {
		e.Lock()
		delete(e.EdgesUp, v)
		e.Unlock()
	}

	// remove parents relationships
	for e := range v.EdgesUp {
		e.Lock()
		delete(e.EdgesDown, v)
		e.Unlock()
	}

	// delete node
	delete(g.Nodes, key)

	return true
}

//GetParents ...nice to have func
func (g *SimpleGraph) GetParents(key string) map[*Node]interface{} {
	g.RLock()
	defer g.RUnlock()

	v := g.Nodes[key]

	if v == nil {
		return nil
	}
	return v.EdgesUp
}

//GetChildren ...nice to have func
func (g *SimpleGraph) GetChildren(key string) map[*Node]interface{} {
	g.RLock()
	defer g.RUnlock()

	v := g.Nodes[key]

	if v == nil {
		return nil
	}
	return v.EdgesDown
}
