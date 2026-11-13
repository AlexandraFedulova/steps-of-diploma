package graph

import "math"

// NodeCommunity — назначение вершины в сообщество Louvain.
type NodeCommunity struct {
	NodeID      string `json:"nodeId"`
	CommunityID int    `json:"communityId"`
}

// LouvainCommunity — группа вершин с одним community id.
type LouvainCommunity struct {
	ID    int      `json:"id"`
	Nodes []string `json:"nodes"`
}

// LouvainResult — результат community detection.
type LouvainResult struct {
	Modularity  float64            `json:"modularity"`
	Assignments []NodeCommunity    `json:"assignments"`
	Communities []LouvainCommunity `json:"communities"`
}

// LouvainCommunities выполняет Louvain modularity optimization на взвешенном графе.
//
// Ориентированный граф приводится к неориентированному: w(u,v) = |w(u→v)| + |w(v→u)|.
// Нулевые связи игнорируются.
func (g *Graph) LouvainCommunities() *LouvainResult {
	if g.NodeCount() == 0 {
		return &LouvainResult{}
	}

	state := newLouvainState(g)
	if state.m2 == 0 {
		return state.trivialResult()
	}

	for {
		improved := state.localMoving()
		if !improved {
			break
		}
		if !state.aggregate() {
			break
		}
	}

	return state.buildResult()
}

type louvainState struct {
	origLabels []string
	members    [][]int // текущий узел → индексы исходных вершин

	adj       map[int]map[int]float64
	degree    []float64
	m2        float64
	community []int
	commTot   map[int]float64
}

func newLouvainState(g *Graph) *louvainState {
	n := g.NodeCount()
	index := make(map[string]int, n)
	origLabels := make([]string, n)

	for i, id := range g.order {
		origLabels[i] = id
		index[id] = i
	}

	adj := make(map[int]map[int]float64, n)
	degree := make([]float64, n)
	members := make([][]int, n)

	for i := 0; i < n; i++ {
		adj[i] = make(map[int]float64)
		members[i] = []int{i}
	}

	for from, edges := range g.adj {
		i := index[from]
		for _, e := range edges {
			j, ok := index[e.To]
			if !ok {
				continue
			}
			w := math.Abs(e.Weight)
			if w == 0 {
				continue
			}
			adj[i][j] += w
			if i != j {
				adj[j][i] += w
			}
			degree[i] += w
			if i != j {
				degree[j] += w
			}
		}
	}

	m2 := totalEdgeWeight(adj, n)

	community := make([]int, n)
	commTot := make(map[int]float64, n)
	for i := 0; i < n; i++ {
		community[i] = i
		commTot[i] = degree[i]
	}

	return &louvainState{
		origLabels: origLabels,
		members:    members,
		adj:        adj,
		degree:     degree,
		m2:         m2,
		community:  community,
		commTot:    commTot,
	}
}

func (s *louvainState) size() int {
	return len(s.members)
}

func (s *louvainState) trivialResult() *LouvainResult {
	assignments := make([]NodeCommunity, len(s.origLabels))
	communities := make([]LouvainCommunity, len(s.origLabels))
	for i, id := range s.origLabels {
		assignments[i] = NodeCommunity{NodeID: id, CommunityID: i}
		communities[i] = LouvainCommunity{ID: i, Nodes: []string{id}}
	}
	return &LouvainResult{
		Assignments: assignments,
		Communities: communities,
	}
}

func (s *louvainState) localMoving() bool {
	improved := false
	n := s.size()

	for {
		passImproved := false

		for i := 0; i < n; i++ {
			current := s.community[i]
			ki := s.degree[i]

			s.removeFromCommunity(i, current)

			bestComm := current
			bestGain := 0.0

			neighborComms := make(map[int]float64)
			for j, w := range s.adj[i] {
				if w == 0 {
					continue
				}
				c := s.community[j]
				if c < 0 {
					continue
				}
				neighborComms[c] += w
			}

			for comm, kiIn := range neighborComms {
				sigmaTot := s.commTot[comm]
				gain := kiIn/s.m2 - (sigmaTot*ki)/(2*s.m2*s.m2)
				if gain > bestGain+1e-12 {
					bestGain = gain
					bestComm = comm
				}
			}

			s.addToCommunity(i, bestComm)

			if bestComm != current {
				passImproved = true
				improved = true
			}
		}

		if !passImproved {
			break
		}
	}

	return improved
}

func (s *louvainState) removeFromCommunity(node, comm int) {
	s.commTot[comm] -= s.degree[node]
	s.community[node] = -1
}

func (s *louvainState) addToCommunity(node, comm int) {
	s.community[node] = comm
	s.commTot[comm] += s.degree[node]
}

func (s *louvainState) aggregate() bool {
	n := s.size()
	unique := make(map[int]int)
	for _, c := range s.community {
		if c < 0 {
			continue
		}
		if _, ok := unique[c]; !ok {
			unique[c] = len(unique)
		}
	}

	if len(unique) >= n {
		return false
	}

	oldToNew := make(map[int]int, len(unique))
	for old, newID := range unique {
		oldToNew[old] = newID
	}

	newN := len(unique)
	newMembers := make([][]int, newN)
	newAdj := make(map[int]map[int]float64, newN)
	newDegree := make([]float64, newN)
	newCommunity := make([]int, newN)
	newCommTot := make(map[int]float64, newN)

	for i := 0; i < newN; i++ {
		newAdj[i] = make(map[int]float64)
		newCommunity[i] = i
	}

	for i := 0; i < n; i++ {
		nc := oldToNew[s.community[i]]
		newMembers[nc] = append(newMembers[nc], s.members[i]...)
	}

	for i := 0; i < n; i++ {
		ci := oldToNew[s.community[i]]
		for j, w := range s.adj[i] {
			if j <= i || w == 0 {
				continue
			}
			cj := oldToNew[s.community[j]]
			if ci == cj {
				continue
			}
			newAdj[ci][cj] += w
			newAdj[cj][ci] += w
			newDegree[ci] += w
			newDegree[cj] += w
		}
	}

	newM2 := totalEdgeWeight(newAdj, newN)

	for i := 0; i < newN; i++ {
		newCommTot[i] = newDegree[i]
	}

	s.members = newMembers
	s.adj = newAdj
	s.degree = newDegree
	s.m2 = newM2
	s.community = newCommunity
	s.commTot = newCommTot

	return true
}

func (s *louvainState) modularity() float64 {
	if s.m2 == 0 {
		return 0
	}
	m := s.m2 / 2
	q := 0.0

	commNodes := make(map[int][]int)
	for i, c := range s.community {
		commNodes[c] = append(commNodes[c], i)
	}

	for _, nodes := range commNodes {
		var lc, dc float64
		for _, i := range nodes {
			dc += s.degree[i]
		}
		for ai := 0; ai < len(nodes); ai++ {
			for aj := ai; aj < len(nodes); aj++ {
				i, j := nodes[ai], nodes[aj]
				lc += s.adj[i][j]
			}
		}
		q += lc/m - (dc/(2*m))*(dc/(2*m))
	}

	return q
}

func (s *louvainState) buildResult() *LouvainResult {
	idRemap := make(map[int]int)
	groups := make(map[int][]string)

	for i := 0; i < s.size(); i++ {
		cid := s.community[i]
		if _, ok := idRemap[cid]; !ok {
			idRemap[cid] = len(idRemap)
		}
		newID := idRemap[cid]
		for _, origIdx := range s.members[i] {
			groups[newID] = append(groups[newID], s.origLabels[origIdx])
		}
	}

	assignments := make([]NodeCommunity, 0, len(s.origLabels))
	communities := make([]LouvainCommunity, 0, len(groups))

	for cid := 0; cid < len(idRemap); cid++ {
		nodes := groups[cid]
		communities = append(communities, LouvainCommunity{ID: cid, Nodes: nodes})
		for _, nodeID := range nodes {
			assignments = append(assignments, NodeCommunity{
				NodeID:      nodeID,
				CommunityID: cid,
			})
		}
	}

	// стабильный порядок assignments по порядку вершин в графе
	byID := make(map[string]int, len(idRemap))
	for _, c := range communities {
		byID[c.Nodes[0]] = c.ID // fallback
		for _, n := range c.Nodes {
			byID[n] = c.ID
		}
	}
	ordered := make([]NodeCommunity, len(s.origLabels))
	for i, id := range s.origLabels {
		ordered[i] = NodeCommunity{NodeID: id, CommunityID: byID[id]}
	}

	return &LouvainResult{
		Modularity:  s.modularity(),
		Assignments: ordered,
		Communities: communities,
	}
}

func totalEdgeWeight(adj map[int]map[int]float64, n int) float64 {
	var m2 float64
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			m2 += adj[i][j]
		}
	}
	return m2
}
