import "bootstrap/dist/css/bootstrap.min.css";
import "./App.css";
import { useEffect, useState } from "react";
import axios, { AxiosError, AxiosResponse } from "axios";
import { useNavigate } from "react-router-dom";
import Accounts from "./Accounts";

type User = {
  id?: string;
  email?: string;
};

function App() {
  const [user, setUser] = useState<User>({});
  const [fetching, setFetching] = useState<boolean>(true);

  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("token");

    axios
      .get(`${import.meta.env.VITE_API_BASE_URL}/api/v1/users/iam/`, {
        headers: {
          Authorization: token,
        },
      })
      .then((res: AxiosResponse) => {
        console.log(res.status);

        setUser(res.data);

        navigate("/");
      })
      .catch((err: AxiosError) => {
        if (err.response?.status === 401) {
          navigate("/login");
        }
      })
      .finally(() => {
        setFetching(false);
      });
  }, [navigate]);

  if (!fetching && user.id) {
    return (
      <div>
        <div className="m-2">Hello {user.email}</div>
        <div className="m-2">
          <Accounts />
        </div>
      </div>
    );
  }
}

export default App;
