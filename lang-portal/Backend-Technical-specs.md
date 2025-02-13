# Backend Server Technical Specs

## Business Goal: 
A language learning school wants to build a prototype of learning portal which will act as three things:
- Inventory of possible vocabulary that can be learned
- Act as a  Learning record store (LRS), providing correct and wrong score on practice vocabulary
- A unified launchpad to launch different learning apps

## Core Functionalities
- The backend will be built using go
- The backend will be responsible for serving the frontend
- The database will be SQLlite3
- The API will be built using Gin framework
- The API will always return JSON
- The will not be authentication or authorization single user

## Database Schema

We have the following tables:
- words - store the vocabulary words
- words groups - join table between words and groups many-to-many
- groups - thematic groups of words
- study_sessions - records of study sessions grouping word_review_items
- study_activities - a specific instance of a study activity, linking a study_session to a group
- word_review_items - a record of word practive, determining if the word was correctly recalled or not

### Words
- id integer primary key autoincrement
- word text not null
- meaning text not null
- group_id integer not null
- json_data text not null

### Words Groups
- id integer primary key autoincrement
- words_ids text not null   
- group_id integer not null

### Groups
- id integer primary key autoincrement
- name text not null
- description text null

### Study Sessions
- id integer primary key autoincrement
- group_id integer not null
- created_at datetime not null 
- study_activity_id integer not null    

### Study Activities
- id integer primary key autoincrement
- study_session_id integer not null
- group_id integer not null
- created_at datetime not null

### Word Review Items
- id integer primary key autoincrement
- word_id integer not null
- study_session_id integer not null
- is_correct boolean not null
- created_at datetime not null

## API

### Dashboard
- GET /api/dashboard/last-study-sessions
    - Returns:
        - study_sessions: array of recent study sessions
        ```json
        {
          "study_sessions": [
            {
              "id": 1,
              "group_id": 1,
              "group_name": "Basic Verbs",
              "created_at": "2024-03-20T15:30:00Z",
              "study_activity_id": 1
            }
          ]
        }
        ```

- GET /api/dashboard/study-progress
    - Returns:
        - total_known: integer
        - total_available_words: integer
        ```json
        {
          "total_known": 150,
          "total_available_words": 200
        }
        ```

- GET /api/dashboard/quick-stats
    - Returns:
        - success_rate: float
        - total_words_studied: integer
        - total_words_correct: integerç
        - total_study_sessions: integer
        - total_active_groups: integer
        - current_streak: integer
        - longest_streak: integer
        ```json
        {
          "success_rate": 85.5,
          "total_words_studied": 500,
          "total_words_correct": 427,
          "total_study_sessions": 25,
          "total_active_groups": 5,
          "current_streak": 3,
          "longest_streak": 7,
          "last_study_date": "2024-03-20T15:30:00Z"
        }
        ```

### Study Activities
- GET /api/study-activities
    - Query Params:
        - page: integer
        - limit: integer
    - Returns:
        - study_activities: array of study activities
        - pagination: object
            - page_number: integer
            - page_size: integer
            - total_pages: integer
            - total_items: integer
        ```json
        {
          "study_activities": [
            {
              "id": 1,
              "name": "Flashcards",
              "group_id": 1,
              "group_name": "Basic Verbs",
              "created_at": "2024-03-20T15:30:00Z",
              "completed": true,
              "success_rate": 85.5,
              "total_words": 20,
              "correct_count": 17,
              "incorrect_count": 3
            }
          ],
          "pagination": {
            "page_number": 1,
            "page_size": 10,
            "total_pages": 5,
            "total_items": 45
          }
        }
        ```

- GET /api/study-activities/:id 
    - Returns:
        - study_activity: object
            - id: integer
            - name: string
            - group_id: integer
            - created_at: datetime
            - word_reviews: array of word_review_items
        ```json
        {
          "study_activity": {
            "id": 1,
            "name": "Flashcards",
            "group_id": 1,
            "group_name": "Basic Verbs",
            "created_at": "2024-03-20T15:30:00Z",
            "completed": true,
            "success_rate": 85.5,
            "word_reviews": [
              {
                "id": 1,
                "word": {
                  "id": 1,
                  "word": "run",
                  "meaning": "to move at a speed faster than walking"
                },
                "is_correct": true,
                "created_at": "2024-03-20T15:31:00Z"
              }
            ]
          }
        }
        ```

- GET /api/study-activities/:id/study-session
    - Returns:
        - session_data: object
            - words: array of words to study
            - settings: object with activity settings
        ```json
        {
          "session_data": {
            "words": [
              {
                "id": 1,
                "word": "run",
                "meaning": "to move at a speed faster than walking",
                "previous_attempts": {
                  "correct_count": 5,
                  "incorrect_count": 1,
                  "last_studied": "2024-03-19T10:30:00Z"
                }
              }
            ],
            "settings": {
              "review_mode": "flashcards",
              "words_per_session": 20,
              "time_limit_minutes": 15,
              "show_hints": true
            }
          }
        }
        ```

- POST /api/study-activities/
    - Required params:
        - group_id: integer
        - study_activity_id: integer
    - Request body:
        ```json
        {
          "group_id": 1,
          "study_activity_id": 1
        }
        ```
    - Returns:
        - study_activity: newly created study activity object
        ```json
        {
          "study_activity": {
            "id": 2,
            "name": "Flashcards",
            "group_id": 1,
            "group_name": "Basic Verbs",
            "created_at": "2024-03-20T15:30:00Z",
            "completed": false,
            "settings": {
              "review_mode": "flashcards",
              "words_per_session": 20,
              "time_limit_minutes": 15,
              "show_hints": true
            }
          }
        }
        ```

### Words
- GET /api/words
    - Query Params:
        - page: integer
        - limit: integer
        - search: string
        - group_id: integer
        - meaning: string
        - correct_count: integer
        - incorrect_count: integer
    - Returns:
        - words: array of words
        - pagination: object
            - page_number: integer
            - page_size: integer
            - total_pages: integer
            - total_items: integer
        ```json
        {
          "words": [
            {
              "id": 1,
              "word": "run",
              "meaning": "to move at a speed faster than walking",
              "json_data": {
                "pronunciation": "rʌn",
                "examples": [
                  "She runs every morning",
                  "The train runs on time"
                ],
                "part_of_speech": "verb"
              },
              "study_statistics": {
                "correct_count": 15,
                "incorrect_count": 3,
                "success_rate": 83.3,
                "last_studied": "2024-03-20T15:30:00Z"
              },
              "groups": [
                {
                  "id": 1,
                  "name": "Basic Verbs"
                }
              ]
            }
          ],
          "pagination": {
            "page_number": 1,
            "page_size": 10,
            "total_pages": 20,
            "total_items": 195
          }
        }
        ```

- GET /api/words/:id
    - Returns:
        - word: object
            - id: integer
            - word: string
            - meaning: string
            - study_statistics: object
                - correct_count: integer
                - incorrect_count: integer
                - total_study_sessions: integer
                - last_study_session_date: datetime
            - groups: array of group objects
        ```json
        {
          "word": {
            "id": 1,
            "word": "run",
            "meaning": "to move at a speed faster than walking",
            "json_data": {
              "pronunciation": "rʌn",
              "examples": [
                "She runs every morning",
                "The train runs on time"
              ],
              "part_of_speech": "verb",
              "difficulty_level": "A1",
              "tags": ["action", "movement"]
            },
            "study_statistics": {
              "correct_count": 15,
              "incorrect_count": 3,
              "total_study_sessions": 6,
              "success_rate": 83.3,
              "last_study_session_date": "2024-03-20T15:30:00Z",
              "study_history": [
                {
                  "date": "2024-03-20T15:30:00Z",
                  "is_correct": true,
                  "study_session_id": 123
                }
              ]
            },
            "groups": [
              {
                "id": 1,
                "name": "Basic Verbs",
                "description": "Common everyday verbs"
              }
            ]
          }
        }
        ```

### Words Groups
- GET /api/words-groups
    - Query Params:
        - page: integer
        - limit: integer
    - Returns:
        - groups: array of group objects
            - id: integer
            - name: string
            - words_count: integer
        - pagination: object
            - page_number: integer
            - page_size: integer
            - total_pages: integer
            - total_items: integer
        ```json
        {
          "groups": [
            {
              "id": 1,
              "name": "Basic Verbs",
              "description": "Common everyday verbs",
              "words_count": 50,
              "study_statistics": {
                "total_sessions": 10,
                "success_rate": 85.5,
                "last_studied": "2024-03-20T15:30:00Z"
              }
            }
          ],
          "pagination": {
            "page_number": 1,
            "page_size": 10,
            "total_pages": 3,
            "total_items": 25
          }
        }
        ```

- GET /api/words-groups/:id
    - Returns:
        - group: object
            - id: integer
            - name: string
            - description: string
            - words_count: integer
            - words: array of word objects
        ```json
        {
          "group": {
            "id": 1,
            "name": "Basic Verbs",
            "description": "Common everyday verbs",
            "words_count": 50,
            "study_statistics": {
              "total_sessions": 10,
              "success_rate": 85.5,
              "last_studied": "2024-03-20T15:30:00Z"
            },
            "words": [
              {
                "id": 1,
                "word": "run",
                "meaning": "to move at a speed faster than walking",
                "study_statistics": {
                  "correct_count": 15,
                  "incorrect_count": 3,
                  "success_rate": 83.3
                }
              }
            ]
          }
        }
        ```

- GET /api/words-groups/:id/words
    - Query Params:
        - page: integer
        - limit: integer
    - Returns:
        - words: array of word objects
        - pagination: object
            - page_number: integer
            - page_size: integer
            - total_pages: integer
            - total_items: integer
        ```json
        {
          "words": [
            {
              "id": 1,
              "word": "run",
              "meaning": "to move at a speed faster than walking",
              "study_statistics": {
                "correct_count": 15,
                "incorrect_count": 3,
                "success_rate": 83.3
              }
            }
          ],
          "pagination": {
            "page_number": 1,
            "page_size": 10,
            "total_pages": 5,
            "total_items": 50
          }
        }
        ```

- GET /api/words-groups/:id/study-sessions
    - Query Params:
        - page: integer
        - limit: integer
    - Returns:
        - study_sessions: array of study session objects
        - pagination: object
            - page_number: integer
            - page_size: integer
            - total_pages: integer
            - total_items: integer
        ```json
        {
          "study_sessions": [
            {
              "id": 1,
              "created_at": "2024-03-20T15:30:00Z",
              "activity_name": "Flashcards",
              "statistics": {
                "total_words": 20,
                "correct_count": 17,
                "incorrect_count": 3,
                "success_rate": 85.0,
                "duration_minutes": 15
              }
            }
          ],
          "pagination": {
            "page_number": 1,
            "page_size": 10,
            "total_pages": 3,
            "total_items": 25
          }
        }
        ```

### Study Sessions
- GET /api/study-sessions
    - Query Params:
        - page: integer
        - limit: integer
    - Returns:
        - study_sessions: array of study sessions
        - pagination: object
            - page_number: integer
            - page_size: integer
            - total_pages: integer
            - total_items: integer
        ```json
        {
          "study_sessions": [
            {
              "id": 1,
              "group_id": 1,
              "group_name": "Basic Verbs",
              "activity_name": "Flashcards",
              "created_at": "2024-03-20T15:30:00Z",
              "completed_at": "2024-03-20T15:45:00Z",
              "statistics": {
                "total_words": 20,
                "correct_count": 17,
                "incorrect_count": 3,
                "success_rate": 85.0,
                "duration_minutes": 15
              }
            }
          ],
          "pagination": {
            "page_number": 1,
            "page_size": 10,
            "total_pages": 5,
            "total_items": 45
          }
        }
        ```

- GET /api/study-sessions/:id
    - Returns:
        - study_session: object
            - id: integer
            - activity_name: string
            - group_name: string
            - start_time: datetime
            - end_time: datetime
            - review_items_count: integer
            - review_items: array of objects
                - id: integer
                - word: object
                    - id: integer
                    - word: string
                    - meaning: string
                - is_correct: boolean
                - created_at: datetime
            - statistics: object
                - correct_count: integer
                - incorrect_count: integer
                - success_rate: float
            - group: object
                - id: integer
                - name: string
            - study_activity: object
                - id: integer
                - name: string
        ```json
        {
          "study_session": {
            "id": 1,
            "activity_name": "Flashcards",
            "group_name": "Basic Verbs",
            "start_time": "2024-03-20T15:30:00Z",
            "end_time": "2024-03-20T15:45:00Z",
            "review_items_count": 20,
            "review_items": [
              {
                "id": 1,
                "word": {
                  "id": 1,
                  "word": "run",
                  "meaning": "to move at a speed faster than walking"
                },
                "is_correct": true,
                "created_at": "2024-03-20T15:31:00Z",
                "response_time_ms": 2500
              }
            ],
            "statistics": {
              "correct_count": 17,
              "incorrect_count": 3,
              "success_rate": 85.0,
              "average_response_time_ms": 3000
            },
            "group": {
              "id": 1,
              "name": "Basic Verbs"
            },
            "study_activity": {
              "id": 1,
              "name": "Flashcards",
              "settings": {
                "review_mode": "flashcards",
                "words_per_session": 20,
                "time_limit_minutes": 15
              }
            }
          }
        }
        ```

### Word Review Items
- GET /api/word-review-items
    - Query Params:
        - study_session_id: integer
        - word_id: integer
        - page: integer
        - limit: integer
    - Returns:
        - word_review_items: array of review items
        - pagination: object
            - page_number: integer
            - page_size: integer
            - total_pages: integer
            - total_items: integer
        ```json
        {
          "word_review_items": [
            {
              "id": 1,
              "word": {
                "id": 1,
                "word": "run",
                "meaning": "to move at a speed faster than walking"
              },
              "study_session_id": 1,
              "is_correct": true,
              "created_at": "2024-03-20T15:31:00Z",
              "response_time_ms": 2500,
              "session_details": {
                "activity_name": "Flashcards",
                "group_name": "Basic Verbs"
              }
            }
          ],
          "pagination": {
            "page_number": 1,
            "page_size": 10,
            "total_pages": 5,
            "total_items": 45
          }
        }
        ```

### Settings
- POST /api/settings
    - Required params:
        - theme: string (light|dark|system)
    - Request body:
        ```json
        {
          "theme": "dark",
          "words_per_session": 20,
          "show_hints": true,
          "default_study_time_minutes": 15
        }
        ```
    - Returns:
        - settings: object
            - theme: string
            - updated_at: datetime
        ```json
        {
          "settings": {
            "theme": "dark",
            "words_per_session": 20,
            "show_hints": true,
            "default_study_time_minutes": 15,
            "updated_at": "2024-03-20T15:30:00Z"
          }
        }
        ```

- POST /api/settings/reset-progress
    - Returns:
        - success: boolean
        - message: string
        ```json
        {
          "success": true,
          "message": "All progress has been reset successfully",
          "reset_timestamp": "2024-03-20T15:30:00Z",
          "affected_items": {
            "study_sessions": 25,
            "word_reviews": 500,
            "statistics": "reset"
          }
        }
        ```

- POST /api/settings/reload-seed-data
    - Returns:
        - success: boolean
        - message: string
        ```json
        {
          "success": true,
          "message": "Seed data reloaded successfully",
          "loaded_items": {
            "words": 200,
            "groups": 10,
            "relationships": 250
          },
          "timestamp": "2024-03-20T15:30:00Z"
        }
        ```

## Doc

