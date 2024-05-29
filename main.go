package main

/*

This is Main Backend for EVAL (eBPFVirtualAILab ) project

Backend's main purpose is connecting from k8s pods to the frontend

Mainly used packages:
- gin-gonic/gin
- kubernetes.io/client-go

*/

import (
	"net/http"

	global "eval/global"
	k8s "eval/pkg/k8s"

	"github.com/gin-gonic/gin"
)

func main() {

	global.SetEnv()

	// Create a gin router
	r := gin.Default()

	// ############################################################

	// db api group
	dbGroup := r.Group("/api/db")
	// get user data
	dbGroup.GET("/user", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "api test",
		})
	})

	// ############################################################

	// k8s api group
	k8sGroup := r.Group("/api/k8s")

	// fixme: namespace is fixed = "eval"
	// pod list
	k8sGroup.GET("/pod", k8s.GetPodsHandler())
	// get pod log
	k8sGroup.GET("/pod/log", k8s.GetPodLogHandler())
	// get pod connection
	k8sGroup.GET("/pod/connection", k8s.PodConnectionHandler())

	// ############################################################
	// todo: AuthN api
	// todo: AuthZ api
	// todo: etc...

	// ############################################################

	// Run the server
	r.Run(global.HOST + ":" + global.PORT) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
