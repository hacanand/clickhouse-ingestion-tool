import { useState } from "react";
// import useToken from "./hooks/useToken";
import TokenInput from "./components/TokenInput";
import FileUpload from "./components/FileUpload";
import CSVPreview from "./components/CSVPreview";
import ColumnSelector from "./components/ColumnSelector";
import { exportToCSV } from "./services/api";
// import apiClient from "./utils/apiClient";
import axios from "axios";

function App() {
  // const [token, setToken] = useToken();
  const [file, setFile] = useState(null);
  const [previewData, setPreviewData] = useState([]);
  const [table, setTable] = useState("");
  const [selectedCols, setSelectedCols] = useState([]);
  const [exportedFile, setExportedFile] = useState("");

  const previewCSV = (file) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const lines = e.target.result.split("\n").slice(0, 5);
      setPreviewData(lines.map((line) => line.split(",")));
    };
    reader.readAsText(file);
  };

  const handleExport = async () => {
    const form = new FormData();
    form.append("host", "clickhouse");
    form.append("port", "9000");
    form.append("database", "default");
    form.append("user", "default");
    // form.append("jwt_token", token);
    form.append("table", table);
    selectedCols.forEach((col) => form.append("columns[]", col));

    try {
      const res = await exportToCSV(form);
      setExportedFile(res.data.file);
    } catch (e) {
      alert("Export failed: " + e.response?.data?.error);
    }
  };
  const handleLogin = async () => {
    try {
      const response = await axios.post('http://localhost:3000/api/login')
      alert("Login successful");
      console.log("Login response:", response);
    } catch (error) {
      alert("Login failed: " + error.response?.data?.error);
      console.error("Login error:", error);
    }
  }
  return (
    <div className="p-6 bg-gray-100 min-h-screen">
      <div className="bg-white max-w-xl mx-auto p-6 rounded shadow">
        <h1 className="text-2xl font-bold mb-4">🔄 Ingestion Tool</h1>

         <button
          onClick={() => handleLogin()}
          className="bg-red-600 text-white px-4 py-2 mb-4 rounded"
        >Login</button>

        <hr className="my-4" />
        <h2 className="text-lg font-semibold mb-2">📤 Upload CSV</h2>
        <FileUpload file={file} setFile={setFile} previewCSV={previewCSV} />
        <CSVPreview previewData={previewData} />

        <hr className="my-4" />
        <h2 className="text-lg font-semibold mb-2">
          📥 Export ClickHouse Table
        </h2>
        <input
          value={table}
          onChange={(e) => setTable(e.target.value)}
          placeholder="Enter table name"
          className="border px-2 py-1 w-full mb-2"
        />
        <ColumnSelector
          table={table}
          selectedCols={selectedCols}
          setSelectedCols={setSelectedCols}
        />
        <button
          onClick={handleExport}
          className="bg-green-600 text-white px-4 py-2 mt-2 rounded"
        >
          Export
        </button>

        {exportedFile && (
          <div className="mt-2 text-sm">
            ✅ Exported:{" "}
            <a
              href={`http://localhost:8080/${exportedFile}`}
              className="text-blue-600 underline"
              target="_blank"
            >
              Download
            </a>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
