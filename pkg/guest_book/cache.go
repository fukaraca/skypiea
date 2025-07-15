package guest_book

import (
	"iter"
	"sync"
	"time"

	"github.com/fukaraca/skypiea/pkg/cache"
)

var GuestBook *VisitorMetric

// VisitorMetric stores metrics for visitors
type VisitorMetric struct {
	storage *cache.Storage
}

type Stats struct {
	mu                  sync.Mutex
	TotalHits           int
	Endpoints           map[string]int
	FirstSeen, LastSeen time.Time
}

func New() *VisitorMetric {
	return &VisitorMetric{
		storage: cache.New(),
	}
}

func (m *VisitorMetric) RegisterGuest(ip, endpoint string) {
	stat := m.GetStat(ip)

	stat.mu.Lock()
	defer stat.mu.Unlock()

	stat.Endpoints[endpoint]++
	stat.LastSeen = time.Now()
	stat.TotalHits++
}

// FIXME: not safe against same ip origined requests but we can live with it
func (m *VisitorMetric) GetStat(ip string) *Stats {
	var stat *Stats
	stat, ok := m.storage.Get(ip).(*Stats)
	if !ok {
		stat = &Stats{
			Endpoints: make(map[string]int),
			FirstSeen: time.Now(),
		}
		m.storage.Set(ip, stat)
	}
	return stat
}

func (m *VisitorMetric) IPs() iter.Seq[string] {
	return m.storage.Keys()
}

type Visitor struct {
	IP        string
	FirstSeen time.Time
	LastSeen  time.Time
	TotalHits int
	Endpoints map[string]int // endpoint → hit count
}

func (m *VisitorMetric) DumpVisitorMetric() []*Visitor {
	out := make([]*Visitor, 0, m.storage.Len())
	for s := range m.IPs() {
		stat := m.GetStat(s)
		out = append(out, &Visitor{
			IP:        s,
			FirstSeen: stat.FirstSeen,
			LastSeen:  stat.LastSeen,
			TotalHits: stat.TotalHits,
			Endpoints: stat.Endpoints,
		})
	}
	return out
}
