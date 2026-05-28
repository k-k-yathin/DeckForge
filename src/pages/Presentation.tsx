import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import {
  getPresentation,
  downloadExport,
} from '../api/presentations';
import { getErrorMessage } from '../api/client';
import { Button } from '../components/ui/Button';
import { SlideCard } from '../components/slides/SlideCard';
import { LoadingSpinner } from '../components/ui/LoadingSpinner';
import { useProtectedRoute } from '../hooks/useProtectedRoute';
import type { Presentation } from '../types';

export function PresentationPage() {
  const { id } = useParams<{ id: string }>();
  const { user, isLoading: authLoading } = useProtectedRoute();
  const [presentation, setPresentation] = useState<Presentation | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [exporting, setExporting] = useState<'pptx' | 'pdf' | null>(null);

  useEffect(() => {
    if (!user || !id) return;
    getPresentation(id)
      .then(setPresentation)
      .catch((err) => setError(getErrorMessage(err)))
      .finally(() => setLoading(false));
  }, [user, id]);

  async function handleExport(format: 'pptx' | 'pdf') {
    if (!id || !presentation) return;
    setExporting(format);
    try {
      await downloadExport(
        id,
        format,
        `${presentation.title.replace(/\s+/g, '-')}.${format}`
      );
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setExporting(null);
    }
  }

  if (authLoading || !user) return <LoadingSpinner />;
  if (loading) return <LoadingSpinner message="Loading presentation..." />;
  if (error || !presentation) {
    return (
      <div className="text-center py-16">
        <p className="text-red-600">{error || 'Presentation not found'}</p>
        <Link to="/dashboard" className="mt-4 inline-block text-brand-600">
          Back to dashboard
        </Link>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-4xl animate-fade-in">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <Link to="/dashboard" className="text-sm text-brand-600 hover:underline">
            ← Dashboard
          </Link>
          <h1 className="mt-2 text-2xl font-bold text-slate-900">{presentation.title}</h1>
          {presentation.source_summary && (
            <p className="mt-2 text-sm text-slate-500 line-clamp-2">
              {presentation.source_summary}
            </p>
          )}
        </div>
        <div className="flex gap-2">
          <Button
            variant="secondary"
            loading={exporting === 'pptx'}
            onClick={() => handleExport('pptx')}
            disabled={presentation.status !== 'completed'}
          >
            Export PPTX
          </Button>
          <Button
            variant="secondary"
            loading={exporting === 'pdf'}
            onClick={() => handleExport('pdf')}
            disabled={presentation.status !== 'completed'}
          >
            Export PDF
          </Button>
        </div>
      </div>

      <div className="mt-8 grid gap-6">
        {presentation.slides?.map((slide, i) => (
          <SlideCard key={slide.id} slide={slide} index={i} />
        ))}
      </div>
    </div>
  );
}
