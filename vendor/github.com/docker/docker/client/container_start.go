package client // import "github.com/docker/docker/client"

import (
	"context"
	"net/url"
	"k8s.io/klog" //add
	"regexp" //add
	// "os"//add
	// "time"//add
	"github.com/otiai10/copy"
	"github.com/docker/docker/api/types"
)

// ContainerStart sends a request to the docker daemon to start a container.
func (cli *Client) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	query := url.Values{}
	// klog.V(3).Infof("so-ta: ctx=%v, containerID=%v, opt=%v", ctx, containerID, options)
	klog.V(3).Infof("so-ta: 1 container_start.go")
	klog.V(3).Infof("so-ta: containerID + , = %v", containerID)

	req := regexp.MustCompile(`^.*,`)
	kService := req.ReplaceAllString(containerID, "")
	req = regexp.MustCompile(`(,.*$)`)
	containerID = req.ReplaceAllString(containerID, "")
	
	if kService == "helloworld-go" {
		klog.V(3).Infof("so-ta: checkpoint stop!")
		// time.Sleep(30 * time.Second)
		klog.V(3).Infof("so-ta: checkpoint start!")
		newname := "/var/lib/docker/containers/" + containerID + "/checkpoints/cp-helloworld-go"
		err := copy.Copy("/cp/cp-helloworld-go/", newname)
		klog.V(3).Infof("so-ta: err=%v",err)
		// klog.V(3).Infof("so-ta: containerID=%v", containerID)
		// klog.V(3).Infof("so-ta: kService =%v", kService)
		options.CheckpointID = "cp-helloworld-go"
		// options.CheckpointDir = "/cp/"
	}
	if len(options.CheckpointID) != 0 {
		query.Set("checkpoint", options.CheckpointID)
	}
	if len(options.CheckpointDir) != 0 {
		query.Set("checkpoint-dir", options.CheckpointDir)
	}

	

	resp, err := cli.post(ctx, "/containers/"+containerID+"/start", query, nil, nil)
	klog.V(3).Infof("so-ta: err=%v", err)
	ensureReaderClosed(resp)
	return err
}
