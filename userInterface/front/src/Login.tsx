import axios, { AxiosResponse } from "axios";
import React, { useState, FormEvent } from "react";
import { Button } from "react-bootstrap";
import { Link, useNavigate } from "react-router-dom";
import { useAppDispatch } from "./store/hooks";
import { set } from "./store/userReducer";

const Login: React.FC = () => {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");

  const navigate = useNavigate();

  const dispatch = useAppDispatch();

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault();

    axios
      .post(`${import.meta.env.VITE_API_BASE_URL}/api/v1/users/login`, {
        email,
        password,
      })
      .then((response: AxiosResponse) => {
        dispatch(set(response.data.token));

        localStorage.setItem('token', response.data.token);

        navigate("/");
      });
  };

  return (
    <div className="container mt-5 max-w-screen-md">
      <form onSubmit={handleSubmit}>
        <div className="mb-3">
          <label htmlFor="email" className="form-label">
            Email address
          </label>
          <input
            type="email"
            className="form-control"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div className="mb-3">
          <label htmlFor="password" className="form-label">
            Password
          </label>
          <input
            type="password"
            className="form-control"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <Button type="submit" variant="primary">
          Login
        </Button>
      </form>
      <div className="mt-3">
        <span>Don't have an account? </span>
        <Link to="/registration">
          <Button variant="success">Register here</Button>
        </Link>
      </div>
    </div>
  );
};

export default Login;
