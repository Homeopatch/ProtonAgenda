# TODO:
Develop API
Split API generated code into several files

# Vue.js + Go Starter Project

This is a simple starter template for a project with a Go backend and a Vue.js frontend.

## Prerequisites

- **Go**: Install from [golang.org](https://golang.org/).
- **Node.js & npm**: Install from [nodejs.org](https://nodejs.org/).

## Setup Instructions

### Backend (Go)
1. Navigate to the `backend` directory:
   ```bash
   cd backend

2. Run the Go server:

   ```
   go run main.go
   ```

3. The backend will run on http://localhost:8080.

### Frontend (Vue.js)

1. Navigate to the frontend directory:

   ```
   cd frontend
   ```

2. Install dependencies:

   ```
   npm install
   ```

3. Start the development server:
   ```
   npm run dev
   ```
    
4. Access the app at http://localhost:5173.

Notes

The frontend development server proxies API requests to the backend.
To build the frontend for production:
```
npm run build
```

