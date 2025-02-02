package cmd

import (
	"fmt"
	"os"

	"path/filepath"

	"os/signal"
	"syscall"

	"time"

	"k8s.io/client-go/util/homedir"

	"github.com/ismferd/serf-publisher/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spotahome/kooper/operator/controller"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	zapLogger, _ = zap.NewProduction(zap.AddCallerSkip(1))
	logger       = pkg.NewLogger(*zapLogger.Sugar())
	stopC        = make(chan struct{})
	finishC      = make(chan error)
	signalC      = make(chan os.Signal, 1)
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "serf-publisher",
	Short: "Kubernetes controller that automatically adds annotations in Pods to they can assume AWS Roles.",
	Long:  `Kubernetes controller that automatically adds annotations in Pods to they can assume AWS Roles.`,
	Run: func(cmd *cobra.Command, args []string) {
		defer zapLogger.Sync()
		signal.Notify(signalC, syscall.SIGTERM, syscall.SIGINT)

		k8sCli, err := getKubernetesClient(viper.GetString("kubeconfig"), logger)
		if err != nil {
			logger.Errorf("Can't create k8s client: %s", err)
			os.Exit(1)
		}

		go func() {
			finishC <- getController(k8sCli, logger).Run(stopC)
		}()

		select {
		case err := <-finishC:
			if err != nil {
				logger.Errorf("error running controller: %s", err)
				os.Exit(1)
			}
		case <-signalC:
			logger.Info("Signal captured, exiting...")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		namespace := viper.GetString("namespace")
		if len(namespace) == 0 {
			logger.Error("Error: required flag \"namespace\" or environment variable \"NAMESPACE\" not set")
			os.Exit(1)
		}
	})

	rootCmd.Flags().String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "Path to the kubeconfig file")
	_ = viper.BindPFlag("kubeconfig", rootCmd.Flags().Lookup("kubeconfig"))

	rootCmd.Flags().String("namespace", "", "kubernetes namespace where this app is running")
	_ = viper.BindPFlag("namespace", rootCmd.Flags().Lookup("namespace"))

	rootCmd.Flags().Int("resync-seconds", 30, "The number of seconds the controller will resync the resources")
	_ = viper.BindPFlag("resync-seconds", rootCmd.Flags().Lookup("resync-seconds"))

	viper.AutomaticEnv()
}

func getKubernetesClient(kubeconfig string, logger pkg.Logger) (kubernetes.Interface, error) {
	var err error
	var cfg *rest.Config

	cfg, err = rest.InClusterConfig()
	if err != nil {
		logger.Warningf("Falling back to using kubeconfig file: %s", err)
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	return kubernetes.NewForConfig(cfg)
}

func getController(k8sCli kubernetes.Interface, logger pkg.Logger) controller.Controller {
	return controller.NewSequential(
		time.Duration(viper.GetInt("resync-seconds"))*time.Second,
		pkg.NewHandler(pkg.NewSerfPublisher(k8sCli, logger)),
		pkg.NewDeploymentRetrieve(viper.GetString("namespace"), k8sCli),
		nil,
		logger,
	)
}
