package util

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Workload struct {
    APIVersion string `yaml:"apiVersion"`
    Kind       string `yaml:"kind"`
    Metadata   struct {
        Name      string `yaml:"name"`
        Namespace string `yaml:"namespace"`
    } `yaml:"metadata"`
    Spec struct {
        Template struct {
            Spec struct {
                Containers []struct {
                    Resources struct {
                        Limits struct {
                            CPU    string `yaml:"cpu"`
                            Memory string `yaml:"memory"`
                        } `yaml:"limits"`
                        Requests struct {
                            CPU    string `yaml:"cpu"`
                            Memory string `yaml:"memory"`
                        } `yaml:"requests"`
                    } `yaml:"resources"`
                } `yaml:"containers"`
            } `yaml:"spec"`
        } `yaml:"template"`
    } `yaml:"spec"`
}

// GenerateWorkloads walks through YAML files in the specified folderPath,
// decodes them into Workload structs, and returns a slice of Workload structs.
func GenerateWorkloads(folderPath string) ([]map[string]interface{}, error) {
	var workloads []map[string]interface{}

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml")) {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			decoder := yaml.NewDecoder(strings.NewReader(string(fileContent)))
			for {
				var workload Workload
				err := decoder.Decode(&workload)
				if err != nil {
					if err.Error() == "EOF" {
						break
					}
					return err
				}

				if workload.APIVersion == "apps/v1" && (workload.Kind == "Deployment" || workload.Kind == "StatefulSet") {
					for _, container := range workload.Spec.Template.Spec.Containers {
						workloadInfo := map[string]interface{}{
							"metadata-name":           workload.Metadata.Name,
							"metadata-namespace":      workload.Metadata.Namespace,
							"resources-limit-cpu":     getStringOrDefault(container.Resources.Limits.CPU, "not defined"),
							"resources-limit-memory":  getStringOrDefault(container.Resources.Limits.Memory, "not defined"),
							"resources-requests-cpu":  getStringOrDefault(container.Resources.Requests.CPU, "not defined"),
							"resources-requests-memory": getStringOrDefault(container.Resources.Requests.Memory, "not defined"),
						}
						workloads = append(workloads, workloadInfo)
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return workloads, nil
}

// getStringOrDefault fix the issue when resources are empty/nul
func getStringOrDefault(value, defaultValue string) string {
    if value == "" {
        return defaultValue
    }
    return value
}
