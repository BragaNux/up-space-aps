import { useState, useRef, useEffect } from 'react'
import { ChevronDown, Search, X } from 'lucide-react'

// dropdown com busca embutida; options e um array de { value, label, hint?, badge? }, e ao escolher chama onChange(value)
export default function SearchableSelect({ options, value, onChange, placeholder = 'Buscar...', emptyMessage = 'Nenhum resultado' }) {
  const [open, setOpen] = useState(false)
  const [query, setQuery] = useState('')
  const containerRef = useRef(null)

  useEffect(() => {
    // fecha o dropdown se clicar fora dele
    function handleClickOutside(e) {
      if (containerRef.current && !containerRef.current.contains(e.target)) {
        setOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const selected = options.find((o) => o.value === value)
  const filtered = options.filter((o) => o.label.toLowerCase().includes(query.toLowerCase()))

  return (
    <div ref={containerRef} className="relative">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        className="flex w-full items-center justify-between gap-2 rounded-lg border border-slate-200 bg-white px-3 py-2 text-left text-sm outline-none transition focus:border-emerald-400"
      >
        <span className={selected ? 'text-slate-800' : 'text-slate-400'}>
          {selected ? selected.label : placeholder}
        </span>
        <span className="flex items-center gap-1">
          {selected && (
            <X
              className="h-3.5 w-3.5 text-slate-300 hover:text-slate-500"
              onClick={(e) => { e.stopPropagation(); onChange(null) }}
            />
          )}
          <ChevronDown className={`h-4 w-4 text-slate-400 transition-transform ${open ? 'rotate-180' : ''}`} />
        </span>
      </button>

      {open && (
        <div className="animate-pop-in absolute z-20 mt-1.5 w-full overflow-hidden rounded-xl border border-slate-200 bg-white shadow-lg">
          <div className="flex items-center gap-2 border-b border-slate-100 px-3 py-2">
            <Search className="h-3.5 w-3.5 text-slate-400" />
            <input
              autoFocus
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Digite para filtrar..."
              className="w-full text-sm outline-none"
            />
          </div>
          <div className="max-h-48 overflow-y-auto">
            {filtered.length === 0 ? (
              <p className="px-3 py-3 text-center text-xs text-slate-400">{emptyMessage}</p>
            ) : (
              filtered.map((o) => (
                <button
                  key={o.value}
                  type="button"
                  onClick={() => { onChange(o.value); setOpen(false); setQuery('') }}
                  className={`flex w-full flex-col items-start px-3 py-2 text-left text-sm transition hover:bg-emerald-50 ${
                    o.value === value ? 'bg-emerald-50 text-emerald-700' : 'text-slate-700'
                  }`}
                >
                  <span className="flex w-full items-center justify-between gap-2">
                    {o.label}
                    {o.badge && (
                      <span className="rounded-full bg-amber-100 px-2 py-0.5 text-[10px] font-semibold text-amber-700">
                        {o.badge}
                      </span>
                    )}
                  </span>
                  {o.hint && <span className="text-xs text-slate-400">{o.hint}</span>}
                </button>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  )
}
