package rss

import (
	"testing"
)

func TestRssToStruct(t *testing.T) {
	ups, err := RssToStruct("https://habr.com/ru/rss/hub/go/all/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	if len(ups) == 0 {
		t.Fatal("данные не раскодированы")
	}
	t.Logf("получено %d новостей\n%+v", len(ups), ups)

	ups, err = RssToStruct("https://habr.com/ru/rss/best/daily/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	if len(ups) == 0 {
		t.Fatal("данные не раскодированы")
	}
	t.Logf("получено %d новостей\n%+v", len(ups), ups)
}
