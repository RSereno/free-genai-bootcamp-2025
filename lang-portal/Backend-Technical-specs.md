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
- GET /api/dashboard/study-progress
    - Returns:
        - total_known: integer
        - total_unknown: integer
        - progress_percentage: float
- GET /api/dashboard/quick-stats
    - Returns:
        - success_rate: float
        - total_words_studied: integer
        - total_words_correct: integer
        - total_words_incorrect: integer
        - total_study_sessions: integer
        - total_active_groups: integer
        - current_streak: integer
        - longest_streak: integer

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

- GET /api/study-activities/:id 
    - Returns:
        - study_activity: object
            - id: integer
            - name: string
            - group_id: integer
            - created_at: datetime
            - word_reviews: array of word_review_items

- GET /api/study-activities/:id/study-session
    - Returns:
        - session_data: object
            - words: array of words to study
            - settings: object with activity settings

- POST /api/study-activities/
    - Required params:
        - group_id: integer
        - study_activity_id: integer
    - Returns:
        - study_activity: newly created study activity object

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

- GET /api/words-groups/:id
    - Returns:
        - group: object
            - id: integer
            - name: string
            - description: string
            - words_count: integer
            - words: array of word objects

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

### Settings
- POST /api/settings
    - Required params:
        - theme: string (light|dark|system)
    - Returns:
        - settings: object
            - theme: string
            - updated_at: datetime

- POST /api/settings/reset-progress
    - Returns:
        - success: boolean
        - message: string

- POST /api/settings/reload-seed-data
    - Returns:
        - success: boolean
        - message: string

## Doc

