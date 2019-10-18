# simplegraph
 A simple thread-safe graph data structure in Go


How to use it:
```
import "github.com/rambutan-soft/simplegraph"
g := SimpleGraph.New()
g.Add("a", 1)
g.Add("b", 2)
g.Connect("a", "b", "ok")

n, _ := g.Get("a")

c := g.GetParents("b")

f := g.GetChildren("a")

g.Disconnect("a", "b")

g.Delete("b")
```