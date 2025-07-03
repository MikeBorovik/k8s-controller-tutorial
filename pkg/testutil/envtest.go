package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

// SetupEnv starts envtest, creates a clientset, populates the cluster with sample Deployments, and returns env, clientset, and cleanup.
// It also creates necessary resources for testing metrics and leader election.
func SetupEnv(t *testing.T) (*envtest.Environment, *kubernetes.Clientset, func()) {
	t.Helper()
	ctx := context.Background()
	env := &envtest.Environment{}

	cfg, err := env.Start()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Write kubeconfig to /tmp/envtest.kubeconfig
	kubeconfig := clientcmdapi.NewConfig()
	kubeconfig.Clusters["envtest"] = &clientcmdapi.Cluster{
		Server:                   cfg.Host,
		CertificateAuthorityData: cfg.CAData,
	}
	kubeconfig.AuthInfos["envtest-user"] = &clientcmdapi.AuthInfo{
		ClientCertificateData: cfg.CertData,
		ClientKeyData:         cfg.KeyData,
	}
	kubeconfig.Contexts["envtest-context"] = &clientcmdapi.Context{
		Cluster:  "envtest",
		AuthInfo: "envtest-user",
	}
	kubeconfig.CurrentContext = "envtest-context"

	kubeconfigBytes, err := clientcmd.Write(*kubeconfig)
	require.NoError(t, err)
	err = os.WriteFile("/tmp/envtest.kubeconfig", kubeconfigBytes, 0644)
	require.NoError(t, err)

	clientset, err := kubernetes.NewForConfig(cfg)
	require.NoError(t, err)

	// Create default namespace if it doesn't exist
	_, err = clientset.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: "default"},
	}, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		require.NoError(t, err)
	}

	// Create sample Deployments
	for i := 1; i <= 2; i++ {
		dep := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("sample-deployment-%d", i),
				Namespace: "default",
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "test"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "test"}},
					Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "nginx", Image: "nginx"}}},
				},
			},
		}
		_, err := clientset.AppsV1().Deployments("default").Create(ctx, dep, metav1.CreateOptions{})
		require.NoError(t, err)
	}
	
	// Create a ConfigMap for leader election testing
	_, err = clientset.CoreV1().ConfigMaps("default").Create(ctx, &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-leader-election",
			Namespace: "default",
		},
		Data: map[string]string{},
	}, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		require.NoError(t, err)
	}

	cleanup := func() {
		_ = env.Stop()
	}
	return env, clientset, cleanup
}

// NewObjectMeta creates a new ObjectMeta with the given name and namespace.
func NewObjectMeta(name, namespace string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}
}

// NewDeploymentSpec creates a new DeploymentSpec with the given replicas, labels, and container image.
func NewDeploymentSpec(replicas int32, labels map[string]string, image string) appsv1.DeploymentSpec {
	return appsv1.DeploymentSpec{
		Replicas: int32Ptr(replicas),
		Selector: &metav1.LabelSelector{
			MatchLabels: labels,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: labels},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Name: "container", Image: image},
				},
			},
		},
	}
}

func int32Ptr(i int32) *int32 { return &i }
