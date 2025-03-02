import chromadb
from pathlib import Path

# setup Chroma in-memory, for easy prototyping. Can add persistence easily!
client = chromadb.Client()

# Create collection. get_collection, get_or_create_collection, delete_collection also available!
collection = client.create_collection("listening-collection")

def load_documents(directory_path):
    documents = []
    metadatas = []
    ids = []
    
    directory = Path(directory_path)
    for i, file_path in enumerate(directory.glob('*.txt')):  # Adjust file pattern as needed
        with open(file_path, 'r', encoding='utf-8') as file:
            content = file.read()
            documents.append(content)
            metadatas.append({"source": str(file_path)})
            ids.append(f"doc_{i}")
    
    return documents, metadatas, ids

# Example usage:
documents, metadatas, ids = load_documents("path/to/your/documents")
collection.add(
    documents=documents,
    metadatas=metadatas,
    ids=ids
)

# Query/search 2 most similar results. You can also .get by id
results = collection.query(
    query_texts=["This is a query document"],
    n_results=2,
    # where={"metadata_field": "is_equal_to_this"}, # optional filter
    # where_document={"$contains":"search_string"}  # optional filter
)