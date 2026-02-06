import axios from 'axios';
import type {
  LoginResponse,
  RegisterResponse,
  LoginRequest,
  RegisterRequest,
  Group,
  Transaction,
  CreateTransactionRequest,
  CreateTransactionResponse,
  GroupMember,
  Debt,
} from '../types';

const API_BASE_URL = '/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = localStorage.getItem('refresh_token');
        if (refreshToken) {
          const response = await axios.post(`${API_BASE_URL}/token/refresh`, {
            refresh_token: refreshToken,
          });
          const { token } = response.data;
          localStorage.setItem('token', token);
          originalRequest.headers.Authorization = `Bearer ${token}`;
          return api(originalRequest);
        }
      } catch (refreshError) {
        localStorage.removeItem('token');
        localStorage.removeItem('refresh_token');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

// Auth APIs
export const authAPI = {
  login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post('/login', credentials);
    return response.data;
  },

  register: async (userData: RegisterRequest): Promise<RegisterResponse> => {
    const response = await api.post('/register', userData);
    return response.data;
  },

  refreshToken: async (): Promise<{ token: string }> => {
    const response = await api.post('/token/refresh');
    return response.data;
  },
};

// Group APIs
export const groupAPI = {
  getMyGroups: async (): Promise<Group[]> => {
    const response = await api.get('/groups');
    return response.data;
  },

  getAllGroups: async (): Promise<Group[]> => {
    const response = await api.get('/groups/all');
    return response.data;
  },

  createGroup: async (name: string): Promise<Group> => {
    const response = await api.post('/groups', { name });
    return response.data;
  },

  joinGroup: async (name: string): Promise<Group> => {
    const response = await api.post('/groups/join', { name });
    return response.data;
  },

  getMembers: async (groupId: string): Promise<GroupMember[]> => {
    const response = await api.get(`/groups/members/${groupId}`);
    return response.data;
  },
};

// Transaction APIs
export const transactionAPI = {
  create: async (
    data: CreateTransactionRequest
  ): Promise<CreateTransactionResponse> => {
    const response = await api.post('/transaction', data);
    return response.data;
  },

  getTransactions: async (params?: {
    group_id?: string;
    author_id?: string;
  }): Promise<Transaction[]> => {
    const response = await api.get('/transactions', { params });
    return response.data;
  },

  getTransactionDetails: async (
    transactionIds: string[]
  ): Promise<{ [key: string]: Debt[] }> => {
    const response = await api.get('/transactions/details', {
      params: { transactions: transactionIds.join(',') },
    });
    return response.data;
  },
};

export default api;
