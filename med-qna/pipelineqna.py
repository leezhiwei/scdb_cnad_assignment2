# Use a pipeline as a high-level helper
from transformers import pipeline

# Question for model
messages = [
    {"role": "user", "content": "If i fall, what should i do?"},
]

# Pipeline to Med42 Llama3 model
pipe = pipeline("text-generation", model="m42-health/Llama3-Med42-8B")
pipe(messages)