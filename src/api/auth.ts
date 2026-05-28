import { api } from './client';
import type { AuthResponse } from '../types';

export async function register(
  email: string,
  password: string,
  fullName: string
): Promise<AuthResponse> {
  const { data } = await api.post<AuthResponse>('/register', {
    email,
    password,
    full_name: fullName,
  });
  return data;
}

export async function login(email: string, password: string): Promise<AuthResponse> {
  const { data } = await api.post<AuthResponse>('/login', { email, password });
  return data;
}
