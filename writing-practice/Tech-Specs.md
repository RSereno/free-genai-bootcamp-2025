# Technical Specs

## Initialization Step
When the app first initializes it needs to the following:
Fetch from the GET localhost:{serverPort}/api/words_groups/:id/study_sessions/raw (currently missing), this will return a collection of words in a json structure. It will have portuguese words with their english translation. We need to store this collection of words in memory

## API Endpoints
The application requires the following API endpoints:

### GET /api/words_groups/:id/study_sessions/raw
Returns a collection of Portuguese-English word pairs.
Response format:
```json
{
  "words": [
    {
      "portuguese": "string",
      "english": "string"
    }
  ]
}
```

## Page States

Page states describes the state the single page application should behaviour from a user's perspective. 

### Setup State
When a user first's start up the app.
They will only see a button called "Generate Sentence"
When they press the button, The app will generate a sentence using
the Sentence Genreator LLM, and the state will move to Practice State

### Practice State
When a user in is a practice state,
they will see an english sentence,
They will also see an upload field under the english sentence
They will see a button called "Submit for Review"
When they press the Submit for Review Button an uploaded image
will be passed to the Grading System and then will transition to the Review State

### Review State
 When a user in is the review review state,
 The user will still see the english sentence.
 The upload field will be gone.
 The user will now see a review of the output from the Grading System:
- Transcription of Image
- Translation of Transcription
- Grading
  - a letter score using the S Rank to score
  - a description of whether the attempt was accurarte to the english sentence and suggestions.
There will be a button called "Next Question" when clicked
it will it generate a new question and place the app into Practice State

## State Management
The application should maintain the following state:
- Current word collection (from API)
- Current sentence
- Current application state (Setup, Practice, or Review)
- Current submission data (during Review state)

## Error Handling
The application should handle:
- API fetch failures
- Image upload failures
- OCR transcription failures
- Invalid image formats
- LLM service failures

## Performance Requirements
- Image upload size should be limited to 5MB
- LLM responses should timeout after 10 seconds
- Word collection should be cached in memory

## Sentence Generator LLM Prompt
Generate a simple sentence using the following word: {{word}}
The grammar should be scoped to A1 level.
You can use the following vocabulary to construct a simple sentence:
- Simple objects (e.g., book, car, ramen, sushi)
- Simple verbs (e.g., to drink, to eat, to meet)
- Simple time expressions (e.g., tomorrow, today, yesterday)

## Grading System
The Grading System will:
1. Transcribe the uploaded image using OCR
2. Use an LLM to translate the transcribed Portuguese text to English
3. Use a separate LLM to evaluate and grade the translation
4. Return the combined results to the frontend application

## Grading System Details
Grade scale and criteria:
- S: Perfect match - Translation is exactly as expected
- A: Minor differences - Meaning preserved with slight variations
- B: Some inaccuracies - Generally correct but with minor errors
- C: Significant differences - Core meaning partially preserved
- D: Major errors - Meaning significantly altered
- F: Failed - Completely incorrect or processing error

## Technical Notes
1. The `/api/words_groups/:id/study_sessions/raw` endpoint needs to be implemented in the Go backend
2. Expected response format is already defined in the API Endpoints section
3. This endpoint should return a curated list of word pairs for study sessions

## Techincal Uncertainity
 Q: I dont seem to have the "raw" endpoint in the GO backend will have to make the decision if I quickly implement it or fall back to the flask version.