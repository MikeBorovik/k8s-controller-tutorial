package ctrl_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/MikeBorovik/k8s-controller-tutorial/pkg/ctrl"
	"github.com/MikeBorovik/k8s-controller-tutorial/pkg/testutil" // Import your custom envtest package
)

func TestDeploymentReconciler(t *testing.T) {
	// Use the custom SetupEnv function
	env, clientset, cleanup := testutil.SetupEnv(t)
	defer cleanup()

	// Add the Deployment API to the scheme
	err := appsv1.AddToScheme(scheme.Scheme)
	require.NoError(t, err)

	// Create a manager
	mgr, err := manager.New(env.Config, manager.Options{
		Scheme: scheme.Scheme,
	})
	require.NoError(t, err)

	// Add the DeploymentReconciler to the manager
	err = ctrl.AddDeploymentController(mgr)
	require.NoError(t, err)

	// Start the manager in a separate goroutine
	ctx := context.Background()
	go func() {
		require.NoError(t, mgr.Start(ctx))
	}()

	// Create a sample Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: testutil.NewObjectMeta("test-deployment", "default"),
		Spec:       testutil.NewDeploymentSpec(1, map[string]string{"app": "test"}, "nginx"),
	}
	_, err = clientset.AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
	require.NoError(t, err)

	// Verify that the Reconcile method is called
	reconciledDeployment := &appsv1.Deployment{}
	require.Eventually(t, func() bool {
		err := mgr.GetClient().Get(ctx, types.NamespacedName{
			Name:      "test-deployment",
			Namespace: "default",
		}, reconciledDeployment)
		return err == nil
	}, 10*time.Second, 500*time.Millisecond)

	// Verify the Deployment exists in the cluster
	require.Equal(t, "test-deployment", reconciledDeployment.Name)
	require.Equal(t, "default", reconciledDeployment.Namespace)
}
