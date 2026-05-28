import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  uploadFile,
  generateFromFile,
  generateFromText,
} from '../api/presentations';
import { getErrorMessage } from '../api/client';
import { Button } from '../components/ui/Button';
import { Card } from '../components/ui/Card';
import { DropZone } from '../components/upload/DropZone';
import { LoadingSpinner } from '../components/ui/LoadingSpinner';
import { useProtectedRoute } from '../hooks/useProtectedRoute';

type Tab = 'file' | 'text';

export function Upload() {
  const { user, isLoading: authLoading } = useProtectedRoute();
  const navigate = useNavigate();
  const [tab, setTab] = useState<Tab>('file');
  const [text, setText] = useState('');
  const [fileId, setFileId] = useState<string | null>(null);
  const [fileName, setFileName] = useState('');
  const [step, setStep] = useState<'idle' | 'uploading' | 'generating'>('idle');
  const [error, setError] = useState('');

  if (authLoading || !user) return <LoadingSpinner />;

  async function handleFile(file: File) {
    setError('');
    setStep('uploading');
    try {
      const { file: uploaded } = await uploadFile(file);
      setFileId(uploaded.id);
      setFileName(uploaded.original_name);
      setStep('idle');
    } catch (err) {
      setError(getErrorMessage(err));
      setStep('idle');
    }
  }

  async function handleGenerate() {
    setError('');
    setStep('generating');
    try {
      const result =
        tab === 'text'
          ? await generateFromText(text)
          : fileId
            ? await generateFromFile(fileId)
            : null;

      if (!result) {
        setError('Please upload a file or paste text first');
        setStep('idle');
        return;
      }

      navigate(`/presentation/${result.presentation.id}`);
    } catch (err) {
      setError(getErrorMessage(err));
      setStep('idle');
    }
  }

  const busy = step !== 'idle';

  return (
    <div className="mx-auto max-w-2xl animate-fade-in">
      <h1 className="text-2xl font-bold text-slate-900">Create a new deck</h1>
      <p className="mt-1 text-slate-500">
        Upload a document or paste text — AI will build a 7-slide pitch deck.
      </p>

      <div className="mt-6 flex gap-2 rounded-lg bg-slate-100 p-1">
        {(['file', 'text'] as Tab[]).map((t) => (
          <button
            key={t}
            onClick={() => setTab(t)}
            className={`flex-1 rounded-md py-2 text-sm font-medium transition ${
              tab === t ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-600'
            }`}
          >
            {t === 'file' ? 'Upload file' : 'Paste text'}
          </button>
        ))}
      </div>

      <Card className="mt-6">
        {step === 'generating' ? (
          <LoadingSpinner message="GPT is crafting your pitch deck... This may take 30–60 seconds." />
        ) : tab === 'file' ? (
          <>
            <DropZone onFileSelect={handleFile} disabled={busy} />
            {step === 'uploading' && (
              <p className="mt-4 text-center text-sm text-slate-500">Uploading...</p>
            )}
            {fileId && (
              <p className="mt-4 text-center text-sm text-emerald-600">
                ✓ Ready: {fileName}
              </p>
            )}
          </>
        ) : (
          <textarea
            className="w-full min-h-[200px] rounded-lg border border-slate-200 p-4 text-sm focus:border-brand-500 focus:outline-none focus:ring-2 focus:ring-brand-500/20"
            placeholder="Paste your business plan, product notes, or any source text..."
            value={text}
            onChange={(e) => setText(e.target.value)}
            disabled={busy}
          />
        )}

        {error && (
          <p className="mt-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p>
        )}

        {step !== 'generating' && (
          <Button
            className="mt-6 w-full"
            onClick={handleGenerate}
            disabled={busy || (tab === 'file' ? !fileId : text.trim().length < 50)}
            loading={step === 'uploading'}
          >
            Generate pitch deck
          </Button>
        )}
        {tab === 'text' && text.length > 0 && text.length < 50 && (
          <p className="mt-2 text-xs text-slate-500">Enter at least 50 characters of text.</p>
        )}
      </Card>
    </div>
  );
}
