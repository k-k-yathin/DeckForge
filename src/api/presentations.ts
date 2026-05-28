import { api } from './client';
import type { Presentation, UploadedFile } from '../types';

/** Parse slide content — backend stores bullets as JSON array in content field */
function normalizePresentation(p: Presentation): Presentation {
  if (!p.slides) return p;
  return {
    ...p,
    slides: p.slides.map((s) => ({
      ...s,
      content: Array.isArray(s.content)
        ? s.content
        : typeof s.content === 'string'
          ? JSON.parse(s.content as unknown as string)
          : [],
    })),
  };
}

export async function uploadFile(file: File): Promise<{ file: UploadedFile }> {
  const form = new FormData();
  form.append('file', file);
  const { data } = await api.post('/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
  return data;
}

export async function generateFromFile(fileId: string): Promise<{ presentation: Presentation }> {
  const { data } = await api.post('/generate', { file_id: fileId });
  return { presentation: normalizePresentation(data.presentation) };
}

export async function generateFromText(text: string): Promise<{ presentation: Presentation }> {
  const { data } = await api.post('/generate', { text });
  return { presentation: normalizePresentation(data.presentation) };
}

export async function listPresentations(): Promise<Presentation[]> {
  const { data } = await api.get<{ presentations: Presentation[] }>('/presentations');
  return data.presentations;
}

export async function getPresentation(id: string): Promise<Presentation> {
  const { data } = await api.get<{ presentation: Presentation }>(`/presentation/${id}`);
  return normalizePresentation(data.presentation);
}

export function exportUrl(id: string, format: 'pptx' | 'pdf'): string {
  const base = import.meta.env.VITE_API_URL || '/api/v1';
  const token = localStorage.getItem('deckforge_token');
  // Browser download with auth: we use fetch in UI, but this builds the path
  return `${base}/presentation/${id}/export/${format}?token=${token}`;
}

export async function downloadExport(id: string, format: 'pptx' | 'pdf', filename: string) {
  const { data } = await api.get(`/presentation/${id}/export/${format}`, {
    responseType: 'blob',
  });
  const url = window.URL.createObjectURL(new Blob([data]));
  const link = document.createElement('a');
  link.href = url;
  link.setAttribute('download', filename);
  document.body.appendChild(link);
  link.click();
  link.remove();
  window.URL.revokeObjectURL(url);
}
