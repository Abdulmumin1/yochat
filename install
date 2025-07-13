#!/bin/bash

# --- Configuration ---
APP_NAME="yochat"
INSTALL_DIR="$HOME/.local/bin"
TEMP_DIR=$(mktemp -d) # Create a temporary directory for downloads and extraction

# --- Styling for messages ---
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# --- Helper Functions ---
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# --- Cleanup function to run on exit ---
cleanup() {
    if [ -d "$TEMP_DIR" ]; then
        log_info "Cleaning up temporary files..."
        rm -rf "$TEMP_DIR"
        log_info "Temporary files removed."
    fi
}

# Register the cleanup function to be called on script exit (success or failure)
trap cleanup EXIT

# --- Pre-checks ---
check_prerequisites() {
    if ! command -v curl &> /dev/null; then
        log_error "curl is not installed. Please install curl to proceed (e.g., 'sudo apt install curl' on Debian/Ubuntu, 'sudo yum install curl' on Fedora/CentOS, or 'brew install curl' on macOS)."
    fi
    if ! command -v tar &> /dev/null; then
        log_error "tar is not installed. This should be available by default on Linux/macOS."
    fi
    log_info "All necessary tools found."
}

# --- Main Script Logic ---
main() {
    check_prerequisites

    # 1. GitHub repository details (hardcoded based on user's request)
    GITHUB_USER="Abdulmumin1"
    GITHUB_REPO="yochat"

    # 2. Determine OS and Architecture
    OS_TYPE=""
    ARCH_TYPE=""
    case "$(uname -s)" in
        Linux*)
            OS_TYPE="linux"
            ;;
        Darwin*)
            OS_TYPE="darwin"
            ;;
        *)
            log_error "Unsupported operating system: $(uname -s). This script supports Linux and macOS."
            ;;
    esac

    case "$(uname -m)" in
        x86_64)
            ARCH_TYPE="amd64"
            ;;
        arm64)
            ARCH_TYPE="arm64"
            ;;
        aarch64) # Common for ARM Linux
            ARCH_TYPE="arm64"
            ;;
        *)
            log_error "Unsupported CPU architecture: $(uname -m). This script supports amd64 and arm64."
            ;;
    esac

    log_info "Detected OS: $OS_TYPE, Architecture: $ARCH_TYPE"

    # 3. Fetch latest release tag
    log_info "Fetching latest release information from GitHub..."
    API_URL="https://api.github.com/repos/${GITHUB_USER}/${GITHUB_REPO}/releases/latest"
    RELEASE_INFO=$(curl -s "$API_URL")

    if [ -z "$RELEASE_INFO" ]; then
        log_error "Failed to fetch release information from $API_URL. Please check GitHub username/repository name and your internet connection."
    fi

    # Extract tag_name using grep and sed (robust parsing)
    LATEST_TAG=$(echo "$RELEASE_INFO" | grep '"tag_name":' | head -n 1 | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')

    if [ -z "$LATEST_TAG" ]; then
        log_error "Could not find the latest release tag. Ensure the repository has public releases."
    fi
    log_info "Latest release tag found: $LATEST_TAG"

    # 4. Construct download filename and URL
    DOWNLOAD_FILENAME="${APP_NAME}-${OS_TYPE}-${ARCH_TYPE}.tar.gz"
    DOWNLOAD_URL="https://github.com/${GITHUB_USER}/${GITHUB_REPO}/releases/download/${LATEST_TAG}/${DOWNLOAD_FILENAME}"

    log_info "Attempting to download: $DOWNLOAD_URL"

    # 5. Download the application
    cd "$TEMP_DIR" || log_error "Failed to change to temporary directory."
    curl -L -o "$DOWNLOAD_FILENAME" "$DOWNLOAD_URL"
    if [ $? -ne 0 ]; then
        log_error "Failed to download $APP_NAME from $DOWNLOAD_URL. Please verify the URL and that the asset exists for your OS/architecture."
    fi
    log_info "Downloaded $DOWNLOAD_FILENAME to $TEMP_DIR"

    # 6. Extract the application
    log_info "Extracting $APP_NAME..."
    tar -xzf "$DOWNLOAD_FILENAME"
    if [ $? -ne 0 ]; then
        log_error "Failed to extract $DOWNLOAD_FILENAME. The archive might be corrupted or in an unexpected format."
    fi

    # Find the executable (it might be in a subdirectory after extraction)
    # Search for an executable file named $APP_NAME within the extracted contents
    EXECUTABLE_PATH=$(find . -type f -name "$APP_NAME")

    if [ -z "$EXECUTABLE_PATH" ]; then
        log_error "Could not find the '$APP_NAME' executable after extraction. Please check the archive structure or manually locate it."
    fi
    log_info "Found executable at: $EXECUTABLE_PATH"

    # 7. Create ~/.local/bin if it doesn't exist
    log_info "Ensuring installation directory exists: $INSTALL_DIR"
    mkdir -p "$INSTALL_DIR"
    if [ $? -ne 0 ]; then
        log_error "Failed to create directory $INSTALL_DIR. Check user permissions for your home directory."
    fi

    # 8. Move executable to ~/.local/bin and grant permissions
    log_info "Moving $APP_NAME to $INSTALL_DIR and setting execute permissions..."
    mv "$EXECUTABLE_PATH" "$INSTALL_DIR/$APP_NAME"
    if [ $? -ne 0 ]; then
        log_error "Failed to move $APP_NAME to $INSTALL_DIR. Check permissions for $INSTALL_DIR."
    fi
    chmod +x "$INSTALL_DIR/$APP_NAME"
    log_info "$APP_NAME successfully moved and made executable."

    # 9. Update PATH environment variable
    log_info "Checking and updating PATH environment variable..."
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        log_warn "$INSTALL_DIR is not currently in your PATH."
        SHELL_CONFIG_FILE=""
        CURRENT_SHELL=$(basename "$SHELL")

        case "$CURRENT_SHELL" in
            bash)
                if [ -f "$HOME/.bashrc" ]; then
                    SHELL_CONFIG_FILE="$HOME/.bashrc"
                elif [ -f "$HOME/.bash_profile" ]; then
                    SHELL_CONFIG_FILE="$HOME/.bash_profile"
                fi
                ;;
            zsh)
                if [ -f "$HOME/.zshrc" ]; then
                    SHELL_CONFIG_FILE="$HOME/.zshrc"
                fi
                ;;
            *)
                log_warn "Unknown shell ($CURRENT_SHELL). Please manually add 'export PATH=\"\$HOME/.local/bin:\$PATH\"' to your shell's configuration file."
                ;;
        esac

        if [ -n "$SHELL_CONFIG_FILE" ]; then
            log_info "Adding '$INSTALL_DIR' to PATH in $SHELL_CONFIG_FILE"
            echo -e "\n# Add $APP_NAME to PATH for user-level executables" >> "$SHELL_CONFIG_FILE"
            echo "export PATH=\"\$HOME/.local/bin:\$PATH\"" >> "$SHELL_CONFIG_FILE"
            log_warn "PATH updated. You MAY need to run 'source $SHELL_CONFIG_FILE' or open a NEW terminal for changes to take effect."
        else
            log_warn "Could not automatically detect shell configuration file. Please manually add '$INSTALL_DIR' to your PATH."
        fi
    else
        log_info "$INSTALL_DIR is already in your PATH. No changes needed."
    fi

    log_info "Installation complete!"
    log_info "To verify, please open a NEW terminal window and type: ${GREEN}$APP_NAME${NC}"
    log_info "If you encounter any issues, you can follow the manual Installation guide"
    log_info "you can also add alias for  ${YELLOW}yochat${NC} as ${YELLOW}chat${NC}"
}

# --- Run the main function ---
main
