import { useState } from "react";
import axios from "axios";

function App() {
  const [file, setFile] = useState(null);
  const [exportedFile, setExportedFile] = useState(null);
  const [columns, setColumns] = useState("");
  const [table, setTable] = useState("");

  const baseURL = "http://localhost:8080/api";

  const handleUpload = async () => {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("host", "clickhouse");
    formData.append("port", "9000");
    formData.append("database", "default");
    formData.append("user", "default");
    formData.append("jwt_token", "your_token");
    formData.append("target_table", "test");

    try {
      const res = await axios.post(`${baseURL}/file-to-clickhouse`, formData);
      alert(`✅ Uploaded ${res.data.rows} rows successfully!`);
    } catch (e) {
      alert("❌ Upload failed: " + e.response?.data?.error);
    }
  };

  const handleExport = async () => {
    try {
      const res = await axios.post(`${baseURL}/clickhouse-to-file`, {
        host: "clickhouse",
        port: "9000",
        database: "default",
        user: "default",
        jwt_token: "your_token",
        table,
        columns: columns.split(",").map((c) => c.trim()),
      });
      setExportedFile(res.data.file);
    } catch (e) {
      alert("❌ Export failed: " + e.response?.data?.error);
    }
  };

  return (
    <div className="min-h-screen p-6 bg-gray-100">
      <div className="max-w-xl mx-auto bg-white shadow-md rounded-lg p-6">
        <h1 className="text-2xl font-bold mb-4">
          📦 ClickHouse Ingestion Tool
        </h1>

        <h2 className="font-semibold text-lg mt-4 mb-2">
          📤 Upload CSV to ClickHouse
        </h2>
        <input
          type="file"
          onChange={(e) => setFile(e.target.files[0])}
          className="mb-2"
        />
        <button
          onClick={handleUpload}
          className="bg-blue-600 text-white px-4 py-2 rounded"
        >
          Upload
        </button>

        <h2 className="font-semibold text-lg mt-6 mb-2">
          📥 Export ClickHouse Table to CSV
        </h2>
        <input
          type="text"
          placeholder="Table name"
          value={table}
          onChange={(e) => setTable(e.target.value)}
          className="border px-2 py-1 w-full mb-2"
        />
        <input
          type="text"
          placeholder="Columns (comma separated)"
          value={columns}
          onChange={(e) => setColumns(e.target.value)}
          className="border px-2 py-1 w-full mb-2"
        />
        <button
          onClick={handleExport}
          className="bg-green-600 text-white px-4 py-2 rounded"
        >
          Export
        </button>

        {exportedFile && (
          <div className="mt-4">
            ✅ Exported:{" "}
            <a
              href={`http://localhost:8080/${exportedFile}`}
              target="_blank"
              className="text-blue-600 underline"
            >
              Download CSV
            </a>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
