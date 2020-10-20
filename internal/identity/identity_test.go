package identity

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type fakeGcpClient struct {

}

func (fgc *fakeGcpClient) KeyForServiceAccount(saEmail string, recreateKey bool) (string, error) {
	return "testKey", nil
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestCreateNewKey(t *testing.T) {
	fgc := fakeGcpClient{}
	testSaEmail :="testAccount@gcloudserviceaccount.com"
	gcpAuth := NewGcpAuthMaterial(testSaEmail, &fgc)
	key := gcpAuth.ObtainKey()
	assert.Equal(t, "testKey", key)
}

//func TestObtainExistingKey(t *testing.T) {
//
//}
