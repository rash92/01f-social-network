
# Social Network

A **Facebook-inspired platform** that allows users to connect, share posts, and interact in real time. Developed as part of the 01 Founders curriculum, this project was a collaborative effort by a team of four students, showcasing features such as user authentication, group management, and real-time communication.

## Features

### Core Functionality
- **User Authentication**:  
  Secure registration, login, and logout using cookies for session management.
- **Profiles**:  
  Users can create and manage public or private profiles, complete with personal information and activity.
- **Posts**:  
  Users can create, edit, and manage posts with customizable privacy settings (public, private, or restricted to followers).
- **Followers**:  
  Users can follow or unfollow others, with requests required for private profiles.

### Real-Time Features
- **Group Chat**:  
  Members of a group can communicate in real-time through a shared chat room.
- **Notifications**:  
  Users receive real-time updates for events like group invitations, follow requests, and more.

### Group Management
- Create groups with titles and descriptions.
- Invite others to join groups or request membership.
- Organize events within groups, including RSVP options.

## Technologies Used

### Frontend
- **React**: For building the user interface and handling user interactions.
- **CSS**: For styling and responsive design.

### Backend
- **golang**: For building the server-side logic and APIs.
- **WebSockets**: For real-time features like group chat and notifications.

### Database
- **SQLite**: To store user data, posts, and group information.

### Tools
- **Docker**: For containerizing the frontend and backend for deployment.
- **Git**: For version control and team collaboration.



## Installation

### docker

from the root directory run ```docker compose up```

### If you do not want to use docker, or can't get docker to work

cd into the backend directory, and run ```go run server.go```

in a separate terminal, cd into the frontend directory and run ```npm install``` and then ```npm start```


### run in browsers

Open the application in your browser at `http://localhost:3000`. Separate instances can be run in separate browsers, and different users may log in and interact with each other through chats, reacting to posts and comments, joining groups etc.

## Future Improvements
- Add automated tests for both frontend and backend.
- Enhance database migrations and seeding.
- Add support for more advanced privacy controls.

## Contributors
- **Rashid**
- **Daisy**
- **Bilaal** 
- **Peter**
