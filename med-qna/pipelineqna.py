# Import libraries
import transformers
import torch

# Med42 Llama3 Model
model_name_or_path = "m42-health/Llama3-Med42-8B"

# Setup the pipeline for text generating 
pipeline = transformers.pipeline(
    # Set text generation task
    "text-generation",
    # Load Med42 Llama3 model
    model=model_name_or_path,
    # bfloat 16 bit to represent dynamic range of larger and smaller numbers
    # bfloat16 used to optimise model performance
    torch_dtype=torch.bfloat16,
    # Automatically map the model to available CPU or GPU
    device_map="auto",
)

# Message content for Med42 Llama3 model
messages = [
    # Set behaviour of the model
    {
        "role": "system",
        "content": (
            "Always answer as helpfully as possible, while being safe. "
            "Your answers should not include any harmful, unethical, racist, sexist, toxic, dangerous, or illegal content. "
            "Please ensure that your responses are socially unbiased and positive in nature. If a question does not make any sense, or is not factually coherent, explain why instead of answering something not correct. "
            "If you don’t know the answer to a question, please don’t share false information."
        ),
    },
    # Message for the model
    {"role": "user", "content": "If i fall, what should i do?"},
]

# Prepare the prompt for model input using pipeline
prompt = pipeline.tokenizer.apply_chat_template(
    # Message conversation without tokenizing or generation prompt
    messages, tokenize=False, add_generation_prompt=False
)

# Stop tokens to end the model output
stop_tokens = [
    # End-of-sequence token for stopping generation
    pipeline.tokenizer.eos_token_id,
    # Token for end of text
    pipeline.tokenizer.convert_tokens_to_ids("<|eot_id|>"),
]

# Set pipeline with the specified parameters
outputs = pipeline(
    # Input prompt for text generation
    prompt,
    # Limit the generated response to 512 tokens
    max_new_tokens=512,
    # Stop the generation if stop tokens are encountered
    eos_token_id=stop_tokens,
    # Enable sampling for more diverse result
    do_sample=True,
    # Control randomness of the output. Lower value makes allow deterministic response
    temperature=0.4,
    # Limit sampling to top 150 tokens
    top_k=150,
    # Nucleus sampling for coherent text from model
    top_p=0.75,
)

# Display the generated response excluding the input prompt
print(outputs[0]["generated_text"][len(prompt) :])