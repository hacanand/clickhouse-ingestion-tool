import { getColumns } from "../services/api";
import { useState } from "react";
export default function ColumnSelector({
  table,
  selectedCols,
  setSelectedCols,
}) {
  const [fetchedCols, setFetchedCols] = useState([]);

  const fetchColumns = async () => {
    const form = new FormData();
    form.append("host", "clickhouse");
    form.append("port", "9000");
    form.append("database", "default");
    form.append("user", "default");
    form.append("jwt_token", localStorage.getItem("jwt_token") || "");
    form.append("table", table);

    const res = await getColumns(form);
    setFetchedCols(res.data.columns);
    setSelectedCols(res.data.columns);
  };

  return (
    <>
      <button
        onClick={fetchColumns}
        className="bg-gray-600 text-white px-2 py-1 rounded mt-2"
      >
        Load Columns
      </button>
      <div className="mt-2">
        {fetchedCols.map((col) => (
          <label key={col} className="block">
            <input
              type="checkbox"
              checked={selectedCols.includes(col)}
              onChange={() =>
                setSelectedCols((prev) =>
                  prev.includes(col)
                    ? prev.filter((c) => c !== col)
                    : [...prev, col]
                )
              }
            />{" "}
            {col}
          </label>
        ))}
      </div>
    </>
  );
}
