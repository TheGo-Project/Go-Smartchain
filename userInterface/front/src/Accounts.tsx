import axios, { AxiosResponse } from "axios";
import React, { useEffect, useState } from "react";
import {
  Button,
  FormControl,
  FormGroup,
  FormLabel,
  Modal,
} from "react-bootstrap";
import { Form, Link } from "react-router-dom";

type Account = {
  id: string;
  user_id: string;
  ext_id: string;
};

const Accounts: React.FC = () => {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [newAccount, setNewAccount] = useState<{
    address?: string;
    password?: string;
  }>({});
  const [newAccountModalShow, setNewAccountModalShow] = useState(false);

  const [faucetModalShow, setFaucetModalShow] = useState(false);
  const [faucetRequestAddress, setFaucetRequestAddress] = useState("");
  const [faucetRequestPassword, setFaucetRequestPassword] = useState("");

  useEffect(() => {
    const token = localStorage.getItem("token");

    axios
      .get(`${import.meta.env.VITE_API_BASE_URL}/api/v1/accounts`, {
        headers: {
          Authorization: token,
        },
      })
      .then((res: AxiosResponse) => {
        console.log(res.status);

        if (Array.isArray(res.data)) {
          setAccounts(res.data);
        } else {
          console.error("Unexpected data format:", res.data);
        }
      });
  }, []);

  const createAccountOnClick = () => {
    const token = localStorage.getItem("token");

    axios
      .post(
        `${import.meta.env.VITE_API_BASE_URL}/api/v1/accounts`,
        {},
        {
          headers: {
            Authorization: token,
          },
        }
      )
      .then((res: AxiosResponse) => {
        if (res.data.address && res.data.password) {
          setNewAccount({
            address: res.data.address,
            password: res.data.password,
          });
          setNewAccountModalShow(true);
        }
      })
      .catch(console.error);
  };

  const handleClose = () => {
    setNewAccountModalShow(false);
    setNewAccount({});
  };

  const faucetOnClick = () => {
    setFaucetModalShow(true);
  };

  const handleFaucetClose = () => {
    setFaucetModalShow(false);
  };

  const faucetSubmit = () => {
    console.log(faucetRequestAddress, faucetRequestPassword);

    const token = localStorage.getItem("token");

    axios
      .post(
        `${import.meta.env.VITE_API_BASE_URL}/api/v1/faucet`,
        { address: faucetRequestAddress, password: faucetRequestPassword },
        {
          headers: {
            Authorization: token,
          },
        }
      )
      .then((res: AxiosResponse) => {
        console.log(res.status);
        console.log(res.data);
      })
      .catch(console.error);
  };

  return (
    <div>
      <div className="flex flex-row">
        <div className="m-2">
          <Button onClick={createAccountOnClick}>Create Account</Button>
        </div>
        <div className="m-2">
          <Button onClick={faucetOnClick}>Faucet</Button>
        </div>
      </div>
      <div>Accounts</div>
      <ul>
        {accounts.map((account) => (
          <div key={account.id}>
            <Link to={`/accounts/${account.id}`}>{account.ext_id}</Link>
          </div>
        ))}
      </ul>

      <Modal show={newAccountModalShow} onHide={handleClose}>
        <Modal.Header closeButton>
          <Modal.Title>Save this information</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <div className="font-bold">Address</div>
          <div>{newAccount.address}</div>
          <div className="font-bold">Password</div>
          <div>{newAccount.password}</div>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="danger" onClick={handleClose}>
            Close
          </Button>
        </Modal.Footer>
      </Modal>

      <Modal show={faucetModalShow} onHide={handleFaucetClose}>
        <Modal.Header closeButton>
          <Modal.Title>Faucet</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form>
            <FormGroup className="mb-3" controlId="formBasicAddress">
              <FormLabel>Account address</FormLabel>
              <FormControl
                placeholder="Enter account address"
                value={faucetRequestAddress}
                onChange={(e) => setFaucetRequestAddress(e.target.value)}
              />
            </FormGroup>

            <FormGroup className="mb-3" controlId="formBasicPassword">
              <FormLabel>Password</FormLabel>
              <FormControl
                type="password"
                placeholder="Password"
                value={faucetRequestPassword}
                onChange={(e) => setFaucetRequestPassword(e.target.value)}
              />
            </FormGroup>
            <Button variant="primary" type="submit" onClick={faucetSubmit}>
              Submit
            </Button>
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="danger" onClick={handleFaucetClose}>
            Close
          </Button>
        </Modal.Footer>
      </Modal>
    </div>
  );
};

export default Accounts;
