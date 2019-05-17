package server

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/base32"
	"net/http"
	"sync"
	"time"

	"github.com/cloudflare/cfssl/certdb"
	"github.com/cloudflare/cfssl/certdb/sql"
	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/ocsp"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/mongo"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/demosdemon/super-potato/pkg/app"
	"github.com/demosdemon/super-potato/pkg/pki"
)

const (
	OneYear  = 60 * 60 * 24 * 365
	OneMonth = 60 * 60 * 24 * 30
)

type Server struct {
	*app.App      `flag:"-"`
	SessionCookie string `flag:"session-cookie" desc:"The name of the session cookie." env:"PKI_SESSION_COOKIE"`

	once       sync.Once
	start      time.Time
	engine     *gin.Engine
	db         *sqlx.DB
	accessor   certdb.Accessor
	signer     signer.Signer
	ocspSigner ocsp.Signer
}

func (s *Server) Use() string {
	return "serve"
}

func (s *Server) Args(cmd *cobra.Command, args []string) error {
	return cobra.NoArgs(cmd, args)
}

func (s *Server) Run(cmd *cobra.Command, args []string) error {
	return s.Serve()
}

func (s *Server) Init() {
	s.once.Do(s.init)
}

func (s *Server) init() {
	s.start = time.Now().Truncate(time.Second)
	s.engine = gin.New()

	var err error

	s.db, err = s.getDB()
	if err != nil {
		logrus.WithError(err).Panic("unable to get database connection")
	}

	s.accessor = sql.NewAccessor(s.db)
	s.signer, err = s.getSigner()
	if err != nil {
		logrus.WithError(err).Panic("unable to get certificate signer")
	}
	s.signer.SetDBAccessor(s.accessor)

	s.ocspSigner, err = s.getOCSPSigner()
	if err != nil {
		logrus.WithError(err).Panic("unable to get OCSP signer")
	}

	s.register(s.engine)
}

func (s *Server) getDB() (*sqlx.DB, error) {
	rels, err := s.Relationships()
	if err != nil {
		return nil, errors.Wrap(err, "unable to locate relationships")
	}

	dbOpen, err := rels.Postgresql("database")
	if err != nil {
		return nil, errors.Wrap(err, "unable to get database connection string")
	}

	db, err := sqlx.Open("postgres", dbOpen)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to postgres")
	}

	return db, nil
}

func (s *Server) getSigner() (signer.Signer, error) {
	rootPem, ok := s.Lookup("PKI_ROOT_CERTIFICATE")
	if !ok {
		return nil, errors.New("PKI_ROOT_CERTIFICATE not found in environment")
	}

	intermediatePem, ok := s.Lookup("PKI_INTERMEDIATE_CERTIFICATE")
	if !ok {
		return nil, errors.New("PKI_INTERMEDIATE_CERTIFICATE not found in environment")
	}

	intermediateKeyPem, ok := s.Lookup("PKI_INTERMEDIATE_PRIVATE_KEY")
	if !ok {
		return nil, errors.New("PKI_INTERMEDIATE_PRIVATE_KEY not found in environment")
	}

	pool := x509.CertPool{}
	if !pool.AppendCertsFromPEM([]byte(rootPem)) {
		return nil, errors.New("failed adding root certs to pool")
	}

	policy := &config.Signing{
		Profiles: map[string]*config.SigningProfile{},
		Default:  config.DefaultConfig(),
	}
	policy.SetRemoteCAs(&pool)

	b := pki.Bundle{}

	if err := b.Cert.UnmarshalText([]byte(intermediatePem)); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal intermediate certificate")
	}

	if err := b.Key.UnmarshalText([]byte(intermediateKeyPem)); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal private key")
	}

	return local.NewSigner(b.Key, b.Cert.Certificate, x509.SHA256WithRSA, policy)
}

func (s *Server) getOCSPSigner() (ocsp.Signer, error) {
	rootPem, ok := s.Lookup("PKI_ROOT_CERTIFICATE")
	if !ok {
		return nil, errors.New("PKI_ROOT_CERTIFICATE not found in environment")
	}

	intermediatePem, ok := s.Lookup("PKI_INTERMEDIATE_CERTIFICATE")
	if !ok {
		return nil, errors.New("PKI_INTERMEDIATE_CERTIFICATE not found in environment")
	}

	intermediateKeyPem, ok := s.Lookup("PKI_INTERMEDIATE_PRIVATE_KEY")
	if !ok {
		return nil, errors.New("PKI_INTERMEDIATE_PRIVATE_KEY not found in environment")
	}

	root, err := helpers.ParseCertificatePEM([]byte(rootPem))
	if err != nil {
		return nil, err
	}

	intermediate, err := helpers.ParseCertificatePEM([]byte(intermediatePem))
	if err != nil {
		return nil, err
	}

	intermediateKey, err := helpers.ParsePrivateKeyPEM([]byte(intermediateKeyPem))
	if err != nil {
		return nil, err
	}

	return ocsp.NewSigner(root, intermediate, intermediateKey, time.Hour)
}

func (s *Server) Serve() error {
	s.Init()

	l, err := s.Listener()
	if err != nil {
		return errors.Wrap(err, "unable to open listener")
	}
	defer l.Close()

	done := make(chan error)
	defer close(done)

	go func() {
		srv := http.Server{Handler: s.engine}
		go func() {
			done <- srv.Serve(l)
		}()

		<-s.Done()
		logrus.WithError(s.Err()).Debug("context done")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logrus.WithError(err).Warn("error shutting down server")
		}
	}()

	err = <-done
	if err != nil {
		logrus.WithError(err).Warning("server shutdown")
	}
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (s *Server) GetSecret() []byte {
	entropy, err := s.ProjectEntropy()
	if err == nil {
		logrus.WithField("entropy", entropy).Debug("found project entropy")
		if rv, err := base32.StdEncoding.DecodeString(entropy); err == nil {
			return rv
		} else {
			logrus.WithField("err", err).WithField("entropy", entropy).Warn("error decoding project entropy")
		}
	} else {
		logrus.WithField("err", err).Warn("project entropy not found")
	}

	secret := make([]byte, 40)
	_, _ = rand.Read(secret)
	logrus.WithField("secret", base32.StdEncoding.EncodeToString(secret)).Warn("using random secret")
	return secret
}

func (s *Server) GetSessionStore() sessions.Store {
	secret := s.GetSecret() // TODO: rotate secret?
	var store sessions.Store
	if rels, err := s.Relationships(); err == nil {
		db, err := rels.MongoDB("sessions")
		if err == nil {
			col := db.C("sessions")
			store = mongo.NewStore(col, OneYear, true, secret)
		} else {
			logrus.WithError(err).Warn("unable to connect to mongo server")
		}
	} else {
		logrus.WithField("err", err).Warn("unable to determine relationships")
	}
	if store == nil {
		logrus.Warn("using cookie session store")
		store = cookie.NewStore(secret)
	}
	store.Options(sessions.Options{
		MaxAge: OneMonth,
		Secure: true,
	})
	return store
}
