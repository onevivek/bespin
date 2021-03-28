# datadev

Data Developers Examples and Tutorials


## MacOS Homebrew Setup

  - Update homebrew first and then upgrade formula for pyenv
    ```
      brew update
    ```
  - Check if brew prefix shows '/usr/local'
    ```
      brew --prefix
    ```


## Python Version Upgrade (need atleast 3.9.x for GCloud SDK)

  - Upgrade formula for pyenv
    ```
      brew upgrade pyenv
    ```
  - Install python version and check installation
    ```
      which pyenv
      pyenv list
      pyenv --version
      pyenv install 3.9.2
      pyenv --version
      pyenv local 3.9.2
    ```


## Install and Configure Google Cloud SDK

  - Step 1a: Install Google Cloud SDK via Homebrew
    ```
      brew install --cask google-cloud-sdk
    ```
  - Or Step 1a: Re-install Google Cloud SDK via Homebrew
    ```
      brew reinstall google-cloud-sdk
    ```
  - Step 2: Check if gcloud is installed
    ```
      which gcloud
    ```
  - Step 3: Add path variables to bash profile
    ```
      source "$(brew --prefix)/Caskroom/google-cloud-sdk/latest/google-cloud-sdk/path.bash.inc"
      source "$(brew --prefix)/Caskroom/google-cloud-sdk/latest/google-cloud-sdk/completion.bash.inc"
    ```
  - Step 4: Authenticate with Google Credentials (this will launch a browser and prompt login)
    ```
      gcloud auth login
    ```
  - Step 5: Run gcloud command to set the project
    - Note: Get the project ID from the Google Cloud Console (select the project you want)
            The project ID is the text ID of the project.
    ```
      gcloud config set project PROJECT-ID
    ```
  - Step 6: List the GCP regions and zones (this tests that the gcloud is setup correctly)
    ```
      gcloud compute zones list
      gcloud compute regions list
    ```


## Terraform Installation and Configuraiton

  - Install Terraform using Homebrew
    ```
      brew tap hashicorp/tap
      brew install hashicorp/tap/terraform
      brew upgrade hashicorp/tap/terraform
    ```

  - Check Terraform Installation
    ```
      terraform --help
    ```

  - Install Terraform Autocomplete for bash or zsh
    ```
      terraform -install-autocomplete
    ```


## Terraform Installation and Checking Installation

  - Install Terraform using Homebrew
    ```
      brew tap hashicorp/tap
      brew install hashicorp/tap/terraform
      brew upgrade hashicorp/tap/terraform
    ```
  - Check installation by running through the tutorial with Docker
    - Note: This tutorial requires Docker to be installed
    - See [Terraform Tutorial: https://learn.hashicorp.com/tutorials/terraform/install-cli](https://learn.hashicorp.com/tutorials/terraform/install-cli)
