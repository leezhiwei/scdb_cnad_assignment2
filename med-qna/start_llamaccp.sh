# Github repo for llama.cpp to expose Med42 Llama3 as API but might be intensive for CPU
git clone https://github.com/ggerganov/llama.cpp
# Change directory to the cloned repo
cd llama.cpp
# Compile the source to build llama.cpp
make
# Start the Llama.cpp server to expose the model as an API with context length of 2048 tokens and port 8000
llama-server -m models/Llama3-Med42-8B-Q5_K_M.gguf -c 2048 --host 0.0.0.0 --port 8000 -ngl 33

# Testing with curl
curl.exe -X POST http://192.168.2.108:8000/v1/chat/completions -H "Content-Type: applications/json" --data "{\"model\":\"m42-health/Llama3-Med42-8B\",\"messages\":[{\"role\":\"user\",\"content\":\"If I fall, what should I do?\"}]}"