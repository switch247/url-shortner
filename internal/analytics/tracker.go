package analytics

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"tinyurl-analytics/internal/db"
	"tinyurl-analytics/internal/models"

	_ "embed"

	"github.com/mssola/useragent"
	"github.com/oschwald/geoip2-golang"
)

//go:embed GeoLite2-City.mmdb
var geoDBBytes []byte // ← this will now contain the real file

type Tracker struct {
	reader *geoip2.Reader
}

func NewTracker() *Tracker {
	reader, err := geoip2.FromBytes(geoDBBytes)
	if err != nil {
		panic(err)
	}
	return &Tracker{reader: reader}
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return strings.Split(ip, ",")[0]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

func hashIP(ip string) string {
	h := sha1.Sum([]byte(ip + "salt"))
	return hex.EncodeToString(h[:])[:16]
}

func (t *Tracker) Track(r *http.Request, code string) {
	ip := getIP(r)
	ua := useragent.New(r.UserAgent())

	record, err := t.reader.City(net.ParseIP(ip))
	country := "Unknown"
	city := "Unknown"
	utms := r.URL.Query().Get("utm_source")
	log.Println("Tracking click from IP:", ip, "Country:", country, "City:", city, "UTMs:", utms)
	if err == nil {
		country = record.Country.IsoCode
		if record.City.Names != nil {
			city = record.City.Names["en"]
		}
	}

	click := models.Click{
		ShortCode: code,
		IP:        hashIP(ip),
		Country:   country,
		City:      city,
		UserAgent: r.UserAgent(),
		Referrer:  r.Referer(),
		OS:        ua.OS(),
		Device:    ua.Platform(),
		Browser:   func() string { n, _ := ua.Browser(); return n }(),
		Utms:      utms,
		CreatedAt: time.Now(),
	}

	go db.GORM.Create(&click)
	go db.Redis.Incr(context.Background(), "clickcount:"+code)
}
