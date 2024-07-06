import axios, { AxiosError, AxiosResponse } from "axios";
import React, { useEffect, useState } from "react";
import { Button, Form } from "react-bootstrap";
import NavMenu from "./NavMenu";

const Admin: React.FC = () => {
  const [faucetContractAddress, setFaucetContractAddress] =
    useState<string>("");
  const [faucetContractAddressRO, setFaucetContractAddressRO] =
    useState<boolean>(true);

  //   const [coinbaseAddress, setCoinbaseAddress] = useState<string>("");
  //   const [coinbaseAddressRO, setCoinbaseAddressRO] = useState<boolean>(true);

  useEffect(() => {
    const token = localStorage.getItem("token");

    axios
      .get(
        `${
          import.meta.env.VITE_API_BASE_URL
        }/api/v1/params/faucet-contract-address`,
        {
          headers: {
            Authorization: token,
          },
        }
      )
      .then((res: AxiosResponse) => {
        setFaucetContractAddress(res.data.value);
      })
      .catch((err: AxiosError) => {
        console.log(err);
      });

    // axios
    //   .get(
    //     `${import.meta.env.VITE_API_BASE_URL}/api/v1/params/coinbase-address`,
    //     {
    //       headers: {
    //         Authorization: token,
    //       },
    //     }
    //   )
    //   .then((res: AxiosResponse) => {
    //     setCoinbaseAddress(res.data.value);
    //   })
    //   .catch((err: AxiosError) => {
    //     console.log(err);
    //   });
  }, []);

  const saveParam = (key: string, value: string): Promise<string> => {
    return new Promise((resolve, reject) => {
      const token = localStorage.getItem("token");

      let url = "";
      if (key === "faucetContractAddress") {
        url = "/api/v1/params/faucet-contract-address";
      } else if (key === "coinbaseAddress") {
        url = "/api/v1/params/coinbase-address";
      }

      axios
        .post(
          `${import.meta.env.VITE_API_BASE_URL}${url}`,
          {
            value: value,
          },
          {
            headers: {
              Authorization: token,
            },
          }
        )
        .then((res: AxiosResponse) => {
          console.log(res.status);

          resolve(value);
        })
        .catch((err: AxiosError) => {
          console.log(err);

          reject(err);
        });
    });
  };

  const saveFaucetContractAddress = () => {
    saveParam("faucetContractAddress", faucetContractAddress).then(() => {
      setFaucetContractAddressRO(true);
    });
  };

  //   const saveCoinbaseAddress = () => {
  //     saveParam("coinbaseAddress", coinbaseAddress).then(() => {
  //       setCoinbaseAddressRO(true);
  //     });
  //   };

  const onClickDelete = () => {
    const token = localStorage.getItem("token");

    axios.delete(`${import.meta.env.VITE_API_BASE_URL}/api/v1/accounts/all`, {
      headers: {
        Authorization: token,
      },
    });
  };

  return (
    <>
      <NavMenu />
      <div className="m-2">
        <Form>
          <Form.Group className="mb-3" controlId="exampleForm.ControlInput1">
            <Form.Label>Faucet contract address</Form.Label>
            <Form.Control
              type="text"
              value={faucetContractAddress}
              readOnly={faucetContractAddressRO}
              onChange={(e) => setFaucetContractAddress(e.target.value)}
            />
            <div className="flex flex-row">
              <div className="mx-1 my-2">
                <Button
                  variant="primary"
                  onClick={() => setFaucetContractAddressRO(false)}
                  disabled={!faucetContractAddressRO}
                >
                  Edit
                </Button>
              </div>
              {!faucetContractAddressRO && (
                <div className="mx-1 my-2">
                  <Button
                    variant="success"
                    onClick={() => saveFaucetContractAddress()}
                  >
                    Save
                  </Button>
                </div>
              )}
            </div>
          </Form.Group>

          {/* <Form.Group className="mb-3" controlId="exampleForm.ControlInput1">
      <Form.Label>Coinbase address</Form.Label>
      <Form.Control
        type="text"
        value={coinbaseAddress}
        readOnly={coinbaseAddressRO}
        onChange={(e) => setCoinbaseAddress(e.target.value)}
      />
      <div className="flex flex-row">
        <div className="mx-1 my-2">
          <Button
            variant="primary"
            onClick={() => setCoinbaseAddressRO(false)}
            disabled={!coinbaseAddressRO}
          >
            Edit
          </Button>
        </div>
        {!coinbaseAddressRO && (
          <div className="mx-1 my-2">
            <Button variant="success" onClick={() => saveCoinbaseAddress()}>
              Save
            </Button>
          </div>
        )}
      </div>
    </Form.Group> */}

          <div>
            <Button variant="danger" onClick={() => onClickDelete()}>
              Delete all accounts
            </Button>
          </div>
        </Form>
      </div>
    </>
  );
};

export default Admin;
