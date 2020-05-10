/*

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

package eks

import (
	"github.com/keikoproj/instance-manager/api/v1alpha1"
	"github.com/keikoproj/instance-manager/controllers/common"
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
)

func (ctx *EksInstanceGroupContext) Delete() error {
	var (
		instanceGroup = ctx.GetInstanceGroup()
		state         = ctx.GetDiscoveredState()
		role          = state.GetRole()
		roleARN       = aws.StringValue(role.Arn)
	)

	instanceGroup.SetState(v1alpha1.ReconcileDeleting)
	// delete scaling group
	err := ctx.DeleteScalingGroup()
	if err != nil {
		return errors.Wrap(err, "failed to delete scaling group")
	}

	// if scaling group is deleted, defer removal from aws-auth
	defer common.RemoveAuthConfigMap(ctx.KubernetesClient.Kubernetes, []string{roleARN})

	// delete launchconfig
	err = ctx.DeleteLaunchConfiguration()
	if err != nil {
		return errors.Wrap(err, "failed to delete launch configuration")
	}

	// delete the managed IAM role if one was created
	err = ctx.DeleteManagedRole()
	if err != nil {
		return errors.Wrap(err, "failed to delete scaling group role")
	}

	return nil
}

func (ctx *EksInstanceGroupContext) DeleteScalingGroup() error {
	var (
		state         = ctx.GetDiscoveredState()
		scalingGroup  = state.GetScalingGroup()
		instanceGroup = ctx.GetInstanceGroup()
		asgName       = aws.StringValue(scalingGroup.AutoScalingGroupName)
	)

	if !state.HasScalingGroup() {
		return nil
	}

	err := ctx.AwsWorker.DeleteScalingGroup(asgName)
	if err != nil {
		return err
	}
	ctx.Log.Info("deleted scaling group", "instancegroup", instanceGroup.GetName(), "scalinggroup", asgName)
	return nil
}

func (ctx *EksInstanceGroupContext) DeleteLaunchConfiguration() error {
	var (
		state         = ctx.GetDiscoveredState()
		instanceGroup = ctx.GetInstanceGroup()
		lcName        = state.GetActiveLaunchConfigurationName()
	)

	if !state.HasLaunchConfiguration() {
		return nil
	}

	err := ctx.AwsWorker.DeleteLaunchConfig(lcName)
	if err != nil {
		return err
	}
	ctx.Log.Info("deleted launch config", "instancegroup", instanceGroup.GetName(), "launchconfig", lcName)
	return nil
}

func (ctx *EksInstanceGroupContext) DeleteManagedRole() error {
	var (
		instanceGroup      = ctx.GetInstanceGroup()
		configuration      = instanceGroup.GetEKSConfiguration()
		state              = ctx.GetDiscoveredState()
		additionalPolicies = configuration.GetManagedPolicies()
		role               = state.GetRole()
		roleName           = aws.StringValue(role.RoleName)
	)

	if !state.HasRole() || configuration.HasExistingRole() {
		return nil
	}

	managedPolicies := ctx.GetManagedPoliciesList(additionalPolicies)

	err := ctx.AwsWorker.DeleteScalingGroupRole(roleName, managedPolicies)
	if err != nil {
		return err
	}
	ctx.Log.Info("deleted scaling group role", "instancegroup", instanceGroup.GetName(), "iamrole", roleName)
	return nil
}
