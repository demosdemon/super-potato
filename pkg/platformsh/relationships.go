package platformsh

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type (
	Relationships map[string][]Relationship

	Relationship struct {
		Cluster  string     `json:"cluster"`
		Fragment string     `json:"fragment"`
		Host     string     `json:"host"`
		Hostname string     `json:"hostname"`
		IP       string     `json:"ip"`
		Password string     `json:"password"`
		Path     string     `json:"path"`
		Port     int        `json:"port"`
		Public   bool       `json:"public"`
		Query    JSONObject `json:"query"`
		Rel      string     `json:"rel"`
		Scheme   string     `json:"scheme"`
		Service  string     `json:"service"`
		SSL      JSONObject `json:"ssl"`
		Type     string     `json:"type"`
		Username string     `json:"username"`
	}
)

// ElasticSearch
// InfluxDB
// Kafka
// Memcached
// MongoDB
// MySQL
// PostgreSQL
// Redis
// RabbitMQ
// Solr

func (r Relationship) URL(user, query bool) string {
	var b strings.Builder
	b.WriteString(r.Scheme)
	b.WriteString("://")

	if user && r.Username != "" {
		b.WriteString(r.Username)
		if r.Password != "" {
			b.WriteString(":")
			b.WriteString(r.Password)
		}
		b.WriteString("@")
	}

	b.WriteString(r.Host)
	if r.Port > 0 {
		fmt.Fprintf(&b, ":%d", r.Port)
	}

	if r.Path != "" {
		b.WriteString("/")
		b.WriteString(r.Path)
	}

	if query {
		b.WriteString("?")
		first := true
		for k, v := range r.Query {
			if !first {
				b.WriteString("&")
			}
			first = false
			fmt.Fprintf(&b, "%s=%v", k, v)
		}
	}

	return b.String()
}

func (r Relationships) MongoDB(name string) (*mgo.Database, error) {
	rels, ok := r[name]
	if !ok {
		return nil, errors.New("missing relationship")
	}

	if len(rels) == 0 {
		return nil, errors.New("empty relationship")
	}

	hosts := make([]string, len(rels))
	for idx, rel := range rels {
		hosts[idx] = net.JoinHostPort(rel.Host, strconv.Itoa(rel.Port))
	}

	var b strings.Builder
	b.WriteString("mongodb://")
	b.WriteString(strings.Join(hosts, ","))
	b.WriteString("/")
	b.WriteString(rels[0].Path)

	url := b.String()

	mgo.SetLogger(log.New(os.Stderr, "mongo ", log.LstdFlags))

	for count := 0; count < 10; count++ {
		sess, err := mgo.Dial(url)
		if err == nil {
			db := sess.DB(rels[0].Path)
			if err := db.Login(rels[0].Username, rels[0].Password); err != nil {
				logrus.WithError(err).Warn("error logging into mongo database")
			}

			return db, nil
		}
		logrus.WithError(err).WithField("attempt", count+1).Warn("failed to connect to mongo server")
	}

	return nil, errors.New("failed to connect to mongo server")
}

func (r Relationships) Postgresql(name string) (string, error) {
	rels, ok := r[name]
	if !ok {
		return "", errors.New("missing relationship")
	}

	if len(rels) == 0 {
		return "", errors.New("empty relationship")
	}

	for len(rels) > 0 {
		rand.Shuffle(len(rels), func(i, j int) {
			rels[i], rels[j] = rels[j], rels[i]
		})

		dbURL := rels[0].URL(true, false)
		dbURL = strings.Replace(dbURL, "pgsql://", "postgresql://", 1)
		dbOpen, err := pq.ParseURL(dbURL)
		if err == nil {
			dbOpen += " sslmode=disable"
			return dbOpen, nil
		}

		rels = rels[1:]
	}

	return "", errors.New("error parsing postgres url")
}
