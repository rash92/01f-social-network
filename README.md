
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

1. Clone the repository:
   ```bash
   git clone https://github.com/bilaal441/social-network.git
   cd social-network
   ```

2. Set up the backend:
   ```bash
   cd backend

   go run .
   ```

3. Set up the frontend:
   ```bash
   cd frontend
   npm install
   npm start
   ```

4. Open the application in your browser at `http://localhost:3000`.

## Future Improvements
- Add automated tests for both frontend and backend.
- Enhance database migrations and seeding.
- Add support for more advanced privacy controls.

## Contributors
- **Bilal** 
- **Peter** 
- **daisy** 
- **rashid** 
