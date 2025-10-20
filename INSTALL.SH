#!/bin/bash

# Iran Intelligence Enterprise - Installation Script
set -e

echo "🔄 Starting Iran Intelligence Enterprise Installation..."
echo "📍 Detecting operating system..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Detect OS
OS="unknown"
case "$(uname -s)" in
    Linux*)     OS="linux";;
    Darwin*)    OS="macos";;
    CYGWIN*)    OS="windows";;
    MINGW*)     OS="windows";;
    *)          OS="unknown"
esac

# Detect if WSL
if [[ "$OS" == "linux" ]] && grep -q Microsoft /proc/version; then
    OS="wsl"
fi

echo -e "${BLUE}Detected OS: $OS${NC}"

# Function to check command existence
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to install dependencies
install_dependencies() {
    echo -e "${YELLOW}Installing dependencies for $OS...${NC}"
    
    case $OS in
        "linux"|"wsl")
            if command_exists apt; then
                # Ubuntu/Debian
                sudo apt update
                sudo apt install -y docker.io docker-compose git curl wget build-essential
            elif command_exists yum; then
                # CentOS/RHEL
                sudo yum install -y docker git curl wget
                sudo systemctl start docker
                sudo systemctl enable docker
            elif command_exists pacman; then
                # Arch Linux
                sudo pacman -S docker docker-compose git curl wget
                sudo systemctl start docker.service
                sudo systemctl enable docker.service
            fi
            ;;
        "macos")
            if ! command_exists brew; then
                /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
            fi
            brew install docker docker-compose git curl wget
            ;;
        "windows")
            echo "Please install Docker Desktop from https://docker.com"
            echo "Please install Git from https://git-scm.com"
            read -p "Press Enter after installing Docker Desktop and Git..."
            ;;
    esac
}

# Function to setup Docker
setup_docker() {
    if ! command_exists docker; then
        echo -e "${RED}Docker not found. Please install Docker first.${NC}"
        exit 1
    fi
    
    # Start Docker service on Linux
    if [[ "$OS" == "linux" || "$OS" == "wsl" ]]; then
        sudo systemctl start docker || true
        sudo systemctl enable docker || true
    fi
    
    # Test Docker
    if ! docker ps >/dev/null 2>&1; then
        echo -e "${RED}Docker is not running. Please start Docker Desktop (Windows/macOS) or Docker service (Linux).${NC}"
        exit 1
    fi
}

# Function to clone and setup project
setup_project() {
    echo -e "${YELLOW}Setting up project...${NC}"
    
    # Create project directory
    PROJECT_DIR="$HOME/iran-intel-enterprise"
    if [ ! -d "$PROJECT_DIR" ]; then
        mkdir -p "$PROJECT_DIR"
    fi
    
    cd "$PROJECT_DIR"
    
    # Copy necessary files (assuming script is run from project root)
    if [ -f "docker-compose.yml" ]; then
        cp docker-compose.yml "$PROJECT_DIR/"
    fi
    
    if [ -d "config" ]; then
        cp -r config "$PROJECT_DIR/"
    fi
    
    # Create necessary directories
    mkdir -p {data,logs,proxies,models}
}

# Function to setup environment
setup_environment() {
    echo -e "${YELLOW}Setting up environment...${NC}"
    
    cd "$PROJECT_DIR"
    
    # Create .env file if not exists
    if [ ! -f ".env" ]; then
        cat > .env << EOF
# Iran Intelligence Enterprise Configuration
NODE_ENV=production
IRAN_DEPLOYMENT=true

# API Configuration
API_PORT=8080
UI_PORT=3000
AI_PORT=5000

# Database
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_DB=iran_intel
POSTGRES_USER=intel_user
POSTGRES_PASSWORD=$(openssl rand -base64 32)

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# Security
ENCRYPTION_KEY=$(openssl rand -base64 32)
JWT_SECRET=$(openssl rand -base64 32)

# Iranian Operators
MCI_API_ENABLED=true
MTN_API_ENABLED=true
RIGHTEL_API_ENABLED=true
SHATEL_API_ENABLED=true

# Platform Limits
MAX_CONCURRENT_AGENTS=50
MAX_REQUESTS_PER_MINUTE=1000
EOF
        echo -e "${GREEN}Environment file created${NC}"
    fi
}

# Function to build and start services
start_services() {
    echo -e "${YELLOW}Building and starting services...${NC}"
    
    cd "$PROJECT_DIR"
    
    # Pull latest images
    docker-compose pull || true
    
    # Build services
    docker-compose build
    
    # Start services
    docker-compose up -d
    
    # Wait for services to be healthy
    echo -e "${YELLOW}Waiting for services to start...${NC}"
    sleep 30
    
    # Check service status
    if docker-compose ps | grep -q "Up"; then
        echo -e "${GREEN}✅ All services are running!${NC}"
    else
        echo -e "${RED}❌ Some services failed to start${NC}"
        docker-compose logs
        exit 1
    fi
}

# Function to display success message
show_success() {
    echo -e "${GREEN}"
    echo "███████ ███████ ██████  ███████ ███████ ████████ ███████ ██████  "
    echo "██      ██      ██   ██ ██      ██         ██    ██      ██   ██ "
    echo "█████   █████   ██████  █████   ███████    ██    █████   ██████  "
    echo "██      ██      ██   ██ ██           ██    ██    ██      ██   ██ "
    echo "██      ███████ ██   ██ ███████ ███████    ██    ███████ ██   ██ "
    echo -e "${NC}"
    
    echo -e "${BLUE}🎉 Installation Completed Successfully!${NC}"
    echo ""
    echo -e "${YELLOW}📊 Access Points:${NC}"
    echo -e "  Web UI: ${GREEN}http://localhost:3000${NC}"
    echo -e "  API: ${GREEN}http://localhost:8080${NC}"
    echo -e "  Admin: ${GREEN}http://localhost:3000/admin${NC}"
    echo ""
    echo -e "${YELLOW}🔧 Management Commands:${NC}"
    echo -e "  View logs: ${GREEN}docker-compose logs -f${NC}"
    echo -e "  Stop services: ${GREEN}docker-compose down${NC}"
    echo -e "  Restart services: ${GREEN}docker-compose restart${NC}"
    echo ""
    echo -e "${YELLOW}📚 Next Steps:${NC}"
    echo -e "  1. Open web UI and configure your license"
    echo -e "  2. Add proxy configurations"
    echo -e "  3. Configure Iranian operator APIs"
    echo -e "  4. Start intelligence gathering"
    echo ""
    echo -e "${RED}🔒 IMPORTANT: Keep your .env file secure!${NC}"
}

# Main installation process
main() {
    echo -e "${BLUE}🚀 Iran Intelligence Enterprise Installer${NC}"
    echo ""
    
    install_dependencies
    setup_docker
    setup_project
    setup_environment
    start_services
    show_success
}

# Run main function
main "$@"
