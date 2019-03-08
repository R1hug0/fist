package main

import (
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//Heartbeater is
type Heartbeater interface {
	//terminal deployment is in sealyun-tty namespace
	UpdateTimestamp(clientSet *kubernetes.Clientset) error
	//need delete deployment and service in sealyun-tty if it timeout
	CleanTerminalJob(clientSet *kubernetes.Clientset) error
}

type terHeartbeater struct {
	terminalID string
	namespace  string
}

func (t *terHeartbeater) CleanTerminalJob(clientSet *kubernetes.Clientset) error {
	//get deploy of terminalId
	deploymentsClient := clientSet.AppsV1().Deployments(t.namespace)
	deploymentResult, err := deploymentsClient.Get(PrefixDeploy+t.terminalID, metav1.GetOptions{})
	if err != nil {
		return err
	}
	//get svc of terminalId
	serviceClient := clientSet.CoreV1().Services(t.namespace)
	_, err = serviceClient.Get(PrefixSvc+t.terminalID, metav1.GetOptions{})
	if err != nil {
		return err
	}
	lastHeartbeat, _ := strconv.Atoi(deploymentResult.Annotations["fist.seayun.com/lastHeartbeat"])
	heartbeatTime := time.Now().Unix() - int64(lastHeartbeat)
	if heartbeatTime > 600 { // time.Minute() * 10 = 600
		deletePolicy := metav1.DeletePropagationForeground
		//delete deploy
		if err := deploymentsClient.Delete(PrefixDeploy+t.terminalID, &metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			return err
		}
		//delete svc
		if err := serviceClient.Delete(PrefixSvc+t.terminalID, &metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (t *terHeartbeater) UpdateTimestamp(clientSet *kubernetes.Clientset) error {
	deploymentsClient := clientSet.AppsV1().Deployments(t.namespace)
	result, err := deploymentsClient.Get(PrefixDeploy+t.terminalID, metav1.GetOptions{})
	if err != nil {
		return err
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	result.Annotations["fist.seayun.com/lastHeartbeat"] = timestamp
	_, err = deploymentsClient.Update(result)
	if err != nil {
		return err
	}
	return nil
}