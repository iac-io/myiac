package identity

import (
	"fmt"
	"github.com/dfernandezm/myiac/internal/gcp"
)

type Iam interface {
	AuthMaterialFactory(provider string) AuthMaterial
}

type AuthMaterial interface {
	ObtainKey() string
}

type GcpClient interface {
	KeyForServiceAccount(saEmail string, recreateKey bool) (string, error)
}

type gcpClient struct {
}

func NewGcpClient() *gcpClient {
	return &gcpClient{}
}

// KeyForServiceAccount generates a new JSON key from a given service account email. The key is returned as string,
// and it's cached for the given email unless 'recreateKey' has been set to 'true'
func (gc *gcpClient) KeyForServiceAccount(saEmail string, recreateKey bool) (string, error) {
	return gcp.KeyForServiceAccount(saEmail, recreateKey)
}

type GcpAuthMaterial struct {
	serviceAccountEmail string
	gcpClient GcpClient
}

func NewGcpAuthMaterial(serviceAccountEmail string, gcpClient GcpClient) *GcpAuthMaterial {
	return &GcpAuthMaterial{
		serviceAccountEmail:serviceAccountEmail,
		gcpClient:gcpClient,
	}
}

func (gam *GcpAuthMaterial) ObtainKey() string {
	fmt.Printf("Obtaining service account key for email %s", gam.serviceAccountEmail)
	jsonkey, err := gam.gcpClient.KeyForServiceAccount(gam.serviceAccountEmail, false)
	if err != nil {
		fmt.Errorf("trying to regenerate key for email %s as error occurred %v",
			gam.serviceAccountEmail, err)
		jsonkey, err2 := gcp.KeyForServiceAccount(gam.serviceAccountEmail, true)

		if err2 != nil {
			errReturn := fmt.Errorf("cannot obtain key from: %s %w", gam.serviceAccountEmail, err2)
			panic(errReturn)
		}

		return jsonkey
	}

	return jsonkey
}