package docker

import (
	"fmt"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fpgeek/docker-gc/env"
	"github.com/fsouza/go-dockerclient"
)

var (
	dockerClient *docker.Client
)

func init() {
	var err error
	if env.DockerCertPath != "" {
		dockerClient, err = docker.NewTLSClient(env.DockerHost,
			fmt.Sprintf("%s/cert.pem", env.DockerCertPath),
			fmt.Sprintf("%s/key.pem", env.DockerCertPath),
			fmt.Sprintf("%s/ca.pem", env.DockerCertPath),
		)
	} else {
		dockerClient, err = docker.NewClient(env.DockerHost)
	}
	if err != nil {
		log.Fatal(err)
	}
}

type dockerImages []docker.APIImages

func (d dockerImages) Len() int {
	return len(d)
}

func (d dockerImages) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d dockerImages) Less(i, j int) bool {
	return d[i].Created > d[j].Created
}

// RemoveExistedContainers removes containers that existed
func RemoveExistedContainers() error {
	existedContainers, err := dockerClient.ListContainers(docker.ListContainersOptions{
		Filters: map[string][]string{
			"status": {"exited"},
		},
	})
	if err != nil {
		return err
	}
	for _, c := range existedContainers {
		if err := dockerClient.RemoveContainer(docker.RemoveContainerOptions{
			ID: c.ID,
		}); err != nil {
			return err
		}
	}
	return nil
}

// RemoveStaleImages removes stale images
func RemoveStaleImages() error {
	runningImageIDsSet, err := listRunningImageIDsSet()
	if err != nil {
		return err
	}
	allImages, err := dockerClient.ListImages(docker.ListImagesOptions{})
	if err != nil {
		return err
	}
	staleImageSets := map[string]dockerImages{}
	for _, image := range allImages {
		if !runningImageIDsSet[image.ID] {
			for _, fullRepoName := range image.RepoTags {
				repoNames := strings.Split(fullRepoName, ":")
				repo := repoNames[0]
				if images, ok := staleImageSets[repo]; ok {
					images = append(images, image)
					staleImageSets[repo] = images
				} else {
					staleImageSets[repo] = dockerImages{image}
				}
			}
		}
	}

	staleImageList := dockerImages{}
	for _, images := range staleImageSets {
		sort.Sort(images)
		if len(images) > env.KeepStaleImageCount {
			staleImageList = append(staleImageList, images[env.KeepStaleImageCount:]...)
		}
	}

	for _, image := range staleImageList {
		log.WithField("image", fmt.Sprintf("%#+v", image)).Info("remove image")
		if err := dockerClient.RemoveImage(image.ID); err != nil {
			log.WithField("image", err.Error()).Error("removing image has falied")
		}
	}

	return nil
}

func listRunningImageIDsSet() (map[string]bool, error) {
	runningContainers, err := dockerClient.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return nil, err
	}
	imageIDsSet := map[string]bool{}
	for _, c := range runningContainers {
		container, err := dockerClient.InspectContainer(c.ID)
		if err != nil {
			return nil, err
		}
		log.WithField("container", fmt.Sprintf("%#+v", container)).Debug("running container")
		imageIDsSet[container.Image] = true
	}
	return imageIDsSet, nil
}
