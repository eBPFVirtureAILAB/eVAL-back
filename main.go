package main

/*

This is Main Backend for EVAL (eBPFVirtualAILab ) project

Backend's main purpose is connecting from k8s pods to the frontend

Mainly used packages:
- gin-gonic/gin
- kubernetes.io/client-go
- entgo.io/ent

*/

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"eval/ent"
	global "eval/global"
	k8s "eval/pkg/k8s"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {

	fmt.Println("############################################################")
	fmt.Println("EVAL Backend! at time : ", time.Now())
	fmt.Println("############################################################")

	// ############################################################

	global.SetEnv()

	// ############################################################

	openPostgresql := "host=" + global.Postgres.Host + " port=" + global.Postgres.Port + " user=" + global.Postgres.User + " dbname=" + global.Postgres.DBName + " password=" + global.Postgres.Password + " sslmode=disable"
	fmt.Println("openPostgresql : ", openPostgresql)

	dbClient, err := ent.Open("postgres", openPostgresql)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer dbClient.Close()
	// Run the auto migration tool.
	if err := dbClient.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// ############################################################

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
	// save kube config to db
	k8sGroup.POST("/config", k8s.SaveKubeConfigHandler())
	// pod list
	k8sGroup.GET("/pods", k8s.GetPodsHandler())
	// get pod log
	k8sGroup.GET("/pods/log", k8s.GetPodLogHandler())
	// get pod connection
	k8sGroup.GET("/pods/connection", k8s.PodConnectionHandler())

	// ############################################################
	// todo: AuthN api
	// todo: AuthZ api
	// todo: etc...

	// ############################################################

	// Run the server
	r.Run(global.HOST + ":" + global.PORT) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
