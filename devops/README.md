# Configuration

Configuring the local environment and tools for python, gcp, and terraform setup.


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
      pyenv --version
      pyenv install 3.9.2
      pyenv local 3.9.2
    ```


## Install and Configure Google Cloud SDK

  - Step 1a: Install Google Cloud SDK via Homebrew
    ```
      brew install --cask google-cloud-sdk
    ```
  - Or Step 1b: Re-install Google Cloud SDK via Homebrew
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
    - Note: Get the project ID from the Google Cloud Console Home -> Dashboard
            The project ID is the text ID of the project.
    ```
      gcloud config set project PROJECT-ID
    ```
  - Step 6: List the GCP regions and zones
    - Note: this tests that the gcloud is setup correctly
    ```
      gcloud compute zones list
      gcloud compute regions list
    ```
  - Step 7: Create a local settings directory for GCP
    ```
      mkdir ~/.gcp
    ```
  - Step 8: Save the creedentials file and project environment file to the local settings directory
    - Login to the GCP cloud console and create a credentials file and download it as a JSON file
    - Place the downloaded JSON credentials file to the local settings directory for GCP (~/.gcp)
    - Create an environment file (name it as PROJECTID-env.sh) in the local settings directory for GCP (~/.gcp)
        and place the following in there (select a GCP region and zone from the command output above).
      Note: substitute the CREDENTIALS_FILE, PROJECTID, GCP_REGION, and GCP_ZONE with actual values below.
    ```
      export TF_VAR_gcp_credentials_file=${HOME}/.gcp/CREDENTIALS_FILE.json
      export TF_VAR_gcp_project=PROJECTID
      export TF_VAR_gcp_region=GCP_REGION
      export TF_VAR_gcp_zone=GCP_ZONE
    ```


## GCP gcloud commands

  - Get list of available GKE cluster versions by region
    ```
      gcloud container get-server-config --region us-west1
    ```

  - Get list of available GKE cluster versions by zone
    ```
      gcloud container get-server-config --zone us-west1-a
    ```

  - Get list of running GKE clusters
    ```
      gcloud container clusters list 
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
      terraform --version
    ```
    Note: make sure that the Terraform version is v0.14.9

  - Install Terraform Autocomplete for bash or zsh
    ```
      terraform -install-autocomplete
    ```

  - Check installation by running through the tutorial with Docker
    - Note: This tutorial requires Docker to be installed
    - See [Terraform Tutorial: https://learn.hashicorp.com/tutorials/terraform/install-cli](https://learn.hashicorp.com/tutorials/terraform/install-cli)


## GCP Terraform Modules

  - See [https://github.com/terraform-google-modules](https://github.com/terraform-google-modules)
  - For GKE and with several examples:
    - See [https://github.com/terraform-google-modules/terraform-google-kubernetes-engine](https://github.com/terraform-google-modules/terraform-google-kubernetes-engine)
