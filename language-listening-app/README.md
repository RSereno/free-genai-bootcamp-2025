# Language Listening Comprehension App

## Project Overview
**Difficulty Level:** 300

A Language Listening Comprehension App that generates practice exercises based on YouTube content, designed for language learning test preparation.

## Business Context
As an Applied AI Engineer, this project aims to create an automated system that:
- Extracts content from YouTube language learning videos
- Generates listening comprehension exercises
- Provides interactive practice sessions

## Technical Stack
- **Frontend:** Streamlit
- **Database:** SQLite3 with vector store capabilities
- **AI/ML Components:**
  - LLM with Agent capabilities
  - Speech-to-Text (ASR) - Optional (e.g., Amazon Transcribe, OpenWhisper)
  - Text-to-Speech (TTS) - (e.g., Amazon Polly)
  - YouTube Transcript API
- **Development Tools:**
  - AI Coding Assistant (e.g., GitHub Copilot, Amazon CodeWhisperer)
  - Guardrails implementation

## Technical Challenges
- Vector store integration with SQLite3
- TTS/ASR availability and quality for target languages
- YouTube transcript accessibility
- Language model accuracy and context understanding

## Core Features
1. YouTube transcript extraction
2. Vector store data formatting
3. Topic-based question retrieval
4. Dynamic question generation
5. Audio generation for listening exercises

## Setup Instructions

### Environment Configuration
```bash
# Copy the example environment file
cp .env.example .env

# Edit .env with your actual API keys
nano .env
```

### API Configuration
```bash
export GROQ_API_KEY='your-api-key'
```

### Venv start
```bash
source .venv/bin/activate
```

### Running the Application
```bash
streamlit run frontend/main.py
```


### Technical uncertainity
Q: The transcript of some of the videos I tried dont seem to bring the Portuguese text on them. Or its completely gibberish will have to take another route. Keeping the Japanese version for now.

A: There seems to be one of the attributes that defines the language to get extracted.