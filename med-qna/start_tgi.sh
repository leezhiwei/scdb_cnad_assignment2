# Run docker container to expose Med42 Llama3-Med42-8B as API with Hugging Face's Text Generation Inference (TGI)
docker run \
    # Enable GPU acceleration
    --gpus all \
    # Shared memory size increased
    --shm-size 12g \
    # Map to host on port 80 inside the container
    -p 8080:80 \
    # Disable the Hugging Face Hub transfer feature via environment variable = 0
    -e HF_HUB_ENABLE_HF_TRANSFER=0 \
    # Mount local volume to store model data
    -v $PWD/modeldata:/data \
    # Use latest TGI image
    ghcr.io/huggingface/text-generation-inference:latest \
    # Med42 Llama3 model path
    --model-id "m42-health/Llama3-Med42-8B" \
    # Limit max batch size to 5 process in single batch
    --max-batch-size 5 \
    # Limit max tokens (chunks of text) to 512
    --max-total-tokens 512 


# Run docker container to expose Llama3-Med42-8B-4bit  as API with Hugging Face's Text Generation Inference (TGI)
docker run \
	# Enable GPU acceleration
	--gpus all \
	# Shared memory size increased
	--shm-size 12g \
	# Map to host on port 80 inside the container
	-p 8080:80 \
	# Disable the Hugging Face Hub transfer feature via environment variable = 0
	-e HF_HUB_ENABLE_HF_TRANSFER=0 \
	# Mount local volume
	-v $PWD/modeldata:/data \
	# Use latest TGI image
	ghcr.io/huggingface/text-generation-inference:latest \
	# Med42 Llama3 model path
	--model-id "emircanerol/Llama3-Med42-8B-4bit" \
        # Limit max input length
        --max-input-length 2048 \
        # Limit max tokens (chunks of text)
        --max-total-tokens 4096 