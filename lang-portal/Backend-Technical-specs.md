# Backend Server Technical Specs

## Business Goal: 
A language learning school wants to build a prototype of learning portal which will act as three things:
- Inventory of possible vocabulary that can be learned
- Act as a  Learning record store (LRS), providing correct and wrong score on practice vocabulary
- A unified launchpad to launch different learning apps

## Core Functionalities
- The backend will be built using Go
- The backend will be responsible for serving the frontend
- The database will be SQLlite3
- The API will be built using Gin framework - Mage is a task runner for Go.
- The API will always return JSON
- The will not be authentication or authorization
- Everything will be treated as a single user

## Database Schema
Our database will be a single sqlite database called ´words.db´ that will be in the root of the project folder of ´backend_go´

We have the following tables:
- words - store the vocabulary words
- words_groups - join table between words and groups many-to-many
- groups - thematic groups of words
- study_sessions - records of study sessions grouping word_review_items
- study_activities - a specific instance of a study activity, linking a study_session to a group
- word_review_items - a record of word practive, determining if the word was correctly recalled or not

### Words
- id integer primary key autoincrement
- english text not null
- portuguese text not null
- parts text not null

### Words Groups
- id integer primary key autoincrement
- word_id integer not null   
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

#### GET /api/dashboard/last_study_sessions
- Description: Retrieves the most recent study sessions
- Query Parameters:
    - page: integer (default: 1)
    - limit: integer (default: 100)
- Returns: Array of recent study sessions with pagination

Response: 200 OK
```json
{
  "items": [
    {
      "id": 123,
      "group_id": 456,
      "created_at": "2024-03-20T15:30:00Z",
      "study_activity_id": 789,
      "group_name": "Basic Greetings"
    }
  ],
  "pagination": {
    "page_number": 1,
    "page_size": 100,
    "total_pages": 5,
    "total_items": 450
  }
}
```

#### GET /api/dashboard/study_progress
- Description: Returns the total count of studied words versus available words
- Returns: Basic progress statistics

Response: 200 OK
```json
{
  "data": {
    "total_words_studied": 150,
    "total_available_words": 200
  }
}
```

#### GET /api/dashboard/quick_stats
- Description: Returns key performance indicators and statistics
- Returns: Basic statistics about learning progress

Response: 200 OK
```json
{
  "data": {
    "success_rate": 85.5,
    "total_words_studied": 500,
    "total_words_correct": 427,
    "total_study_sessions": 25,
    "total_active_groups": 5,
    "current_streak": 3,
    "longest_streak": 7,
    "last_study_date": "2024-03-20T15:30:00Z"
  }
}
```

### Study Activities

#### GET /api/study_activities
- Description: Returns a paginated list of study activities
- Query Parameters:
    - page: integer (default: 1)
    - limit: integer (default: 100)
- Returns: List of study activities with pagination info

Response: 200 OK
```json
{
  "items": [
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
    "page_size": 100,
    "total_pages": 5,
    "total_items": 450
  }
}
```

#### GET /api/study_activities/:id
- Description: Returns detailed information about a specific study activity
- Returns: Study activity object with basic information

Response: 200 OK
```json
{
  "data": {
    "id": 1,
    "name": "Vocabulary Quiz",
    "thumbnail_url": "https://example.com/thumbnail.jpg",
    "description": "Practice your vocabulary with flashcards"
  }
}
```

#### GET /api/study_activities/:id/study_session
- Description: Returns a list of study sessions for a specific activity with pagination
- Query Parameters:
    - page: integer (default: 1)
    - limit: integer (default: 100)
- Returns: List of study sessions with pagination info

Response: 200 OK
```json
{
  "items": [
    {
      "id": 123,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "start_time": "2024-03-20T15:30:00Z",
      "end_time": "2024-03-20T15:45:00Z",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "page_number": 1,
    "page_size": 100,
    "total_pages": 5,
    "total_items": 450
  }
}
```

- POST /api/study_activities/
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
    - Description: Returns a paginated list of words
    - Query Parameters:
        - page: integer (default: 1)
        - limit: integer (default: 100)
    - Returns: List of words with pagination info

Response: 200 OK
```json
{
  "items": [
    {
      "id": 123,
      "portuguese": "correr",
      "english": "run",
      "study_statistics": {
        "correct_count": 5,
        "incorrect_count": 2
      }
    }
  ],
  "pagination": {
    "page_number": 1,
    "page_size": 100,
    "total_pages": 5,
    "total_items": 450
  }
}
```

- GET /api/words/:id
    - Description: Returns detailed information about a specific word
    - Returns: Word object with study statistics and group information

Response: 200 OK
```json
{
  "id": 123,
  "portuguese": "correr",
  "english": "run",
  "study_statistics": {
    "correct_count": 5,
    "incorrect_count": 2
  },
  "groups": [
    {
      "id": 1,
      "name": "Basic Verbs"
    }
  ]
}
```

### Words Groups

#### GET /api/words_groups
- Description: Returns a paginated list of word groups
- Query Parameters:
    - page: integer (default: 1)
    - limit: integer (default: 100)
- Returns: List of groups with pagination info

Response: 200 OK
```json
{
  "items": [
    {
      "id": 1,
      "name": "Basic Verbs",
      "description": "Common everyday verbs",
      "words_count": 50
    }
  ],
  "pagination": {
    "page_number": 1,
    "page_size": 100,
    "total_pages": 5,
    "total_items": 450
  }
}
```

#### GET /api/words_groups/:id
- Description: Returns detailed information about a specific word group
- Returns: Group object with basic information

Response: 200 OK
```json
{
  "data": {
    "id": 1,
    "name": "Basic Verbs",
    "description": "Common everyday verbs",
    "words_count": 50
  }
}
```

#### GET /api/words_groups/:id/words
- Description: Returns a paginated list of words in a specific group
- Query Parameters:
    - page: integer (default: 1)
    - limit: integer (default: 100)
- Returns: List of words with pagination info

Response: 200 OK
```json
{
  "items": [
    {
      "id": 1,
      "english": "run",
      "portuguese": "correr",
      "study_statistics": {
        "correct_count": 15,
        "incorrect_count": 3,
        "success_rate": 83.3
      }
    }
  ],
  "pagination": {
    "page_number": 1,
    "page_size": 100,
    "total_pages": 5,
    "total_items": 450
  }
}
```

#### GET /api/words_groups/:id/study_sessions
- Description: Returns a paginated list of study sessions for a specific group
- Query Parameters:
    - page: integer (default: 1)
    - limit: integer (default: 100)
- Returns: List of study sessions with pagination info

Response: 200 OK
```json
{
  "items": [
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
    "page_size": 100,
    "total_pages": 5,
    "total_items": 450
  }
}
```

### Study Sessions

#### GET /api/study_sessions
- Description: Returns a paginated list of study sessions
- Query Parameters:
    - page: integer (default: 1)
    - limit: integer (default: 100)
- Returns: List of study sessions with pagination info

Response: 200 OK
```json
{
  "items": [
    {
      "id": 123,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "start_time": "2024-03-20T15:30:00Z",
      "end_time": "2024-03-20T15:45:00Z",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "page_number": 1,
    "page_size": 100,
    "total_pages": 5,
    "total_items": 450
  }
}
```

#### GET /api/study_sessions/:id
- Description: Returns detailed information about a specific study session
- Returns: Study session object with basic information

Response: 200 OK
```json
{
  "data": {
    "id": 123,
    "activity_name": "Vocabulary Quiz",
    "group_name": "Basic Greetings",
    "start_time": "2024-03-20T15:30:00Z",
    "end_time": "2024-03-20T15:45:00Z",
    "review_items_count": 20
  }
}
```

#### GET /api/study_sessions/:id/words
- Description: Returns a paginated list of words reviewed in a specific study session
- Query Parameters:
    - page: integer (default: 1)
    - limit: integer (default: 100)
- Returns: List of reviewed words with pagination info

Response: 200 OK
```json
{
  "items": [
    {
      "id": 1,
      "english": "run",
      "portuguese": "correr",
      "correct_count": 5,
      "incorrect_count": 2
    }
  ],
  "pagination": {
    "page_number": 1,
    "page_size": 100,
    "total_pages": 5,
    "total_items": 450
  }
}
```

#### POST /api/study_sessions/:id/words/:word_id/review
- Description: Records a word review result in a study session
- URL Parameters:
    - id: integer (study_session_id)
    - word_id: integer
- Request Body:
```json
{
  "correct": true
}
```
- Returns: Review result information

Response: 200 OK
```json
{
  "success": true,
  "data": {
    "word_id": 1,
    "study_session_id": 123,
    "correct": true,
    "created_at": "2024-03-20T15:30:00Z"
  }
}
```

### Settings

#### GET /api/settings
- Description: Returns current application settings
- Returns: Settings object

Response: 200 OK
```json
{
  "data": {
    "theme": "dark",
    "words_per_session": 20,
    "show_hints": true,
    "default_study_time_minutes": 15,
    "updated_at": "2024-03-20T15:30:00Z"
  }
}
```

#### POST /api/settings
- Description: Updates application settings
- Required params:
    - theme: string (light|dark|system)
- Request Body:
```json
{
  "theme": "dark",
  "words_per_session": 20,
  "show_hints": true,
  "default_study_time_minutes": 15
}
```
- Returns: Updated settings object

Response: 200 OK
```json
{
  "success": true,
  "data": {
    "theme": "dark",
    "words_per_session": 20,
    "show_hints": true,
    "default_study_time_minutes": 15,
    "updated_at": "2024-03-20T15:30:00Z"
  }
}
```

#### POST /api/reset_history
- Description: Resets all study history while maintaining words and groups
- Returns: Success confirmation

Response: 200 OK
```json
{
  "success": true,
  "data": {
    "message": "Study history has been reset"
  }
}
```

#### POST /api/full_reset
- Description: Performs a complete system reset, removing all data
- Returns: Success confirmation

Response: 200 OK
```json
{
  "success": true,
  "data": {
    "message": "System has been fully reset"
  }
}
```

## Task Runner Tasks

### Initialize Database
- Description: Creates the SQLite database file `words.db` in the project root
- Task Name: `initdb`
- Actions:
  - Creates database file if it doesn't exist
  - Verifies database connection

### Migrate Database
- Description: Runs SQL migration files in sequential order
- Task Name: `migrate`
- Migration Files Location: `./migrations/`
- File Naming Convention: 
  - Format: `NNNN_description.sql`
  - Example: 
    - `0001_init.sql`
    - `0002_create_words_table.sql`
- Actions:
  - Runs migrations in order by file name
  - Tracks executed migrations
  - Prevents duplicate migration runs

### Seed Data
- Description: Imports initial data from JSON files
- Task Name: `seed`
- Seed Files Location: `./seeds/`
- File Format: JSON
- Example Seed File:
```json
[
  {
    "portuguese": "pagar",
    "english": "to pay",
    "parts": "verb"
  }
]
```
- Actions:
  - Reads JSON files
  - Maps data to database schema
  - Inserts data into appropriate tables
  - Creates word-group associations