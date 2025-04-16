import axios from "axios";

// Create an Axios instance
const apiClient = axios.create({
  baseURL: "http://localhost:3000/api", // Set your API base URL
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true, // Allow cookies (e.g., access_token, refresh_token) to be sent with requests
});

// Request interceptor to add Authorization token from cookies
apiClient.interceptors.request.use(
  (config) => {
    // Get the token from the cookie (Access Token)
    const token = getCookie("access_token");
    if (token) {
      config.headers["Authorization"] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle token expiration
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    // If the error status is 401 and token expired
    if (
      error.response &&
      error.response.status === 401 &&
      error.response.data.error === "Token expired"
    ) {
      const newToken = await refreshAccessToken();
      if (newToken) {
        // Retry the original request with the new token
        error.config.headers["Authorization"] = `Bearer ${newToken}`;
        return axios(error.config); // Retry the failed request
      }
    }
    return Promise.reject(error);
  }
);

// Helper function to get a cookie value by name
function getCookie(name) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(";").shift();
  return null;
}

// Function to refresh the access token using the refresh token
const refreshAccessToken = async () => {
  const refreshToken = getCookie("refresh_token"); // Retrieve refresh token from cookies
  if (!refreshToken) {
    console.error("Refresh token not found");
    return null;
  }

  try {
    const response = await axios.post(
      "http://localhost:3000/api/refresh-token",
      {
        refresh_token: refreshToken,
      },
      { withCredentials: true }
    );

    const newAccessToken = response.data.access_token;
    document.cookie = `access_token=${newAccessToken}; path=/; secure; HttpOnly; SameSite=Strict`;
    return newAccessToken;
  } catch (error) {
    console.error("Error refreshing token:", error);
    // Handle token refresh failure (e.g., redirect to login)
    window.location.href = "/login"; // Redirect to login page if token refresh fails
    return null;
  }
};

export default apiClient;
