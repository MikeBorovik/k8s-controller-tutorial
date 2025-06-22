package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockKubernetesClient mock для тестирования
type MockKubernetesClient struct {
	DeploymentNames []string
	Error           error
}

func (m *MockKubernetesClient) ListDeployments(namespace string) ([]string, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.DeploymentNames, nil
}

// TestRunListCommand тестирует логику команды list
func TestRunListCommand(t *testing.T) {
	// Сохраняем оригинальную фабрику
	origFactory := clientFactory
	defer func() {
		clientFactory = origFactory
	}()

	t.Run("Success", func(t *testing.T) {
		// Настраиваем мок
		mockClient := &MockKubernetesClient{
			DeploymentNames: []string{"test-deployment-1", "test-deployment-2"},
		}

		// Заменяем фабрику на функцию, возвращающую наш мок
		clientFactory = func(kubeconfigPath string) (KubernetesClient, error) {
			return mockClient, nil
		}

		// Перехватываем вывод
		buf := new(bytes.Buffer)

		// Вызываем тестируемую функцию
		err := runListCommand("test-kubeconfig", buf)

		// Проверяем результат
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "Found 2 deployments")
		assert.Contains(t, output, "test-deployment-1")
		assert.Contains(t, output, "test-deployment-2")
	})

	t.Run("Client creation error", func(t *testing.T) {
		// Заменяем фабрику на функцию, возвращающую ошибку
		clientFactory = func(kubeconfigPath string) (KubernetesClient, error) {
			return nil, fmt.Errorf("failed to create client")
		}

		// Перехватываем вывод
		buf := new(bytes.Buffer)

		// Вызываем тестируемую функцию
		err := runListCommand("test-kubeconfig", buf)

		// Проверяем результат
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create client")
	})

	t.Run("List deployments error", func(t *testing.T) {
		// Настраиваем мок с ошибкой
		mockClient := &MockKubernetesClient{
			Error: fmt.Errorf("failed to list deployments"),
		}

		// Заменяем фабрику
		clientFactory = func(kubeconfigPath string) (KubernetesClient, error) {
			return mockClient, nil
		}

		// Перехватываем вывод
		buf := new(bytes.Buffer)

		// Вызываем тестируемую функцию
		err := runListCommand("test-kubeconfig", buf)

		// Проверяем результат
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to list deployments")
	})
}
