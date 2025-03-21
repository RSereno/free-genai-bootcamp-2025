import chromadb
import json
from fastapi import FastAPI

app = FastAPI()

class DatabaseOperations:
    def __init__(self):
        # Initialize the database client and collection
        self.db_client = chromadb.PersistentClient("./chroma_db")
        self.collection = self.db_client.get_or_create_collection("language_learning")

    def add_scene_to_collection(self, scene):
        """Add a single scene to the collection."""
        document = json.dumps(scene)
        self.collection.add(
            documents=[document],  # Pass the JSON string
            ids=[scene["id"]]
        )

    def clear_collection(self):
        """Clear all documents from the collection."""
        all_ids = self.collection.get()["ids"]
        if all_ids:
            self.collection.delete(ids=all_ids)
            print("Cleared existing story data from the database.")
        else:
            print("No existing story data to clear.")

    def get_scene_by_id(self, scene_id):
        """Retrieve a scene by its ID."""
        results = self.collection.get(ids=[scene_id])
        if results["documents"]:
            return json.loads(results["documents"][0])  # Convert string back to dict
        return {"error": "Scene not found"}

db_operations = DatabaseOperations()

@app.get("/story/{scene_id}")
def get_scene(scene_id: str):
    return db_operations.get_scene_by_id(scene_id)

@app.post("/story")
def add_scene(scene: dict):
    try:
        db_operations.add_scene_to_collection(scene)
        return {"message": "Scene added"}
    except Exception as e:
        return {"error": f"Failed to add scene: {str(e)}"}
