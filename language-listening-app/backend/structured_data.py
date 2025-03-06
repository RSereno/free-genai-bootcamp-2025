import os
import re
from typing import List, Dict
from groq import Groq
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Model ID
MODEL_ID = "deepseek-r1-distill-llama-70b"

class TranscriptProcessor:
    def __init__(self, groq_api_key: str = None):
        self.groq_api_key = groq_api_key
        if groq_api_key:
            print("Groq:")
            self.groq_client = Groq(api_key=groq_api_key)
        else:
            print("No Groq:")
            self.groq_client = None

    def clean_text(self, text: str) -> str:
        """
        Clean the text by removing special characters, extra spaces, etc.
        """
        text = re.sub(r'\n', ' ', text)  # Replace newlines with spaces
        text = re.sub(r'\[.*?\]', '', text)  # Remove timestamps or speaker labels
        text = re.sub(r'\s+', ' ', text).strip()  # Remove extra spaces
        return text

    def extract_dialogue(self, transcript: List[Dict]) -> List[Dict]:
        """
        Extract dialogue from the transcript, assuming each entry is a turn.
        """
        dialogue = []
        for entry in transcript:
            text = self.clean_text(entry['text'])
            dialogue.append({'text': text})
        return dialogue

    def structure_data(self, transcript: List[Dict]) -> Dict:
        """
        Structure the transcript data into sections based on content using Groq.
        """
        cleaned_transcript = [
            {'text': self.clean_text(entry['text'])} for entry in transcript
        ]

        if self.groq_client:
            # Convert cleaned_transcript to a string representation
            transcript_str = "\n".join(entry["text"] for entry in cleaned_transcript)
            
            # Load JSON template
            template_path = os.path.join(os.path.dirname(__file__), 'templates', 'structuredata.json')
            with open(template_path, 'r') as f:
                json_template = f.read()
            
            prompt = f"""You are an expert in structuring transcripts of spoken European Portuguese. Your task is to analyze the given transcript and extract key sections to generate structured question-based learning content.  
            Identify and extract the following sections:  

1. **Introduction**:  
   - Extract the initial part of the transcript where the speaker(s) introduce themselves, the topic, or the purpose of the recording.  

2. **Conversation**:  
   - Extract any conversational setup before questions are asked.  

3. **Questions & Answers**:  
   - Identify questions posed within the transcript.  
   - Extract corresponding multiple-choice answer options.  
   - Ensure each question is properly paired with its answer choices.  

4. **Conclusion**:  
   - Extract the final part of the transcript where the speaker(s) summarize the content or provide closing remarks.  

**Format the extracted content as JSON:**  

{json_template}

#### **Rules:**  
- **Do not translate any Portuguese text.**  
- **Analyze the entire transcript carefully to extract each section.**  
- **If a section type is not clearly identifiable, return it as an empty string (`""`) or an empty list (`[]`).**  
- **Maintain the original structure of the transcript while organizing it into JSON format.**  

Here is the transcript to analyze:  
{transcript_str}
"""
            try:
                chat_completion = self.groq_client.chat.completions.create(
                    messages=[
                        {
                            "role": "user",
                            "content": prompt
                        }
                    ],
                    model=MODEL_ID,
                    temperature=0.0,
                    max_tokens=4096,
                )
                response_text = chat_completion.choices[0].message.content
                # Parse the JSON response from Groq
                import json
                try:
                    structured_data = json.loads(response_text)
                    return structured_data
                except json.JSONDecodeError:
                    print(f"Error decoding JSON from Groq: {response_text}")
                    # Fallback to basic logic if JSON decoding fails
                    return {
                        'introduction': [],
                        'dialogue': [],
                        'topic_discussion': [],
                        'conclusion': []
                    }
            except Exception as e:
                print(f"Error during Groq API call: {e}")
                # Fallback to basic logic if Groq API call fails
                return {
                    'introduction': [],
                    'dialogue': [],
                    'topic_discussion': [],
                    'conclusion': []
                }
        else:
            print("Groq API key not set. Skipping Groq processing.")
            return {
                'introduction': [],
                'dialogue': [],
                'topic_discussion': [],
                'conclusion': []
            }

    def process_transcript(self, transcript_path: str) -> Dict:
        """
        Process the entire transcript: clean, extract dialogue, and structure.
        """
        # Read the transcript from the file
        with open(transcript_path, 'r') as f:
            transcript_text = f.readlines()

        # Parse the transcript into a list of dictionaries
        transcript = [{'text': line.strip()} for line in transcript_text]

        dialogue = self.extract_dialogue(transcript)
        structured_data = self.structure_data(transcript)
        
        return {
            'dialogue': dialogue,
            'structured_data': structured_data
        }

if __name__ == '__main__':
    groq_api_key = os.getenv("GROQ_API_KEY")  
    processor = TranscriptProcessor(groq_api_key=groq_api_key)
    processed_data = processor.process_transcript("transcripts/sX6xBrSb-TU.txt")
    
    print("Processed Data:")
    print(processed_data)
