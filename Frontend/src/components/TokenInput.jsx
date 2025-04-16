export default function TokenInput({ token, setToken }) {
  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        localStorage.setItem("jwt_token", token);
      }}
    >
      <input
        type="text"
        value={token}
        onChange={(e) => setToken(e.target.value)}
        placeholder="JWT Token"
        className="border px-2 py-1 mr-2"
      />
      <button
        type="submit"
        className="bg-purple-600 text-white px-3 py-1 rounded"
      >
        Save Token
      </button>
    </form>
  );
}
