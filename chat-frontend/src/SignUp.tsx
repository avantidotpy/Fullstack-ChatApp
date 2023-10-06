import React, { useState } from "react";
import { Form, Input, Button, Typography, message } from 'antd';
import { Link } from 'react-router-dom';
import axios from "axios";
import { useNavigate } from 'react-router-dom';
import './Login.css'; 

const { Title } = Typography;

const SignUp: React.FC = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const navigate = useNavigate();


  // Function to handle the signup process
  const handleSignUp = async (values: any) => {
    try {
      // Send a POST request to the '/api/signup' endpoint with the entered username and password
      const response = await axios.post("http://localhost:8000/api/signup", { username, password, confirmPassword });
      console.log("User registered successfully:", response.data);
      
      // Route to another page, after successful sign-up.
      message.info("Redirecting to login page...");
      setTimeout(() => {
        navigate("/api/login");
      }, 2000);
    } catch (error) {
      if (error.response && error.response.status === 400) {
        console.error("Error during sign up:", error.response.data);
        // Display an error message to the user
        message.error(error.response.data);
      } else {
        console.error("Error during sign up:", error);
      }
    }
  };

  return (
    <div className="login-container">
      <Title level={2} className="login-heading">
        Sign Up
      </Title>
      <Form onFinish={handleSignUp} className="login-form">
        <Form.Item
          label="Username"
          name="username"
          rules={[{ required: true, message: "Please enter your username!" }]}
        >
          <Input value={username} onChange={(e) => setUsername(e.target.value)} />
        </Form.Item>

        <Form.Item
          label="Password"
          name="password"
          rules={[{ required: true, message: "Please enter your password!" }]}
        >
          <Input.Password
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </Form.Item>

        <Form.Item
          label="Confirm Password"
          name="confirmPassword"
          rules={[
            { required: true, message: "Please confirm your password!" },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue("password") === value) {
                  return Promise.resolve();
                }
                return Promise.reject("Passwords do not match!");
              },
            }),
          ]}
        >
          <Input.Password
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
          />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" className="login-button">
            Sign Up
          </Button>
        </Form.Item>
      </Form>

      <p className="signup-link">
        Already have an account? <Link to="/api/login">Log in</Link>
      </p>
    </div>
  );
};

export default SignUp;
