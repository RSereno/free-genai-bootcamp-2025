import groq
import os
from dotenv import load_dotenv
from typing import Optional, Dict, Any

# Load environment variables
load_dotenv()

# Model ID
MODEL_ID = "deepseek-r1-distill-llama-70b"  # Groq's recommended model

class GroqChat:
    def __init__(self, model_id: str = MODEL_ID):
        """Initialize Groq chat client"""
        self.model_id = model_id
        self.client = groq.Groq(api_key=os.getenv("GROQ_API_KEY"))

    def generate_response(self, message: str, inference_config: Optional[Dict[str, Any]] = None) -> Optional[str]:
        """Generate a response using Groq's API"""
        if inference_config is None:
            inference_config = {"temperature": 0.7}

        try:
            completion = self.client.chat.completions.create(
                model=self.model_id,
                messages=[
                    {"role": "system", "content": "You are a helpful assistant."},
                    {"role": "user", "content": message}
                ],
                temperature=inference_config["temperature"]
            )
            return completion.choices[0].message.content.strip()
            
        except Exception as e:
            #st.error(f"Error generating response: {str(e)}")
            print(f"Error generating response: {str(e)}")
            return None

if __name__ == "__main__":
    chat = GroqChat()
    while True:
        user_input = input("You: ")
        if user_input.lower() == '/exit':
            break
        response = chat.generate_response(user_input)
        print("Bot:", response)