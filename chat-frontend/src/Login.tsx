import React, { useState } from 'react';
import { Form, Input, Button, Typography, message } from 'antd';
import { useNavigate, Link } from 'react-router-dom';
import axios from 'axios';
import './Login.css';

const { Title } = Typography;

const Login: React.FC = () => {
  // React hooks to manage state
  const navigate = useNavigate();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  // Function to handle the login process
  const handleLogin = async () => {
    try {
      // Send a POST request to the '/api/login' endpoint with the entered username and password
      const response = await axios.post('http://localhost:8000/api/login', { username, password });

      // Extract the token from the response data
      const token = response.config.data;
      console.log(token);

      // Navigate to the '/chat' route with the token passed as state
      navigate('/chat', { state: { token } });
    } catch (error) {
      // Handle errors during login
      if (error.response && error.response.status === 401) {
        console.error('Error during login:', error.response.data);
        message.error(error.response.data);
      } else {
        console.error('Error during login:', error);
      }
    }
  };

  return (
    <div className="login-container">
      <Title level={2} className="login-heading">
        Login
      </Title>
      <Form className="login-form">
        <Form.Item
          label="Username"
          name="username"
          rules={[{ required: true, message: 'Please enter your username!' }]}
        >
          <Input value={username} onChange={(e) => setUsername(e.target.value)} />
        </Form.Item>

        <Form.Item
          label="Password"
          name="password"
          rules={[{ required: true, message: 'Please enter your password!' }]}
        >
          <Input.Password value={password} onChange={(e) => setPassword(e.target.value)} />
        </Form.Item>

        <Form.Item>
          <Button type="primary" className="login-button" onClick={handleLogin}>
            Login
          </Button>
        </Form.Item>
      </Form>
      <p className="signup-link">
        Don't have an account? <Link to="/api/signup">Sign up</Link>
      </p>
    </div>
  );
};

export default Login;
