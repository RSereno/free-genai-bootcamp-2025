# Interactive Story Mode - Technical Specifications

## Overview
This project is an **Interactive Language Learning Story Mode** application designed to help users learn a new language by making choices in a branching narrative. The system consists of:
- **Frontend:** Streamlit (for an interactive UI)
- **Backend:** FastAPI (for API endpoints)
- **Database:** ChromaDB (for storing and retrieving story data)

---

## Architecture
### **Frontend (Streamlit)**
- Displays the interactive story to users.
- Retrieves scenes from the FastAPI backend.
- Provides buttons for users to select responses.
- Stores user session state to track progress.

### **Backend (FastAPI)**
- Handles API requests for retrieving and adding story scenes.
- Communicates with ChromaDB for data storage.
- Provides endpoints:
  - `GET /story/{scene_id}`: Retrieves a scene by ID.
  - `POST /story`: Adds a new scene to the database.

### **Database (ChromaDB)**
- Stores story scenes as JSON documents.
- Supports fast retrieval based on scene ID.
- Persistent storage to maintain user progress and stories.

---

## Data Structure
Each story scene follows a structured JSON format:
```json
{
  "id": "cafe_entrance",
  "prompt": "Um empregado de mesa cumprimenta-o. O que você responde?",
  "options": [
    { "text": "Quero um café, por favor.", "nextScene": "waiter_serves_coffee", "points": 10 },
    { "text": "Obrigado, adeus!", "nextScene": "waiter_looks_confused", "points": 0, "hint": "Está a sair demasiado cedo!" }
  ]
}
```

---

## Installation & Setup
### **Prerequisites**
- Python 3.8+
- FastAPI
- Streamlit
- ChromaDB

### **Installation**
1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/interactive-story
   cd interactive-story
   ```
# UV enviroment
https://docs.astral.sh/uv/#getting-started

curl -LsSf https://astral.sh/uv/install.sh | sh

uv venv intStory --python 3.12 --seed

source intStory/bin/activate

2. Install dependencies:
   ```bash
   pip install fastapi streamlit chromadb uvicorn
   ```
3. Run the Streamlit frontend:
   ```bash
   streamlit run main.py
   ```

---

## Future Enhancements
- **Speech Recognition:** Allow users to speak responses instead of selecting.
- **User Progress Tracking:** Save user choices and progression.
- **Multimedia Support:** Add images and audio clips to enrich storytelling.
- **AI-driven Conversations:** Implement an AI chatbot for open-ended interactions.

