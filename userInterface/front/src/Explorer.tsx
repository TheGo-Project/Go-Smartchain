import axios, { AxiosError, AxiosResponse } from "axios";
import React, { useEffect, useState } from "react";
import NavMenu from "./NavMenu";
import { Link } from "react-router-dom";

type Block = {
  number: number;
  hash: string;
  parentHash: string;
  time: string;
};

const Explorer: React.FC = () => {
  const [blocks, setBlocks] = useState<Block[]>([]);

  useEffect(() => {
    const token = localStorage.getItem("token");

    axios
      .get(`${import.meta.env.VITE_API_BASE_URL}/api/v1/blocks`, {
        headers: {
          Authorization: token,
        },
      })
      .then((res: AxiosResponse) => {
        // console.log(res.data);

        setBlocks(res.data.blocks);
      })
      .catch((err: AxiosError) => {
        console.log(err);
      });
  }, []);

  return (
    <>
      <NavMenu />

      <div className="m-2">
        {blocks.map((block) => (
          <div key={block.number} className="flex flex-row">
            <div className="mx-2">
              <Link to={`/explorer/${block.number}`}>{block.number}</Link>
            </div>
            <div className="mx-2">{new Date(block.time).toUTCString()}</div>
            <div className="mx-2">{block.hash}</div>
          </div>
        ))}
      </div>
    </>
  );
};

export default Explorer;
