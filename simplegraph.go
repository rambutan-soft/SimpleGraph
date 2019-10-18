package simplegraph

import (
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

// func (g *SimpleGraph) get(key string) *Node {
// 	return g.Nodes[key]
// }

//Connect ... allow recursive for now, will see
func (g *SimpleGraph) Connect(parentKey string, childKey string, edge interface{}) bool {
	g.RLock()
	defer g.RUnlock()

	// get vertexes and check for validity of keys
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

	// success
	return true
}

/*
func (g *SimpleGraph) Disconnect(key string, otherKey string) bool {
	// recursive edges are forbidden
	if key == otherKey {
		return false
	}

	// lock graph for reading until this method is finished to prevent changes made by other goroutines while this one is running
	g.RLock()
	defer g.RUnlock()

	// get vertexes and check for validity of keys
	v := g.get(key)
	otherV := g.get(otherKey)

	if v == nil || otherV == nil {
		return false
	}

	// delete the edge from both vertexes
	v.Lock()
	otherV.Lock()

	delete(v.Edges, otherV)
	delete(otherV.Edges, v)

	v.Unlock()
	otherV.Unlock()

	return true
}

func (v *Node) GetNeighbors() map[*Node]int {
	if v == nil {
		return nil
	}

	v.RLock()
	neighbors := v.Edges
	v.RUnlock()

	return neighbors
}

func (g *SimpleGraph) Delete(key string) bool {
	// lock graph until this method is finished to prevent changes made by other goroutines while this one is looping etc.
	g.Lock()
	defer g.Unlock()

	// get vertex in question
	v := g.get(key)
	if v == nil {
		return false
	}

	// iterate over neighbors, remove edges from neighboring vertexes
	for neighbor, _ := range v.Edges {
		// delete edge to the to-be-deleted vertex
		neighbor.Lock()
		delete(neighbor.Edges, v)
		neighbor.Unlock()
	}

	// delete vertex
	delete(g.Nodes, key)

	return true
}

func (g *SimpleGraph) Get(key string) (v *Node, err error) {
	g.RLock()
	v = g.get(key)
	g.RUnlock()

	if v == nil {
		err = errors.New("graph: invalid key")
	}

	return
}
*/

// func main() {
// 	var item Student
// 	doSomethinWithThisParam(&item)
// 	fmt.Printf("%+v", item)
// }

// type Student struct {
// 	ID   string
// 	Name string
// }

// func doSomethinWithThisParam(item interface{}) {
// 	switch v := item.(type) {
// 	case *Student:
// 		*v = Student{
// 			ID:   "124",
// 			Name: "Iman Tumorang",
// 		}
// 		// another case
// 	}
// }
