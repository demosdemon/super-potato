// This file is generated - do not edit!

package server

import (
	"net/http"

	gin "github.com/gin-gonic/gin"
	logrus "github.com/sirupsen/logrus"

	platformsh "github.com/demosdemon/super-potato/pkg/platformsh"
)

func addRoutes(group gin.IRoutes) gin.IRoutes {
	group.GET("application", getApplication)
	group.GET("application_name", getApplicationName)
	group.GET("app_name", getApplicationName)
	group.GET("app_command", getAppCommand)
	group.GET("application_command", getAppCommand)
	group.GET("app_dir", getAppDir)
	group.GET("branch", getBranch)
	group.GET("dir", getDir)
	group.GET("document_root", getDocumentRoot)
	group.GET("environment", getEnvironment)
	group.GET("port", getPort)
	group.GET("project", getProject)
	group.GET("project_entropy", getProjectEntropy)
	group.GET("relationships", getRelationships)
	group.GET("routes", getRoutes)
	group.GET("smtp_host", getSMTPHost)
	group.GET("socket", getSocket)
	group.GET("tree_id", getTreeID)
	group.GET("variables", getVariables)
	group.GET("vars", getVariables)
	return group
}

func getApplication(c *gin.Context) {
	logrus.Trace("getApplication")
	obj, err := env.Application()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getApplicationName(c *gin.Context) {
	logrus.Trace("getApplicationName")
	obj, err := env.ApplicationName()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getAppCommand(c *gin.Context) {
	logrus.Trace("getAppCommand")
	obj, err := env.AppCommand()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getAppDir(c *gin.Context) {
	logrus.Trace("getAppDir")
	obj, err := env.AppDir()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getBranch(c *gin.Context) {
	logrus.Trace("getBranch")
	obj, err := env.Branch()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getDir(c *gin.Context) {
	logrus.Trace("getDir")
	obj, err := env.Dir()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getDocumentRoot(c *gin.Context) {
	logrus.Trace("getDocumentRoot")
	obj, err := env.DocumentRoot()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getEnvironment(c *gin.Context) {
	logrus.Trace("getEnvironment")
	obj, err := env.Environment()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getPort(c *gin.Context) {
	logrus.Trace("getPort")
	obj, err := env.Port()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getProject(c *gin.Context) {
	logrus.Trace("getProject")
	obj, err := env.Project()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getProjectEntropy(c *gin.Context) {
	logrus.Trace("getProjectEntropy")
	obj, err := env.ProjectEntropy()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getRelationships(c *gin.Context) {
	logrus.Trace("getRelationships")
	obj, err := env.Relationships()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getRoutes(c *gin.Context) {
	logrus.Trace("getRoutes")
	obj, err := env.Routes()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getSMTPHost(c *gin.Context) {
	logrus.Trace("getSMTPHost")
	obj, err := env.SMTPHost()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getSocket(c *gin.Context) {
	logrus.Trace("getSocket")
	obj, err := env.Socket()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getTreeID(c *gin.Context) {
	logrus.Trace("getTreeID")
	obj, err := env.TreeID()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func getVariables(c *gin.Context) {
	logrus.Trace("getVariables")
	obj, err := env.Variables()
	_, ok := err.(platformsh.MissingEnvironment)
	switch {
	case err == nil:
		c.IndentedJSON(http.StatusOK, obj)
	case ok:
		c.IndentedJSON(http.StatusNotFound, err)
	default:
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}
