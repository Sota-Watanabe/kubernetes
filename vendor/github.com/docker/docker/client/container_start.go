package client // import "github.com/docker/docker/client"

import (
	"context"
	"net/url"
	"k8s.io/klog" //add

	"github.com/docker/docker/api/types"
)

// ContainerStart sends a request to the docker daemon to start a container.
func (cli *Client) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	query := url.Values{}
	klog.V(3).Infof("so-ta: ctx=%v, containerID=%v, opt=%v", ctx, containerID, options)
	klog.V(3).Infof("so-ta: 1 container_start.go")
	if len(options.CheckpointID) != 0 {
		query.Set("checkpoint", options.CheckpointID)
	}
	if len(options.CheckpointDir) != 0 {
		query.Set("checkpoint-dir", options.CheckpointDir)
	}

	resp, err := cli.post(ctx, "/containers/"+containerID+"/start", query, nil, nil)
	ensureReaderClosed(resp)
	return err
}
