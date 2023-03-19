package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const NATURE_URL = "https://pokeapi.co/api/v2/nature/calm"
const GENERATION_URL = "https://pokeapi.co/api/v2/generation/generation-ii"

type NatureResult struct {
	Name       string       `json:"name"`
	StatChange []StatChange `json:"pokeathlon_stat_changes"`
}

type StatChange struct {
	MaxChange float32        `json:"max_change"`
	Stats     PokeathlonStat `json:"pokeathlon_stat"`
}

type PokeathlonStat struct {
	Name string `json:"name"`
}

type GenerationResult struct {
	Name          string           `json:"name"`
	VersionGroups []version_groups `json:"version_groups"`
}

type version_groups struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func timeCheck(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func getNatureData() {
	natureReq, _ := http.NewRequest("GET", NATURE_URL, nil)

	natureRes, _ := http.DefaultClient.Do(natureReq)

	defer natureRes.Body.Close()
	body, _ := ioutil.ReadAll(natureRes.Body)
	var responseObject NatureResult
	json.Unmarshal(body, &responseObject)

	fmt.Println(responseObject.Name)
	fmt.Println(len(responseObject.StatChange))

	for _, StatChange := range responseObject.StatChange {
		fmt.Println(StatChange.Stats.Name)
	}
}

func getGenerationData() {
	generationReq, _ := http.NewRequest("GET", GENERATION_URL, nil)

	generationRes, _ := http.DefaultClient.Do(generationReq)

	defer generationRes.Body.Close()
	body, _ := ioutil.ReadAll(generationRes.Body)
	var responseObject GenerationResult
	json.Unmarshal(body, &responseObject)

	fmt.Println(responseObject.Name)
	fmt.Println(len(responseObject.VersionGroups))

	for _, versionGroup := range responseObject.VersionGroups {
		fmt.Println(versionGroup.Name)
	}
}

func main() {

	defer timeCheck(time.Now(), "Fetching time")

	fmt.Println("Starting concurrent calls...")

	var waitGroup sync.WaitGroup
	waitGroup.Add(3)

	go func() {
		getNatureData()
		waitGroup.Done()
	}()

	go func() {
		getGenerationData()
		waitGroup.Done()
	}()

	waitGroup.Wait()
}
