/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

// KubernetesClient определяет интерфейс для работы с кластером
type KubernetesClient interface {
	ListDeployments(namespace string) ([]string, error)
}

// DefaultKubernetesClient реализует KubernetesClient с реальным API Kubernetes
type DefaultKubernetesClient struct {
	clientset kubernetes.Interface
}

// NewKubernetesClient создает новый клиент Kubernetes
func NewKubernetesClient(kubeconfigPath string) (KubernetesClient, error) {
	var config *rest.Config
	var err error

	// Используем указанный kubeconfig или дефолтный путь
	if kubeconfigPath == "" {
		log.Debug().Msg("Using in-cluster configuration")
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Debug().Msg("In-cluster config failed, falling back to default kubeconfig")
			config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		}
	} else {
		log.Debug().Str("kubeconfig", kubeconfigPath).Msg("Using provided kubeconfig")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	return &DefaultKubernetesClient{
		clientset: clientset,
	}, nil
}

// ListDeployments получает список деплойментов в указанном namespace
func (c *DefaultKubernetesClient) ListDeployments(namespace string) ([]string, error) {
	deployments, err := c.clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	var names []string
	for _, deployment := range deployments.Items {
		names = append(names, deployment.Name)
	}

	return names, nil
}

// Factory для создания клиентов - можно заменить в тестах
var clientFactory = NewKubernetesClient

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List deployments in the default namespace",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().Msg("List command started")
		return runListCommand(kubeconfig, cmd.OutOrStdout())
	},
}

// runListCommand выполняет логику команды list
func runListCommand(kubeconfigPath string, out io.Writer) error {
	// Создаем клиент
	client, err := clientFactory(kubeconfigPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Kubernetes client")
		return err
	}

	// Получаем список деплойментов
	deployments, err := client.ListDeployments("default")
	if err != nil {
		log.Error().Err(err).Msg("Failed to list deployments")
		return err
	}

	// Выводим результат
	fmt.Fprintf(out, "Found %d deployments in 'default' namespace:\n", len(deployments))
	for _, name := range deployments {
		fmt.Fprintln(out, "-", name)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "Path to kubeconfig file")
}
