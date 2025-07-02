package testutil

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestInt32Ptr(t *testing.T) {
	v := int32(42)
	ptr := int32Ptr(v)
	if ptr == nil || *ptr != v {
		t.Errorf("int32Ptr(%d) = %v, want pointer to %d", v, ptr, v)
	}
}

func TestNewObjectMeta(t *testing.T) {
	name := "test-name"
	namespace := "test-namespace"
	meta := NewObjectMeta(name, namespace)

	require.Equal(t, name, meta.Name)
	require.Equal(t, namespace, meta.Namespace)
}

func TestNewDeploymentSpec(t *testing.T) {
	replicas := int32(3)
	labels := map[string]string{"app": "test"}
	image := "nginx"
	spec := NewDeploymentSpec(replicas, labels, image)

	require.Equal(t, replicas, *spec.Replicas)
	require.Equal(t, labels, spec.Selector.MatchLabels)
	require.Equal(t, labels, spec.Template.Labels)
	require.Len(t, spec.Template.Spec.Containers, 1)
	require.Equal(t, image, spec.Template.Spec.Containers[0].Image)
}

func TestSetupEnv(t *testing.T) {
	env, clientset, cleanup := SetupEnv(t)
	defer cleanup()

	// Verify that the environment started successfully
	require.NotNil(t, env)
	require.NotNil(t, clientset)

	// Verify that the sample deployments were created
	ctx := context.Background()
	for i := 1; i <= 2; i++ {
		name := types.NamespacedName{
			Name:      "sample-deployment-" + strconv.Itoa(i),
			Namespace: "default",
		}
		deployment, err := clientset.AppsV1().Deployments(name.Namespace).Get(ctx, name.Name, metav1.GetOptions{})
		require.NoError(t, err)
		require.Equal(t, name.Name, deployment.Name)
		require.Equal(t, name.Namespace, deployment.Namespace)
	}
}
