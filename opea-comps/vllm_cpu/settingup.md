# Setting up vLLM docker
Inside of the "/third_parties/vllm/src" there is some bash file to help construct the image. Let's use them.
The version is locked at 0.6.1 of vLLM we are going to correct that.

## 🚀1. Set up Environment Variables

```bash
export LLM_ENDPOINT_PORT=8008
export host_ip=${host_ip}
export HF_TOKEN=${HF_TOKEN}
export LLM_ENDPOINT="http://${host_ip}:${LLM_ENDPOINT_PORT}"
export LLM_MODEL_ID="Intel/neural-chat-7b-v3-3"
```

Although in the readme it says that the token is only need for gated models it seem that its used in all calls so  is always needed.
The HF_TOKEN its the wrong variable being used across the scripts. Use HUGGINGFACEHUB_API_TOKEN instead

``` bash
export LLM_MODEL_ID="deepseek-ai/DeepSeek-R1-Distill-Qwen-1.5B"
export HUGGINGFACEHUB_API_TOKEN={your_token}
export host_ip=localhost
```

## 🚀2. Set up vLLM Service
### 2.3 vLLM with OpenVINO (on Intel GPU and CPU)

#### Update "build_docker_vllm_openvino.sh"
In the src/build_docker_vllm_openvino.sh its getting the version 0.6.1 tag. Remove it to ge the latest:

``` bash
# before
 cd ./temp/vllm_source/ && git checkout v0.6.1
# after
 cd ./temp/vllm_source/
```

#### Build Docker Image

To build the docker image for Intel CPU, run the command

```bash
bash ./src/build_docker_vllm_openvino.sh
```

Once it successfully builds, you will have the `vllm-openvino` local image. It can be used to spawn a serving container with OpenAI API endpoint or you can work with it interactively via bash shell.

#### Launch vLLM service

If previously defined the variables below just skip to the next step. If not create .env file:

```sh
export TAG=latest
export LLM_ENDPOINT_PORT=8008
export HF_CACHE_DIR=/home/digital/.cache/huggingface
export LLM_MODEL_ID=deepseek-ai/DeepSeek-R1-Distill-Qwen-1.5B
export host_ip=127.0.0.1
export HUGGINGFACEHUB_API_TOKEN= {HF_API_Token}
```

Now we ready to start up the docker with vLLM using CPU and OpenVINO
```sh
docker-compose -f docker-compose.cpu.yaml up
```

**Note:** The provided docker-compose file includes the `--disable-frontend-multiprocessing` flag because multiprocessing was causing startup issues.

### Testing
To test the running service, send a chat completion request using curl:

```bash
curl http://${host_ip}:8008/v1/chat/completions \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "model": "deepseek-ai/DeepSeek-R1-Distill-Qwen-1.5B",
        "messages": [
            {
                "role": "user",
                "content": "What is Deep Learning?"
            }
        ]
    }'
```

### References
- [vLLM AI Accelerator Documentation](https://docs.vllm.ai/en/stable/getting_started/installation/ai_accelerator/index.html#set-up-using-docker)
- [OPEA GenAI Components Repository](https://github.com/opea-project/GenAIComps/tree/main/comps/third_parties/vllm)