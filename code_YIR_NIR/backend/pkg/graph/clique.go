package graph

// CliqueCommunity — устойчивое сообщество-клика (взаимные «+» внутри, без «+» наружу).
type CliqueCommunity struct {
	ID          int      `json:"id"`
	Nodes       []string `json:"nodes"`
	Size        int      `json:"size"`
	Mutual      bool     `json:"mutual"`      // все пары взаимно «+»
	Directed    bool     `json:"directed"`    // все пары u→v «+» (ориентированная клика)
	Closed      bool     `json:"closed"`      // нет положительных связей с внешними
	OutgoingExt int      `json:"outgoingExt"` // число исходящих «+» за пределы группы
	IncomingExt int      `json:"incomingExt"` // число входящих «+» извне
}

// CliqueCommunitiesResult — результат поиска клик-сообществ.
type CliqueCommunitiesResult struct {
	Strict            []CliqueCommunity `json:"strict"`            // closed: нет «+» с внешними, внутри взаимная клика
	MutualCliques     []CliqueCommunity `json:"mutualCliques"`     // макс. клики по взаимному «+»
	DirectedCliques   []CliqueCommunity `json:"directedCliques"`   // макс. ориентированные клики (все u→v «+»)
	AnyPositiveCliques []CliqueCommunity `json:"anyPositiveCliques"` // клики графа «есть u→v или v→u +»
}

// CliqueCommunities ищет сообщества как клики взаимного выбора.
//
// Строгое сообщество: полный взаимный «+» внутри и ни одного «+» между S и остальными.
func (g *Graph) CliqueCommunities() *CliqueCommunitiesResult {
	n := g.NodeCount()
	if n == 0 {
		return &CliqueCommunitiesResult{}
	}

	ids := make([]string, n)
	index := make(map[string]int, n)
	for i, node := range g.Nodes() {
		ids[i] = node.ID
		index[node.ID] = i
	}

	pos := func(i, j int) bool {
		for _, e := range g.adj[ids[i]] {
			if e.To == ids[j] {
				return e.Weight > 0
			}
		}
		return false
	}

	mutual := func(i, j int) bool {
		return pos(i, j) && pos(j, i)
	}

	anyPos := func(i, j int) bool {
		return pos(i, j) || pos(j, i)
	}

	mutualMax := bronKerbosch(n, mutual)
	dirMax := findMaximalDirectedCliques(n, pos)
	anyMax := bronKerbosch(n, anyPos)

	toComm := func(nodes []int, mutualClique, directedClique bool) CliqueCommunity {
		nodeIDs := make([]string, len(nodes))
		for i, v := range nodes {
			nodeIDs[i] = ids[v]
		}
		outExt, inExt := 0, 0
		set := make(map[int]struct{}, len(nodes))
		for _, v := range nodes {
			set[v] = struct{}{}
		}
		for _, u := range nodes {
			for v := 0; v < n; v++ {
				if _, inside := set[v]; inside {
					continue
				}
				if pos(u, v) {
					outExt++
				}
				if pos(v, u) {
					inExt++
				}
			}
		}
		closed := outExt == 0 && inExt == 0
		return CliqueCommunity{
			Nodes:       nodeIDs,
			Size:        len(nodeIDs),
			Mutual:      mutualClique,
			Directed:    directedClique,
			Closed:      closed,
			OutgoingExt: outExt,
			IncomingExt: inExt,
		}
	}

	var mutualCliques, directedCliques, anyCliques, strict []CliqueCommunity
	for i, clique := range mutualMax {
		if len(clique) < 2 {
			continue
		}
		c := toComm(clique, true, false)
		c.ID = i
		mutualCliques = append(mutualCliques, c)
		if c.Closed && c.Mutual {
			strict = append(strict, c)
		}
	}
	for i, clique := range dirMax {
		if len(clique) < 2 {
			continue
		}
		c := toComm(clique, isMutualClique(clique, mutual), true)
		c.ID = i
		directedCliques = append(directedCliques, c)
	}
	for i, clique := range anyMax {
		if len(clique) < 2 {
			continue
		}
		c := toComm(clique, isMutualClique(clique, mutual), isDirectedClique(clique, pos))
		c.ID = i
		anyCliques = append(anyCliques, c)
	}

	return &CliqueCommunitiesResult{
		Strict:             strict,
		MutualCliques:      mutualCliques,
		DirectedCliques:    directedCliques,
		AnyPositiveCliques: anyCliques,
	}
}

// CliqueCommunityAssignments назначает вершинам communityId для визуализации.
//
// Берутся замкнутые взаимные клики (Strict); если их нет — максимальные клики
// взаимного положительного выбора (MutualCliques). Клики обрабатываются по
// убыванию размера: вершина получает id первой клики, в которую входит.
// Вершины вне клик (размер < 2) объединяются в одно сообщество «прочие».
func (g *Graph) CliqueCommunityAssignments() []NodeCommunity {
	n := g.NodeCount()
	if n == 0 {
		return nil
	}

	res := g.CliqueCommunities()
	// По теории используем только клики взаимного положительного выбора:
	// 1) сначала назначаем замкнутые взаимные клики (Strict);
	// 2) затем — остальные максимальные взаимные клики (MutualCliques) для нераспределённых вершин.
	strict := append([]CliqueCommunity(nil), res.Strict...)
	sortCliquesBySizeDesc(strict)
	mutual := append([]CliqueCommunity(nil), res.MutualCliques...)
	sortCliquesBySizeDesc(mutual)

	byNode := make(map[string]int, n)
	nextID := 0
	assignFrom := func(cliques []CliqueCommunity) {
		for _, c := range cliques {
			if len(c.Nodes) < 2 {
				continue
			}
			assignedAny := false
			for _, nodeID := range c.Nodes {
				if _, ok := byNode[nodeID]; !ok {
					byNode[nodeID] = nextID
					assignedAny = true
				}
			}
			if assignedAny {
				nextID++
			}
		}
	}

	assignFrom(strict)
	assignFrom(mutual)

	otherID := nextID
	assignments := make([]NodeCommunity, 0, n)
	for _, nodeID := range g.order {
		cid, ok := byNode[nodeID]
		if !ok {
			cid = otherID
		}
		assignments = append(assignments, NodeCommunity{
			NodeID:      nodeID,
			CommunityID: cid,
		})
	}
	return assignments
}

func sortCliquesBySizeDesc(cliques []CliqueCommunity) {
	for i := 0; i < len(cliques); i++ {
		for j := i + 1; j < len(cliques); j++ {
			if cliques[j].Size > cliques[i].Size {
				cliques[i], cliques[j] = cliques[j], cliques[i]
			}
		}
	}
}

func isDirectedClique(nodes []int, pos func(int, int) bool) bool {
	for _, u := range nodes {
		for _, v := range nodes {
			if u != v && !pos(u, v) {
				return false
			}
		}
	}
	return true
}

func isMutualClique(nodes []int, mutual func(int, int) bool) bool {
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			if !mutual(nodes[i], nodes[j]) {
				return false
			}
		}
	}
	return true
}

func bronKerbosch(n int, edge func(i, j int) bool) [][]int {
	neighbors := make([]map[int]struct{}, n)
	for i := 0; i < n; i++ {
		neighbors[i] = make(map[int]struct{})
		for j := 0; j < n; j++ {
			if i != j && edge(i, j) {
				neighbors[i][j] = struct{}{}
			}
		}
	}

	var result [][]int
	var bk func(r, p, x map[int]struct{})
	bk = func(r, p, x map[int]struct{}) {
		if len(p) == 0 && len(x) == 0 {
			if len(r) >= 2 {
				clique := make([]int, 0, len(r))
				for v := range r {
					clique = append(clique, v)
				}
				result = append(result, clique)
			}
			return
		}
		if len(p) == 0 {
			return
		}
		for v := range p {
			np := intersectSet(p, neighbors[v])
			nx := intersectSet(x, neighbors[v])
			nr := copySet(r)
			nr[v] = struct{}{}
			bk(nr, np, nx)
			delete(p, v)
			x[v] = struct{}{}
		}
	}
	all := make(map[int]struct{}, n)
	for i := 0; i < n; i++ {
		all[i] = struct{}{}
	}
	bk(nil, all, make(map[int]struct{}))
	return maximalBySize(result)
}

func copySet(s map[int]struct{}) map[int]struct{} {
	out := make(map[int]struct{}, len(s))
	for k := range s {
		out[k] = struct{}{}
	}
	return out
}

func intersectSet(a map[int]struct{}, b map[int]struct{}) map[int]struct{} {
	out := make(map[int]struct{})
	for k := range a {
		if _, ok := b[k]; ok {
			out[k] = struct{}{}
		}
	}
	return out
}

func findMaximalDirectedCliques(n int, pos func(i, j int) bool) [][]int {
	var all [][]int
	for size := n; size >= 2; size-- {
		combinations(n, size, func(combo []int) {
			for _, u := range combo {
				for _, v := range combo {
					if u != v && !pos(u, v) {
						return
					}
				}
			}
			all = append(all, append([]int(nil), combo...))
		})
	}
	return maximalBySize(all)
}

func combinations(n, k int, fn func([]int)) {
	combo := make([]int, k)
	var dfs func(start, depth int)
	dfs = func(start, depth int) {
		if depth == k {
			fn(combo)
			return
		}
		for i := start; i <= n-(k-depth); i++ {
			combo[depth] = i
			dfs(i+1, depth+1)
		}
	}
	dfs(0, 0)
}

func maximalBySize(sets [][]int) [][]int {
	if len(sets) == 0 {
		return nil
	}
	contains := func(a, b []int) bool {
		if len(a) > len(b) {
			return false
		}
		in := make(map[int]struct{}, len(b))
		for _, v := range b {
			in[v] = struct{}{}
		}
		for _, v := range a {
			if _, ok := in[v]; !ok {
				return false
			}
		}
		return true
	}
	var out [][]int
	for _, s := range sets {
		subset := false
		for _, t := range sets {
			if len(t) > len(s) && contains(s, t) {
				subset = true
				break
			}
		}
		if !subset {
			out = append(out, s)
		}
	}
	return out
}

