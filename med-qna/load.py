# Load Med42 Llama3 model directly using transformer
from transformers import AutoTokenizer, AutoModelForCausalLM
# Tokenizer and Model
tokenizer = AutoTokenizer.from_pretrained("m42-health/Llama3-Med42-8B")
model = AutoModelForCausalLM.from_pretrained("m42-health/Llama3-Med42-8B")