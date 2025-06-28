import axios from '../api/axios';

const AuthService = {
  login: async (username, password) => {
    const response = await axios.post('/auth/login', { username, password });
    return response.data;
  },

  register: async (userData) => {
    const response = await axios.post('/auth/register', userData);
    return response.data;
  },

  getProfile: async () => {
    const response = await axios.get('/auth/me');
    return response.data;
  },

  changePassword: async (oldPassword, newPassword) => {
    const response = await axios.post(
      '/auth/change-password',
      { oldPassword, newPassword }
    );
    return response.data;
  }
};

export default AuthService;
