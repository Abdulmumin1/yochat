# YoChat

YoChat is a super duper fast command-line chat tool that leverages the Gemini 2.5 Flash Lite model to provide direct and concise answers to your questions. It supports both text-based queries and multimodal input, allowing you to ask questions about files (documents, images, audio, and video).


## Usage

> fun tip: add alias for 'yochat' as "chat" 😅

so you be like: chat "how do i compress a video file"

### Setting Your API Key

Before you can use YoChat, you need to set your Gemini API key. You can obtain a free API key from [Google AI Studio](https://aistudio.google.com).

To set your API key, use the `set` command:

```bash
yochat set <YOUR_API_KEY>
````

For example:

```bash
yochat set AbCdEfGhIjKlMnOpQrStUvWxYz1234567890
```

This will save your API key to a configuration file, typically located in your user's configuration directory. The exact path will be displayed after successful setup.

### Asking Questions

You can ask YoChat questions in several ways:

#### 1. Questions

Simply type your question after the `yochat` command:

```bash
yochat "how to convert file to mp4"
```

```bash
yochat "how to revert to a previous commit"
```


#### 3\. Analyzing Files (Multimodal Input)

YoChat can analyze various file types, including PDFs, audio, images, and videos. When you provide a file, you can also ask a specific question about its content.

Use the `--file` flag followed by the path to your file.


To ask a specific question about the file, combine `--file` with the `-q` flag:

```bash
yochat --file ./image.jpg -q "give me an ocr of this file"
```

# **Installation**

This guide explains how to download, install, and run the yochat executable on various operating systems, including making it accessible globally from any terminal.
Our application is built using Go, which allows us to provide a single executable file for each platform, making installation straightforward.

## **1\. Download the Application**

The compiled application binaries are available on our GitHub Releases page.

1. Go to the [Releases page of this repository](https://www.google.com/search?q=https://github.com/abdulmumin1/yochat/releases)
2. Find the latest release (e.g., v1.0).
3. Under the "Assets" section, download the archive file that matches your operating system and CPU architecture.
   * **For Windows (64-bit Intel/x86):** Look for yochat-windows-x86.zip
   * **For macOS (Intel-based):** Look for yochat-darwin-x86.zip
   * **For macOS (Apple Silicon \- M1/M2/M3 chips):** Look for yochat-darwin-arm64.zip
   * **For Linux (64-bit Intel/AMD):** Look for yochat-linux-x86.zip
   * **For Linux (64-bit ARM):** Look for yochat-linux-arm64.zip

## **2\. Extract the Application**

Once you've downloaded the correct archive, you need to extract its contents.

### **For Windows (.zip file)**

1. Locate the downloaded .zip file (e.g., yochat-windows-amd64.zip).
2. Extract the zip file.
4. Choose a destination folder where you want to keep the application. A good practice is to create a dedicated folder like `C:\Program Files\yochat` or `C:\yochat`

Making an application globally accessible on Windows typically involves adding its directory to the system's Path environment variable.

1. **Move the extracted yochat.exe** to a permanent location, for example, C:\\yochat.
2. **Add the directory to your System PATH:**

   * Search for "Environment Variables" in the Windows search bar and select "Edit the system environment variables".
   * Create a "New" env variable and add the full path to the directory where you placed yochat.exe (e.g., C:\\yochat).

3. **Verify the installation** by opening a **new Command Prompt or PowerShell window** and typing:


```bash
> yochat
```

You should see the application's output. (Note: Existing terminal windows might not reflect the PATH changes until restarted).

### **For macOS and Linux**

#### Auto insall with:

```bash
curl -fsSL https://yochat.yaqeen.me/install | bash
```

#### Manually

1. Open your **Terminal** application.
2. Navigate to the directory where you downloaded the .zip file using the cd command.

```bash
cd ~/Downloads
  ```

3. Extract the archive using the tar command. Replace your-downloaded-file.zip with the actual filename:

```bash
unzip your-downloaded-file.zip
```

This will extract the yochat executable file into the current directory.

To run yochat by simply typing yochat in any terminal, you need to place its executable in a directory that is part of your system's PATH environment variable.

### **For macOS and Linux**

The easiest way to make yochat globally accessible is to move it to /usr/local/bin, which is typically included in your system's PATH.

1. Open your **Terminal** application.
2. Navigate to the directory where you extracted the yochat executable (e.g., `~/Downloads`).
3. **Grant execute permissions** to the file (if you haven't already):

```bash
chmod \+x yochat
```

4. **Move the executable** to `/home/{YOUR_USER}/.local/bin/` using `sudo` (administrator privileges are required for this location):

```bash
mv yochat /home/{YOUR_USER}/.local/bin
```

5. **Verify the installation** by opening a **new terminal window** and typing:

```bash
yochat
```

You should see the application's output.

If you encounter any issues, please refer to the project's GitHub repository for more information or to open an issue.
