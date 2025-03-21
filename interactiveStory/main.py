import streamlit as st
import chromadb
import json
from fastapi import FastAPI

# Initialize ChromaDB client
db_client = chromadb.PersistentClient(path="./chroma_db")
collection = db_client.get_or_create_collection("language_learning")

# Load story data from JSON
def load_story_data(file_name="story_data.json"):
    with open(file_name, "r", encoding="utf-8") as file:
        story_scenes = json.load(file)
        for scene in story_scenes:
            # Convert the scene dictionary to a JSON string
            document = json.dumps(scene)

            # Add the scene to the collection
            collection.add(
                documents=[document],  # Pass the JSON string
                ids=[scene["id"]]
            )

# FastAPI Backend
app = FastAPI()

@app.get("/story/{scene_id}")
def get_scene(scene_id: str):
    results = collection.get(ids=[scene_id])
    if results["documents"]:
        return json.loads(results["documents"][0])  # Convert string back to dict
    return {"error": "Scene not found"}

@app.post("/story")
def add_scene(scene: dict):
    try:
        # Convert the scene dictionary to a JSON string
        document = json.dumps(scene)

        # Add the scene to the collection
        collection.add(
            documents=[document],  # Pass the JSON string
            ids=[scene["id"]]
        )
        return {"message": "Scene added"}
    except Exception as e:
        return {"error": f"Failed to add scene: {str(e)}"}

# Load initial story data
STORY_FILE = "story_data.json"  # You can change this to your story file path
try:
    load_story_data(STORY_FILE)
    load_story_data("stories/cinema.json")
    print(f"Successfully loaded story from {STORY_FILE}")
except FileNotFoundError:
    print(f"Story file {STORY_FILE} not found. Please create it first.")
except json.JSONDecodeError:
    print(f"Error parsing {STORY_FILE}. Please check if it's valid JSON.")

# Streamlit Frontend
st.title("🌍 Interactive Language Learning Story")

if "scene_id" not in st.session_state:
    st.session_state.scene_id = "start"
if "points" not in st.session_state:
    st.session_state.points = 0

scene = get_scene(st.session_state.scene_id)
st.write("Debug - Current scene_id:", st.session_state.scene_id)
st.write("Debug - Scene content:", scene)

if "error" in scene:
    st.error("Nenhuma cena encontrada.")
elif "prompt" in scene:
    st.write(scene["prompt"])
    options = scene.get("options", [])
    
    for option in options:
        if st.button(option["text"]):
            st.session_state.scene_id = option["nextScene"]
            st.session_state.points += option.get("points", 0)
            st.rerun()
    
    # Display hint if available in the scene or options
    hint = scene.get("hint")
    if not hint:  # Check if hint is inside options
        for option in options:
            if "hint" in option:
                hint = option["hint"]
                break
    if hint:
        st.info(f"Hint: {hint}")
    
    # Display current points
    st.sidebar.write(f"🌟 Points: {st.session_state.points}")
else:
    st.error("Erro: Nenhuma cena encontrada. Verifique o backend. 🚨")