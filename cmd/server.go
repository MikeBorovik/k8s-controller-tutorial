package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/MikeBorovik/k8s-controller-tutorial/pkg/ctrl"
	"github.com/MikeBorovik/k8s-controller-tutorial/pkg/informer"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrlruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var serverPort int
var serverInCluster bool
var serverKubeConfig string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a FastHTTP server and deployment informer",
	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := getServerKubeClient(serverKubeConfig, serverInCluster)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create kubernetes client.")
			os.Exit(1)
		}
		ctx := context.Background()
		go informer.StartDeploymentInformer(ctx, clientset)

		mgr, err := ctrlruntime.NewManager(ctrlruntime.GetConfigOrDie(), manager.Options{})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create controller-runtime manager")
			os.Exit(1)
		}
		if err := ctrl.AddDeploymentController(mgr); err != nil {
			log.Error().Err(err).Msg("Failed to add deployment controller")
			os.Exit(1)
		}
		go func() {
			log.Info().Msg("Starting controller-runtime manager...")
			if err := mgr.Start(cmd.Context()); err != nil {
				log.Error().Err(err).Msg("Error starting controller-runtime manager")
				os.Exit(1)
			}
		}()

		handler := createHandler(&informer.DeploymentInformer{})
		addr := fmt.Sprintf(":%d", serverPort)
		log.Info().Msgf("Starting FastHTTP server on %s", addr)
		if err := fasthttp.ListenAndServe(addr, handler); err != nil {
			log.Error().Err(err).Msg("Error starting FastHTTP server")
			os.Exit(1)
		}
	},
}

func createHandler(lister informer.DeploymentLister) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		requestID := uuid.New().String()
		ctx.Response.Header.Set("X-Request-ID", requestID)
		logger := log.With().Str("request_id", requestID).Logger()
		switch string(ctx.Path()) {
		case "/deployments":
			ctx.Response.Header.Set("Content-Type", "application/json")
			deployments := lister.GetDeploymentsNames()
			logger.Info().Msgf("Deployments: %v", deployments)
			ctx.SetStatusCode(200)
			ctx.Write([]byte("["))
			for i, name := range deployments {
				ctx.WriteString("\"")
				ctx.WriteString(name)
				ctx.WriteString("\"")
				if i < len(deployments)-1 {
					ctx.WriteString(",")
				}
			}
			ctx.Write([]byte("]"))
			return
		default:
			logger.Info().Msg("Default path received")
			fmt.Fprintf(ctx, "Hello from FastHTTP!")
		}
	}
}

func getServerKubeClient(kubeconfigPath string, inCluster bool) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if inCluster {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVar(&serverPort, "port", 8080, "Port to run the server on")
	serverCmd.Flags().StringVar(&serverKubeConfig, "kubeconfig", "", "Path to kubeconfig")
	serverCmd.Flags().BoolVar(&serverInCluster, "in-cluster", false, "Use in-cluster kubeconfg")
}
