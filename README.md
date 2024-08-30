# Beego_Project

This project implements a web application similar to The Cat API using Beego for the backend and JavaScript for frontend interactions.
# Features
- Display random cat images
- Filter cat images by breed
- Vote on favorite cat images
- Upload cat images

# Technologies Used

- Backend: Beego (Go web framework)
- Frontend: Vanilla JS
- API: The Cat API

# Prerequisites
- Go (version 1.16 or later)
- Beego framework
- Vanilla JavaScript

# Setup and Installation

1. Clone the repository:
   ```
   git clone https://github.com/YasinRafin01/Beego_Project.git
   cd Beego_Project
   ```
2. Install Beego:
   ```
   go get github.com/beego/beego/v2@latest
   go mod tidy
   ```
3. Set up your Cat API key in conf/app.conf:
   ```
   cat_api_key =your_api_key
   ```
4. Run the application:
   ```
   bee run
   ```
5. Open your browser and navigate to http://localhost:8080


# Project Structure

cat-api-project/

├── conf/

│   └── app.conf

├── controllers/

│   └── default.go

├── models/

│   └── cat.go

├── routers/

│   └── router.go

├── static/

│   ├── css/

│   └── js/

├── views/

│   └── index.tpl

├── main.go

└── README.md

# Key Implementation Details

1. Beego Controller: Handles requests and renders templates.
2. JavaScript Interaction: Manages frontend interactions and dynamic content updates.
3. Go Channels: Used for concurrent API calls to improve performance.
4. Beego Config: Stores and manages the API key and other configuration settings.

# API Integration
This project uses The Cat API for fetching cat images and related data. API calls are made using Go's HTTP client and processed concurrently using channels.
# Contributing
Contributions are welcome! Please feel free to submit a Pull Request.
   
