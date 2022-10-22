package utils

import (
	"context"
	"log"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
)

// GetClusterConfig provides the cluster config
func GetClusterConfig() (*rest.Config, error) {
	if len(os.Getenv("KUBECONFIG")) != 0 {
		conf, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
		if err != nil {
			log.Fatalf(err.Error())
			return nil, err
		}
		return conf, nil
	}
	conf, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("error in getting cluster config err: %v", err)
		return nil, err
	}
	conf.RateLimiter = flowcontrol.NewFakeAlwaysRateLimiter()
	return conf, nil
}

func GenerateClientSet(cfg *rest.Config) kubernetes.Interface {
	return kubernetes.NewForConfigOrDie(cfg)
}

// ListPods Lists down the pods in a namespace
func ListPods(ctx context.Context, cs kubernetes.Interface, namespace string, listOptions metav1.ListOptions) (*v1.PodList, error) {
	pods, err := cs.CoreV1().Pods(namespace).List(ctx, listOptions)
	if err != nil {
		return nil, err
	}
	return pods, nil
}

// ReplicatePods replicate pods from one namespace to another
func ReplicatePods(ctx context.Context, cs kubernetes.Interface, sourceNamespace string, destNamespace string, listOptions metav1.ListOptions) error {
	pods, err := cs.CoreV1().Pods(sourceNamespace).List(ctx, listOptions)
	if err != nil {
		log.Printf("error in getting pods for namespace %s err: %v", sourceNamespace, err)
		return err
	}
	for item := range pods.Items {
		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: pods.Items[item].Name,
			},
			Spec:   pods.Items[item].Spec,
			Status: v1.PodStatus{},
		}
		pod, err := cs.CoreV1().Pods(destNamespace).Create(ctx, pod, metav1.CreateOptions{})
		if err != nil {
			log.Printf("error in creating pods %v for namespace %s err: %v", pod, destNamespace, err)
			continue
		}
	}
	return nil
}
