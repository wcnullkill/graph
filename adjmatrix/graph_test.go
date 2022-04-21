package adjmatrix

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initGraph(g *Graph, typ GraphType, matrix [][]int) {
	verList := make([]*Vertex, len(matrix))

	for i := 0; i < len(verList); i++ {
		verList[i] = NewVertex("V" + strconv.Itoa(i+1))
		g.InsertVex(verList[i])
	}

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			if matrix[i][j] > 0 {
				v, w := verList[i], verList[j]
				if typ == UDN || typ == DN {
					arc := NewArc(v, w, matrix[i][j])
					g.InsertArc(arc)
				} else {
					arc := NewArc(v, w, 0)
					g.InsertArc(arc)
				}
			}
		}
	}

}
func TestCreateGraph(t *testing.T) {
	assert := assert.New(t)
	data := []struct {
		typ GraphType
		err error
	}{
		{UDG, nil},
		{UDN, nil},
		{DG, nil},
		{DN, nil},
		{-1, graphTypeErr},
		{10, graphTypeErr},
	}
	for _, v := range data {
		g, err := NewGraph(v.typ)
		if v.err == nil {
			assert.Nil(err)
			assert.NotNil(g)
		} else {
			assert.ErrorIs(err, v.err)
			assert.Nil(g)
		}
	}
}

func TestInsert(t *testing.T) {
	assert := assert.New(t)

	data := []struct {
		typ    GraphType
		matrix [][]int
		vnum   int // 顶点个数
		enum   int // 边或者弧条数
	}{
		{
			typ: UDG,
			matrix: [][]int{
				{0, 1, 0, 1, 0},
				{1, 0, 1, 0, 1},
				{0, 1, 0, 1, 1},
				{1, 0, 1, 0, 0},
				{0, 1, 1, 0, 0},
			},
			vnum: 5,
			enum: 6,
		},
		{
			typ: UDN,
			matrix: [][]int{
				{0, 10, 0, 11, 0},
				{10, 0, 12, 0, 15},
				{0, 12, 0, 14, 13},
				{11, 0, 14, 0, 0},
				{0, 15, 13, 0, 0},
			},
			vnum: 5,
			enum: 6,
		},
		{
			typ: DG,
			matrix: [][]int{
				{0, 1, 1, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 1},
				{1, 0, 0, 0},
			},
			vnum: 4,
			enum: 4,
		},
		{
			typ: DN,
			matrix: [][]int{
				{0, 11, 12, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 13},
				{14, 0, 0, 0},
			},
			vnum: 4,
			enum: 4,
		},
	}

	for _, v := range data {
		g, err := NewGraph(v.typ)
		assert.Nil(err)
		assert.NotNil(g)
		initGraph(g, v.typ, v.matrix)

		assert.Equal(g.nv, v.vnum)
		assert.Equal(g.getArcNum(), v.enum)
		for i := range g.matrix {
			for j := range g.matrix[i] {
				if v.matrix[i][j] == 0 && g.matrix[i][j] != nil {
					assert.Failf("matrix error", "i:%d,j:%d", i, j)
				}
				if v.matrix[i][j] > 0 && g.matrix[i][j] == nil {
					assert.Failf("matrix error", "i:%d,j:%d", i, j)
				}
				if v.typ == UDN || v.typ == DN {
					if v.matrix[i][j] > 0 && g.matrix[i][j] != nil && g.matrix[i][j].info != v.matrix[i][j] {
						assert.Failf("matrix error", "i:%d,j:%d", i, j)
					}
				}
			}
		}
	}
}

func TestDeleteArc(t *testing.T) {
	assert := assert.New(t)

	data := []struct {
		typ    GraphType
		matrix [][]int
		vnum   int // 顶点个数
		enum   int // 边或者弧条数
	}{
		{
			typ: UDG,
			matrix: [][]int{
				{0, 1, 0, 1, 0},
				{1, 0, 1, 0, 1},
				{0, 1, 0, 1, 1},
				{1, 0, 1, 0, 0},
				{0, 1, 1, 0, 0},
			},
			vnum: 5,
			enum: 6,
		},
		{
			typ: UDN,
			matrix: [][]int{
				{0, 10, 0, 11, 0},
				{10, 0, 12, 0, 15},
				{0, 12, 0, 14, 13},
				{11, 0, 14, 0, 0},
				{0, 15, 13, 0, 0},
			},
			vnum: 5,
			enum: 6,
		},
		{
			typ: DG,
			matrix: [][]int{
				{0, 1, 1, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 1},
				{1, 0, 0, 0},
			},
			vnum: 4,
			enum: 4,
		},
		{
			typ: DN,
			matrix: [][]int{
				{0, 11, 12, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 13},
				{14, 0, 0, 0},
			},
			vnum: 4,
			enum: 4,
		},
	}

	for _, v := range data {
		g, _ := NewGraph(v.typ)
		initGraph(g, v.typ, v.matrix)

		for i := range g.matrix {
			for j := range g.matrix[i] {
				if v.matrix[i][j] > 0 {
					vi := NewVertex("V" + strconv.Itoa(i+1))
					vj := NewVertex("V" + strconv.Itoa(j+1))
					arc := NewArc(vi, vj, 0)
					g.DeleteArc(arc)
				}
			}
		}
		for i := range g.matrix {
			for j := range g.matrix[i] {
				assert.Nil(g.matrix[i][j])
			}
		}
		assert.Equal(0, g.getArcNum())
		assert.Equal(v.vnum, g.getVexNum())
	}
}

func TestDeleteVex(t *testing.T) {
	assert := assert.New(t)

	data := []struct {
		typ    GraphType
		matrix [][]int
		vnum   int // 顶点个数
		enum   int // 边或者弧条数
	}{
		{
			typ: UDG,
			matrix: [][]int{
				{0, 1, 0, 1, 0},
				{1, 0, 1, 0, 1},
				{0, 1, 0, 1, 1},
				{1, 0, 1, 0, 0},
				{0, 1, 1, 0, 0},
			},
			vnum: 5,
			enum: 6,
		},
		{
			typ: UDN,
			matrix: [][]int{
				{0, 10, 0, 11, 0},
				{10, 0, 12, 0, 15},
				{0, 12, 0, 14, 13},
				{11, 0, 14, 0, 0},
				{0, 15, 13, 0, 0},
			},
			vnum: 5,
			enum: 6,
		},
		{
			typ: DG,
			matrix: [][]int{
				{0, 1, 1, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 1},
				{1, 0, 0, 0},
			},
			vnum: 4,
			enum: 4,
		},
		{
			typ: DN,
			matrix: [][]int{
				{0, 11, 12, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 13},
				{14, 0, 0, 0},
			},
			vnum: 4,
			enum: 4,
		},
	}

	for _, v := range data {
		g, _ := NewGraph(v.typ)
		initGraph(g, v.typ, v.matrix)

		for i := 0; i < v.vnum; i++ {
			vi := NewVertex("V" + strconv.Itoa(i+1))
			g.DeleteVex(vi)
		}

		assert.Equal(0, g.getArcNum())
		assert.Equal(0, g.getVexNum())
		assert.Equal(0, len(g.matrix))
		assert.Equal(0, len(g.vertics))
	}
}

func TestDFSTraverse(t *testing.T) {
	assert := assert.New(t)
	data := []struct {
		typ    GraphType
		matrix [][]int
		dfs    []string
	}{
		{
			typ: UDG,
			matrix: [][]int{
				{0, 1, 0, 1, 0},
				{1, 0, 1, 0, 1},
				{0, 1, 0, 1, 1},
				{1, 0, 1, 0, 0},
				{0, 1, 1, 0, 0},
			},
			dfs: []string{"V1", "V2", "V3", "V4", "V5"},
		},
		{
			typ: UDN,
			matrix: [][]int{
				{0, 11, 0, 12, 0},
				{11, 0, 13, 0, 14},
				{0, 13, 0, 15, 16},
				{12, 0, 15, 0, 0},
				{0, 14, 16, 0, 0},
			},
			dfs: []string{"V1", "V2", "V3", "V4", "V5"},
		},
		{
			typ: UDG,
			matrix: [][]int{
				{0, 1, 1, 0, 0, 0, 0, 0},
				{1, 0, 0, 1, 1, 0, 0, 0},
				{1, 0, 0, 0, 0, 1, 1, 0},
				{0, 1, 0, 0, 0, 0, 0, 1},
				{0, 1, 0, 0, 0, 0, 0, 1},
				{0, 0, 1, 0, 0, 0, 1, 0},
				{0, 0, 1, 0, 0, 1, 0, 0},
				{0, 0, 0, 1, 1, 0, 0, 0},
			},
			dfs: []string{"V1", "V2", "V4", "V8", "V5", "V3", "V6", "V7"},
		},
		{
			typ: UDN,
			matrix: [][]int{
				{0, 11, 12, 0, 0, 0, 0, 0},
				{11, 0, 0, 13, 14, 0, 0, 0},
				{12, 0, 0, 0, 0, 15, 16, 0},
				{0, 1, 0, 0, 0, 0, 0, 17},
				{0, 1, 0, 0, 0, 0, 0, 18},
				{0, 0, 1, 0, 0, 0, 1, 0},
				{0, 0, 1, 0, 0, 1, 0, 0},
				{0, 0, 0, 1, 1, 0, 0, 0},
			},
			dfs: []string{"V1", "V2", "V4", "V8", "V5", "V3", "V6", "V7"},
		},
	}
	for _, v := range data {
		g, _ := NewGraph(v.typ)
		initGraph(g, v.typ, v.matrix)

		dfs := []string{}
		g.DFSTraverse(func(vex *Vertex) {
			dfs = append(dfs, vex.data)
		})

		assert.Equal(len(dfs), len(v.dfs))
		for i := range dfs {
			assert.Equal(dfs[i], v.dfs[i])
		}
	}
}

func TestBFSTraverse(t *testing.T) {
	assert := assert.New(t)
	data := []struct {
		typ    GraphType
		matrix [][]int
		bfs    []string
	}{
		{
			typ: UDG,
			matrix: [][]int{
				{0, 1, 1, 0, 0, 0, 0, 0},
				{1, 0, 0, 1, 1, 0, 0, 0},
				{1, 0, 0, 0, 0, 1, 1, 0},
				{0, 1, 0, 0, 0, 0, 0, 1},
				{0, 1, 0, 0, 0, 0, 0, 1},
				{0, 0, 1, 0, 0, 0, 1, 0},
				{0, 0, 1, 0, 0, 1, 0, 0},
				{0, 0, 0, 1, 1, 0, 0, 0},
			},
			bfs: []string{"V1", "V2", "V3", "V4", "V5", "V6", "V7", "V8"},
		},
		{
			typ: UDN,
			matrix: [][]int{
				{0, 11, 12, 0, 0, 0, 0, 0},
				{11, 0, 0, 13, 14, 0, 0, 0},
				{12, 0, 0, 0, 0, 15, 16, 0},
				{0, 1, 0, 0, 0, 0, 0, 17},
				{0, 1, 0, 0, 0, 0, 0, 18},
				{0, 0, 1, 0, 0, 0, 1, 0},
				{0, 0, 1, 0, 0, 1, 0, 0},
				{0, 0, 0, 1, 1, 0, 0, 0},
			},
			bfs: []string{"V1", "V2", "V3", "V4", "V5", "V6", "V7", "V8"},
		},
	}
	for _, v := range data {
		g, _ := NewGraph(v.typ)
		initGraph(g, v.typ, v.matrix)

		bfs := []string{}
		g.BFSTraverse(func(vex *Vertex) {
			bfs = append(bfs, vex.data)
		})

		assert.Equal(len(bfs), len(v.bfs))
		for i := range bfs {
			assert.Equal(bfs[i], v.bfs[i])
		}
	}
}
