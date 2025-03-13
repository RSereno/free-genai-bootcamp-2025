import streamlit as st
import random
from services.sentence_generator import SentenceGenerator
from dotenv import load_dotenv
import requests

# Load environment variables
load_dotenv()

# Function to fetch word collection
def fetch_words(id):
    """Fetch words from a specific group via the backend API"""
    try:
        # Call the backend API
        response = requests.get(f"http://localhost:8080/api/words_groups/{id}/words", 
                              params={"page": 1, "limit": 100})
        response.raise_for_status()  # Raise exception for 4xx or 5xx status codes
        
        # Parse the response
        data = response.json()
        
        # Transform the data into the expected format
        word_list = []
        for word in data.get("items", []):
            word_list.append({
                "portuguese": word["portuguese"],
                "english": word["english"]
            })
        
        return word_list

    except requests.RequestException as e:
        print(f"Error fetching words from API: {e}")
        # Return mock data as fallback
        return [
            {"portuguese": "casa", "english": "house"},
            {"portuguese": "carro", "english": "car"},
            {"portuguese": "gato", "english": "cat"}
        ]

# Initialize session state variables
if "app_state" not in st.session_state:
    st.session_state.app_state = "Setup"
if "words" not in st.session_state:
    st.session_state.words = fetch_words(1)  # Assuming id=1
if "current_sentence" not in st.session_state:
    st.session_state.current_sentence = ""
if "submission" not in st.session_state:
    st.session_state.submission = {}

# Function to generate a sentence (mocked)
def generate_sentence():
    if st.session_state.words:
        word = random.choice(st.session_state.words)["english"]
        print("\n=== Starting Sentence Generation ===")
        print(f"Selected word: {word}")
        try:
            generator = SentenceGenerator()
            sentence = generator.generate_sentence(word)
            print(f"Generated sentence: {sentence}")
            st.session_state.current_sentence = sentence
            st.session_state.app_state = "Practice"
        except Exception as e:
            print(f"Error in app.py: {str(e)}")
            sentence = f"I see a {word} today. Mock"  # Fallback to simple sentence
            st.session_state.current_sentence = sentence
            st.session_state.app_state = "Practice"
        print("=== End Sentence Generation ===\n")
    else:
        st.error("No words available to generate a sentence.")

# Function to submit image for review (mocked)
def submit_for_review(image):
    if image is not None:
        # Check image size (max 5MB)
        if len(image.getvalue()) > 5 * 1024 * 1024:
            st.error("Image size exceeds 5MB. Please upload a smaller image.")
        else:
            # Mock grading response
            transcription = "Eu vejo um gato hoje."  # Mock Portuguese translation
            translation = "I see a cat today."     # Mock English translation
            grade = "S"                              # Mock grade
            description = "Perfect match - Translation is exactly as expected"
            st.session_state.submission = {
                "transcription": transcription,
                "translation": translation,
                "grade": grade,
                "description": description
            }
            st.session_state.app_state = "Review"
    else:
        st.error("Please upload an image before submitting.")

# Function to proceed to the next question
def next_question():
    st.session_state.submission = {}  # Clear previous submission
    generate_sentence()               # Generate new sentence and go to Practice

# Main application logic
st.title("Language Learning App")

if st.session_state.app_state == "Setup":
    ### Setup State ###
    st.write("Welcome to the Language Learning App!")
    st.write("Click the button below to generate a sentence and start practicing.")
    st.button("Generate Sentence", on_click=generate_sentence)

elif st.session_state.app_state == "Practice":
    ### Practice State ###
    st.write("Translate the following sentence into Portuguese and upload an image of your handwritten translation:")
    st.write(f"**{st.session_state.current_sentence}**")
    image = st.file_uploader("Upload your handwritten translation", type=["jpg", "png"])
    if st.button("Submit for Review"):
        submit_for_review(image)

elif st.session_state.app_state == "Review":
    ### Review State ###
    st.write("Here is the review of your submission:")
    st.write(f"**Original Sentence:** {st.session_state.current_sentence}")
    st.write(f"**Transcription:** {st.session_state.submission['transcription']}")
    st.write(f"**Translation:** {st.session_state.submission['translation']}")
    st.write(f"**Grade:** {st.session_state.submission['grade']}")
    st.write(f"**Description:** {st.session_state.submission['description']}")
    st.button("Next Question", on_click=next_question)