import { uploadCSV } from "../services/api";

export default function FileUpload({ file, setFile, previewCSV }) {
  const handleChange = (e) => {
    const selected = e.target.files[0];
    setFile(selected);
    previewCSV(selected);
  };

  const handleUpload = async () => {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("host", "clickhouse");
    formData.append("port", "9000");
    formData.append("database", "default");
    formData.append("user", "default");
    formData.append("jwt_token", localStorage.getItem("jwt_token") || "");
    formData.append("target_table", "test");

    try {
      const res = await uploadCSV(formData);
      alert(`✅ Uploaded ${res.data.rows} rows`);
    } catch (e) {
      alert(`❌ Upload failed: ${e.response?.data?.error}`);
    }
  };

  return (
    <>
      <input type="file" onChange={handleChange} className="mb-2" />
      <button
        onClick={handleUpload}
        className="bg-blue-600 text-white px-4 py-2 rounded"
      >
        Upload
      </button>
    </>
  );
}
