import axios, { AxiosError, AxiosResponse } from "axios";
import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

const Account: React.FC = () => {
  const { accountId } = useParams<{ accountId: string }>();

  const [account, setAccount] = useState<{ extId?: string; balance?: string }>(
    {}
  );

  useEffect(() => {
    const token = localStorage.getItem("token");

    axios
      .get(
        `${
          import.meta.env.VITE_API_BASE_URL
        }/api/v1/accounts/${accountId}/balance`,
        {
          headers: {
            Authorization: token,
          },
        }
      )
      .then((res: AxiosResponse) => {
        if (res.data.ext_id && res.data.balance) {
          setAccount({ extId: res.data.ext_id, balance: res.data.balance });
        }
      })
      .catch((err: AxiosError) => {
        console.log(err);
      });
  }, [accountId]);

  return (
    <div className="m-2">
      <div className="m-2">Address: {account.extId}</div>
      <div className="m-2">Balance: {account.balance}</div>
    </div>
  );
};

export default Account;
