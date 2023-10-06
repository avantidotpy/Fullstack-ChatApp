import React from 'react';
import { Link } from 'react-router-dom';
import SignUp from './SignUp.tsx';
import Login from './Login.tsx';
import './App.css';
//import Chat from './Chat.tsx';

function App() {
  return (
    <div className="app-container">
      <h1 className="app-heading">Welcome to ChatApp</h1>
      <p className="app-subheading">Thank you for visiting! Get started by signing up or logging in.</p>
      <div className="app-links">
        <Link to="/api/signup" className="app-link">Sign Up</Link>
        <Link to="/api/login" className="app-link">Log In</Link>
      </div>
    </div>
  );
}

export default App;
