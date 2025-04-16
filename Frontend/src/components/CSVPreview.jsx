export default function CSVPreview({ previewData }) {
  if (!previewData.length) return null;

  return (
    <div className="bg-white p-2 mt-2 border">
      <strong>Preview:</strong>
      <table className="table-auto mt-1 text-sm">
        <tbody>
          {previewData.map((row, i) => (
            <tr key={i}>
              {row.map((cell, j) => (
                <td key={j} className="border px-2">
                  {cell}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
