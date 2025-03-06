import os
import re
from typing import List, Dict
from groq import Groq

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
            # Use Groq to determine the structure
            prompt = f"""
            You are an expert in structuring transcripts of Japanese Language Proficiency Test (JLPT) listening tests. Your task is to analyze the given transcript and extract specific sections based on their question format.

            Here are the section extraction rules:

            For 問題1:
            Extract questions where the answer can be determined solely from the conversation without needing visual aids.

            ONLY include questions that meet these criteria:
            - The answer can be determined purely from the spoken dialogue.
            - No spatial/visual information is needed (like locations, layouts, or physical appearances).
            - No physical objects or visual choices need to be compared.

            For example, INCLUDE questions about:
            - Times and dates
            - Numbers and quantities
            - Spoken choices or decisions
            - Clear verbal directions

            DO NOT include questions about:
            - Physical locations that need a map or diagram
            - Visual choices between objects
            - Spatial arrangements or layouts
            - Physical appearances of people or things

            Format each question exactly like this:

            <question>
            Introduction:
            [the situation setup in japanese]

            Conversation:
            [the dialogue in japanese]

            Question:
            [the question being asked in japanese]

            Options:
            1. [first option in japanese]
            2. [second option in japanese]
            3. [third option in japanese]
            4. [fourth option in japanese]
            </question>

            Rules:
            - Only extract questions from the 問題1 section
            - Only include questions where answers can be determined from dialogue alone
            - Ignore any practice examples (marked with 例)
            - Do not translate any Japanese text
            - Do not include any section descriptions or other text
            - Output questions one after another with no extra text between them

            For 問題2:
            Extract questions where the answer can be determined solely from the conversation without needing visual aids.

            ONLY include questions that meet these criteria:
            - The answer can be determined purely from the spoken dialogue
            - No spatial/visual information is needed (like locations, layouts, or physical appearances)
            - No physical objects or visual choices need to be compared

            For example, INCLUDE questions about:
            - Times and dates
            - Numbers and quantities
            - Spoken choices or decisions
            - Clear verbal directions

            DO NOT include questions about:
            - Physical locations that need a map or diagram
            - Visual choices between objects
            - Spatial arrangements or layouts
            - Physical appearances of people or things

            Format each question exactly like this:

            <question>
            Introduction:
            [the situation setup in japanese]

            Conversation:
            [the dialogue in japanese]

            Question:
            [the question being asked in japanese]
            </question>

            Rules:
            - Only extract questions from the 問題2 section
            - Only include questions where answers can be determined from dialogue alone
            - Ignore any practice examples (marked with 例)
            - Do not translate any Japanese text
            - Do not include any section descriptions or other text
            - Output questions one after another with no extra text between them

            For 問題3:
            Extract all questions.
            Format each question exactly like this:

            <question>
            Situation:
            [the situation in japanese where a phrase is needed]

            Question:
            何と言いますか
            </question>

            Rules:
            - Only extract questions from the 問題3 section
            - Ignore any practice examples (marked with 例)
            - Do not translate any Japanese text
            - Do not include any section descriptions or other text
            - Output questions one after another with no extra text between them

            Here is the transcript to analyze:
            {cleaned_transcript}

            Return a JSON object with keys 'mondai1', 'mondai2', and 'mondai3', each containing a list of extracted questions for the corresponding section. If a section is not present or no questions are extracted, the value should be an empty list.
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
                        'mondai1': [],
                        'mondai2': [],
                        'mondai3': []
                    }
            except Exception as e:
                print(f"Error during Groq API call: {e}")
                # Fallback to basic logic if Groq API call fails
                return {
                    'mondai1': [],
                    'mondai2': [],
                    'mondai3': []
                }
        else:
            print("Groq API key not set. Skipping Groq processing.")
            return {
                'mondai1': [],
                'mondai2': [],
                'mondai3': []
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
    processed_data = processor.process_transcript("transcripts/sY7L5cfCWno.txt")
    
    print("Processed Data:")
    print(processed_data)
