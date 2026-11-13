package graph

import (
	"os"
	"strings"
	"testing"

)

func TestCliqueCommunityAssignments_threeGroups(t *testing.T) {
	f, err := os.Open("../sociomatrix/testdata/communities.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	g, err := ParseCSV(f)
	if err != nil {
		t.Fatal(err)
	}

	assignments := g.CliqueCommunityAssignments()
	byNode := make(map[string]int, len(assignments))
	for _, a := range assignments {
		byNode[a.NodeID] = a.CommunityID
	}

	// Три взаимные тройки: Аня–Борис–Вика, Глеб–Дарья–Егор, Жанна–Иван–Кира
	same := func(a, b string) bool { return byNode[a] == byNode[b] }
	diff := func(a, b string) bool { return byNode[a] != byNode[b] }

	for _, trio := range [][3]string{
		{"Аня", "Борис", "Вика"},
		{"Глеб", "Дарья", "Егор"},
		{"Жанна", "Иван", "Кира"},
	} {
		if !same(trio[0], trio[1]) || !same(trio[1], trio[2]) {
			t.Errorf("expected one community for %v, got ids %d %d %d",
				trio, byNode[trio[0]], byNode[trio[1]], byNode[trio[2]])
		}
	}
	if !diff("Аня", "Глеб") || !diff("Глеб", "Жанна") {
		t.Fatalf("expected different communities across groups: %v", byNode)
	}
}

func TestCliqueCommunityAssignments_bridgeWeight1DoesNotMergeMutualTriangles(t *testing.T) {
	// Две взаимные тройки A–B–C и D–E–F.
	// Между тройками есть одиночная положительная дуга веса 1 (не взаимная) — сообщества не должны склеиваться.
	csv := `,A,B,C,D,E,F
A,0,1,1,0,0,0
B,1,0,1,0,0,0
C,1,1,0,0,0,0
D,1,0,0,0,1,1
E,0,0,0,1,0,1
F,0,0,0,1,1,0`

	g, err := ParseCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}

	assignments := g.CliqueCommunityAssignments()
	byNode := make(map[string]int, len(assignments))
	for _, a := range assignments {
		byNode[a.NodeID] = a.CommunityID
	}

	same := func(a, b string) bool { return byNode[a] == byNode[b] }
	diff := func(a, b string) bool { return byNode[a] != byNode[b] }

	// Внутри треугольников — одно сообщество.
	for _, trio := range [][3]string{
		{"A", "B", "C"},
		{"D", "E", "F"},
	} {
		if !same(trio[0], trio[1]) || !same(trio[1], trio[2]) {
			t.Fatalf("expected one community for %v, got ids %d %d %d",
				trio, byNode[trio[0]], byNode[trio[1]], byNode[trio[2]])
		}
	}

	// Между треугольниками — разные сообщества, несмотря на A→D = +1.
	if !diff("A", "D") {
		t.Fatalf("expected different communities across triangles, got: %v", byNode)
	}
}

func TestCliqueCommunityAssignments_threeTriangles_withOneWayBridgeKeepsThreeCommunities(t *testing.T) {
	// Три взаимные тройки: (1,2,3), (4,5,6), (7,8,9).
	// И одна дуга 9→1 веса 1 (не взаимная) — (1,2,3) и (7,8,9) должны остаться разными сообществами.
	csv := `,1,2,3,4,5,6,7,8,9
1,0,2,2,0,0,0,0,0,0
2,2,0,2,0,0,0,0,0,0
3,2,2,0,0,0,0,0,0,0
4,0,0,0,0,2,2,0,0,0
5,0,0,0,2,0,2,0,0,0
6,0,0,0,2,2,0,0,0,0
7,0,0,0,0,0,0,0,2,2
8,0,0,0,0,0,0,2,0,2
9,1,0,0,0,0,0,2,2,0`

	g, err := ParseCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}

	assignments := g.CliqueCommunityAssignments()
	byNode := make(map[string]int, len(assignments))
	for _, a := range assignments {
		byNode[a.NodeID] = a.CommunityID
	}

	same := func(a, b string) bool { return byNode[a] == byNode[b] }
	diff := func(a, b string) bool { return byNode[a] != byNode[b] }

	for _, trio := range [][3]string{
		{"1", "2", "3"},
		{"4", "5", "6"},
		{"7", "8", "9"},
	} {
		if !same(trio[0], trio[1]) || !same(trio[1], trio[2]) {
			t.Fatalf("expected one community for %v, got ids %d %d %d",
				trio, byNode[trio[0]], byNode[trio[1]], byNode[trio[2]])
		}
	}

	if !diff("1", "4") || !diff("4", "7") || !diff("1", "7") {
		t.Fatalf("expected three different communities, got: %v", byNode)
	}
}
