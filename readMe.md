# YAML to JSON and CSV Converter

This Go script converts YAML files containing Kubernetes deployment configurations into JSON and CSV formats. It walks through a specified folder, extracts relevant workload information (like metadata and resource specifications), and saves them into separate JSON and CSV files.

## Dependencies

- Go (Golang)
- `gopkg.in/yaml.v2` package for YAML parsing (included in the Go standard library)

## How to run

- Create a folder "templates" and move all .yaml file there to be converted. 
- Run the go script: 'go run main.go'

## Output files

- JSON file (workloads.json) containing extracted workload information.
- CSV file (workloads.csv) containing workload information in tabular format.
