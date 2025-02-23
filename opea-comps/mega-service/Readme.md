## Building opea-comps from Source
```sh
git clone https://github.com/opea-project/GenAIComps
cd GenAIComps
pip install -e .
´´´

## Python Package Installation

To install the package locally, first install the `python3-full` package, then create and activate a virtual environment:

```sh
# Install python3-full package
sudo apt-get update
sudo apt-get install python3-full

# Create a virtual environment
python3 -m venv .venv

# Activate the virtual environment
# On Linux/macOS:
source .venv/bin/activate

# Install the package in editable mode
pip install -e .
```

Note: The virtual environment helps isolate project dependencies from the system Python installation.

Note: On Ubuntu/Debian systems, ensure you have the necessary system dependencies installed.

## Mega Service
Documentation link https://github.com/opea-project/GenAIComps

## Setup

To ensure that the `comps` module can be imported correctly, you need to add the `GenAIComps` directory (the folder where comps was build) to your Python path. You can do this by modifying the `sys.path` at runtime or by setting the `PYTHONPATH` environment variable.

### Modify `sys.path` at Runtime
In your `example.py` file, add the following code to include the `GenAIComps` directory in your Python path:

```python
import sys
import os

# Add the path to the GenAIComps directory
sys.path.append(os.path.expanduser('~/code/bootcamp/GenAIComps'))
```

### Set PYTHONPATH Environment Variable
Alternatively, you can set the PYTHONPATH environment variable to include the GenAIComps directory before running your script. You can do this in the terminal:

```sh
export PYTHONPATH=$PYTHONPATH:~/code/bootcamp/GenAIComps
```

# Technical Uncertainty
 Q: To be able to implement this need to understand the concepts in example provided. Mainly EMBEDDING_SERVICE and LLM_SERVICE. The LLM Service is the LLM it self but what is the Embedding one.

 A: Found my answer here https://github.com/opea-project/GenAIComps/blob/main/comps/embeddings/src/README.md