package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	humanize "github.com/dustin/go-humanize"
	"github.com/prometheus/client_golang/prometheus"
)

// Cli Environment client for Docker
var cli *client.Client
var r = rand.New(rand.NewSource(99)) // DEBUG
// Data to export
type Data struct {
	Infos []ContainerInfos
}

// ContainerInfos informations
type ContainerInfos struct {
	Image         string // repository:tag
	ContainerID10 string // unique tag
	State         string // running, exited, ...
	RawMemory     uint64 // usage
	HumanMemory   string
	Created       int64
	HumanCreated  string
}

func (ci *ContainerInfos) stats(containerID string) {
	readCloser, err := cli.ContainerStats(context.Background(), containerID, false)
	defer readCloser.Body.Close()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(readCloser.Body)
	for scanner.Scan() {
		var stats types.StatsJSON
		err = json.NewDecoder(strings.NewReader(scanner.Text())).Decode(&stats)
		if err != nil {
			panic(err)
		}
		ci.RawMemory = stats.MemoryStats.Usage
		ci.HumanMemory = humanize.Bytes(ci.RawMemory)
	}
}

// ListContainers list events for containers
func ListContainers() []ContainerInfos {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}
	var description []ContainerInfos
	description = make([]ContainerInfos, len(containers))
	for i, container := range containers {
		infos := ContainerInfos{}
		infos.Image = container.Image
		infos.ContainerID10 = container.ID[:10]
		infos.State = container.State
		infos.Created = container.Created
		infos.HumanCreated = humanize.Time(time.Unix(container.Created, 0))
		infos.stats(container.ID)
		description[i] = infos
	}
	return description
}

// ImageInfos informations
type ImageInfos struct {
	Repository     string
	Tag            string
	RawSize        int64
	HumanSize      string
	DateOfCreation string
	RawDuration    int64
	HumanDuration  string
}

// ListImages list events for images
func ListImages() []ImageInfos {
	//cli, _ = client.NewEnvClient() // DEBUG
	cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	images, _ := cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	var description []ImageInfos
	description = make([]ImageInfos, len(images))
	for i, image := range images {
		response, _, err := cli.ImageInspectWithRaw(context.Background(), image.ID)
		if err != nil {
			panic(err)
		}
		s := strings.Split(image.RepoTags[0], ":")
		repo, tag := s[0], s[1]
		size := image.Size
		date := response.Created
		duration := image.Created
		description[i] = ImageInfos{repo, tag, size, humanize.Bytes(uint64(image.Size)), date, duration, humanize.Time(time.Unix(image.Created, 0))}
		fmt.Printf("%s:%s s%d d%d\n", repo, tag, duration, size)
	}
	return description
}

// InitScraping scraping
func InitScraping() error {
	var err error
	cli, err = client.NewEnvClient() // client pour les metriques docker
	return err
}

//------------Exporters----------------//

func (e *Exporter) processMetrics(data *Data, ch chan<- prometheus.Metric) error {
	for _, d := range data.Infos {
		e.setMetrics(d.Image, d.ContainerID10, d.RawMemory)
	}
	fmt.Printf("Process\n")
	return nil
}

func (e *Exporter) gatherData(ch chan<- prometheus.Metric) (*Data, error) {
	var description = ListContainers()
	data := Data{description}
	fmt.Printf("Gather\n")
	return &data, nil
}
