package kube

import (
	"context"

	. "gopkg.in/check.v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/testing"
)

type PodRunnerTestSuite struct{}

var _ = Suite(&PodRunnerTestSuite{})

func (s *PodRunnerTestSuite) TestPodRunnerContextCanceled(c *C) {
	ctx, cancel := context.WithCancel(context.Background())
	cli := fake.NewSimpleClientset()
	cli.PrependReactor("create", "pods", func(action testing.Action) (handled bool, ret runtime.Object, err error) {
		return false, nil, nil
	})
	cli.PrependReactor("get", "pods", func(action testing.Action) (handled bool, ret runtime.Object, err error) {
		p := &v1.Pod{
			Status: v1.PodStatus{
				Phase: v1.PodRunning,
			},
		}
		return true, p, nil
	})
	deleted := make(chan struct{})
	cli.PrependReactor("delete", "pods", func(action testing.Action) (handled bool, ret runtime.Object, err error) {
		c.Log("Pod deleted due to Context Cancelled")
		close(deleted)
		return true, nil, nil
	})
	pr := NewPodRunner(cli, &PodOptions{})
	returned := make(chan struct{})
	go func() {
		_, err := pr.Run(ctx, makePodRunnerTestFunc(deleted))
		c.Assert(err, IsNil)
		close(returned)
	}()
	cancel()
	<-deleted
	<-returned
}

func makePodRunnerTestFunc(deleted chan struct{}) func(ctx context.Context, pod *v1.Pod) (map[string]interface{}, error) {
	return func(ctx context.Context, pod *v1.Pod) (map[string]interface{}, error) {
		<-deleted
		return nil, nil
	}
}
