import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

import SignUp from './SignUp.tsx';
import App from './App.tsx';
import Login from './Login.tsx';
import Chat from './Chat.tsx';

ReactDOM.render(
  <Router>
    <Routes> {/* Wrap routes with the <Routes> component */}
      <Route path="/api/signup" element={<SignUp />} /> 
      <Route path="/" element={<App />} /> 
      <Route path="/api/login" element ={<Login />} />
      <Route path="/chat" element={<Chat />} />
    </Routes>
  </Router>,
  document.getElementById('root')
);

