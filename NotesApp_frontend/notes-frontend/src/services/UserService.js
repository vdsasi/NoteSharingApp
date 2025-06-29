import axios from '../api/axios';

export const searchUsers = async (query) => {
  const response = await axios.get(`/auth/search-users?query=${encodeURIComponent(query)}`);
  return response.data;
};
