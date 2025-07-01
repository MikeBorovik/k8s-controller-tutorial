package informer

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type DeploymentInformer struct{}

func (d *DeploymentInformer) GetDeploymentsNames() []string {
	return GetDeploymentsNames()
}

type DeploymentLister interface {
	GetDeploymentsNames() []string
}

var informer cache.SharedIndexInformer

func StartDeploymentInformer(ctx context.Context, clientset *kubernetes.Clientset) {
	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset,
		25*time.Second,
		informers.WithNamespace("default"),
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.FieldSelector = fields.Everything().String()
		}),
	)
	informer = factory.Apps().V1().Deployments().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Info().Msgf("Deployment added: %s", getDeploymentName(obj))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			log.Info().Msgf("Deployment updated: %s", getDeploymentName(newObj))
		},
		DeleteFunc: func(obj interface{}) {
			log.Info().Msgf("Deployment deleted: %s", getDeploymentName(obj))
		},
	})
	log.Info().Msg("Starting deployment informer...")

	factory.Start(ctx.Done())
	for t, ok := range factory.WaitForCacheSync(ctx.Done()) {
		if !ok {
			log.Error().Msgf("Failed to sync informer foe %v", t)
			os.Exit(1)
		}
	}
	log.Info().Msg("Deployment informer cache synced. Watching for events...")
	<-ctx.Done()
}

func GetDeploymentsNames() []string {
	var names []string
	if informer == nil {
		return names
	}
	for _, obj := range informer.GetStore().List() {
		if deployment, ok := obj.(*appsv1.Deployment); ok {
			names = append(names, deployment.Name)
		}
	}
	return names
}

func getDeploymentName(obj any) string {
	if deployment, ok := obj.(metav1.Object); ok {
		return deployment.GetName()
	}
	return "unknown"
}
