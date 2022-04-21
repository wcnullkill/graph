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
	data string
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

// 创建 n*n nil矩阵
func createMatrix(n int) [][]*Arc {
	matrix := make([][]*Arc, n)
	list := make([]*Arc, n)
	for i := 0; i < n; i++ {
		matrix[i] = list
	}
	return matrix
}

func NewGraph(typ GraphType) (*Graph, error) {
	switch typ {
	case DG, UDG, DN, UDN:
		return &Graph{
			typ:     typ,
			matrix:  make([][]*Arc, 0),
			vertics: make([]*Vertex, 0),
		}, nil
	default:
		return nil, graphTypeErr
	}
}
func NewVertex(data string) *Vertex {
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

// 返回顶点个数
func (g *Graph) getVexNum() int {
	return g.nv
}

func (g *Graph) print() {
	for _, v := range g.vertics {
		fmt.Print(v.data + "	")
	}
	fmt.Println()
	fmt.Println("-----")
	for _, row := range g.matrix {
		for _, v := range row {
			data := 0

			if v != nil {
				if g.typ == UDN || g.typ == DN {
					data = v.info
				} else {
					data = 1
				}
			}
			fmt.Print(data)
			fmt.Print(" ")
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
		// matrix 扩容
		g.grow(1)
		return true
	}
	// 已经存在了
	return false
}

// 扩容
func (g *Graph) grow(num int) bool {
	if num <= 0 {
		return false
	}
	newNum := len(g.matrix) + num
	appendList := make([]*Arc, num)
	for i := range g.matrix {
		g.matrix[i] = append(g.matrix[i], appendList...)
	}
	nilList := make([]*Arc, newNum)
	for j := 0; j < num; j++ {
		g.matrix = append(g.matrix, nilList)
	}
	g.nv++
	return true
}

// 删除v及相关的弧
func (g *Graph) DeleteVex(v *Vertex) bool {
	index := g.LocateVex(v)
	if index < 0 {
		return false
	}
	g.shrink(index)
	return true
}

// 收缩i行i列
func (g *Graph) shrink(index int) bool {
	oldlen := len(g.matrix)
	if oldlen == 0 {
		return false
	}
	if oldlen == 1 {
		g.matrix = make([][]*Arc, 0)
		g.vertics = make([]*Vertex, 0)
		g.ne = 0
		g.nv = 0
		return true
	}

	if index >= oldlen && index < 0 {
		return false
	}
	var ne int
	// 处理matrix
	newMatrix := make([][]*Arc, 0)
	for i := 0; i < oldlen; i++ {
		newRow := make([]*Arc, 0)
		for j := 0; j < oldlen; j++ {
			if i == index || j == index {
				if g.matrix[i][j] != nil {
					ne++
				}
				continue
			}
			newRow = append(newRow, g.matrix[i][j])
		}
		if len(newRow) > 0 {
			newMatrix = append(newMatrix, newRow)
		}
	}
	g.matrix = newMatrix
	// 处理vertics
	newVertics := make([]*Vertex, 0)
	newVertics = append(newVertics, g.vertics[:index]...)
	newVertics = append(newVertics, g.vertics[index+1:]...)
	g.vertics = newVertics
	g.ne -= ne
	g.nv--
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
	if g.matrix[iv][iw] == nil { // 避免重复
		g.matrix[iv][iw] = arc
		g.ne++
	}
	if g.typ == UDG || g.typ == UDN {
		newArc := NewArc(w, v, arc.info)
		if g.matrix[iw][iv] == nil { // 避免重复
			g.matrix[iw][iv] = newArc
			g.ne++
		}
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
	if g.matrix[iv][iw] != nil {
		g.matrix[iv][iw] = nil
		g.ne--
	}

	if g.typ == UDG || g.typ == UDN {
		if g.matrix[iw][iv] != nil {
			g.matrix[iw][iv] = nil
			g.ne--
		}
	}

	return true
}

// 深度优先遍历
func (g *Graph) DFSTraverse(visit func(*Vertex)) {
	visited := make([]int, len(g.vertics))
	for i, v := range g.vertics {
		if visited[i] == 0 {
			dfs(visited, g, v, visit)
		}
	}
}
func dfs(visited []int, g *Graph, v *Vertex, visit func(*Vertex)) {
	// visit
	vi := g.LocateVex(v)
	visited[vi] = 1
	visit(v)

	w, exist := g.FirstAdjVex(v)
	for exist {
		wi := g.LocateVex(w)
		if visited[wi] == 0 {
			dfs(visited, g, w, visit)
		}
		w, exist = g.NextAdjVex(v, w)
	}
}

// 广度优先遍历
func (g *Graph) BFSTraverse(visit func(*Vertex)) {
	visited := make([]int, len(g.vertics))
	queue := make([]*Vertex, 0, len(g.vertics))
	for i, v := range g.vertics {
		if visited[i] == 0 {
			// visit
			visit(v)
			visited[i] = 1
			queue = append(queue, v)
		}
		for len(queue) > 0 {
			item := queue[0]
			queue = queue[1:]
			w, exits := g.FirstAdjVex(item)
			for exits {
				wi := g.LocateVex(w)
				if visited[wi] == 0 {
					// visit
					visit(w)
					visited[wi] = 1
					queue = append(queue, w)
				}
				w, exits = g.NextAdjVex(item, w)
			}
		}
	}
}

// 孩子兄弟链表法表示树
type CSNode struct {
	v           *Vertex
	firstChild  *CSNode // 第一个孩子节点
	nextSibling *CSNode // 下一个兄弟节点
}

// 通过深度遍历+孩子兄弟链表法 生成森林
func (g *Graph) DFSForest() []*CSNode {
	visited := make([]int, len(g.vertics))

	forest := make([]*CSNode, 0)

	for i := 0; i < len(g.vertics); i++ {
		if visited[i] == 0 {
			head := &CSNode{}
			forest = append(forest, head)
			dfsTree(visited, g, g.vertics[i], head)
		}
	}
	return forest
}

func dfsTree(visited []int, g *Graph, v *Vertex, node *CSNode) {
	vi := g.LocateVex(v)
	visited[vi] = 1
	node.v = v
	node.firstChild = &CSNode{}
	node = node.firstChild
	w, exist := g.FirstAdjVex(v)
	first := true
	for exist {
		wi := g.LocateVex(w)
		if visited[wi] == 0 {
			visited[wi] = 1
			node.nextSibling = &CSNode{}
			node = node.nextSibling
			dfsTree(visited, g, w, node)
		}
		w, exist = g.NextAdjVex(v, w)
		first = false
	}
}
