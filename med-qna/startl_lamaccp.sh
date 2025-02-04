# Github repo for llama.cpp to expose Med42 Llama3 as API but might be intensive for CPU
git clone https://github.com/ggerganov/llama.cpp
# Change directory to the cloned repo
cd llama.cpp
# Compile the source to build llama.cpp
make
# Start the Llama.cpp server to expose the model as an API with context length of 2048 tokens and port 8000
./server -m med42/Llama3-med42.GGUF -c 2048 --host 0.0.0.0 --port 8000

# Testing with curl
curl --request POST \
    --url http://localhost:8080/completion \
    --header "Content-Type: application/json" \
    --data '{"prompt": "If I fall, what should I do?","n_predict": 128}'