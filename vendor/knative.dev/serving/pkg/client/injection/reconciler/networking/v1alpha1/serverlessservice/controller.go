/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by injection-gen. DO NOT EDIT.

package serverlessservice

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	watch "k8s.io/apimachinery/pkg/watch"
	scheme "k8s.io/client-go/kubernetes/scheme"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	record "k8s.io/client-go/tools/record"
	client "knative.dev/pkg/client/injection/kube/client"
	controller "knative.dev/pkg/controller"
	logging "knative.dev/pkg/logging"
	versionedscheme "knative.dev/serving/pkg/client/clientset/versioned/scheme"
	injectionclient "knative.dev/serving/pkg/client/injection/client"
	serverlessservice "knative.dev/serving/pkg/client/injection/informers/networking/v1alpha1/serverlessservice"
)

const (
	defaultControllerAgentName = "serverlessservice-controller"
	defaultFinalizerName       = "serverlessservice"
)

func NewImpl(ctx context.Context, r Interface) *controller.Impl {
	logger := logging.FromContext(ctx)

	serverlessserviceInformer := serverlessservice.Get(ctx)

	recorder := controller.GetEventRecorder(ctx)
	if recorder == nil {
		// Create event broadcaster
		logger.Debug("Creating event broadcaster")
		eventBroadcaster := record.NewBroadcaster()
		watches := []watch.Interface{
			eventBroadcaster.StartLogging(logger.Named("event-broadcaster").Infof),
			eventBroadcaster.StartRecordingToSink(
				&v1.EventSinkImpl{Interface: client.Get(ctx).CoreV1().Events("")}),
		}
		recorder = eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: defaultControllerAgentName})
		go func() {
			<-ctx.Done()
			for _, w := range watches {
				w.Stop()
			}
		}()
	}

	c := &reconcilerImpl{
		Client:        injectionclient.Get(ctx),
		Lister:        serverlessserviceInformer.Lister(),
		Recorder:      recorder,
		FinalizerName: defaultFinalizerName,
		reconciler:    r,
	}
	impl := controller.NewImpl(c, logger, "serverlessservices")

	return impl
}

func init() {
	versionedscheme.AddToScheme(scheme.Scheme)
}
