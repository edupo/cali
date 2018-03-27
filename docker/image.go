package docker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
	"gopkg.in/cheggaaa/pb.v1"
)

// ProgressDetail records the progress achieved downloading an image
type ProgressDetail struct {
	Current int `json:"current,omitempty"`
	Total   int `json:"total,omitempty"`
}

// CreateResponse is the response from Docker API when pulling an image
type CreateResponse struct {
	ID             string         `json:"id"`
	Status         string         `json:"status"`
	ProgressDetail ProgressDetail `json:"progressDetail"`
	Progress       string         `json:"progress,omitempty"`
}

// ImageExists determines if an image exist locally
func (c *Client) ImageExists(image string) bool {
	_, _, err := c.Cli.ImageInspectWithRaw(context.Background(), image)

	// TODO: Safe assumption?
	if err != nil {
		log.WithField("image", image).
			Debug(err)
		return false
	}
	return true
}

// PullImage performs an image pull if that image does not exists locally.
// TODO: Image autoupdate as a parameter
func (c *Client) PullImage(image string) error {

	// TODO: Check for changes in the remote
	if !c.ImageExists(image) {
		log.WithFields(log.Fields{
			"image": image,
		}).Info("Pulling image layers... please wait")

		resp, err := c.Cli.ImagePull(context.Background(), image, types.ImagePullOptions{})

		if err != nil {
			return fmt.Errorf("API could not fetch \"%s\": %s", image, err)
		}
		scanner := bufio.NewScanner(resp)
		var cr CreateResponse
		bar := pb.New(1)
		// Send progress bar to stderr to keep stdout clean when piping
		bar.Output = os.Stderr
		bar.ShowCounters = true
		bar.ShowTimeLeft = false
		bar.ShowSpeed = false
		bar.Prefix("          ")
		bar.Postfix("          ")
		started := false

		for scanner.Scan() {
			txt := scanner.Text()
			byt := []byte(txt)

			if err := json.Unmarshal(byt, &cr); err != nil {
				return fmt.Errorf("Error decoding json from create image API: %s", err)
			}

			if cr.Status == "Downloading" {

				if !started {
					fmt.Print("\n")
					bar.Total = int64(cr.ProgressDetail.Total)
					bar.Start()
					started = true
				}
				bar.Total = int64(cr.ProgressDetail.Total)
				bar.Set(cr.ProgressDetail.Current)
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("Failed to get logs: %s", err)
		}
		bar.Finish()
		fmt.Print("\n")
	}
	return nil
}
