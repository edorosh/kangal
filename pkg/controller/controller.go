package controller

import (
	"context"
	"fmt"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/hellofresh/kangal/pkg/core/observability"
	clientSetV "github.com/hellofresh/kangal/pkg/kubernetes/generated/clientset/versioned"
	"github.com/hellofresh/kangal/pkg/kubernetes/generated/informers/externalversions"
	"go.uber.org/zap"
	kubeInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

// Runner encapsulates all Kangal Controller dependencies
type Runner struct {
	Exporter       *prometheus.Exporter
	Logger         *zap.Logger
	KubeClient     kubernetes.Interface
	KangalClient   *clientSetV.Clientset
	KubeInformer   kubeInformers.SharedInformerFactory
	KangalInformer externalversions.SharedInformerFactory
	StatsReporter  observability.StatsReporter
}

// Run runs an instance of kubernetes kubeController
func Run(ctx context.Context, cfg Config, rr Runner) error {
	stopCh := make(chan struct{})

	c := NewController(cfg, rr.KubeClient, rr.KangalClient, rr.KubeInformer, rr.KangalInformer, rr.StatsReporter, rr.Logger)

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	rr.KangalInformer.Start(stopCh)
	rr.KubeInformer.Start(stopCh)

	if err := RunMetricsServer(ctx, cfg, rr, stopCh); err != nil {
		return fmt.Errorf("could not initialise Metrics Server: %w", err)
	}

	if err := c.Run(1, stopCh); err != nil {
		return fmt.Errorf("error running kubeController: %w", err)
	}
	return nil
}
