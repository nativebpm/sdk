package nativebpm

import (
	"encoding/json"
	"testing"
)

func TestSDK_GalaxyBuilder(t *testing.T) {
	g := NewGalaxy("sre_cluster", "Production SRE Galactic Cluster")

	// Edge Galaxy
	g.StarSystem("edge_galaxy", "Cloudflare Edge").
		WithDomainType(DomainEdge).
		AsSupernova("cf_anycast", "Cloudflare Anycast Core").
		WithOrbit(140.0, 0.02).
		Planet("tunnel_1", "Argo QUIC Hamina", "fin-1").AsNetworkNode(IconNetArgoTunnel).
		Planet("tunnel_2", "Argo QUIC Hamina", "fin-2").AsNetworkNode(IconNetArgoTunnel)

	// Compute Galaxy
	g.StarSystem("compute_galaxy", "FastCGI Compute Cluster").
		WithDomainType(DomainCompute).
		AsStar("nginx_gw", "Nginx L7 Ingress Gateway").
		WithOrbit(120.0, -0.025).
		Planet("worker_1", "PHP-FPM Worker Pool 1", "fin-1").AsComputeNode(IconComputeFastCGI).Moon("fpm_child_1", "Worker-01")

	// Cache Galaxy
	g.StarSystem("cache_galaxy", "KeyDB Multi-Master Mesh").
		WithDomainType(DomainCache).
		AsStar("keydb_seed", "KeyDB Mesh Coordinator").
		WithOrbit(130.0, 0.018).
		Planet("keydb_1", "KeyDB Master 1", "fin-1").AsCacheNode(IconDBCachedKey)

	// Storage Galaxy
	g.StarSystem("storage_galaxy", "TiKV Multi-Raft & TiDB Regions").
		WithDomainType(DomainStorageRaft).
		AsRaftLeader("tikv_fin_1", "TiKV Master Leader :20160").
		WithQuorumShield(3, 5).
		WithOrbit(220.0, 0.015).
		Planet("tikv_fin_2", "TiKV Follower", "fin-2").AsDatabaseNode(IconDBTiKV).
		WithRegionAsteroidBelt(500, 96.0)

	// Hyperlanes
	g.Hyperlane("hl_edge_compute", "edge_galaxy", "compute_galaxy").
		WithProtocol("QUIC :7844 (0 WAN Drop)").
		WithTelemetryRate(15000.0).
		WithLatency(0.4)

	universe := g.Build()
	if universe.ID != "sre_cluster" {
		t.Fatalf("expected 'sre_cluster', got %s", universe.ID)
	}

	if len(universe.StarSystems) != 4 {
		t.Fatalf("expected 4 star systems, got %d", len(universe.StarSystems))
	}

	jsonBytes, err := g.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	var parsed GalaxyUniverse
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if len(parsed.StarSystems) != 4 {
		t.Fatalf("expected 4 parsed star systems, got %d", len(parsed.StarSystems))
	}

	// Verify custom icons
	var foundTiKV bool
	for _, sys := range parsed.StarSystems {
		for _, track := range sys.Orbits {
			for _, body := range track.Bodies {
				if body.ID == "tikv_fin_2" && body.Icon == IconDBTiKV && body.CustomBadge == "DB" {
					foundTiKV = true
				}
			}
		}
	}

	if !foundTiKV {
		t.Fatal("expected tikv_fin_2 to have IconDBTiKV and DB badge")
	}

	t.Log("SDK GalaxyBuilder verified successfully!")
}
