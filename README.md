# Full stack - Chat Application

This is a simple chat application consisting of a backend and frontend. The backend is built using Go, and the frontend is built using React Typescript. Users can sign up, log in, and participate in real-time chat conversations.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Frontend-Setup](#Frontend-Setup)
- [Backend-Setup](#Backend-Setup)
- [Deployment](#deployment)

## Features
- User Registration: Users can sign up for a new account by providing a username and password.
- User Authentication: Users can log in with their registered credentials to access the chat application.
- Real-time Chat: Users can participate in real-time chat conversations with other users who are currently logged in.
- Message History: Users can view the chat history and see previous messages exchanged in the chat room.
- Upvoting and Downvoting: Users can upvote or downvote messages posted by other users, allowing for community-driven moderation of content.
- User Presence: Users can see the list of online users who are currently active in the chat room.
- Logout: Users can log out of the application to end their session.

## Prerequisites

Before running the application, ensure that the following software is installed on your machine:

- Go (1.16 or later)
- Node.js (v14 or later)
- npm (v7 or later)
- Docker (for containerization)
- Minikube (for Kubernetes deployment)

## Frontend-Setup

Navigate to the `chat-frontend` directory and start the frontend development server:
```
npm start
```
The frontend application should now be running on http://localhost:4000.

## Backend-Setup

Install the Go dependencies and drivers and start the backend code:
```
go run main.go
```

## Deployment

Build the Docker images and deploy kubernetes. These commands can be used to deploy kubernetes:
```
kubectl apply -f backend-deployment.yaml
kubectl apply -f frontend-deployment.yaml
```





