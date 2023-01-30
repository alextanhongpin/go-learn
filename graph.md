```golang
package main

import (
	"os"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

//go:generate dot -Tpng -O graph.gv
func main() {
	g := graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())
	g.AddVertex("start")
	g.AddVertex("completed")
	g.AddVertex("compensated")
	g.AddVertex("booking_created")
	g.AddVertex("booking_cancelled")
	g.AddVertex("payment_made")
	g.AddVertex("payment_refunded")
	g.AddVertex("hotel_reserved")
	g.AddVertex("hotel_unavailable")

	red := graph.EdgeAttribute("color", "red")
	green := graph.EdgeAttribute("color", "green")
	g.AddEdge("start", "booking_created", graph.EdgeAttribute("label", "make_booking"), green)
	g.AddEdge("booking_created", "payment_made", graph.EdgeAttribute("label", "wait_payment"), green)
	g.AddEdge("payment_made", "hotel_reserved", graph.EdgeAttribute("label", "reserve_hotel"), green)
	g.AddEdge("hotel_reserved", "completed", graph.EdgeAttribute("label", "complete"), green)

	g.AddEdge("booking_created", "booking_cancelled", graph.EdgeAttribute("label", "reject_payment"), red)
	g.AddEdge("payment_made", "payment_refunded", graph.EdgeAttribute("label", "refund_payment"), red)
	//g.AddEdge("payment_made", "hotel_unavailable", graph.EdgeAttribute("label", "refund_payment"), red)
	g.AddEdge("payment_refunded", "booking_cancelled", graph.EdgeAttribute("label", "cancel_booking"), red)
	g.AddEdge("booking_cancelled", "compensated", graph.EdgeAttribute("label", "compensate"), red)

	//m, err := g.AdjacencyMap()
	//if err != nil {
	//panic(err)
	//}

	file, _ := os.Create("./graph.gv")
	_ = draw.DOT(g, file)
}
```
