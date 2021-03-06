package volume

import (
	"context"
	"path/filepath"
	"reflect"
	"testing"

	. "gopkg.in/check.v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test(t *testing.T) { TestingT(t) }

type TestVolSuite struct{}

var _ = Suite(&TestVolSuite{})

func (s *TestVolSuite) TestCreatePVC(c *C) {
	// Create PVC
	ctx := context.Background()
	pvcSize := int64(1)
	ns := "kanister-pvc-test"
	targetVolID := "testVolID"
	annotations := map[string]string{"a1": "foo"}
	cli := fake.NewSimpleClientset()
	pvcName, err := CreatePVC(ctx, cli, ns, NoPVCNameSpecified, pvcSize, targetVolID, annotations)
	c.Assert(err, IsNil)
	pvc, err := cli.CoreV1().PersistentVolumeClaims(ns).Get(pvcName, metav1.GetOptions{})
	c.Assert(err, IsNil)

	c.Assert(len(pvc.Spec.AccessModes) >= 1, Equals, true)
	accessMode := pvc.Spec.AccessModes[0]
	c.Assert(accessMode, Equals, v1.ReadWriteOnce)
	capacity, ok := pvc.Spec.Resources.Requests[v1.ResourceStorage]
	c.Assert(ok, Equals, true)
	c.Assert(capacity.Value() >= int64(pvcSize*1024*1024*1024), Equals, true)
	eq := reflect.DeepEqual(annotations, pvc.ObjectMeta.Annotations)
	c.Assert(eq, Equals, true)
	c.Assert(len(pvc.Spec.Selector.MatchLabels) >= 1, Equals, true)
	label := pvc.Spec.Selector.MatchLabels[pvMatchLabelName]
	c.Assert(label, Equals, filepath.Base(targetVolID))
}
