package simplegraph

import (
	"errors"
	"sync"
)

//Node ... type node defintion
type Node struct {
	Key   string
	Value interface{}
	Heads map[*Node]interface{}
	Tails map[*Node]interface{}
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

	val.Tails[childVal] = edge
	childVal.Heads[val] = edge

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

	delete(val.Tails, childVal)
	delete(childVal.Heads, val)

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
	for e := range v.Tails {
		e.Lock()
		delete(e.Heads, v)
		e.Unlock()
	}

	// remove parents relationships
	for e := range v.Heads {
		e.Lock()
		delete(e.Tails, v)
		e.Unlock()
	}

	// delete node
	delete(g.Nodes, key)

	return true
}

//GetEdge ...
func (g *SimpleGraph) GetEdge(headKey, tailKey string) interface{} {
	g.RLock()
	defer g.RUnlock()

	head := g.Nodes[headKey]

	if head == nil {
		return nil
	}

	tail := g.Nodes[tailKey]

	if tail == nil {
		return nil
	}
	var edge interface{}

	//hope this goof for performance
	if len(head.Tails) > len(tail.Heads) {
		edge = tail.Heads[head]
	} else {
		edge = head.Tails[tail]
	}
	if edge == nil {
		return nil
	}
	return edge
}

//GetParents ...nice to have func
func (g *SimpleGraph) GetParents(key string) map[*Node]interface{} {
	g.RLock()
	defer g.RUnlock()

	v := g.Nodes[key]

	if v == nil {
		return nil
	}
	return v.Heads
}

//GetChildren ...nice to have func
func (g *SimpleGraph) GetChildren(key string) map[*Node]interface{} {
	g.RLock()
	defer g.RUnlock()

	v := g.Nodes[key]

	if v == nil {
		return nil
	}
	return v.Tails
}

//TailConnected ...nice to have func
func (g *SimpleGraph) TailConnected(headKey, tailKey string) bool {
	if headKey == tailKey {
		return true
	}
	g.RLock()
	defer g.RUnlock()

	head := g.Nodes[headKey]
	if head == nil || len(head.Tails) == 0 {
		return false
	}

	tail := g.Nodes[tailKey]
	if tail == nil || len(tail.Heads) == 0 {
		return false
	}
	head.RLock()
	defer head.RUnlock()
	tail.RLock()
	defer tail.RUnlock()

	_, ok := head.Tails[tail]
	if ok {
		return true
	} else {
		for t := range head.Tails {
			if g.TailConnected(t.Key, tailKey) {
				return true
			}
		}
	}

	return false
}

//HeadConnected ...nice to have func
func (g *SimpleGraph) HeadConnected(tailKey, headKey string) bool {
	if headKey == tailKey {
		return true
	}
	g.RLock()
	defer g.RUnlock()

	tail := g.Nodes[tailKey]
	if tail == nil || len(tail.Heads) == 0 {
		return false
	}

	head := g.Nodes[headKey]
	if head == nil || len(head.Tails) == 0 {
		return false
	}

	tail.RLock()
	defer tail.RUnlock()
	head.RLock()
	defer head.RUnlock()

	_, ok := tail.Heads[head]
	if ok {
		return true
	} else {
		for t := range tail.Heads {
			if g.HeadConnected(tailKey, t.Key) {
				return true
			}
		}
	}

	return false
}


//Size ...Size
func (g *SimpleGraph) Size() int {
	g.RLock()
	defer g.RUnlock()
	return len(g.Nodes)
}
