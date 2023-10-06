import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { Button, message } from 'antd';
import { UpOutlined, DownOutlined } from '@ant-design/icons';
import './Chat.css';

const Chat = () => {
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const navigate = useNavigate();

  // Extract the token from the current location state
  const location = useLocation();
  const token = location.state && location.state.token;
  const currToken = JSON.parse(token);
  const strUsername = currToken.username;

  // Function to fetch the history messages from the server
  const fetchHistoryMessages = async () => {
    try {
      const response = await axios.get('http://localhost:8000/api/messages/history', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      setMessages(response.data);
    } catch (error) {
      console.error('Error fetching history messages:', error);
    }
  };

  const fetchMessages = () => {
    fetchHistoryMessages();
  };

  useEffect(() => {
    if (!token) {
      // Redirect the user to the login page if the token is not available
      navigate('/api/login', { replace: true });
    } else {
      fetchHistoryMessages();
      const interval = setInterval(fetchMessages, 1000); // Fetch messages every 5 seconds
      return () => clearInterval(interval); // Clean up the interval when the component unmounts
    }
  }, [token, navigate]);

  // Function to handle sending a new message
  const handleSendMessage = async () => {
    try {
      const message = {
        content: newMessage,
        username: strUsername,
      };

      await axios.post('http://localhost:8000/api/messages', message, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      setNewMessage('');
      fetchHistoryMessages();
    } catch (error) {
      console.error('Error sending message:', error);
    }
  };

  // Function to handle upvoting a message
  const handleUpvote = async (messageId:any) => {
    console.log(messageId);
    try {
      await axios.post(
        `http://localhost:8000/api/messages/${messageId}/upvote`,
        {},
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      fetchHistoryMessages();
    } catch (error) {
      console.error('Error upvoting message:', error);
      if (error.response && error.response.data === 'User has already voted for this message') {
          message.warning('You have already voted for this message');
    }
  }
  };

  // Function to handle downvoting a message
  const handleDownvote = async (messageId:any) => {
    try {
      await axios.post(
        `http://localhost:8000/api/messages/${messageId}/downvote`,
        {},
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      fetchHistoryMessages();
    } catch (error) {
      console.error('Error downvoting message:', error);
      if (error.response && error.response.data === 'User has already voted for this message') {
          message.warning('You have already voted for this message');
    }
  }
  };

  // Function to handle the logout process
  const handleLogout = async () => {
    try {
      await axios.post('http://localhost:8000/api/logout', null, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      navigate('/api/login', { replace: true });
    } catch (error) {
      console.error('Error logging out:', error);
    }
  };

  return (
    <div className="chat-container">
      <h2>Welcome to the Chatroom!</h2>
      <div className="message-container">
        {messages &&
          messages.map((message:any) => (
            <div key={message._id} className="message">
              <p>{message.content}</p>
              <div className="message-stats">
                <Button
                  className="vote-button"
                  shape="circle"
                  icon={<UpOutlined />}
                  onClick={() => handleUpvote(message.ID)}
                />
                <span className="vote-count">{message.upvotes}</span>
                <Button
                  className="vote-button"
                  shape="circle"
                  icon={<DownOutlined />}
                  onClick={() => handleDownvote(message.ID)}
                />
                <span className="vote-count">{message.downvotes}</span>
              </div>
            </div>
          ))}
      </div>
      <div className="message-input">
        <input
          type="text"
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          placeholder="Type your message..."
        />
        <Button type="primary" onClick={handleSendMessage}>
          Send
        </Button>
      </div>
      <Button className="logout-button" onClick={handleLogout}>
        Logout
      </Button>
    </div>
  );
};

export default Chat;
