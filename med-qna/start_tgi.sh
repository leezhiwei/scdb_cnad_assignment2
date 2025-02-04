# Run docker container to expose Med42 Llama3 as API with Hugging Face's Text Generation Inference (TGI)
docker run \
    # Enable GPU acceleration
    --gpus all \
    # Shared memory size increased
    --shm-size 1g \
    # Map to host on port 80 inside the container
    -p 8080:80 \
    # Use latest TGI image
    ghcr.io/huggingface/text-generation-inference:latest \
    # Med42 Llama3 model path
    --model-id "med42/Llama3-med42"