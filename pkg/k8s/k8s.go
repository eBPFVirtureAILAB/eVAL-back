package k8s

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// Connect to the k8s cluster with the kubeconfig file.
func Connect(kubeconfig string) (*kubernetes.Clientset, *rest.Config, error) {

	// Connect to the k8s cluster with the kubeconfig file.

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientset, config, nil
}

// Get the pods in the eval namespace.
func GetPodsInEval(clientset *kubernetes.Clientset) ([]string, error) {
	// Get the pods in the eval namespace.

	pods, err := clientset.CoreV1().Pods("eval").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// change the return type to []string
	var podNames []string
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}

	return podNames, nil
}

// Get the logs of the pod.
func GetPodLogs(clientset *kubernetes.Clientset, podName string) ([]byte, error) {

	line := int64(10)

	podLogOptions := &apiv1.PodLogOptions{TailLines: &line}
	podLog := clientset.CoreV1().Pods("eval").GetLogs(podName, podLogOptions).Do(context.TODO())

	log, err := podLog.Raw()
	if err != nil {
		return nil, err
	}

	return log, nil
}

type wsStream struct {
	ws *websocket.Conn
}

func (w wsStream) Read(p []byte) (size int, err error) {
	_, r, err := w.ws.ReadMessage()
	if err != nil {
		return 0, err
	}

	// Add newline to the end of the message
	// 추후 프론트엔드 측 터미널 라이브러리에 맞게 수정 필요
	r = []byte(string(r) + "\n")

	size = copy(p, r)

	return size, nil
}

func (w wsStream) Write(p []byte) (size int, err error) {
	err = w.ws.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func PodConnection(clientset *kubernetes.Clientset, config *rest.Config, podName string, ws *websocket.Conn) error {

	// Exec 요청 준비
	namespace := "eval"
	command := "sh"

	req := clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		VersionedParams(&apiv1.PodExecOptions{
			Command: []string{command},
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)

	// Executor 생성 및 실행 : NewSPDYExecutor는 SPDY 프로토콜을 사용하여 원격 명령을 실행하는 Executor를 생성
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		fmt.Println("remotecommand.NewSPDYExecutor", err)
		return err
	}

	// wsStream 구조체 생성
	wsStream := wsStream{ws}

	ws.WriteMessage(websocket.TextMessage, []byte("Connected to the pod : "+podName))

	// 원격 명령 실행
	ctx := context.TODO()
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  wsStream,
		Stdout: wsStream,
		Stderr: wsStream,
		Tty:    true,
	})
	if err != nil {
		fmt.Println("exec.StreamWithContex", err)
		return err
	}

	return nil
}
