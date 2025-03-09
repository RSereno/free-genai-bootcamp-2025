import sys
import os

# Update path to point to backend-flask directory
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../lang-portal/backend-flask')))

import pytest
from app import create_app  # Import create_app instead of app directly
from datetime import datetime

@pytest.fixture
def app():
    app = create_app({
        'TESTING': True,
        'DATABASE': ':memory:'  # Use in-memory database for testing
    })
    return app

@pytest.fixture
def client(app):
    # Create test client
    test_client = app.test_client()
    
    # Set up application context
    ctx = app.app_context()
    ctx.push()
    
    # Set up database
    cursor = app.db.cursor()
    
    # Drop tables if they exist
    cursor.execute('DROP TABLE IF EXISTS word_review_items')
    cursor.execute('DROP TABLE IF EXISTS study_sessions')
    cursor.execute('DROP TABLE IF EXISTS words')
    cursor.execute('DROP TABLE IF EXISTS groups')
    cursor.execute('DROP TABLE IF EXISTS study_activities')
    
    # Create tables
    cursor.execute('''
        CREATE TABLE groups (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            words_count INTEGER DEFAULT 0
        )
    ''')
    
    cursor.execute('''
        CREATE TABLE study_activities (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            url TEXT NOT NULL,
            preview_url TEXT
        )
    ''')
    
    cursor.execute('''
        CREATE TABLE study_sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            group_id INTEGER NOT NULL,
            study_activity_id INTEGER NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (group_id) REFERENCES groups(id),
            FOREIGN KEY (study_activity_id) REFERENCES study_activities(id)
        )
    ''')

    cursor.execute('''
        CREATE TABLE word_review_items (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            study_session_id INTEGER NOT NULL,
            word_id INTEGER NOT NULL,
            correct BOOLEAN NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (study_session_id) REFERENCES study_sessions(id)
        )
    ''')
    
    # Insert test data
    cursor.execute('INSERT INTO groups (id, name) VALUES (1, "Test Group")')
    cursor.execute('''
        INSERT INTO study_activities (id, name, url, preview_url) 
        VALUES (1, "Test Activity", "http://test.com", "http://test.com/preview")
    ''')
    
    app.db.commit()
    
    yield test_client
    
    # Cleanup
    ctx.pop()

@pytest.fixture
def runner(app):
    return app.test_cli_runner()

def test_create_study_session_success(client):
    """Test successful creation of a study session"""
    response = client.post('/api/study-sessions', json={
        'group_id': 1,
        'study_activity_id': 1
    })
    
    # Add debug output
    print("\nResponse status:", response.status_code)
    print("Response data:", response.get_json())
    
    assert response.status_code == 201
    data = response.get_json()
    assert 'id' in data
    assert data['group_id'] == 1
    assert data['activity_id'] == 1
    assert 'group_name' in data
    assert 'activity_name' in data
    assert 'start_time' in data
    assert 'end_time' in data
    assert 'review_items_count' in data

def test_create_study_session_missing_fields(client):
    """Test validation of required fields"""
    response = client.post('/api/study-sessions', json={})
    assert response.status_code == 400
    assert response.get_json()['error'] == 'Missing required fields'

def test_create_study_session_invalid_group(client):
    """Test handling of invalid group_id"""
    response = client.post('/api/study-sessions', json={
        'group_id': 99999,
        'study_activity_id': 1
    })
    assert response.status_code == 404
    assert response.get_json()['error'] == 'Group not found'

def test_create_study_session_invalid_activity(client):
    """Test handling of invalid study_activity_id"""
    response = client.post('/api/study-sessions', json={
        'group_id': 1,
        'study_activity_id': 99999
    })
    assert response.status_code == 404
    assert response.get_json()['error'] == 'Study activity not found' 