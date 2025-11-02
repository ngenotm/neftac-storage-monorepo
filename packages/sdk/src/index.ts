import axios from 'axios';

export class NeftacStorage {
  private token: string;
  private baseURL = 'https://s3.airsoko.com';

  constructor(token: string) {
    this.token = token;
  }

  async upload(bucket: string, key: string, file: File) {
    await axios.put(`${this.baseURL}/v1/buckets/${bucket}/objects/${key}`, file, {
      headers: { Authorization: `Bearer ${this.token}` }
    });
  }

  async list(bucket: string) {
    const res = await axios.get(`${this.baseURL}/v1/buckets/${bucket}/objects`, {
      headers: { Authorization: `Bearer ${this.token}` }
    });
    return res.data;
  }
}
