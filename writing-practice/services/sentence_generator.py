from groq import Groq
import os
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

class SentenceGenerator:
    def __init__(self):
        self.client = Groq(api_key=os.getenv('GROQ_API_KEY'))

    def generate_sentence(self, word: str) -> str:
        print("\n=== SentenceGenerator Service ===")
        print(f"Generating sentence for word: {word}")
        
        prompt = f"""Generate a simple sentence using the following word: {word}
The grammar should be scoped to A1 level.
You can use the following vocabulary to construct a simple sentence:
- Simple objects (e.g., book, car, ramen, sushi)
- Simple verbs (e.g., to drink, to eat, to meet)
- Simple time expressions (e.g., tomorrow, today, yesterday)

 Please provide the response in this format:
        Portuguese: [sentence in Portuguese]
        English: [English translation]"""

        try:
            print("Calling Groq API...")
            response = self.client.chat.completions.create(
                model="qwen-2.5-32b",
                messages=[
                    {"role": "system", "content": "You are a basic sentence generator. Please generate a simple sentence."},
                    {"role": "user", "content": prompt.format(word)}
                ],
                max_tokens=200,
                temperature=0.3
            )
            
            print("\n=== API Response ===")
            content = response.choices[0].message.content.strip()
            print(f"Raw response: {content}")
            
            # Extract English sentence
            sentences = {}
            for line in content.split('\n'):
                line = line.strip()
                if line.startswith('Portuguese:'):
                    sentences['portuguese'] = line.replace('Portuguese:', '').strip()
                elif line.startswith('English:'):
                    sentences['english'] = line.replace('English:', '').strip()
            
            # Use English sentence or fallback
            sentence = sentences.get('english', f"I have a {word}. fallback")
            print(f"Extracted English sentence: {sentence}")
            
            print("=== End SentenceGenerator ===\n")
            return sentence
            
        except Exception as e:
            print(f"\nError in SentenceGenerator: {str(e)}\n")
            return f"I have a {word}. error"
