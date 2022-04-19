// 邻接矩阵实现graph
package adjmatrix

import (
	"errors"
	"fmt"
)

type GraphType int
type VRType int

const (
	// 有向图
	DG = iota
	// 有向网
	DN
	// 无向图
	UDG
	// 无向网
	UDN

	// 两个顶点是否有关联
	VRType0 = 0
	VRType1 = 1
)

type Vertex struct {
	data byte
}
type Arc struct {
	v, w *Vertex
	// vrType VRType
	info int // 如果是网，则是权重
}
type Graph struct {
	matrix  [][]*Arc
	vertics []*Vertex
	nv      int // 顶点个数
	ne      int // matrix中不为空的个数
	typ     GraphType
}

var vertexErr = errors.New("vertext error")
var edgeErr = errors.New("edge error")
var graphTypeErr = errors.New("graph type error")

func NewGraph(typ GraphType) (*Graph, error) {
	switch typ {
	case DG, UDG, DN, UDN:
		return &Graph{typ: typ}, nil
	default:
		return nil, graphTypeErr
	}
}
func NewVertex(data byte) *Vertex {
	return &Vertex{data: data}
}
func NewArc(v, w *Vertex, info int) *Arc {
	return &Arc{
		v:    v,
		w:    w,
		info: info,
	}
}

// 返回弧或者边的个数
func (g *Graph) getArcNum() int {
	if g.typ == UDG || g.typ == UDN {
		return g.ne / 2
	}
	return g.ne
}

func (g *Graph) print() {
	for _, v := range g.vertics {
		fmt.Print(v.data, ',')
	}
	fmt.Println()
	fmt.Println("-----")
	for _, row := range g.matrix {
		for _, v := range row {
			data := 0
			if v != nil {
				data = 1
			} else if g.typ == UDN || g.typ == DN {
				data = v.info
			}
			fmt.Print(data, ',')
		}
		fmt.Println()
	}
}

// 返回v的下标，若不存在，返回-1
func (g *Graph) LocateVex(v *Vertex) int {
	for i := range g.vertics {
		if v.data == g.vertics[i].data {
			return i
		}
	}
	return -1
}

// 返回v的第一个邻接点
func (g *Graph) FirstAdjVex(v *Vertex) (*Vertex, bool) {
	vi := g.LocateVex(v)
	if vi < 0 {
		return nil, false
	}
	for i := 0; i < g.nv; i++ {
		if g.matrix[vi][i] != nil {
			return g.vertics[i], true
		}
	}
	return nil, false
}

// 返回v相对于w的下一个邻接点
func (g *Graph) NextAdjVex(v, w *Vertex) (*Vertex, bool) {
	vi, wi := g.LocateVex(v), g.LocateVex(w)
	if vi < 0 || wi < 0 {
		return nil, false
	}
	for i := wi + 1; i < g.nv; i++ {
		if g.matrix[vi][i] != nil {
			return g.vertics[i], true
		}
	}
	return nil, false
}

// 添加新的v
func (g *Graph) InsertVex(v *Vertex) bool {
	if index := g.LocateVex(v); index < 0 {
		g.vertics = append(g.vertics, v)
		g.nv++
		return true
	}
	// 已经存在了
	return false
}

// 删除v及相关的弧
func (g *Graph) DeleteVex(v *Vertex) bool {
	index := g.LocateVex(v)
	if index < 0 {
		return false
	}
	var ne int
	for i, row := range g.matrix {
		for j := range row {
			if i == index && g.matrix[i][j] != nil {
				g.matrix[i][j] = nil
				ne++
			}
			if j == index && g.matrix[i][j] != nil {
				g.matrix[i][j] = nil
				ne++
			}
		}
	}
	g.nv--
	g.ne -= ne
	return true
}

// 添加新的弧v->w，如果g是无向，还要添加w->v
func (g *Graph) InsertArc(arc *Arc) bool {
	v, w := arc.v, arc.w
	iv := g.LocateVex(v)
	iw := g.LocateVex(w)
	if iv < 0 || iw < 0 {
		return false
	}
	g.matrix[iv][iw] = arc
	if g.typ == UDG || g.typ == UDN {
		newArc := NewArc(w, v, arc.info)
		g.matrix[iw][iv] = newArc
	}

	return true
}

// 删除弧v->w,如果g是无向的，还要删除w->v
func (g *Graph) DeleteArc(arc *Arc) bool {
	v, w := arc.v, arc.w
	iv := g.LocateVex(v)
	iw := g.LocateVex(w)
	if iv < 0 || iw < 0 {
		return false
	}
	g.matrix[iv][iw] = nil
	if g.typ == UDG || g.typ == UDN {
		g.matrix[iw][iv] = nil
	}

	return true
}

// 深度优先遍历
func (g *Graph) DFSTraverse() {
	visited := make([]int, len(g.vertics))
	v := g.vertics[0]
	dfs(visited, g, v)
}

// 广度优先遍历
func (g *Graph) BFSTraverse() {

}
func dfs(visited []int, g *Graph, v *Vertex) {
	// visit
	vi := g.LocateVex(v)
	visited[vi] = 1
	fmt.Println(v.data)

	// todo
	for {
		w, exist := g.FirstAdjVex(v)
		if exist {
			wi := g.LocateVex(w)
			if visited[wi] == 0 {
				dfs(visited, g, w)
			}
		}
	}
}
