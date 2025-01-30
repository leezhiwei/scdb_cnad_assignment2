# Deploy the docker
docker run --runtime nvidia --gpus all \
	--name my_vllm_container \
	-v ~/.cache/huggingface:/root/.cache/huggingface \
 	--env "HUGGING_FACE_HUB_TOKEN=<secret>" \
	-p 8000:8000 \
	--ipc=host \
	vllm/vllm-openai:latest \
	--model m42-health/Llama3-Med42-8B

# Load and run the Med42 Llama3 model:
echo "Starting the Med42 model in the container..."
docker exec -it my_vllm_container bash -c "vllm serve m42-health/Llama3-Med42-8B"

# Call the server using curl:
curl -X POST "http://localhost:8000/v1/chat/completions" \
	-H "Content-Type: application/json" \
	--data '{
		"model": "m42-health/Llama3-Med42-8B",
		"messages": [
			{
				"role": "user",
				"content": "If I fall, what should I do?"
			}
		]
	}'