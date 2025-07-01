package informer

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	testutil "github.com/MikeBorovik/k8s-controller-tutorial/pkg/testutil"
)

func TestStartDeploymentInformer(t *testing.T) {
	_, clientset, cleanup := testutil.SetupEnv(t)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	// Patch log to write to a buffer or just rely on test output
	added := make(chan string, 2)

	// Patch event handler for test
	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset,
		30*time.Second,
		informers.WithNamespace("default"),
	)
	informer := factory.Apps().V1().Deployments().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			if d, ok := obj.(metav1.Object); ok {
				added <- d.GetName()
			}
		},
	})

	go func() {
		defer wg.Done()
		factory.Start(ctx.Done())
		factory.WaitForCacheSync(ctx.Done())
		<-ctx.Done()
	}()

	// Wait for events
	found := map[string]bool{}
	for range 2 {
		select {
		case name := <-added:
			found[name] = true
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for deployment add events")
		}
	}

	require.True(t, found["sample-deployment-1"])
	require.True(t, found["sample-deployment-2"])

	cancel()
	wg.Wait()

	//t.Log("Sleeping for 5 minutes to allow manual kubectl inspection of envtest cluster...")
	//time.Sleep(5 * time.Minute)
}

func TestGetDeploymentName(t *testing.T) {
	dep := &metav1.PartialObjectMetadata{}
	dep.SetName("my-deployment")
	name := getDeploymentName(dep)
	if name != "my-deployment" {
		t.Errorf("expected 'my-deployment', got %q", name)
	}
	name = getDeploymentName("not-an-object")
	if name != "unknown" {
		t.Errorf("expected 'unknown', got %q", name)
	}
}

func TestStartDeploymentInformer_CoversFunction(t *testing.T) {
	_, clientset, cleanup := testutil.SetupEnv(t)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run StartDeploymentInformer in a goroutine
	go func() {
		StartDeploymentInformer(ctx, clientset)
	}()

	// Give the informer some time to start and process events
	time.Sleep(1 * time.Second)
	cancel()
}

func TestGetDeploymentsNames(t *testing.T) {
	_, clientset, cleanup := testutil.SetupEnv(t)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset,
		30*time.Second,
		informers.WithNamespace("default"),
	)
	informer := factory.Apps().V1().Deployments().Informer()
	deployment1 := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample-deployment-1",
			Namespace: "default",
		},
	}
	deployment2 := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample-deployment-2",
			Namespace: "default",
		},
	}

	informer.GetStore().Add(deployment1)
	informer.GetStore().Add(deployment2)

	go func() {
		factory.Start(ctx.Done())
		factory.WaitForCacheSync(ctx.Done())
	}()

	names := GetDeploymentsNames()

	require.ElementsMatch(t, []string{"sample-deployment-1", "sample-deployment-2"}, names)
}
