# Github repo for llama.cpp to expose Med42 Llama3 as API but might be intensive for CPU
git clone https://github.com/ggerganov/llama.cpp
# Change directory to the cloned repo
cd llama.cpp
# Compile the source to build llama.cpp
make
# Start a new screen session
screen -s llama_med42_server_session
# Start the Llama.cpp server to expose the model as an API with context length of 2048 tokens and port 8000
llama-server -m models/Llama3-Med42-8B-Q5_K_M.gguf -c 2048 --host 0.0.0.0 --port 8000 -ngl 33
# Testing with curl
curl.exe -X POST http://192.168.2.108:8000/v1/chat/completions -H "Content-Type: applications/json" --data "{\"model\":\"m42-health/Llama3-Med42-8B\",\"messages\":[{\"role\":\"user\",\"content\":\"If I fall, what should I do?\"}]}"
# Curl with system role to set guidelines on responses
curl.exe -X POST http://192.168.2.108:8000/v1/chat/completions -H "Content-Type: application/json" --data "{\"model\":\"m42-health/Llama3-Med42-8B\",\"messages\":[{\"role\":\"system\",\"content\":\"Always answer as helpfully as possible, while being safe. Your answers should not include any harmful, unethical, racist, sexist, toxic, dangerous, or illegal content. Please ensure that your responses are socially unbiased and positive in nature. If a question does not make any sense, or is not factually coherent, explain why instead of answering something not correct. If you don’t know the answer to a question, please don’t share false information.\"},{\"role\":\"user\",\"content\":\"If I fall, what should I do?\"}]}"
