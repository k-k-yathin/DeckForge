// Shared TypeScript types matching the Go API JSON responses

export interface User {
  id: string;
  email: string;
  full_name: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}

export interface UploadedFile {
  id: string;
  original_name: string;
  file_type: string;
  file_size: number;
  extracted_text?: string;
}

export interface Slide {
  id: string;
  presentation_id: string;
  slide_order: number;
  slide_type: string;
  title: string;
  subtitle?: string;
  content: string[]; // bullet points parsed from JSONB
}

export interface Presentation {
  id: string;
  user_id: string;
  uploaded_file_id?: string;
  title: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  source_summary?: string;
  created_at: string;
  updated_at: string;
  slides?: Slide[];
}
