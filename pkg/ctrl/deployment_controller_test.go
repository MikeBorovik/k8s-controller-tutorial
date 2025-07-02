package ctrl_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/MikeBorovik/k8s-controller-tutorial/pkg/ctrl"
)

func TestDeploymentReconciler(t *testing.T) {
	// Set up the test environment
	ctx := context.Background()
	testEnv := &envtest.Environment{}
	cfg, err := testEnv.Start()
	require.NoError(t, err)
	defer func() { _ = testEnv.Stop() }()

	// Add the Deployment API to the scheme
	err = appsv1.AddToScheme(scheme.Scheme)
	require.NoError(t, err)

	// Create a manager
	mgr, err := manager.New(cfg, manager.Options{
		Scheme: scheme.Scheme,
	})
	require.NoError(t, err)

	// Add the DeploymentReconciler to the manager
	err = ctrl.AddDeploymentController(mgr)
	require.NoError(t, err)

	// Start the manager in a separate goroutine
	go func() {
		require.NoError(t, mgr.Start(ctx))
	}()

	// Create a client to interact with the test cluster
	k8sClient := mgr.GetClient()

	// Create a sample Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-deployment",
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "test"},
			},
			Template: corev1.PodTemplateSpec{ // Use corev1.PodTemplateSpec
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "test"}},
				Spec: corev1.PodSpec{ // Use corev1.PodSpec
					Containers: []corev1.Container{ // Use corev1.Container
						{Name: "nginx", Image: "nginx"},
					},
				},
			},
		},
	}
	err = k8sClient.Create(ctx, deployment)
	require.NoError(t, err)

	// Verify that the Reconcile method is called
	reconciledDeployment := &appsv1.Deployment{}
	require.Eventually(t, func() bool {
		err := k8sClient.Get(ctx, types.NamespacedName{
			Name:      "test-deployment",
			Namespace: "default",
		}, reconciledDeployment)
		return err == nil
	}, 10*time.Second, 500*time.Millisecond)

	// Verify the Deployment exists in the cluster
	require.Equal(t, "test-deployment", reconciledDeployment.Name)
	require.Equal(t, "default", reconciledDeployment.Namespace)
}

func int32Ptr(i int32) *int32 { return &i }
