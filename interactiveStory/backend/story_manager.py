import os
import json
from backend.backend import DatabaseOperations

db_ops = DatabaseOperations()

def load_story_data(file_name="story_data.json"):
    with open(file_name, "r", encoding="utf-8") as file:
        story_scenes = json.load(file)
        for scene in story_scenes:
            db_ops.add_scene_to_collection(scene)

def load_story_data_from_folder(folder_path="stories"):
    try:
        for file_name in os.listdir(folder_path):
            if file_name.endswith(".json"):
                file_path = os.path.join(folder_path, file_name)
                with open(file_path, "r", encoding="utf-8") as file:
                    story_scenes = json.load(file)
                    for scene in story_scenes:
                        db_ops.add_scene_to_collection(scene)
        print(f"Successfully loaded stories from folder: {folder_path}")
    except Exception as e:
        print(f"Error loading stories from folder {folder_path}: {str(e)}")

def refresh_story_data(folder_path="stories"):
    try:
        db_ops.clear_collection()
        load_story_data_from_folder(folder_path)
        print(f"Successfully refreshed stories from folder: {folder_path}")
    except Exception as e:
        print(f"Error refreshing stories from folder {folder_path}: {str(e)}")
