import axios, { AxiosRequestConfig } from 'axios';

class ApiClient {
  private readonly baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  public async get(path: string, config?: AxiosRequestConfig) {
    return await axios.get(`${this.baseUrl}${path}`, config);
  }
}

const apiClient = new ApiClient(process.env.ApiBaseUrl);

export default apiClient;
