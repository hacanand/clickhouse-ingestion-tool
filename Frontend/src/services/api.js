import axios from "axios";
const baseURL = "http://localhost:8080/api";

export const uploadCSV = (formData) =>
  axios.post(`${baseURL}/file-to-clickhouse`, formData);

export const exportToCSV = (formData) =>
  axios.post(`${baseURL}/clickhouse-to-file`, formData);

export const getColumns = (formData) =>
  axios.post(`${baseURL}/get-columns`, formData);
