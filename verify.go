package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

func main() {

	buildpack := libcnb.Buildpack{}

	_, err := toml.DecodeFile("buildpack.toml", &buildpack)
	if err != nil {
		panic(err)
	}

	context := libcnb.BuildContext{}
	context.Buildpack = buildpack

	logger := bard.NewLogger(os.Stderr)

	cr, err := libpak.NewConfigurationResolver(buildpack, &logger)
	if err != nil {
		panic(err)
	}

	defaultVersion, _ := cr.Resolve("BP_NEWRELIC_AGENT_VERSION")

	dr, err := libpak.NewDependencyResolver(context)
	if err != nil {
		panic(err)
	}

	err = VerifyDefaultVersion(defaultVersion, &dr.Dependencies)
	if err != nil {
		panic(err)
	}

	err = VerifyDependencyChecksums(&dr.Dependencies)
	if err != nil {
		panic(err)
	}

}

func VerifyDefaultVersion(defaultVersion string, depedencies *[]libpak.BuildpackDependency) error {
	containsDefaultVersion := false

	for _, dependency := range *depedencies {
		if dependency.Version == defaultVersion {
			containsDefaultVersion = true
		}
	}

	if !containsDefaultVersion {
		return fmt.Errorf("illegal default version :%s", defaultVersion)
	}

	return nil
}

func VerifyDependencyChecksums(deps *[]libpak.BuildpackDependency) error {
	var wg sync.WaitGroup
	for _, dependency := range *deps {
		wg.Add(1)
		go func(dep libpak.BuildpackDependency) {
			defer wg.Done()
			err := VerifyDependencyChecksum(&dep)
			if err != nil {
				panic(err)
			}

		}(dependency)
	}
	wg.Wait()
	return nil
}

func VerifyDependencyChecksum(dep *libpak.BuildpackDependency) error {
	tmp, err := os.CreateTemp("/tmp", fmt.Sprintf("%s_%s", dep.ID, dep.Version))
	if err != nil {
		return err
	}

	defer os.Remove(tmp.Name())

	err = DownloadFile(tmp.Name(), dep.URI)
	if err != nil {
		return err
	}

	sha256, err := CalculateChecksum(tmp.Name())
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Verify dependency checksum: %s:%s\n", dep.ID, dep.Version)
	fmt.Fprintf(os.Stderr, "\tExpect sha256: %s\n", dep.SHA256)
	fmt.Fprintf(os.Stderr, "\tActual sha256: %s\n", sha256)

	if sha256 != dep.SHA256 {
		return fmt.Errorf("checksum mismatched.\n\texpect sha256: %s\n\tactual sha256: %s", dep.SHA256, sha256)
	}

	return nil
}

func DownloadFile(filepath string, url string) error {

	fmt.Fprintf(os.Stderr, "Download file: %s\n", url)

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)

	if res.StatusCode > 299 {
		return fmt.Errorf("response failed with status code: %d", res.StatusCode)
	}

	fmt.Fprintf(os.Stderr, "Download Succeeded: %d, URL: %s\n", res.StatusCode, url)

	return err
}

func CalculateChecksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
