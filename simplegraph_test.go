package simplegraph

import (
	"testing"
)

func TestGraph(t *testing.T) {
	g := New()
	g.Add("a", 1)
	g.Add("b", 2)
	g.Connect("a", "b", "ok")
	if g.Nodes["a"].Tails[g.Nodes["b"]].(string) != "ok" &&
		g.Nodes["b"].Heads[g.Nodes["a"]].(string) != "ok" && len(g.Nodes) != 2 {
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
	if g.Nodes["a"].Tails[g.Nodes["b"]] != nil &&
		g.Nodes["b"].Heads[g.Nodes["a"]] != nil && len(g.Nodes) != 2 {
		t.Error("Expected connected, but not")
	}

	g.Delete("b")
	_, error := g.Get("b")
	if error == nil {
		t.Error("delete wrong")
	}
}

func TestGraph2(t *testing.T) {
	//        a
	//     b      c
	//   c   d   e   f

	g := New()
	g.Add("a", 1)
	g.Add("b", 2)
	g.Add("c", 2)
	g.Add("d", 2)
	g.Add("e", 2)
	g.Add("f", 2)
	g.Connect("a", "b", "1")
	g.Connect("a", "c", "2")
	g.Connect("b", "c", "3")
	g.Connect("b", "d", "4")
	g.Connect("c", "e", "5")
	g.Connect("c", "f", "6")

	if g.GetEdge("a", "c").(string) != "2" {
		t.Errorf("a->c edge is not 2, is %s\n", g.GetEdge("a", "c").(string))
	}
	if g.GetEdge("a", "f") != nil {
		t.Errorf("a->f edge is not nil, is %s\n", g.GetEdge("a", "f").(string))
	}

	if !g.TailConnected("a", "c") || !g.TailConnected("a", "f") {
		t.Error("a tail nodes not right")
	}

	if g.TailConnected("a", "g") || g.TailConnected("c", "a") {
		t.Error("a tail nodes not right")
	}

	if g.HeadConnected("a", "c") || g.HeadConnected("a", "f") {
		t.Error("a Head nodes not right")
	}

	if !g.HeadConnected("c", "a") || !g.HeadConnected("d", "a") {
		t.Error("c,d Head nodes not right")
	}
}
