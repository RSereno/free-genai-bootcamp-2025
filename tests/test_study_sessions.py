import pytest
from app import app

def test_create_study_session_success(client):
    """Test successful creation of a study session"""
    response = client.post('/api/study-sessions', json={
        'group_id': 1,
        'study_activity_id': 1
    })
    
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