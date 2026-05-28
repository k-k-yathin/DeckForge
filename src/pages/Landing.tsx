import { Link } from 'react-router-dom';
import { Button } from '../components/ui/Button';

export function Landing() {
  return (
    <div className="min-h-screen hero-grid">
      <header className="mx-auto flex max-w-6xl items-center justify-between px-6 py-6">
        <div className="flex items-center gap-2">
          <span className="flex h-10 w-10 items-center justify-center rounded-xl bg-brand-600 text-white font-bold">
            D
          </span>
          <span className="text-xl font-bold">DeckForge</span>
        </div>
        <div className="flex gap-3">
          <Link to="/login">
            <Button variant="ghost">Sign in</Button>
          </Link>
          <Link to="/register">
            <Button>Get started</Button>
          </Link>
        </div>
      </header>

      <section className="mx-auto max-w-4xl px-6 py-24 text-center animate-fade-in">
        <span className="rounded-full bg-brand-50 px-4 py-1 text-sm font-medium text-brand-700">
          AI-powered pitch decks
        </span>
        <h1 className="mt-6 text-5xl font-bold tracking-tight text-slate-900 sm:text-6xl">
          Turn documents into{' '}
          <span className="bg-gradient-to-r from-brand-600 to-purple-600 bg-clip-text text-transparent">
            investor-ready decks
          </span>
        </h1>
        <p className="mx-auto mt-6 max-w-2xl text-lg text-slate-600">
          Upload a PDF, Word doc, or paste your notes. DeckForge extracts the content,
          uses GPT to structure a 7-slide pitch deck, and lets you export to PPTX or PDF.
        </p>
        <div className="mt-10 flex flex-wrap justify-center gap-4">
          <Link to="/register">
            <Button size="lg">Start for free</Button>
          </Link>
          <Link to="/login">
            <Button size="lg" variant="secondary">
              View demo account
            </Button>
          </Link>
        </div>
      </section>

      <section className="mx-auto grid max-w-5xl gap-6 px-6 pb-24 sm:grid-cols-3">
        {[
          {
            title: 'Upload anything',
            desc: 'PDF, DOCX, or plain text — we extract the content automatically.',
            icon: '📤',
          },
          {
            title: 'AI structuring',
            desc: 'GPT builds title, problem, solution, market, features, roadmap & conclusion slides.',
            icon: '🤖',
          },
          {
            title: 'Export & share',
            desc: 'Download professional PPTX or PDF files in one click.',
            icon: '📥',
          },
        ].map((f) => (
          <div
            key={f.title}
            className="rounded-xl border border-slate-200 bg-white p-6 shadow-sm"
          >
            <div className="text-3xl">{f.icon}</div>
            <h3 className="mt-4 font-semibold text-slate-900">{f.title}</h3>
            <p className="mt-2 text-sm text-slate-600">{f.desc}</p>
          </div>
        ))}
      </section>
    </div>
  );
}
