package client // import "github.com/docker/docker/client"

import (
	"context"
	"net/url"
	"k8s.io/klog" //add
	"regexp" //add
	"os"//add
	"bufio"//add
	"github.com/otiai10/copy"
	"github.com/docker/docker/api/types"
)
const file = "checkpoint-list.dat"
// ContainerStart sends a request to the docker daemon to start a container.
func (cli *Client) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	query := url.Values{}
	klog.V(3).Infof("so-ta: 1 container_start.go")
	// klog.V(3).Infof("so-ta: containerID + , = %v", containerID)
	dir, err := os.Getwd()
	klog.V(3).Infof("so-ta: dir = %v", dir)
	req := regexp.MustCompile(`^.*,`)
	kService := req.ReplaceAllString(containerID, "")
	req = regexp.MustCompile(`(,.*$)`)
	containerID = req.ReplaceAllString(containerID, "")
	
	if checkCP(kService, file) {
		newname := "/var/lib/docker/containers/" + containerID + "/checkpoints/" + kService
		err := copy.Copy("/cp/" + kService, newname)
		klog.V(3).Infof("so-ta: err=%v",err)
		// klog.V(3).Infof("so-ta: containerID=%v", containerID)
		// klog.V(3).Infof("so-ta: kService =%v", kService)
		options.CheckpointID = kService
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

func checkCP(kService string, file string) bool{
	klog.V(3).Infof("so-ta: kService: %v", kService)
	fp, err := os.Open(file)
    if err != nil {
        klog.V(3).Infof("so-ta: No such file")
    }
    defer fp.Close()

    scanner := bufio.NewScanner(fp)

    for scanner.Scan() {
        if scanner.Text() == kService {
			klog.V(3).Infof("so-ta: match! checkpoint: %v", kService)
			return true
		}
    }
	return false
}