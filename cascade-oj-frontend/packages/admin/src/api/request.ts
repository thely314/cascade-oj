import axios from 'axios';
import router from '@/router';

const service = axios.create({
  /*baseURL: import.meta.env.VITE_API_BASE_URL || '/api',*/
  baseURL: '/api', //临时
  timeout: 10000
});

service.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('cascade_token');
    if (token) {
      config.headers['token'] = token;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

service.interceptors.response.use(
  (response) => {
    // 直接返回 `response.data`，简化调用方处理
    return response.data;
  },
  (error) => {
    if (error.response) {
      const status = error.response.status;
      if (status === 401) {
        localStorage.removeItem('cascade_token');
        if (router.currentRoute.value.name !== 'Login') {
          router.push({ name: 'Login', params: { sourceApp: 'admin' } });
        }
      }
    }
    // 对于非 401 错误，直接将原始错误抛出，让调用方处理
    return Promise.reject(error);
  }
);

export default service;
