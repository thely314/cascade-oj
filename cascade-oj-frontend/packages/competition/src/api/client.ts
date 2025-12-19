import axios from 'axios'

const baseURL = import.meta.env.VITE_API_BASE || '/api'

const client = axios.create({
  baseURL,
  withCredentials: true,
  timeout: 15000,
})

client.interceptors.response.use(
  (resp) => resp.data,
  (error) => {
    // 尽量返回后端的错误结构，否则抛原错误
    throw error?.response?.data ?? error
  }
)

export default client
