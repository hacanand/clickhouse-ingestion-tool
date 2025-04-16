import { useState } from "react";

export default function useToken() {
  const [token, setTokenState] = useState(
    localStorage.getItem("jwt_token") || ""
  );

  const setToken = (val) => {
    setTokenState(val);
    localStorage.setItem("jwt_token", val);
  };

  return [token, setToken];
}
