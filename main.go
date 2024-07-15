package main

import (
	"cloudResource/util"

	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
    folderPath := "./templates-ias" // Change to your folder path
    logsFolderPath := "./logs"

    // unMarshall and generate workloads from YAML
    workloads, err := util.GenerateWorkloads(folderPath)
    if err != nil {
        log.Fatalf("error walking the path %v: %v", folderPath, err)
    }

    // Marshall into JSON
    jsonData, err := json.MarshalIndent(workloads, "", "  ")
    if err != nil {
        log.Fatalf("error marshaling to json: %v", err)
    }

    // Current time
    t := time.Now()
    timeNow := fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())

    // make sure there is a log folder to save the files. If there is, skip
    if _, err := os.Stat(logsFolderPath); os.IsNotExist(err) {
        if err := os.Mkdir(logsFolderPath, 0755); err != nil {
            log.Fatalf("error creating logs folder: %v", err)
        }
    }

    // Save JSON
    logFilePath := filepath.Join(logsFolderPath, fmt.Sprintf("workloads_%s.json", timeNow))
    if err := os.WriteFile(logFilePath, jsonData, 0644); err != nil {
        log.Fatalf("error writing json file: %v", err)
    }
    fmt.Printf("Workloads information saved to %s\n", logFilePath)

    // Save CSV
    csvFilePath := filepath.Join(logsFolderPath, fmt.Sprintf("workloads_%s.csv", timeNow))
    csvFile, err := os.Create(csvFilePath)
    if err != nil {
        log.Fatalf("error creating csv file: %v", err)
    }
    defer csvFile.Close()

    csvWriter := csv.NewWriter(csvFile)
    defer csvWriter.Flush()

    // Write CSV headers
    headers := []string{
        "metadata-name", 
        "metadata-namespace", 
        "resources-limit-cpu", 
        "resources-limit-memory", 
        "resources-requests-cpu", 
        "resources-requests-memory",
    }
    if err := csvWriter.Write(headers); err != nil {
        log.Fatalf("error writing csv headers: %v", err)
    }

    // Write CSV records
    for _, workload := range workloads {
        record := []string{
            workload["metadata-name"].(string),
            workload["metadata-namespace"].(string),
            workload["resources-limit-cpu"].(string),
            workload["resources-limit-memory"].(string),
            workload["resources-requests-cpu"].(string),
            workload["resources-requests-memory"].(string),
        }
        if err := csvWriter.Write(record); err != nil {
            log.Fatalf("error writing csv record: %v", err)
        }
    }

    fmt.Printf("Workloads information saved to %s\n", csvFilePath)
}
