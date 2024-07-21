# Installnation

Autobase is a Go application that can be installed on any supported platform in the release section. The installation process is simple and straightforward. The following steps will guide you through the installation process.

## FADE

The most recommended way to install Autobase is to install through FADE. FADE stands for Fastre Autobase Development Environment. FADE ensures that you have the latest version of Autobase installed on your system, along with all the necessary dependencies.

>[!NOTE]
>FADE is tested and intended to work on Windows operating systems. However, it should work on any operating system that supports .NET Runtime 8.0 or higher.

### Install FADE

- Download the latest version of FADE from [FADE Releases](https://github.com/mvishok/fade/releases)
- Install FADE by running the downloaded setup file
- Open a fresh elevated command prompt (Run as Administrator) and run the following command to install Autobase

```bash
fade install autobase
```

## Github Releases

Autobase binaries are available in the [Releases](https://github.com/mvishok/autobase/releases) section. You can download the latest version of Autobase for your operating system from the releases page.

### Install Autobase

- Download the latest version of Autobase from the above link based on your operating system
- Place the downloaded binary in a new folder
- Add the folder to the system PATH

#### Add to PATH

##### On Windows

- Search for `Environment Variables` in the start menu
- Click on `Environment Variables`
- Under `System Variables`, find and select the `Path` variable
- Click `Edit`
- Click `New` and add the path to the folder where the Autobase binary is placed
- Click `OK` to save the changes
- Autobase is now installed on your system

##### On Linux

- Open the terminal
- Run the following command to open the `.bashrc` file
```bash
nano ~/.bashrc
```
- Add the following line to the end of the file
```bash
export PATH=$PATH:/path/to/autobase/folder
```
- Press `Ctrl + X` to exit the editor
- Press `Y` to save the changes
- Press `Enter` to confirm the file name
- Run the following command to apply the changes
```bash
source ~/.bashrc
```
- Autobase is now installed on your system


## Running From Source

You can also run Autobase from the source code. This method is recommended for developers who want to contribute to the project or test the latest features.

### Prerequisites

- Go 1.22 or higher
- Git

### Clone the Repository

- Clone the Autobase repository to your local machine
```bash
git clone https://github.com/mvishok/autobase.git
```
- Change the directory to the cloned repository
```bash
cd autobase
```
- Run the following command to build the project
```bash
go build
```
- Run the following command to start Autobase
```bash
./autobase
```
- Autobase is now running on your system

# Usage

Autobase takes a configuration file as input and starts the server based on the configuration. The configuration file is a JSON file that contains the details of the data sources and the API endpoints.

```bash
autobase config config.json
```
