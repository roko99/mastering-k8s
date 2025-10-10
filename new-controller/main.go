// main.go
package main

import (
	"context"
	"flag"
	"net/http"

	// k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"

	newv1 "github.com/roko99/mastering-k8s/new-controller/api/v1alpha1"
	"github.com/roko99/mastering-k8s/new-controller/controllers"
)

func main() {
	var (
		metricsAddr          string
		enableLeaderElection bool
	)

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false, "Enable leader election for controller manager.")
	flag.Parse()

	scheme := runtime.NewScheme()
	utilruntime.Must(newv1.AddToScheme(scheme))

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Create a variable to hold the client that will be set after manager creation
	var k8sClient client.Client

	// Create health check handler that uses the client
	healthCheck := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		list := &newv1.NewResourceList{}
		if err := k8sClient.List(ctx, list); err != nil {
			// if k8serrors.IsNotFound(err) || meta.IsNoMatchError(err) {
			if meta.IsNoMatchError(err) {
				http.Error(w, "NewResource CRD not found", http.StatusServiceUnavailable)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
		Scheme: scheme,
		Metrics: server.Options{
			BindAddress: metricsAddr,
			ExtraHandlers: map[string]http.Handler{
				"/healthz": healthCheck,
			},
		},
		HealthProbeBindAddress: "", // Disable separate health probe server
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "newresource-controller",
	})
	if err != nil {
		panic(err)
	}

	// Set the client after manager is created
	k8sClient = mgr.GetClient()

	if err := (&controllers.NewResourceReconciler{
		Client: mgr.GetClient(),
	}).SetupWithManager(mgr); err != nil {
		panic(err)
	}

	// Start the manager (this is blocking)
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		panic(err)
	}
}
