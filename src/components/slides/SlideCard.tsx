import type { Slide } from '../../types';

const typeLabels: Record<string, string> = {
  title: 'Title',
  problem: 'Problem',
  solution: 'Solution',
  market: 'Market',
  features: 'Features',
  roadmap: 'Roadmap',
  conclusion: 'Conclusion',
};

const typeColors: Record<string, string> = {
  title: 'from-indigo-600 to-purple-700',
  problem: 'from-rose-500 to-orange-500',
  solution: 'from-emerald-500 to-teal-600',
  market: 'from-blue-500 to-cyan-500',
  features: 'from-violet-500 to-fuchsia-500',
  roadmap: 'from-amber-500 to-yellow-500',
  conclusion: 'from-slate-700 to-slate-900',
};

interface SlideCardProps {
  slide: Slide;
  index: number;
}

export function SlideCard({ slide, index }: SlideCardProps) {
  const gradient = typeColors[slide.slide_type] || 'from-brand-600 to-brand-700';
  const label = typeLabels[slide.slide_type] || slide.slide_type;

  return (
    <article
      className="animate-slide-up overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm"
      style={{ animationDelay: `${index * 50}ms` }}
    >
      <div className={`bg-gradient-to-br ${gradient} px-6 py-8 text-white`}>
        <span className="rounded-full bg-white/20 px-2 py-0.5 text-xs font-medium uppercase tracking-wide">
          {label}
        </span>
        <h3 className="mt-3 text-2xl font-bold">{slide.title}</h3>
        {slide.subtitle && (
          <p className="mt-2 text-sm text-white/90">{slide.subtitle}</p>
        )}
      </div>
      {slide.content && slide.content.length > 0 && (
        <ul className="space-y-2 px-6 py-5">
          {slide.content.map((bullet, i) => (
            <li key={i} className="flex gap-2 text-sm text-slate-600">
              <span className="text-brand-500">•</span>
              {bullet}
            </li>
          ))}
        </ul>
      )}
    </article>
  );
}
