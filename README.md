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

