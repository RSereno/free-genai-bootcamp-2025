import streamlit as st
from backend.story_manager import refresh_story_data  # Import the refactored function
from backend.backend import get_scene
from frontend.utils import display_scene

# Refresh story data on app start
STORY_FOLDER = "stories"  # Specify the folder containing story files
try:
    refresh_story_data(STORY_FOLDER)
except Exception as e:
    print(f"Failed to refresh stories: {str(e)}")

# Streamlit Frontend
st.title("🌍 Interactive Language Learning Story")

if "scene_id" not in st.session_state:
    st.session_state.scene_id = "start"
if "points" not in st.session_state:
    st.session_state.points = 0

scene = get_scene(st.session_state.scene_id)
display_scene(scene)

    