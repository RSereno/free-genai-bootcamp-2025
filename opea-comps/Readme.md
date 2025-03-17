## Running Ollama Third-Party Service

### Choosing a Model

Since we are using Docker image there is information in [Docker Hub](https://hub.docker.com/r/ollama/ollama), we will be using the Nvidia GPU with Docker. Therefore, we need to install the NVIDIA Container Toolkit as described in the instructions.

### NVIDIA Container Toolkit

Since I have WSL with Ubuntu, let's use Apt to install it.

#### Install with Apt

1. Configure the repository
    ```sh
    curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey \
        | sudo gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg
    curl -s -L https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list \
        | sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' \
        | sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list
    sudo apt-get update
    ```

2. Install the NVIDIA Container Toolkit packages
    ```sh
    sudo apt-get install -y nvidia-container-toolkit
    ```

#### Configure Docker to use Nvidia driver
    ```sh
    sudo nvidia-ctk runtime configure --runtime=docker
    sudo systemctl restart docker
    ```

### Start the container

    ```sh
    docker run -d --gpus=all -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
    ```

It's running. Let's put some limitations on it.

    ```sh
    docker run -d \
    --gpus all \
    -v ollama:/root/.ollama \
    -p 11434:11434 \
    --security-opt=no-new-privileges \
    --cap-drop=ALL \
    --cap-add=SYS_NICE \
    --memory=8g \
    --memory-swap=8g \
    --cpus=4 \
    --read-only \
    --name ollama \
    ollama/ollama
    ```

### Run it locally

    ```sh
    docker exec -it ollama ollama run deepseek-r1
    ```

Note: By not defining anything else, it will download the 7b model. Please reference [Ollama GitHub](https://github.com/ollama/ollama).

### Ollama API

Once the Ollama server is running, we can make API calls to the Ollama API.

[Ollama API Documentation](https://github.com/ollama/ollama/blob/main/docs/api.md)

## List running models

    ```sh
    curl http://localhost:11434/api/ps
    ```

By executing the command in the section "Run it locally", the model was downloaded.

## Generate a Request

    ```sh
    curl http://localhost:11434/api/generate -d '{
      "model": "deepseek-r1:latest",
      "prompt": "Why is the sky blue?"
    }'
    ```

# UV enviroment
https://docs.astral.sh/uv/#getting-started

curl -LsSf https://astral.sh/uv/install.sh | sh

uv venv vllm_source --python 3.12 --seed

source vllm_source/bin/activate


# Technical Uncertainty

Q: Since we are putting some limitations on the container, will the API calls work as expected?  
A: No problem so far.

Q: What is the model name that Ollama is expecting since we defined deepseek-r1 only?  
A: By listing the available models, we can see it's "deepseek-r1:latest".

Q: Will this work in later steps of the bootcamp?