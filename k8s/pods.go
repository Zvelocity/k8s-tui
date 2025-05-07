// k8s/pods.go
package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPods returns a list of pods in the current namespace
func (c *Client) GetPods(ctx context.Context) ([]corev1.Pod, error) {
	podList, err := c.clientset.CoreV1().Pods(c.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return podList.Items, nil
}

// GetPodStatus returns a simplified status string for a pod
func GetPodStatus(pod corev1.Pod) string {
	if pod.Status.Phase != corev1.PodRunning {
		return string(pod.Status.Phase)
	}

	// Check container statuses
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			if containerStatus.State.Waiting != nil {
				return containerStatus.State.Waiting.Reason
			}
			return "NotReady"
		}
	}

	return "Running"
}

// GetPodRestarts returns the total number of container restarts in a pod
func GetPodRestarts(pod corev1.Pod) int32 {
	var restarts int32
	for _, containerStatus := range pod.Status.ContainerStatuses {
		restarts += containerStatus.RestartCount
	}
	return restarts
}
