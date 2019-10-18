package simplegraph

import (
	"testing"
)

func TestGraph(t *testing.T) {
	g := New()
	g.Add("a", 1)
	g.Add("b", 2)
	g.Connect("a", "b", "ok")
	if g.Nodes["a"].EdgesDown[g.Nodes["b"]].(string) != "ok" &&
		g.Nodes["b"].EdgesUp[g.Nodes["a"]].(string) != "ok" && len(g.Nodes) != 2 {
		t.Error("Expected connected, but not")
	}

	n, _ := g.Get("a")
	if n.Value.(int) != 1 {
		t.Errorf("Expected 1, but got %d", n.Value.(int))
	}
	c := g.GetParents("b")
	if len(c) != 1 {
		t.Error("get parent wrong")
	}

	f := g.GetChildren("a")
	if len(f) != 1 {
		t.Error("get children wrong")
	}

	g.Disconnect("a", "b")
	if g.Nodes["a"].EdgesDown[g.Nodes["b"]] != nil &&
		g.Nodes["b"].EdgesUp[g.Nodes["a"]] != nil && len(g.Nodes) != 2 {
		t.Error("Expected connected, but not")
	}

	g.Delete("b")
	_, error := g.Get("b")
	if error == nil {
		t.Error("delete wrong")
	}
}
