package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
)

type IamUserPolicyAttachement struct {
	svc        *iam.IAM
	policyArn  string
	policyName string
	roleName   string
}

func (n *IamNuke) ListUserPolicyAttachements() ([]Resource, error) {
	resp, err := n.Service.ListUsers(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, role := range resp.Users {
		resp, err := n.Service.ListAttachedUserPolicies(
			&iam.ListAttachedUserPoliciesInput{
				UserName: role.UserName,
			})
		if err != nil {
			return nil, err
		}

		for _, pol := range resp.AttachedPolicies {
			resources = append(resources, &IamUserPolicyAttachement{
				svc:        n.Service,
				policyArn:  *pol.PolicyArn,
				policyName: *pol.PolicyName,
				roleName:   *role.UserName,
			})
		}
	}

	return resources, nil
}

func (e *IamUserPolicyAttachement) Remove() error {
	_, err := e.svc.DetachUserPolicy(
		&iam.DetachUserPolicyInput{
			PolicyArn: &e.policyArn,
			UserName:  &e.roleName,
		})
	if err != nil {
		return err
	}

	return nil
}

func (e *IamUserPolicyAttachement) String() string {
	return fmt.Sprintf("%s -> %s", e.roleName, e.policyName)
}
