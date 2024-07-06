import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import NavMenu from "./NavMenu";
import axios, { AxiosError, AxiosResponse } from "axios";

type Block = {
  number?: string;
  hash?: string;
  parentHash?: string;
  time?: string;
  transactionsCount?: number;
};

const Block: React.FC = () => {
  const { number } = useParams<{ number: string }>();

  const [block, setBlock] = useState<Block>({});

  useEffect(() => {
    const token = localStorage.getItem("token");

    axios
      .get(`${import.meta.env.VITE_API_BASE_URL}/api/v1/blocks/${number}`, {
        headers: {
          Authorization: token,
        },
      })
      .then((res: AxiosResponse) => {
        console.log(res.data);

        setBlock(res.data);
      })
      .catch((err: AxiosError) => {
        console.log(err);
      });
  }, [number]);

  return (
    <>
      <NavMenu />

      <div className="m-2">
        <div>Number: {block.number}</div>
        <div>Time: {new Date(block.time!).toUTCString()}</div>
        <div>Hash: {block.hash}</div>
        <div>Parent hash: {block.parentHash}</div>
        <div>Transactions: {block.transactionsCount}</div>
      </div>
    </>
  );
};

export default Block;
