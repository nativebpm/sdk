package nativebpm

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// getNodeDimensions returns the standard width and height for a BPMN node based on its type.
func getNodeDimensions(n map[string]interface{}) (width, height float64) {
	nType, _ := n["type"].(string)
	switch nType {
	case "startEvent", "endEvent":
		return 36.0, 36.0
	case "exclusiveGateway", "parallelGateway", "eventBasedGateway":
		return 50.0, 50.0
	default:
		return 100.0, 80.0
	}
}

// buildGraphAdj constructs outgoing and incoming adjacency lists and in-degree counts.
func buildGraphAdj(nodes []map[string]interface{}, flows []map[string]interface{}) (map[string][]string, map[string][]string, map[string]int) {
	outgoing := make(map[string][]string)
	incoming := make(map[string][]string)
	inCount := make(map[string]int)

	for _, f := range flows {
		src, _ := f["source"].(string)
		tgt, _ := f["target"].(string)
		if src != "" && tgt != "" {
			outgoing[src] = append(outgoing[src], tgt)
			incoming[tgt] = append(incoming[tgt], src)
			inCount[tgt]++
		}
	}
	return outgoing, incoming, inCount
}

// computeBFSLevels computes topological BFS depth levels for all nodes, detecting back-edges.
func computeBFSLevels(nodes []map[string]interface{}, outgoing map[string][]string, inCount map[string]int) (map[string]int, int) {
	state := make(map[string]int)
	backEdges := make(map[string]map[string]bool)

	var dfs func(string)
	dfs = func(u string) {
		state[u] = 1
		for _, v := range outgoing[u] {
			if state[v] == 1 {
				if backEdges[u] == nil {
					backEdges[u] = make(map[string]bool)
				}
				backEdges[u][v] = true
			} else if state[v] == 0 {
				dfs(v)
			}
		}
		state[u] = 2
	}

	for _, n := range nodes {
		id, _ := n["id"].(string)
		if id != "" && inCount[id] == 0 {
			dfs(id)
		}
	}
	for _, n := range nodes {
		id, _ := n["id"].(string)
		if id != "" && state[id] == 0 {
			dfs(id)
		}
	}

	var queue []string
	nodeLevels := make(map[string]int)
	queued := make(map[string]bool)

	for _, n := range nodes {
		id, _ := n["id"].(string)
		if id != "" && inCount[id] == 0 {
			queue = append(queue, id)
			queued[id] = true
			nodeLevels[id] = 0
		}
	}
	if len(queue) == 0 && len(nodes) > 0 {
		id, _ := nodes[0]["id"].(string)
		if id != "" {
			queue = append(queue, id)
			queued[id] = true
			nodeLevels[id] = 0
		}
	}

	maxDepth := 0
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		queued[curr] = false

		cLvl := nodeLevels[curr]
		if cLvl > maxDepth {
			maxDepth = cLvl
		}

		for _, child := range outgoing[curr] {
			if backEdges[curr] != nil && backEdges[curr][child] {
				continue
			}
			childLevel := cLvl + 1
			if childLevel > nodeLevels[child] {
				nodeLevels[child] = childLevel
				if !queued[child] {
					queue = append(queue, child)
					queued[child] = true
				}
			}
		}
	}

	for _, n := range nodes {
		id, _ := n["id"].(string)
		if id != "" {
			if _, ok := nodeLevels[id]; !ok {
				nodeLevels[id] = maxDepth + 1
			}
		}
	}

	return nodeLevels, maxDepth
}

// resolveCollisions performs iterative AABB bounding-box collision resolution.
func resolveCollisions(coords map[string]Coords, minGap float64) {
	nodeList := make([]string, 0, len(coords))
	for id := range coords {
		nodeList = append(nodeList, id)
	}
	sort.Strings(nodeList)

	for pass := 0; pass < 12; pass++ {
		collisionFound := false
		for i := 0; i < len(nodeList); i++ {
			idA := nodeList[i]
			cA := coords[idA]
			for j := i + 1; j < len(nodeList); j++ {
				idB := nodeList[j]
				cB := coords[idB]

				hOverlap := (cA.Width+cB.Width)/2.0 + minGap - math.Abs((cA.X+cA.Width/2.0)-(cB.X+cB.Width/2.0))
				vOverlap := (cA.Height+cB.Height)/2.0 + minGap - math.Abs((cA.Y+cA.Height/2.0)-(cB.Y+cB.Height/2.0))

				if hOverlap > 0 && vOverlap > 0 {
					collisionFound = true
					if cB.Y >= cA.Y {
						shift := (cA.Y + cA.Height + minGap) - cB.Y
						cB.Y += shift
						coords[idB] = cB
					} else {
						shift := (cB.Y + cB.Height + minGap) - cA.Y
						cA.Y += shift
						coords[idA] = cA
					}
				}
			}
		}
		if !collisionFound {
			break
		}
	}
}

// applyCustomCoords overrides calculated coordinates with user-provided custom coordinates.
func applyCustomCoords(coords map[string]Coords, customCoords map[string]Coords) {
	for id, custom := range customCoords {
		if c, exists := coords[id]; exists {
			if custom.Width > 0 {
				c.Width = custom.Width
			}
			if custom.Height > 0 {
				c.Height = custom.Height
			}
			c.X = custom.X
			c.Y = custom.Y
			coords[id] = c
		} else {
			coords[id] = custom
		}
	}
}

// fillUnplacedNodes assigns coordinates to any node not yet placed by a specialized strategy.
func fillUnplacedNodes(coords map[string]Coords, nodes []map[string]interface{}, defaultX, defaultY, rowSpacing float64) {
	unplacedCount := 0
	for _, n := range nodes {
		id, _ := n["id"].(string)
		if _, placed := coords[id]; !placed {
			w, h := getNodeDimensions(n)
			coords[id] = Coords{
				X:      defaultX,
				Y:      defaultY + float64(unplacedCount)*rowSpacing + (80.0-h)/2.0,
				Width:  w,
				Height: h,
			}
			unplacedCount++
		}
	}
}

// StrategyTiered implements the classic Left-to-Right topological pipeline layout.
func StrategyTiered(nodes []map[string]interface{}, flows []map[string]interface{}, opts LayoutOptions) map[string]Coords {
	if opts.ColSpacing <= 0 {
		opts.ColSpacing = 240.0
	}
	if opts.RowSpacing <= 0 {
		opts.RowSpacing = 130.0
	}
	if opts.StartX <= 0 {
		opts.StartX = 150.0
	}
	if opts.StartY <= 0 {
		opts.StartY = 200.0
	}

	outgoing, _, inCount := buildGraphAdj(nodes, flows)
	nodeLevels, _ := computeBFSLevels(nodes, outgoing, inCount)

	// Apply manual tiers if specified
	for id, tier := range opts.NodeTiers {
		nodeLevels[id] = tier
	}

	levelGroups := make(map[int][]string)
	for _, n := range nodes {
		id, _ := n["id"].(string)
		if id != "" {
			lvl := nodeLevels[id]
			levelGroups[lvl] = append(levelGroups[lvl], id)
		}
	}

	nodeMap := make(map[string]map[string]interface{})
	for _, n := range nodes {
		id, _ := n["id"].(string)
		nodeMap[id] = n
	}

	coords := make(map[string]Coords)
	for lvl, nIDs := range levelGroups {
		x := opts.StartX + float64(lvl)*opts.ColSpacing
		count := len(nIDs)
		for i, nID := range nIDs {
			w, h := getNodeDimensions(nodeMap[nID])
			y := opts.StartY + (float64(i)-float64(count-1)/2.0)*opts.RowSpacing + (80.0-h)/2.0
			coords[nID] = Coords{
				X:      x,
				Y:      y,
				Width:  w,
				Height: h,
			}
		}
	}

	resolveCollisions(coords, 25.0)
	applyCustomCoords(coords, opts.CustomCoords)
	return coords
}

// StrategyCenterHub implements the Radial Hub & Parallel Fan-out Corridors layout.
func StrategyCenterHub(nodes []map[string]interface{}, flows []map[string]interface{}, opts LayoutOptions) map[string]Coords {
	startX := opts.StartX
	if startX <= 0 {
		startX = 150.0
	}
	startY := opts.StartY
	if startY <= 0 {
		startY = 200.0
	}
	colSpacing := opts.ColSpacing
	if colSpacing <= 0 {
		colSpacing = 240.0
	}
	rowSpacing := opts.RowSpacing
	if rowSpacing <= 0 {
		rowSpacing = 140.0
	}

	nodeMap := make(map[string]map[string]interface{})
	for _, n := range nodes {
		id, _ := n["id"].(string)
		nodeMap[id] = n
	}

	outgoing, incoming, inCount := buildGraphAdj(nodes, flows)
	hubID := opts.CenterHubID
	if hubID == "" {
		metrics := AnalyzeGraphComplexity(nodes, flows)
		hubID = metrics.DetectedHubID
	}
	if hubID == "" && len(nodes) > 0 {
		maxDeg := -1
		for _, n := range nodes {
			id, _ := n["id"].(string)
			deg := len(outgoing[id]) + len(incoming[id])
			if deg > maxDeg {
				maxDeg = deg
				hubID = id
			}
		}
	}

	coords := make(map[string]Coords)
	hubY := 460.0
	hubX := startX + 500.0

	// 1. Place the focal Hub
	if hubNode, ok := nodeMap[hubID]; ok {
		hw, hh := getNodeDimensions(hubNode)
		coords[hubID] = Coords{
			X:      hubX,
			Y:      hubY + (80.0-hh)/2.0,
			Width:  hw,
			Height: hh,
		}
	}

	// 2. Place upstream/ingress nodes leading into Hub horizontally to the left
	currUp := hubID
	upDist := 1
	visitedUp := map[string]bool{hubID: true}
	for {
		preds := incoming[currUp]
		var nextUp string
		for _, p := range preds {
			if !visitedUp[p] {
				nextUp = p
				break
			}
		}
		if nextUp == "" {
			break
		}
		visitedUp[nextUp] = true
		w, h := getNodeDimensions(nodeMap[nextUp])
		coords[nextUp] = Coords{
			X:      hubX - float64(upDist)*colSpacing,
			Y:      hubY + (80.0-h)/2.0,
			Width:  w,
			Height: h,
		}
		upDist++
		currUp = nextUp
	}

	// Also catch any roots connected upstream
	for _, n := range nodes {
		id, _ := n["id"].(string)
		if inCount[id] == 0 && id != hubID && !visitedUp[id] {
			w, h := getNodeDimensions(n)
			coords[id] = Coords{
				X:      hubX - float64(upDist)*colSpacing,
				Y:      hubY + (80.0-h)/2.0,
				Width:  w,
				Height: h,
			}
			visitedUp[id] = true
			upDist++
		}
	}

	// 3. Downstream fan-out branches radiating from Hub
	branchTargets := outgoing[hubID]
	sort.Strings(branchTargets)
	numBranches := len(branchTargets)
	if numBranches == 0 {
		numBranches = 1
	}

	// Find the convergence node (where multiple branches rejoin, e.g. gw_ingress_join)
	convergenceID := ""
	for _, n := range nodes {
		id, _ := n["id"].(string)
		if id != hubID && len(incoming[id]) >= 2 {
			convergenceID = id
			break
		}
	}

	// Lay out each parallel branch along a symmetrical horizontal corridor/lane
	maxBranchX := hubX + colSpacing
	branchNodes := make(map[string]bool)
	for bIdx, bTarget := range branchTargets {
		laneY := hubY + (float64(bIdx)-float64(numBranches-1)/2.0)*180.0
		curr := bTarget
		step := 1
		for curr != "" && curr != convergenceID && !branchNodes[curr] {
			branchNodes[curr] = true
			w, h := getNodeDimensions(nodeMap[curr])
			nodeX := hubX + float64(step)*colSpacing
			if nodeX > maxBranchX {
				maxBranchX = nodeX
			}
			coords[curr] = Coords{
				X:      nodeX,
				Y:      laneY + (80.0-h)/2.0,
				Width:  w,
				Height: h,
			}

			// Advance to next node in branch
			next := ""
			for _, child := range outgoing[curr] {
				if child != convergenceID && !branchNodes[child] {
					next = child
					break
				}
			}
			curr = next
			step++
		}
	}

	// 4. Place Convergence node on the central axis
	convX := maxBranchX + colSpacing
	if convergenceID != "" {
		cw, ch := getNodeDimensions(nodeMap[convergenceID])
		coords[convergenceID] = Coords{
			X:      convX,
			Y:      hubY + (80.0-ch)/2.0,
			Width:  cw,
			Height: ch,
		}
	}

	// 5. Lay out post-convergence nodes using BFS from convergence node
	startPost := convergenceID
	if startPost == "" {
		startPost = hubID
	}
	postQueue := []string{startPost}
	postVisited := map[string]bool{startPost: true}
	postDist := make(map[string]int)
	postDist[startPost] = 0

	for len(postQueue) > 0 {
		curr := postQueue[0]
		postQueue = postQueue[1:]

		for _, child := range outgoing[curr] {
			if !postVisited[child] && !branchNodes[child] && child != hubID && coords[child].Width == 0 {
				postVisited[child] = true
				postDist[child] = postDist[curr] + 1
				postQueue = append(postQueue, child)
			}
		}
	}

	// Group post-convergence nodes by distance and place them
	postGroups := make(map[int][]string)
	for id, dist := range postDist {
		if id != startPost {
			postGroups[dist] = append(postGroups[dist], id)
		}
	}

	var distKeys []int
	for d := range postGroups {
		distKeys = append(distKeys, d)
	}
	sort.Ints(distKeys)

	for _, d := range distKeys {
		pIDs := postGroups[d]
		sort.Strings(pIDs)
		x := convX + float64(d)*colSpacing
		count := len(pIDs)
		for i, pID := range pIDs {
			w, h := getNodeDimensions(nodeMap[pID])
			y := hubY + (float64(i)-float64(count-1)/2.0)*rowSpacing + (80.0-h)/2.0
			coords[pID] = Coords{
				X:      x,
				Y:      y,
				Width:  w,
				Height: h,
			}
		}
	}

	fillUnplacedNodes(coords, nodes, convX+float64(len(distKeys)+1)*colSpacing, hubY, rowSpacing)
	resolveCollisions(coords, 25.0)
	applyCustomCoords(coords, opts.CustomCoords)
	return coords
}

// StrategyClusterMesh implements the Server Rack Matrix & Interconnect Mesh layout.
func StrategyClusterMesh(nodes []map[string]interface{}, flows []map[string]interface{}, opts LayoutOptions) map[string]Coords {
	startX := opts.StartX
	if startX <= 0 {
		startX = 100.0
	}
	startY := opts.StartY
	if startY <= 0 {
		startY = 180.0
	}
	rowSpacing := opts.RowSpacing
	if rowSpacing <= 0 {
		rowSpacing = 160.0
	}

	nodeMap := make(map[string]map[string]interface{})
	for _, n := range nodes {
		id, _ := n["id"].(string)
		nodeMap[id] = n
	}

	coords := make(map[string]Coords)

	// Server hosts in SRE topology: fin-1, fin-2, fin-4, fin-5
	hostRows := []string{"fin_1", "fin_2", "fin_4", "fin_5"}
	hostYMap := make(map[string]float64)
	for i, h := range hostRows {
		hostYMap[h] = startY + float64(i)*rowSpacing
	}
	centerY := startY + float64(len(hostRows)-1)/2.0*rowSpacing

	// Check if this workflow contains SRE server rack nodes
	hasSRENodes := false
	for id := range nodeMap {
		if strings.Contains(id, "fin_") {
			hasSRENodes = true
			break
		}
	}

	if hasSRENodes {
		// 1. Ingress Edge & Fork
		startID := "start"
		if _, ok := nodeMap[startID]; !ok {
			if _, ok := nodeMap["start_client"]; ok {
				startID = "start_client"
			}
		}
		if n, ok := nodeMap[startID]; ok {
			w, h := getNodeDimensions(n)
			coords[startID] = Coords{X: startX, Y: centerY + (80.0-h)/2.0, Width: w, Height: h}
		}
		if n, ok := nodeMap["task_cf_edge"]; ok {
			w, h := getNodeDimensions(n)
			coords["task_cf_edge"] = Coords{X: startX + 130.0, Y: centerY + (80.0-h)/2.0, Width: w, Height: h}
		}
		if n, ok := nodeMap["gw_argo_fork"]; ok {
			w, h := getNodeDimensions(n)
			coords["gw_argo_fork"] = Coords{X: startX + 270.0, Y: centerY + (80.0-h)/2.0, Width: w, Height: h}
		}

		// 2. 2D Server Rack Matrix (Columns: Tunnel, Nginx, App, KeyDB Master)
		tunnelX := startX + 440.0
		nginxX := startX + 660.0
		appX := startX + 880.0
		keydbX := startX + 1100.0 // All 4 KeyDB masters vertically aligned in column matrix!

		for _, host := range hostRows {
			rowY := hostYMap[host]

			tunnelID := fmt.Sprintf("task_tunnel_%s", host)
			if n, ok := nodeMap[tunnelID]; ok {
				w, h := getNodeDimensions(n)
				coords[tunnelID] = Coords{X: tunnelX, Y: rowY + (80.0-h)/2.0, Width: w, Height: h}
			}

			nginxID := fmt.Sprintf("task_nginx_%s", host)
			if n, ok := nodeMap[nginxID]; ok {
				w, h := getNodeDimensions(n)
				coords[nginxID] = Coords{X: nginxX, Y: rowY + (80.0-h)/2.0, Width: w, Height: h}
			}

			appID := fmt.Sprintf("task_app_%s", host)
			if n, ok := nodeMap[appID]; ok {
				w, h := getNodeDimensions(n)
				coords[appID] = Coords{X: appX, Y: rowY + (80.0-h)/2.0, Width: w, Height: h}
			}

			keydbID := fmt.Sprintf("task_keydb_%s", host)
			if n, ok := nodeMap[keydbID]; ok {
				w, h := getNodeDimensions(n)
				coords[keydbID] = Coords{X: keydbX, Y: rowY + (80.0-h)/2.0, Width: w, Height: h}
			}
		}

		// 3. Post-rack Ingress Join, Cache check, HAProxy, TiDB
		joinX := startX + 1300.0
		cacheX := startX + 1440.0
		if n, ok := nodeMap["gw_ingress_join"]; ok {
			w, h := getNodeDimensions(n)
			coords["gw_ingress_join"] = Coords{X: joinX, Y: centerY + (80.0-h)/2.0, Width: w, Height: h}
		}
		if n, ok := nodeMap["gw_cache_check"]; ok {
			w, h := getNodeDimensions(n)
			coords["gw_cache_check"] = Coords{X: cacheX, Y: centerY + (80.0-h)/2.0, Width: w, Height: h}
		}

		// Cache HIT path (top)
		if n, ok := nodeMap["task_keydb_hit"]; ok {
			w, h := getNodeDimensions(n)
			coords["task_keydb_hit"] = Coords{X: cacheX + 140.0, Y: centerY - 160.0 + (80.0-h)/2.0, Width: w, Height: h}
		}
		if n, ok := nodeMap["end_response_cached"]; ok {
			w, h := getNodeDimensions(n)
			coords["end_response_cached"] = Coords{X: cacheX + 320.0, Y: centerY - 160.0 + (80.0-h)/2.0, Width: w, Height: h}
		}

		// Cache MISS path (bottom/middle)
		missX := cacheX + 140.0
		haproxyX := missX + 180.0
		tidbX := haproxyX + 180.0
		raftRepX := tidbX + 180.0

		if n, ok := nodeMap["task_keydb_miss"]; ok {
			w, h := getNodeDimensions(n)
			coords["task_keydb_miss"] = Coords{X: missX, Y: centerY + (80.0-h)/2.0, Width: w, Height: h}
		}
		if n, ok := nodeMap["task_haproxy_router"]; ok {
			w, h := getNodeDimensions(n)
			coords["task_haproxy_router"] = Coords{X: haproxyX, Y: centerY + (80.0-h)/2.0, Width: w, Height: h}
		}
		if n, ok := nodeMap["task_tidb_master"]; ok {
			w, h := getNodeDimensions(n)
			coords["task_tidb_master"] = Coords{X: tidbX, Y: centerY + (80.0-h)/2.0, Width: w, Height: h}
		}
		if n, ok := nodeMap["gw_raft_replicate"]; ok {
			w, h := getNodeDimensions(n)
			coords["gw_raft_replicate"] = Coords{X: raftRepX, Y: centerY + (80.0-h)/2.0, Width: w, Height: h}
		}

		// 4. TiKV 5-Node 2D Multi-Raft Cluster Grid (2 columns x 3 rows)
		tikvColAX := raftRepX + 120.0
		tikvColBX := raftRepX + 270.0
		quorumX := raftRepX + 410.0

		tikvCoords := map[string]Coords{
			"task_tikv_fin_1": {X: tikvColAX, Y: centerY - 140.0},
			"task_tikv_fin_2": {X: tikvColAX, Y: centerY},
			"task_tikv_fin_4": {X: tikvColAX, Y: centerY + 140.0},
			"task_tikv_fin_5": {X: tikvColBX, Y: centerY - 70.0},
			"task_tikv_fin_3": {X: tikvColBX, Y: centerY + 70.0},
		}
		for id, c := range tikvCoords {
			if n, ok := nodeMap[id]; ok {
				w, h := getNodeDimensions(n)
				c.Width = w
				c.Height = h
				coords[id] = c
			}
		}

		if n, ok := nodeMap["gw_raft_quorum"]; ok {
			w, h := getNodeDimensions(n)
			coords["gw_raft_quorum"] = Coords{X: quorumX, Y: centerY + (80.0-h)/2.0, Width: w, Height: h}
		}

		// 5. Quorum Results: Commit vs Abort
		commitX := quorumX + 140.0
		syncX := commitX + 180.0
		endX := syncX + 180.0

		if n, ok := nodeMap["task_raft_commit"]; ok {
			w, h := getNodeDimensions(n)
			coords["task_raft_commit"] = Coords{X: commitX, Y: centerY - 60.0 + (80.0-h)/2.0, Width: w, Height: h}
		}
		if n, ok := nodeMap["task_keydb_sync"]; ok {
			w, h := getNodeDimensions(n)
			coords["task_keydb_sync"] = Coords{X: syncX, Y: centerY - 60.0 + (80.0-h)/2.0, Width: w, Height: h}
		}
		if n, ok := nodeMap["end_response"]; ok {
			w, h := getNodeDimensions(n)
			coords["end_response"] = Coords{X: endX, Y: centerY - 60.0 + (80.0-h)/2.0, Width: w, Height: h}
		}

		if n, ok := nodeMap["task_raft_abort"]; ok {
			w, h := getNodeDimensions(n)
			coords["task_raft_abort"] = Coords{X: commitX, Y: centerY + 180.0 + (80.0-h)/2.0, Width: w, Height: h}
		}
		if n, ok := nodeMap["end_failed"]; ok {
			w, h := getNodeDimensions(n)
			coords["end_failed"] = Coords{X: syncX, Y: centerY + 180.0 + (80.0-h)/2.0, Width: w, Height: h}
		}
	} else {
		// Generic 2D matrix layout for general workflows
		outgoing, _, inCount := buildGraphAdj(nodes, flows)
		nodeLevels, _ := computeBFSLevels(nodes, outgoing, inCount)

		levelGroups := make(map[int][]string)
		for _, n := range nodes {
			id, _ := n["id"].(string)
			lvl := nodeLevels[id]
			levelGroups[lvl] = append(levelGroups[lvl], id)
		}

		for lvl, ids := range levelGroups {
			sort.Strings(ids)
			x := startX + float64(lvl)*220.0
			for i, id := range ids {
				w, h := getNodeDimensions(nodeMap[id])
				y := startY + float64(i)*rowSpacing + (80.0-h)/2.0
				coords[id] = Coords{X: x, Y: y, Width: w, Height: h}
			}
		}
	}

	fillUnplacedNodes(coords, nodes, startX+3800.0, centerY, rowSpacing)
	resolveCollisions(coords, 25.0)
	applyCustomCoords(coords, opts.CustomCoords)
	return coords
}

// StrategyLinear implements a compact 3-Row Serpentine Wrap (S-Curve Snake) layout.
// Instead of stretching infinitely across 18 columns, it wraps every 6 columns into alternating rows,
// fitting the entire workflow into a compact 1400x750px viewport.
func StrategyLinear(nodes []map[string]interface{}, flows []map[string]interface{}, opts LayoutOptions) map[string]Coords {
	startX := opts.StartX
	if startX <= 0 {
		startX = 120.0
	}
	startY := opts.StartY
	if startY <= 0 {
		startY = 140.0
	}
	colSpacing := opts.ColSpacing
	if colSpacing <= 0 {
		colSpacing = 220.0
	}
	rowSpacing := 360.0

	outgoing, _, inCount := buildGraphAdj(nodes, flows)
	nodeLevels, _ := computeBFSLevels(nodes, outgoing, inCount)

	// Respect manual tiers if defined
	for id, tier := range opts.NodeTiers {
		nodeLevels[id] = tier
	}

	nodeMap := make(map[string]map[string]interface{})
	for _, n := range nodes {
		id, _ := n["id"].(string)
		nodeMap[id] = n
	}

	// Group nodes by their tier/level
	tierGroups := make(map[int][]string)
	for _, n := range nodes {
		id, _ := n["id"].(string)
		if id != "" {
			lvl := nodeLevels[id]
			tierGroups[lvl] = append(tierGroups[lvl], id)
		}
	}

	const colsPerRow = 6 // Wrap every 6 tiers into a new row
	coords := make(map[string]Coords)

	for tier, ids := range tierGroups {
		sort.Strings(ids)
		rowIdx := tier / colsPerRow
		colIdx := tier % colsPerRow

		// Calculate X based on direction: Even rows flow left-to-right, Odd rows flow right-to-left
		var x float64
		if rowIdx%2 == 0 {
			x = startX + float64(colIdx)*colSpacing
		} else {
			x = startX + float64(colsPerRow-1-colIdx)*colSpacing
		}

		baseY := startY + float64(rowIdx)*rowSpacing
		count := len(ids)

		for i, id := range ids {
			w, h := getNodeDimensions(nodeMap[id])
			// Compact vertical spacing for parallel nodes in the same tier
			y := baseY + (float64(i)-float64(count-1)/2.0)*65.0 + (80.0-h)/2.0
			coords[id] = Coords{
				X:      x,
				Y:      y,
				Width:  w,
				Height: h,
			}
		}
	}

	fillUnplacedNodes(coords, nodes, startX+float64(colsPerRow)*colSpacing, startY, 80.0)
	resolveCollisions(coords, 25.0)
	applyCustomCoords(coords, opts.CustomCoords)
	return coords
}
