// This file is generated - do not edit!

package serve

import (
	"net/http"

	gin "github.com/gin-gonic/gin"
	logrus "github.com/sirupsen/logrus"

	platformsh "github.com/demosdemon/super-potato/pkg/platformsh"
)

func (x *API) addGeneratedRoutes() {
	x.routes.GET("application", x.getApplication)
	x.routes.GET("application_name", x.getApplicationName)
	x.routes.GET("app_name", x.getApplicationName)
	x.routes.GET("app_command", x.getAppCommand)
	x.routes.GET("application_command", x.getAppCommand)
	x.routes.GET("app_dir", x.getAppDir)
	x.routes.GET("branch", x.getBranch)
	x.routes.GET("dir", x.getDir)
	x.routes.GET("document_root", x.getDocumentRoot)
	x.routes.GET("environment", x.getEnvironment)
	x.routes.GET("port", x.getPort)
	x.routes.GET("project", x.getProject)
	x.routes.GET("project_entropy", x.getProjectEntropy)
	x.routes.GET("relationships", x.getRelationships)
	x.routes.GET("routes", x.getRoutes)
	x.routes.GET("smtp_host", x.getSMTPHost)
	x.routes.GET("socket", x.getSocket)
	x.routes.GET("tree_id", x.getTreeID)
	x.routes.GET("variables", x.getVariables)
	x.routes.GET("vars", x.getVariables)
	x.routes.GET("x_client_cert", x.getXClientCert)
	x.routes.GET("x_client_dn", x.getXClientDN)
	x.routes.GET("x_client_ip", x.getXClientIP)
	x.routes.GET("x_client_ssl", x.getXClientSSL)
	x.routes.GET("x_client_verify", x.getXClientVerify)
}

func (x *API) getApplication(c *gin.Context) {
	logrus.Trace("getApplication")
	obj, err := x.env.Application()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getApplicationName(c *gin.Context) {
	logrus.Trace("getApplicationName")
	obj, err := x.env.ApplicationName()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getAppCommand(c *gin.Context) {
	logrus.Trace("getAppCommand")
	obj, err := x.env.AppCommand()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getAppDir(c *gin.Context) {
	logrus.Trace("getAppDir")
	obj, err := x.env.AppDir()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getBranch(c *gin.Context) {
	logrus.Trace("getBranch")
	obj, err := x.env.Branch()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getDir(c *gin.Context) {
	logrus.Trace("getDir")
	obj, err := x.env.Dir()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getDocumentRoot(c *gin.Context) {
	logrus.Trace("getDocumentRoot")
	obj, err := x.env.DocumentRoot()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getEnvironment(c *gin.Context) {
	logrus.Trace("getEnvironment")
	obj, err := x.env.Environment()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getPort(c *gin.Context) {
	logrus.Trace("getPort")
	obj, err := x.env.Port()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getProject(c *gin.Context) {
	logrus.Trace("getProject")
	obj, err := x.env.Project()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getProjectEntropy(c *gin.Context) {
	logrus.Trace("getProjectEntropy")
	obj, err := x.env.ProjectEntropy()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getRelationships(c *gin.Context) {
	logrus.Trace("getRelationships")
	obj, err := x.env.Relationships()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getRoutes(c *gin.Context) {
	logrus.Trace("getRoutes")
	obj, err := x.env.Routes()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getSMTPHost(c *gin.Context) {
	logrus.Trace("getSMTPHost")
	obj, err := x.env.SMTPHost()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getSocket(c *gin.Context) {
	logrus.Trace("getSocket")
	obj, err := x.env.Socket()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getTreeID(c *gin.Context) {
	logrus.Trace("getTreeID")
	obj, err := x.env.TreeID()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getVariables(c *gin.Context) {
	logrus.Trace("getVariables")
	obj, err := x.env.Variables()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getXClientCert(c *gin.Context) {
	logrus.Trace("getXClientCert")
	obj, err := x.env.XClientCert()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getXClientDN(c *gin.Context) {
	logrus.Trace("getXClientDN")
	obj, err := x.env.XClientDN()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getXClientIP(c *gin.Context) {
	logrus.Trace("getXClientIP")
	obj, err := x.env.XClientIP()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getXClientSSL(c *gin.Context) {
	logrus.Trace("getXClientSSL")
	obj, err := x.env.XClientSSL()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}

func (x *API) getXClientVerify(c *gin.Context) {
	logrus.Trace("getXClientVerify")
	obj, err := x.env.XClientVerify()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		negotiate(c, http.StatusOK, obj)
	case ok:
		negotiate(c, http.StatusNotFound, err)
	default:
		negotiate(c, http.StatusInternalServerError, err)
	}
}
