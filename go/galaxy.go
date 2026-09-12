package nativebpm

import (
	"encoding/json"
	"fmt"
	"math"
)

// DomainType represents functional clusters in the distributed architecture.
type DomainType string

const (
	DomainEdge           DomainType = "edge"
	DomainCompute        DomainType = "compute"
	DomainCache          DomainType = "cache"
	DomainSQLCoordinator DomainType = "sql_coordinator"
	DomainStorageRaft    DomainType = "storage_raft"
)

// CelestialType represents the astronomical classification of a node.
type CelestialType string

const (
	BodyStar         CelestialType = "star"
	BodySupernova    CelestialType = "supernova"
	BodyPlanet       CelestialType = "planet"
	BodyMoon         CelestialType = "moon"
	BodyRaftLeader   CelestialType = "raft_leader"
	BodyRaftFollower CelestialType = "raft_follower"
)

// Standard Icon constants for Planetary visualization
const (
	// BPMN Semantics
	IconBPMNServiceTask = "bpmn-service-task"
	IconBPMNUserTask    = "bpmn-user-task"
	IconBPMNScriptTask  = "bpmn-script-task"
	IconBPMNGateway     = "bpmn-gateway"
	IconBPMNStartEvent  = "bpmn-start-event"
	IconBPMNEndEvent    = "bpmn-end-event"

	// Specialized Distributed Database & Storage Icons
	IconDBTiKV      = "db-tikv-raft"
	IconDBTiDB      = "db-tidb-sql"
	IconDBYugabyte  = "db-yugabyte"
	IconDBPostgres  = "db-postgres"
	IconDBCachedKey = "db-keydb-mesh"
	IconDBCylinder  = "db-cylinder"

	// Specialized Network Topology & Ingress Icons
	IconNetCloudflare   = "net-cloudflare-edge"
	IconNetArgoTunnel   = "net-argo-quic"
	IconNetLoadBalancer = "net-load-balancer"
	IconNetGateway      = "net-gateway"
	IconNetMeshRouter   = "net-mesh-router"

	// Specialized Compute & FastCGI Workers
	IconComputeFastCGI = "compute-fastcgi"
	IconComputeWorker  = "compute-worker"
	IconComputeDaemon  = "compute-daemon"
	IconComputeDocker  = "compute-container"
)

// TelemetryData encapsulates runtime operational metrics for HUD rendering.
type TelemetryData struct {
	IP       string  `json:"ip,omitempty"`
	Port     int     `json:"port,omitempty"`
	PingMs   float64 `json:"pingMs,omitempty"`
	CPU      float64 `json:"cpu,omitempty"`
	RAM      float64 `json:"ram,omitempty"`
	IOPS     int     `json:"iops,omitempty"`
	Status   string  `json:"status,omitempty"`
	QPS      float64 `json:"qps,omitempty"`
	Socket   string  `json:"socket,omitempty"`
	RegionDC string  `json:"regionDC,omitempty"`
}

// CelestialBody models a node in the planetary cosmos.
type CelestialBody struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Type         CelestialType    `json:"type"`
	Role         string           `json:"role"`
	SubRole      string           `json:"subRole,omitempty"`
	RegionDC     string           `json:"regionDc,omitempty"`
	X            float64          `json:"x"`
	Y            float64          `json:"y"`
	Z            float64          `json:"z,omitempty"`
	OrbitRadius  float64          `json:"orbitRadius"`
	OrbitSpeed   float64          `json:"orbitSpeed"`
	CurrentAngle float64          `json:"currentAngle"`
	Size         float64          `json:"size"`
	Mass         float64          `json:"mass"`
	Color        string           `json:"color"`
	GlowColor    string           `json:"glowColor"`
	Icon         string           `json:"icon,omitempty"`
	IconCategory string           `json:"iconCategory,omitempty"`
	CustomBadge  string           `json:"customBadge,omitempty"`
	Telemetry    TelemetryData    `json:"telemetry"`
	Satellites   []*CelestialBody `json:"satellites,omitempty"`
}

// OrbitalTrack represents a Keplerian harmonic orbit containing celestial bodies.
type OrbitalTrack struct {
	Radius float64          `json:"radius"`
	Speed  float64          `json:"speed"`
	Bodies []*CelestialBody `json:"bodies"`
}

// DataRegion represents a TiDB/TiKV distributed data region/shard.
type DataRegion struct {
	RegionID      uint64   `json:"regionId"`
	StartKey      string   `json:"startKey"`
	EndKey        string   `json:"endKey"`
	LeaderStoreID string   `json:"leaderStoreId"`
	PeerStores    []string `json:"peerStores"`
	SizeMB        float64  `json:"sizeMb"`
	ReadQPS       float64  `json:"readQps"`
	WriteQPS      float64  `json:"writeQps"`
	Angle         float64  `json:"angle"`
	Distance      float64  `json:"distance"`
	Speed         float64  `json:"speed"`
	IsHotspot     bool     `json:"isHotspot"`
}

// DataRegionBelt represents an asteroid/comet belt of distributed data shards.
type DataRegionBelt struct {
	TotalCount  int           `json:"totalCount"`
	SizeMB      float64       `json:"sizeMb"`
	StartRadius float64       `json:"startRadius"`
	EndRadius   float64       `json:"endRadius"`
	BaseSpeed   float64       `json:"baseSpeed"`
	Regions     []*DataRegion `json:"regions"`
}

// QuorumShield represents a Raft consensus energy perimeter.
type QuorumShield struct {
	MinQuorum     int     `json:"minQuorum"`
	TotalNodes    int     `json:"totalNodes"`
	Active        bool    `json:"active"`
	Radius        float64 `json:"radius"`
	PulseRate     float64 `json:"pulseRate"`
	EnergyColor   string  `json:"energyColor"`
	QuorumReached bool    `json:"quorumReached"`
}

// Photon represents a light pulse traveling across an interstellar hyperlane.
type Photon struct {
	Progress float64 `json:"progress"`
	Speed    float64 `json:"speed"`
	Size     float64 `json:"size"`
	Color    string  `json:"color"`
}

// Hyperlane represents an interstellar gravitational communication channel.
type Hyperlane struct {
	ID            string    `json:"id"`
	FromID        string    `json:"fromId"`
	ToID          string    `json:"toId"`
	Protocol      string    `json:"protocol"`
	LatencyMs     float64   `json:"latencyMs"`
	TelemetryRate float64   `json:"telemetryRate"`
	Color         string    `json:"color"`
	ActivePhotons []*Photon `json:"activePhotons,omitempty"`
}

// StarSystem represents a functional domain (e.g. Ingress, Compute, TiKV Storage).
type StarSystem struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	DomainType   DomainType      `json:"domainType"`
	CenterX      float64         `json:"centerX"`
	CenterY      float64         `json:"centerY"`
	CenterStar   *CelestialBody  `json:"centerStar"`
	Orbits       []*OrbitalTrack `json:"orbits"`
	RegionBelt   *DataRegionBelt `json:"regionBelt,omitempty"`
	QuorumShield *QuorumShield   `json:"quorumShield,omitempty"`
}

// GalaxyParams controls the simulation and universe boundaries.
type GalaxyParams struct {
	Width            float64 `json:"width"`
	Height           float64 `json:"height"`
	MinOrbitalDist   float64 `json:"minOrbitalDist"`
	GravitationalK   float64 `json:"gravitationalK"`
	AllowCollisions  bool    `json:"allowCollisions"`
	TimeAcceleration float64 `json:"timeAcceleration"`
}

// GalaxyUniverse is the root domain model for the cosmic cluster visualization.
type GalaxyUniverse struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	StarSystems []*StarSystem `json:"starSystems"`
	Hyperlanes  []*Hyperlane  `json:"hyperlanes"`
	Params      GalaxyParams  `json:"params"`
	Timestamp   int64         `json:"timestamp,omitempty"`
}

// GalaxyBuilder provides a fluent API to define distributed planetary galaxies.
type GalaxyBuilder struct {
	universe *GalaxyUniverse
}

// NewGalaxy creates a new root galaxy builder.
func NewGalaxy(id, name string) *GalaxyBuilder {
	u := &GalaxyUniverse{
		ID:          id,
		Name:        name,
		StarSystems: make([]*StarSystem, 0),
		Hyperlanes:  make([]*Hyperlane, 0),
		Params: GalaxyParams{
			Width:            2400.0,
			Height:           1600.0,
			MinOrbitalDist:   50.0,
			GravitationalK:   1.0,
			AllowCollisions:  false,
			TimeAcceleration: 1.0,
		},
	}
	return &GalaxyBuilder{universe: u}
}

// StarSystem adds a new star system to the universe.
func (g *GalaxyBuilder) StarSystem(id, name string) *StarSystemBuilder {
	sys := &StarSystem{
		ID:     id,
		Name:   name,
		Orbits: make([]*OrbitalTrack, 0),
	}
	g.universe.StarSystems = append(g.universe.StarSystems, sys)
	return &StarSystemBuilder{
		galaxy: g,
		system: sys,
	}
}

// Hyperlane creates an interstellar communication conduit.
func (g *GalaxyBuilder) Hyperlane(id, from, to string) *HyperlaneBuilder {
	lane := &Hyperlane{
		ID:            id,
		FromID:        from,
		ToID:          to,
		Protocol:      "TCP/TLS",
		LatencyMs:     1.5,
		TelemetryRate: 1000.0,
		Color:         "#a855f7",
	}
	g.universe.Hyperlanes = append(g.universe.Hyperlanes, lane)
	return &HyperlaneBuilder{
		galaxy: g,
		lane:   lane,
	}
}

// Build computes coordinates and returns finalized GalaxyUniverse.
func (g *GalaxyBuilder) Build() *GalaxyUniverse {
	g.computeLayout()
	return g.universe
}

// ToJSON serializes the universe into JSON.
func (g *GalaxyBuilder) ToJSON() ([]byte, error) {
	u := g.Build()
	return json.MarshalIndent(u, "", "  ")
}

func (g *GalaxyBuilder) computeLayout() {
	positions := map[DomainType][2]float64{
		DomainEdge:           {450.0, 380.0},
		DomainCompute:        {1050.0, 420.0},
		DomainCache:          {1700.0, 420.0},
		DomainSQLCoordinator: {850.0, 1050.0},
		DomainStorageRaft:    {1550.0, 1100.0},
	}

	for _, s := range g.universe.StarSystems {
		if s.CenterX == 0 && s.CenterY == 0 {
			if pos, ok := positions[s.DomainType]; ok {
				s.CenterX = pos[0]
				s.CenterY = pos[1]
			}
		}

		if s.CenterStar != nil {
			s.CenterStar.X = s.CenterX
			s.CenterStar.Y = s.CenterY
		}

		lastRadius := 50.0
		if s.CenterStar != nil {
			lastRadius = s.CenterStar.Size + 25.0
		}

		for _, track := range s.Orbits {
			if track.Radius <= lastRadius {
				track.Radius = lastRadius + g.universe.Params.MinOrbitalDist
			}
			lastRadius = track.Radius

			numBodies := len(track.Bodies)
			if numBodies == 0 {
				continue
			}

			step := (2.0 * math.Pi) / float64(numBodies)
			for i, body := range track.Bodies {
				if body.CurrentAngle == 0 {
					body.CurrentAngle = float64(i) * step
				}
				body.OrbitRadius = track.Radius
				if body.OrbitSpeed == 0 {
					body.OrbitSpeed = track.Speed
				}
				body.X = s.CenterX + track.Radius*math.Cos(body.CurrentAngle)
				body.Y = s.CenterY + track.Radius*math.Sin(body.CurrentAngle)
			}
		}

		// Initialize Region Belt
		if s.RegionBelt != nil && s.RegionBelt.TotalCount > 0 && len(s.RegionBelt.Regions) == 0 {
			belt := s.RegionBelt
			if belt.StartRadius <= lastRadius {
				belt.StartRadius = lastRadius + 30.0
			}
			if belt.EndRadius <= belt.StartRadius {
				belt.EndRadius = belt.StartRadius + 90.0
			}
			if belt.BaseSpeed == 0 {
				belt.BaseSpeed = 0.008
			}

			radiusSpan := belt.EndRadius - belt.StartRadius
			goldenRatio := 1.61803398875
			count := belt.TotalCount
			if count > 2000 {
				count = 2000
			}

			belt.Regions = make([]*DataRegion, count)
			for i := 0; i < count; i++ {
				angle := float64(i) * goldenRatio * 2.0 * math.Pi
				distFraction := math.Mod(float64(i)*goldenRatio, 1.0)
				dist := belt.StartRadius + distFraction*radiusSpan
				isHotspot := (i%47 == 0)

				belt.Regions[i] = &DataRegion{
					RegionID:  uint64(100000 + i),
					StartKey:  fmt.Sprintf("t_75_r_%08d", i*1000),
					EndKey:    fmt.Sprintf("t_75_r_%08d", (i+1)*1000),
					Angle:     angle,
					Distance:  dist,
					Speed:     belt.BaseSpeed * (1.0 + 0.3*(distFraction-0.5)),
					SizeMB:    belt.SizeMB,
					IsHotspot: isHotspot,
				}
			}
		}
	}

	for _, lane := range g.universe.Hyperlanes {
		if len(lane.ActivePhotons) == 0 {
			lane.ActivePhotons = []*Photon{
				{Progress: 0.1, Speed: 0.008, Size: 3.5, Color: "#ffffff"},
				{Progress: 0.5, Speed: 0.008, Size: 4.0, Color: lane.Color},
				{Progress: 0.85, Speed: 0.008, Size: 3.0, Color: "#38bdf8"},
			}
		}
	}
}

// StarSystemBuilder configures celestial bodies within a star system.
type StarSystemBuilder struct {
	galaxy       *GalaxyBuilder
	system       *StarSystem
	currentOrbit *OrbitBuilder
}

// WithDomainType sets the domain type of the star system.
func (s *StarSystemBuilder) WithDomainType(dt DomainType) *StarSystemBuilder {
	s.system.DomainType = dt
	return s
}

// AsSupernova designates the center node as an exploding energy core (e.g. Cloudflare Anycast).
func (s *StarSystemBuilder) AsSupernova(id, name string) *StarSystemBuilder {
	s.system.CenterStar = &CelestialBody{
		ID:           id,
		Name:         name,
		Type:         BodySupernova,
		Role:         "Anycast Ingress Gateway",
		Size:         34.0,
		Color:        "#f97316",
		GlowColor:    "rgba(249, 115, 22, 0.65)",
		Icon:         IconNetCloudflare,
		IconCategory: "network",
		CustomBadge:  "CF",
		Telemetry: TelemetryData{
			Status: "healthy",
			QPS:    45000.0,
			PingMs: 0.8,
		},
	}
	return s
}

// AsRaftLeader designates the center node as a consensus master leader (e.g. TiKV Leader).
func (s *StarSystemBuilder) AsRaftLeader(id, name string) *StarSystemBuilder {
	s.system.CenterStar = &CelestialBody{
		ID:           id,
		Name:         name,
		Type:         BodyRaftLeader,
		Role:         "Raft Master Leader",
		Size:         32.0,
		Color:        "#10b981",
		GlowColor:    "rgba(16, 185, 129, 0.7)",
		Icon:         IconDBTiKV,
		IconCategory: "database",
		CustomBadge:  "RAFT",
		Telemetry: TelemetryData{
			Status: "leader",
			Port:   20160,
			QPS:    12500.0,
			PingMs: 0.4,
		},
	}
	return s
}

// AsStar designates the center node as a standard bright star.
func (s *StarSystemBuilder) AsStar(id, name string) *StarSystemBuilder {
	s.system.CenterStar = &CelestialBody{
		ID:           id,
		Name:         name,
		Type:         BodyStar,
		Role:         "Coordinator Core",
		Size:         28.0,
		Color:        "#38bdf8",
		GlowColor:    "rgba(56, 189, 248, 0.6)",
		Icon:         IconNetGateway,
		IconCategory: "network",
		CustomBadge:  "CORE",
		Telemetry: TelemetryData{
			Status: "healthy",
			Port:   4000,
			QPS:    8000.0,
		},
	}
	return s
}

// WithIcon sets a custom icon on the system center star.
func (s *StarSystemBuilder) WithIcon(icon string) *StarSystemBuilder {
	if s.system.CenterStar != nil {
		s.system.CenterStar.Icon = icon
	}
	return s
}

// WithCustomBadge sets a custom text badge on the center star.
func (s *StarSystemBuilder) WithCustomBadge(badge string) *StarSystemBuilder {
	if s.system.CenterStar != nil {
		s.system.CenterStar.CustomBadge = badge
	}
	return s
}

// WithQuorumShield configures a Raft consensus perimeter around the system.
func (s *StarSystemBuilder) WithQuorumShield(minQuorum, totalNodes int) *StarSystemBuilder {
	s.system.QuorumShield = &QuorumShield{
		MinQuorum:     minQuorum,
		TotalNodes:    totalNodes,
		Active:        true,
		EnergyColor:   "rgba(16, 185, 129, 0.5)",
		PulseRate:     1.5,
		QuorumReached: true,
	}
	return s
}

// WithOrbit adds a concentric planetary orbital ring.
func (s *StarSystemBuilder) WithOrbit(radius, speed float64) *OrbitBuilder {
	track := &OrbitalTrack{
		Radius: radius,
		Speed:  speed,
		Bodies: make([]*CelestialBody, 0),
	}
	s.system.Orbits = append(s.system.Orbits, track)
	ob := &OrbitBuilder{
		systemBuilder: s,
		track:         track,
	}
	s.currentOrbit = ob
	return ob
}

// WithRegionAsteroidBelt attaches a swarm of TiDB/TiKV distributed data shards.
func (s *StarSystemBuilder) WithRegionAsteroidBelt(count int, sizeMB float64) *StarSystemBuilder {
	s.system.RegionBelt = &DataRegionBelt{
		TotalCount: count,
		SizeMB:     sizeMB,
	}
	return s
}

// StarSystem chains to define another star system.
func (s *StarSystemBuilder) StarSystem(id, name string) *StarSystemBuilder {
	return s.galaxy.StarSystem(id, name)
}

// Hyperlane chains to hyperlane definition.
func (s *StarSystemBuilder) Hyperlane(id, from, to string) *HyperlaneBuilder {
	return s.galaxy.Hyperlane(id, from, to)
}

// Build delegates to root galaxy builder.
func (s *StarSystemBuilder) Build() *GalaxyUniverse {
	return s.galaxy.Build()
}

// ToJSON delegates to root galaxy builder.
func (s *StarSystemBuilder) ToJSON() ([]byte, error) {
	return s.galaxy.ToJSON()
}

// OrbitBuilder configures planetary nodes on a specific orbit.
type OrbitBuilder struct {
	systemBuilder *StarSystemBuilder
	track         *OrbitalTrack
	lastBody      *CelestialBody
}

// Planet adds a major node to the orbit.
func (o *OrbitBuilder) Planet(id, name, region string) *OrbitBuilder {
	body := &CelestialBody{
		ID:          id,
		Name:        name,
		Type:        BodyPlanet,
		Role:        "Node Replica",
		RegionDC:    region,
		OrbitRadius: o.track.Radius,
		OrbitSpeed:  o.track.Speed,
		Size:        16.0,
		Color:       "#60a5fa",
		GlowColor:   "rgba(96, 165, 250, 0.5)",
		Telemetry: TelemetryData{
			Status:   "healthy",
			RegionDC: region,
			PingMs:   0.7,
			CPU:      24.5,
			RAM:      62.0,
			IOPS:     1800,
		},
		Satellites: make([]*CelestialBody, 0),
	}
	o.track.Bodies = append(o.track.Bodies, body)
	o.lastBody = body
	return o
}

// WithIcon sets a custom vector/SVG icon key for the last added planet.
func (o *OrbitBuilder) WithIcon(icon string) *OrbitBuilder {
	if o.lastBody != nil {
		o.lastBody.Icon = icon
	}
	return o
}

// WithCustomBadge sets a custom text badge on the last added planet.
func (o *OrbitBuilder) WithCustomBadge(badge string) *OrbitBuilder {
	if o.lastBody != nil {
		o.lastBody.CustomBadge = badge
	}
	return o
}

// AsDatabaseNode configures the planet as a specialized database node (TiKV, TiDB, Yugabyte, Postgres, etc.).
func (o *OrbitBuilder) AsDatabaseNode(dbType string) *OrbitBuilder {
	if o.lastBody != nil {
		o.lastBody.Icon = dbType
		o.lastBody.IconCategory = "database"
		o.lastBody.CustomBadge = "DB"
		o.lastBody.Color = "#10b981"
		o.lastBody.GlowColor = "rgba(16, 185, 129, 0.55)"
		o.lastBody.Role = "Distributed Database Store"
	}
	return o
}

// AsNetworkNode configures the planet as a network topology node (Argo Tunnel, QUIC, Gateway, Load Balancer).
func (o *OrbitBuilder) AsNetworkNode(netType string) *OrbitBuilder {
	if o.lastBody != nil {
		o.lastBody.Icon = netType
		o.lastBody.IconCategory = "network"
		o.lastBody.CustomBadge = "NET"
		o.lastBody.Color = "#38bdf8"
		o.lastBody.GlowColor = "rgba(56, 189, 248, 0.55)"
		o.lastBody.Role = "Network Ingress Tunnel"
	}
	return o
}

// AsComputeNode configures the planet as a compute worker node (FastCGI, PHP-FPM, Go worker, Docker).
func (o *OrbitBuilder) AsComputeNode(computeType string) *OrbitBuilder {
	if o.lastBody != nil {
		o.lastBody.Icon = computeType
		o.lastBody.IconCategory = "compute"
		o.lastBody.CustomBadge = "FPM"
		o.lastBody.Color = "#f59e0b"
		o.lastBody.GlowColor = "rgba(245, 158, 11, 0.55)"
		o.lastBody.Role = "Compute Worker Pool"
	}
	return o
}

// AsCacheNode configures the planet as an in-memory cache/mesh node (KeyDB, Redis).
func (o *OrbitBuilder) AsCacheNode(cacheType string) *OrbitBuilder {
	if o.lastBody != nil {
		o.lastBody.Icon = cacheType
		o.lastBody.IconCategory = "cache"
		o.lastBody.CustomBadge = "MEM"
		o.lastBody.Color = "#ec4899"
		o.lastBody.GlowColor = "rgba(236, 72, 153, 0.55)"
		o.lastBody.Role = "In-Memory Replicated Store"
	}
	return o
}

// AsBPMNTask configures the planet to render standard BPMN task semantics.
func (o *OrbitBuilder) AsBPMNTask(taskType string) *OrbitBuilder {
	if o.lastBody != nil {
		o.lastBody.Icon = taskType
		o.lastBody.IconCategory = "bpmn"
		o.lastBody.CustomBadge = "TASK"
	}
	return o
}

// Moon adds a satellite to the most recently added planet.
func (o *OrbitBuilder) Moon(id, name string) *OrbitBuilder {
	if o.lastBody == nil {
		return o
	}
	moon := &CelestialBody{
		ID:        id,
		Name:      name,
		Type:      BodyMoon,
		Role:      "Satellite Worker",
		Size:      7.0,
		Color:     "#cbd5e1",
		GlowColor: "rgba(203, 213, 225, 0.4)",
		Telemetry: TelemetryData{
			Status: "healthy",
			CPU:    15.0,
		},
	}
	o.lastBody.Satellites = append(o.lastBody.Satellites, moon)
	return o
}

// WithOrbit chains to another orbit in the current system.
func (o *OrbitBuilder) WithOrbit(radius, speed float64) *OrbitBuilder {
	return o.systemBuilder.WithOrbit(radius, speed)
}

// WithRegionAsteroidBelt attaches data regions to the parent system.
func (o *OrbitBuilder) WithRegionAsteroidBelt(count int, sizeMB float64) *StarSystemBuilder {
	return o.systemBuilder.WithRegionAsteroidBelt(count, sizeMB)
}

// StarSystem chains to another star system.
func (o *OrbitBuilder) StarSystem(id, name string) *StarSystemBuilder {
	return o.systemBuilder.StarSystem(id, name)
}

// Hyperlane chains to hyperlane definition.
func (o *OrbitBuilder) Hyperlane(id, from, to string) *HyperlaneBuilder {
	return o.systemBuilder.Hyperlane(id, from, to)
}

// Build delegates to root galaxy builder.
func (o *OrbitBuilder) Build() *GalaxyUniverse {
	return o.systemBuilder.Build()
}

// ToJSON delegates to root galaxy builder.
func (o *OrbitBuilder) ToJSON() ([]byte, error) {
	return o.systemBuilder.ToJSON()
}

// HyperlaneBuilder configures interstellar channels.
type HyperlaneBuilder struct {
	galaxy *GalaxyBuilder
	lane   *Hyperlane
}

// WithProtocol specifies the transport protocol.
func (h *HyperlaneBuilder) WithProtocol(protocol string) *HyperlaneBuilder {
	h.lane.Protocol = protocol
	return h
}

// WithTelemetryRate sets packet transfer rate.
func (h *HyperlaneBuilder) WithTelemetryRate(rate float64) *HyperlaneBuilder {
	h.lane.TelemetryRate = rate
	return h
}

// WithLatency sets channel latency in ms.
func (h *HyperlaneBuilder) WithLatency(latency float64) *HyperlaneBuilder {
	h.lane.LatencyMs = latency
	return h
}

// WithColor sets the laser stream color.
func (h *HyperlaneBuilder) WithColor(color string) *HyperlaneBuilder {
	h.lane.Color = color
	return h
}

// StarSystem chains to a star system.
func (h *HyperlaneBuilder) StarSystem(id, name string) *StarSystemBuilder {
	return h.galaxy.StarSystem(id, name)
}

// Hyperlane chains to define another hyperlane.
func (h *HyperlaneBuilder) Hyperlane(id, from, to string) *HyperlaneBuilder {
	return h.galaxy.Hyperlane(id, from, to)
}

// Build delegates to root galaxy builder.
func (h *HyperlaneBuilder) Build() *GalaxyUniverse {
	return h.galaxy.Build()
}

// ToJSON delegates to root galaxy builder.
func (h *HyperlaneBuilder) ToJSON() ([]byte, error) {
	return h.galaxy.ToJSON()
}
