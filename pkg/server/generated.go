// This file is generated - do not edit!

package server

import (
	"net/http"

	gin "github.com/gin-gonic/gin"
	logrus "github.com/sirupsen/logrus"

	platformsh "github.com/demosdemon/super-potato/pkg/platformsh"
)

func (s *Server) registerGeneratedRoutes(r gin.IRoutes) {
	r.GET("application", s.getApplication)
	r.GET("application_name", s.getApplicationName)
	r.GET("app_name", s.getApplicationName)
	r.GET("app_command", s.getAppCommand)
	r.GET("application_command", s.getAppCommand)
	r.GET("app_dir", s.getAppDir)
	r.GET("branch", s.getBranch)
	r.GET("dir", s.getDir)
	r.GET("document_root", s.getDocumentRoot)
	r.GET("environment", s.getEnvironment)
	r.GET("port", s.getPort)
	r.GET("project", s.getProject)
	r.GET("project_entropy", s.getProjectEntropy)
	r.GET("relationships", s.getRelationships)
	r.GET("routes", s.getRoutes)
	r.GET("smtp_host", s.getSMTPHost)
	r.GET("socket", s.getSocket)
	r.GET("tree_id", s.getTreeID)
	r.GET("variables", s.getVariables)
	r.GET("vars", s.getVariables)
	r.GET("x_client_cert", s.getXClientCert)
	r.GET("x_client_dn", s.getXClientDN)
	r.GET("x_client_ip", s.getXClientIP)
	r.GET("x_client_ssl", s.getXClientSSL)
	r.GET("x_client_verify", s.getXClientVerify)
}

func (s *Server) getApplication(c *gin.Context) {
	logrus.Trace("getApplication")
	obj, err := s.Environment.Application()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getApplicationName(c *gin.Context) {
	logrus.Trace("getApplicationName")
	obj, err := s.Environment.ApplicationName()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getAppCommand(c *gin.Context) {
	logrus.Trace("getAppCommand")
	obj, err := s.Environment.AppCommand()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getAppDir(c *gin.Context) {
	logrus.Trace("getAppDir")
	obj, err := s.Environment.AppDir()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getBranch(c *gin.Context) {
	logrus.Trace("getBranch")
	obj, err := s.Environment.Branch()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getDir(c *gin.Context) {
	logrus.Trace("getDir")
	obj, err := s.Environment.Dir()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getDocumentRoot(c *gin.Context) {
	logrus.Trace("getDocumentRoot")
	obj, err := s.Environment.DocumentRoot()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getEnvironment(c *gin.Context) {
	logrus.Trace("getEnvironment")
	obj, err := s.Environment.Environment()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getPort(c *gin.Context) {
	logrus.Trace("getPort")
	obj, err := s.Environment.Port()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getProject(c *gin.Context) {
	logrus.Trace("getProject")
	obj, err := s.Environment.Project()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getProjectEntropy(c *gin.Context) {
	logrus.Trace("getProjectEntropy")
	obj, err := s.Environment.ProjectEntropy()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getRelationships(c *gin.Context) {
	logrus.Trace("getRelationships")
	obj, err := s.Environment.Relationships()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getRoutes(c *gin.Context) {
	logrus.Trace("getRoutes")
	obj, err := s.Environment.Routes()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getSMTPHost(c *gin.Context) {
	logrus.Trace("getSMTPHost")
	obj, err := s.Environment.SMTPHost()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getSocket(c *gin.Context) {
	logrus.Trace("getSocket")
	obj, err := s.Environment.Socket()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getTreeID(c *gin.Context) {
	logrus.Trace("getTreeID")
	obj, err := s.Environment.TreeID()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getVariables(c *gin.Context) {
	logrus.Trace("getVariables")
	obj, err := s.Environment.Variables()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getXClientCert(c *gin.Context) {
	logrus.Trace("getXClientCert")
	obj, err := s.Environment.XClientCert()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getXClientDN(c *gin.Context) {
	logrus.Trace("getXClientDN")
	obj, err := s.Environment.XClientDN()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getXClientIP(c *gin.Context) {
	logrus.Trace("getXClientIP")
	obj, err := s.Environment.XClientIP()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getXClientSSL(c *gin.Context) {
	logrus.Trace("getXClientSSL")
	obj, err := s.Environment.XClientSSL()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}

func (s *Server) getXClientVerify(c *gin.Context) {
	logrus.Trace("getXClientVerify")
	obj, err := s.Environment.XClientVerify()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		s.negotiate(c, http.StatusOK, obj)
	case ok:
		s.negotiate(c, http.StatusNotFound, err)
	default:
		s.negotiate(c, http.StatusInternalServerError, err)
	}
}
