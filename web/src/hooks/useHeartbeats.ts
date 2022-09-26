import axios from "axios";
import { useEffect, useState } from "react";
import { useNetworkContext } from "../contexts/NetworkContext";
import { NumberLong } from "../utils/longToDate";

export type HeartbeatResponse = {
  boottimestamp: NumberLong;
  counter: number;
  createdAt: string;
  features: string[] | null;
  guardianaddr: string;
  networks: {
    contractaddress: string;
    errorcount: number;
    height: number;
    id: number;
  }[];
  nodename: string;
  timestamp: NumberLong;
  updatedAt: string;
  version: string;
  _id: string;
};

function useHeartbeats(): HeartbeatResponse[] {
  const { currentNetwork } = useNetworkContext();
  const [heartbeats, setHeartbeats] = useState<HeartbeatResponse[]>([]);
  useEffect(() => {
    setHeartbeats([]);
  }, [currentNetwork]);
  useEffect(() => {
    let cancelled = false;
    (async () => {
      while (!cancelled) {
        const response = await axios.get<HeartbeatResponse[]>(
          "http://localhost:3000/api/heartbeats"
        );
        if (!cancelled) {
          setHeartbeats(
            response.data.sort(
              (a, b) => a.nodename.localeCompare(b.nodename || "") || -1
            )
          );
          await new Promise((resolve) => setTimeout(resolve, 1000));
        }
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [currentNetwork]);
  return heartbeats;
}
export default useHeartbeats;