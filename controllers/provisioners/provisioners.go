package provisioners

import (
	"github.com/go-logr/logr"
	"github.com/keikoproj/instance-manager/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	awsprovider "github.com/keikoproj/instance-manager/controllers/providers/aws"
	kubeprovider "github.com/keikoproj/instance-manager/controllers/providers/kubernetes"
)

var (
	log = ctrl.Log.WithName("provisioners")
)

const (
	TagClusterName            = "instancegroups.keikoproj.io/ClusterName"
	TagInstanceGroupName      = "instancegroups.keikoproj.io/InstanceGroup"
	TagInstanceGroupNamespace = "instancegroups.keikoproj.io/Namespace"
	TagClusterOwnershipFmt    = "kubernetes.io/cluster/%s"
	TagKubernetesCluster      = "KubernetesCluster"
)

type ProvisionerInput struct {
	AwsWorker     awsprovider.AwsWorker
	Kubernetes    kubeprovider.KubernetesClientSet
	InstanceGroup *v1alpha1.InstanceGroup
	Configuration *corev1.ConfigMap
	Log           logr.Logger
}
