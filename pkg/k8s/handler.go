package k8s

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// get the absolute path of the kubeconfig file
// This folder is the default folder for the kubeconfig file.
// If the kubeconfig file is not in this folder, you can find the kubeconfig file in /home/user/.kube
func kubeconfigPath() string {
	absPath, _ := os.Getwd()

	if _, err := os.Stat(absPath + "/kubeconfig"); os.IsNotExist(err) {
		absPath = os.Getenv("HOME") + "/.kube"
	}

	return absPath + "/kubeconfig"
}

func GetPodsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		clientset, _, err := Connect(kubeconfigPath())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		podNames, err := GetPodsInEval(clientset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"pods": podNames,
		})
	}
}

func GetPodLogHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		podName := c.Query("name")
		fmt.Println(podName)

		clientset, _, err := Connect(kubeconfigPath())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		logs, err := GetPodLogs(clientset, podName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"logs": string(logs),
		})
	}
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// websocket handler for connect to the pod like kubectl exec
func PodConnectionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		podName := c.Query("name")
		fmt.Println(podName)

		clientset, config, err := Connect(kubeconfigPath())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		// Use gorilla websocket
		ws, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.Writer.WriteHeader(http.StatusOK)

		defer ws.Close()

		err = PodConnection(clientset, config, podName, ws)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		}

	}
}

func SaveKubeConfigHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		err = c.SaveUploadedFile(file, kubeconfigPath())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "kubeconfig file is saved",
		})
	}
}
