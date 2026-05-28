import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { listPresentations } from '../api/presentations';
import { getErrorMessage } from '../api/client';
import { Button } from '../components/ui/Button';
import { Card } from '../components/ui/Card';
import { LoadingSpinner } from '../components/ui/LoadingSpinner';
import { useProtectedRoute } from '../hooks/useProtectedRoute';
import type { Presentation } from '../types';

const statusStyles: Record<string, string> = {
  completed: 'bg-emerald-100 text-emerald-700',
  processing: 'bg-amber-100 text-amber-700',
  failed: 'bg-red-100 text-red-700',
  pending: 'bg-slate-100 text-slate-600',
};

export function Dashboard() {
  const { user, isLoading: authLoading } = useProtectedRoute();
  const [presentations, setPresentations] = useState<Presentation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!user) return;
    listPresentations()
      .then(setPresentations)
      .catch((err) => setError(getErrorMessage(err)))
      .finally(() => setLoading(false));
  }, [user]);

  if (authLoading || !user) return <LoadingSpinner message="Checking session..." />;
  if (loading) return <LoadingSpinner message="Loading your decks..." />;

  return (
    <div className="mx-auto max-w-5xl animate-fade-in">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">
            Welcome back, {user.full_name.split(' ')[0]}
          </h1>
          <p className="text-slate-500">Manage your AI-generated pitch decks</p>
        </div>
        <Link to="/upload">
          <Button>+ New presentation</Button>
        </Link>
      </div>

      {error && (
        <p className="mt-4 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700">{error}</p>
      )}

      <div className="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {presentations.length === 0 ? (
          <Card className="col-span-full text-center py-12">
            <p className="text-slate-500">No presentations yet.</p>
            <Link to="/upload" className="mt-4 inline-block">
              <Button variant="secondary">Create your first deck</Button>
            </Link>
          </Card>
        ) : (
          presentations.map((p) => (
            <Link key={p.id} to={`/presentation/${p.id}`}>
              <Card className="h-full transition hover:border-brand-300 hover:shadow-md cursor-pointer">
                <div className="flex items-start justify-between gap-2">
                  <h3 className="font-semibold text-slate-900 line-clamp-2">{p.title}</h3>
                  <span
                    className={`shrink-0 rounded-full px-2 py-0.5 text-xs font-medium ${
                      statusStyles[p.status] || statusStyles.pending
                    }`}
                  >
                    {p.status}
                  </span>
                </div>
                <p className="mt-2 text-xs text-slate-500">
                  {new Date(p.created_at).toLocaleDateString()}
                </p>
              </Card>
            </Link>
          ))
        )}
      </div>
    </div>
  );
}
