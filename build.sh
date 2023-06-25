#!/bin/bash

# Define the project directory
project_dir=$(pwd)

# Define the output directory
output_dir=$project_dir/bin

# Define the entry points of the applications
apps=("apps/api" "apps/worker")

# Create the output directory
mkdir -p $output_dir

# Function to build the app
build() {
    CGO_ENABLED=1 GOOS=$1 GOARCH=$2 go build -o $output_dir/$3.$1.$2
}

# Loop over each app and build it
for app in ${apps[@]}; do
    app_name=$(basename $app)

    # Change to app directory
    cd $project_dir/$app

    # Build for local environment
    build $(go env GOOS) $(go env GOARCH) $app_name
done

# Change back to project directory
cd $project_dir

echo "Compilation completed. The binaries can be found in the $output_dir directory."

