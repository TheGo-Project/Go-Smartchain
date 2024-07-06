import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import Login from "./Login.tsx";
import Register from "./Register.tsx";
import { Provider } from "react-redux";
import { store } from "./store/store.ts";
import Account from "./Account.tsx";
import Admin from "./Admin.tsx";
import Explorer from "./Explorer.tsx";
import Block from "./Block.tsx";

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
  },
  {
    path: "/login",
    element: <Login />,
  },
  {
    path: "/registration",
    element: <Register />,
  },
  {
    path: "/accounts/:accountId",
    element: <Account />,
  },
  {
    path: "/admin",
    element: <Admin />,
  },
  {
    path: "/explorer",
    element: <Explorer />,
  },
  {
    path: "/explorer/:number",
    element: <Block />,
  },
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <Provider store={store}>
      <RouterProvider router={router} />
    </Provider>
  </React.StrictMode>
);
