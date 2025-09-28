#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

print_status "Generating environment files from templates..."

# Create .env files from templates
SERVICES=("inventory" "order" "payment")

for service in "${SERVICES[@]}"; do
    template_file="$SCRIPT_DIR/${service}.env.template"
    env_file="$PROJECT_ROOT/${service}/.env"
    
    if [ -f "$template_file" ]; then
        print_status "Creating $env_file from $template_file"
        cp "$template_file" "$env_file"
    else
        print_warning "Template file $template_file not found, skipping $service"
    fi
done

# Create main .env file
main_template="$SCRIPT_DIR/env.template"
main_env="$PROJECT_ROOT/.env"

if [ -f "$main_template" ]; then
    print_status "Creating $main_env from $main_template"
    cp "$main_template" "$main_env"
else
    print_warning "Main template file $main_template not found"
fi

print_status "Environment files generated successfully!"
print_status "You can now modify the .env files to customize your configuration."