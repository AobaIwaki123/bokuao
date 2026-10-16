package schedule_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AobaIwaki123/bokuao/pkg/schedule"
)

func TestCheckIsLive(t *testing.T) {
	tests := []struct {
		title    string
		category string
		wantLive bool
	}{
		{
			title:    "「GIGA•GIGA SONIC」",
			category: "LIVE/EVENT",
			wantLive: true,
		},
		{
			title:    "『雲組～B♡BBLE～公演』第2回戦",
			category: "LIVE/EVENT",
			wantLive: true,
		},
		{
			title:    "僕青フリーライブ&握手会 【東京】ベルサール汐留",
			category: "LIVE/EVENT",
			wantLive: true,
		},
		{
			title:    "個別握手会【東京】新宿住友ホール",
			category: "LIVE/EVENT",
			wantLive: false,
		},
		{
			title:    "9th SINGLE発売記念 オンラインお話し会",
			category: "LIVE/EVENT",
			wantLive: false,
		},
		{
			title:    "9th SINGLE発売記念 オンラインサイン会",
			category: "LIVE/EVENT",
			wantLive: false,
		},
		{
			title:    "TBS「SASUKEアイドル予選会2026ー出場権争奪バトルー」",
			category: "LIVE/EVENT",
			wantLive: false,
		},
		{
			title:    "ミュージックステーション",
			category: "TV",
			wantLive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			got := schedule.CheckIsLive(tt.title, tt.category)
			if got != tt.wantLive {
				t.Errorf("CheckIsLive(%q, %q) = %v; want %v", tt.title, tt.category, got, tt.wantLive)
			}
		})
	}
}

func TestParseHTML(t *testing.T) {
	mockHTML := `
<ul class="list--schedule">
	<li data-delighter>
		<p class="date date--event">
			<strong>11.14</strong><span class="week">SAT</span>
		</p>
		<div class="list__txt">
			<a href="/schedule/detail/5986">
				<p class="category category-8">LIVE/EVENT</p>
				<p class="time">13:00<span>〜</span>18:00</p>
				<p class="tit">僕青フリーライブ&amp;握手会 【東京】ベルサール汐留</p>
			</a>
		</div>
	</li>
	<li data-delighter>
		<p class="date date--event">
			<strong>11.15</strong><span class="week">SUN</span>
		</p>
		<div class="list__txt">
			<a href="/schedule/detail/6008">
				<p class="category category-8">LIVE/EVENT</p>
				<p class="time">00:00<span>〜</span>23:59</p>
				<p class="tit">オンラインお話し会</p>
			</a>
		</div>
	</li>
</ul>
`

	events, err := schedule.ParseHTML(mockHTML, "https://bokuao.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	ev1 := events[0]
	if ev1.Date != "11.14" || ev1.Week != "SAT" || ev1.Category != "LIVE/EVENT" {
		t.Errorf("unexpected event 1 metadata: %+v", ev1)
	}
	if !ev1.IsLive {
		t.Errorf("expected event 1 to be marked as live")
	}
	if ev1.URL != "https://bokuao.com/schedule/detail/5986" {
		t.Errorf("unexpected URL: %s", ev1.URL)
	}

	ev2 := events[1]
	if ev2.IsLive {
		t.Errorf("expected event 2 (online chat) NOT to be marked as live")
	}
}

func TestClient_Fetch(t *testing.T) {
	mockHTML := `
<ul class="list--schedule">
	<li data-delighter>
		<p class="date"><strong>11.26</strong><span class="week">THU</span></p>
		<a href="/schedule/detail/5970">
			<p class="category category-8">LIVE/EVENT</p>
			<p class="tit">「どぼフェス」</p>
		</a>
	</li>
</ul>`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/schedule/list/2026/11/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(mockHTML))
	}))
	defer srv.Close()

	client := &schedule.Client{
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
	}

	events, err := client.Fetch(context.Background(), 2026, 11)
	if err != nil {
		t.Fatalf("fetch error: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Title != "「どぼフェス」" || !events[0].IsLive {
		t.Errorf("unexpected event: %+v", events[0])
	}
}
