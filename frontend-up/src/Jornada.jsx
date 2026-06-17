import { useState, useEffect } from 'react'
import {
  Footprints,
  Blocks,
  MessageCircle,
  Users,
  Hand,
  Palette,
  CheckCircle2,
  Clock,
  Trophy,
  Plus,
} from 'lucide-react'
import { getMilestones, createMilestone } from './api'

const categoryOptions = ['Motor', 'Linguagem', 'Social', 'Cognitivo']

const categoryStyles = {
  Motor: { bg: 'bg-amber-100', fg: 'text-amber-600', chip: 'bg-amber-50 text-amber-600' },
  Linguagem: { bg: 'bg-indigo-100', fg: 'text-indigo-600', chip: 'bg-indigo-50 text-indigo-600' },
  Social: { bg: 'bg-pink-100', fg: 'text-pink-600', chip: 'bg-pink-50 text-pink-600' },
  Cognitivo: { bg: 'bg-emerald-100', fg: 'text-emerald-600', chip: 'bg-emerald-50 text-emerald-600' },
}

const categoryIcons = {
  Motor: [Footprints, Blocks, Hand],
  Linguagem: [MessageCircle],
  Social: [Users],
  Cognitivo: [Palette],
}

// escolhe um icone pro marco, alternando entre os icones da categoria
function iconFor(category, index) {
  const icons = categoryIcons[category] ?? [CheckCircle2]
  return icons[index % icons.length]
}

// formata "mes de ano" capitalizado, usado pra agrupar os marcos por mes
function monthLabel(dateStr) {
  if (!dateStr) return 'Sem data'
  const date = new Date(dateStr)
  const label = date.toLocaleDateString('pt-BR', { month: 'long', year: 'numeric' })
  return label.charAt(0).toUpperCase() + label.slice(1)
}

// cabecalho com a barra de progresso geral da jornada (quantos marcos foram conquistados)
function ProgressHeader({ total, done }) {
  const pct = total > 0 ? Math.round((done / total) * 100) : 0
  return (
    <section className="animate-slide-up rounded-3xl bg-white p-6 shadow-sm ring-1 ring-slate-100">
      <div className="flex items-center gap-4">
        <span className="grid h-14 w-14 shrink-0 place-items-center rounded-2xl bg-emerald-500 text-white shadow-md shadow-emerald-200 transition-transform duration-300 hover:scale-105">
          <Trophy className="h-7 w-7" />
        </span>
        <div className="flex-1">
          <h2 className="text-xl font-bold text-slate-800">Jornada UP</h2>
          <p className="text-sm text-slate-400">Marcos de desenvolvimento</p>
        </div>
        <div className="text-right">
          <p className="text-2xl font-extrabold text-emerald-600">
            {done}
            <span className="text-base font-medium text-slate-400">/{total}</span>
          </p>
          <p className="text-xs text-slate-400">conquistados</p>
        </div>
      </div>

      <div className="mt-5">
        <div className="h-2.5 w-full overflow-hidden rounded-full bg-slate-100">
          <div className="h-full rounded-full bg-emerald-500 transition-all duration-700 ease-out" style={{ width: `${pct}%` }} />
        </div>
        <p className="mt-2 text-xs text-slate-400">{pct}% da jornada concluída</p>
      </div>
    </section>
  )
}

// card de um marco na linha do tempo, mostrando categoria e se ja foi conquistado
function MilestoneCard({ item, isLast, index }) {
  const c = categoryStyles[item.category] ?? categoryStyles.Cognitivo
  const Icon = iconFor(item.category, item.id)
  return (
    <div className={`animate-slide-in-left stagger-${Math.min(index + 1, 6)} flex gap-4`}>
      <div className="flex flex-col items-center">
        <span className={`grid h-11 w-11 shrink-0 place-items-center rounded-full ${c.bg} ${c.fg} transition-transform duration-200 hover:scale-110`}>
          <Icon className="h-5 w-5" strokeWidth={2} />
        </span>
        {!isLast && <span className="my-1 w-0.5 flex-1 bg-slate-100" />}
      </div>

      <div className="mb-5 flex-1 rounded-2xl bg-white p-4 shadow-sm ring-1 ring-slate-100 transition-all duration-200 hover:-translate-y-0.5 hover:shadow-md">
        <div className="flex flex-wrap items-center gap-2">
          <span className={`rounded-full px-2.5 py-0.5 text-xs font-semibold ${c.chip}`}>{item.category}</span>
          {item.done ? (
            <span className="inline-flex items-center gap-1 text-xs font-semibold text-emerald-600">
              <CheckCircle2 className="h-3.5 w-3.5" />
              Conquistado
            </span>
          ) : (
            <span className="inline-flex items-center gap-1 text-xs font-semibold text-amber-500">
              <Clock className="h-3.5 w-3.5" />
              Em progresso
            </span>
          )}
          <span className="ml-auto text-xs text-slate-400">
            {item.achieved_at ? new Date(item.achieved_at).toLocaleDateString('pt-BR', { day: '2-digit', month: 'short' }) : ''}
          </span>
        </div>
        <h4 className="mt-2 font-bold text-slate-800">{item.title}</h4>
        <p className="mt-1 text-sm leading-relaxed text-slate-500">{item.description}</p>
      </div>
    </div>
  )
}

// formulario de cadastrar um marco novo (titulo, categoria, descricao e se ja foi conquistado)
function NewMilestoneForm({ onCreate, onClose }) {
  const [title, setTitle] = useState('')
  const [category, setCategory] = useState(categoryOptions[0])
  const [description, setDescription] = useState('')
  const [done, setDone] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  // valida titulo/descricao e cria o marco
  async function submit(e) {
    e.preventDefault()
    if (!title.trim() || !description.trim()) {
      setError('Preencha todos os campos.')
      return
    }
    setSaving(true)
    setError('')
    try {
      await onCreate({
        title,
        category,
        description,
        done,
        achieved_at: done ? new Date().toISOString() : null,
      })
      setTitle('')
      setDescription('')
      onClose()
    } catch (err) {
      setError(err.message)
    } finally {
      setSaving(false)
    }
  }

  return (
    <form onSubmit={submit} className="animate-slide-up space-y-2 rounded-2xl bg-white p-5 shadow-sm ring-1 ring-slate-100">
      <h3 className="font-bold text-slate-800">Novo marco</h3>
      <input
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="Título do marco"
        className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
      />
      <select
        value={category}
        onChange={(e) => setCategory(e.target.value)}
        className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
      >
        {categoryOptions.map((c) => (
          <option key={c} value={c}>{c}</option>
        ))}
      </select>
      <textarea
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        placeholder="Descrição"
        rows={2}
        className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
      />
      <label className="flex items-center gap-2 text-sm text-slate-600">
        <input type="checkbox" checked={done} onChange={(e) => setDone(e.target.checked)} />
        Já conquistado
      </label>
      {error && <p className="text-xs text-rose-500">{error}</p>}
      <div className="flex justify-end gap-2">
        <button type="button" onClick={onClose} className="rounded-full px-3 py-1.5 text-sm text-slate-500 hover:bg-slate-100">
          Cancelar
        </button>
        <button
          type="submit"
          disabled={saving}
          className="rounded-full bg-emerald-500 px-4 py-1.5 text-sm font-semibold text-white hover:bg-emerald-600 disabled:opacity-60"
        >
          {saving ? 'Salvando...' : 'Adicionar'}
        </button>
      </div>
    </form>
  )
}

// tela da jornada de desenvolvimento de uma crianca: progresso geral e marcos agrupados por mes
export default function Jornada({ studentId, canEdit }) {
  const [milestones, setMilestones] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showForm, setShowForm] = useState(false)

  useEffect(() => {
    if (!studentId) return
    setLoading(true)
    getMilestones(studentId)
      .then(setMilestones)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [studentId])

  // cria o marco e recarrega a lista pra atualizar progresso e agrupamento
  async function handleCreateMilestone(data) {
    await createMilestone(studentId, data)
    const updated = await getMilestones(studentId)
    setMilestones(updated)
  }

  const groups = milestones.reduce((acc, item) => {
    const key = monthLabel(item.achieved_at)
    if (!acc[key]) acc[key] = []
    acc[key].push(item)
    return acc
  }, {})

  const total = milestones.length
  const done = milestones.filter((m) => m.done).length

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <ProgressHeader total={total} done={done} />

      {canEdit && !showForm && (
        <button
          onClick={() => setShowForm(true)}
          className="animate-pop-in flex items-center gap-1.5 rounded-full bg-emerald-500 px-4 py-2 text-sm font-semibold text-white shadow-md shadow-emerald-200 transition-all duration-200 hover:bg-emerald-600 hover:shadow-lg active:scale-95"
        >
          <Plus className="h-4 w-4" />
          Novo marco
        </button>
      )}
      {showForm && <NewMilestoneForm onCreate={handleCreateMilestone} onClose={() => setShowForm(false)} />}

      {error && <p className="rounded-2xl bg-rose-50 p-4 text-sm text-rose-600">{error}</p>}
      {loading && <p className="text-sm text-slate-400">Carregando...</p>}
      {!loading && total === 0 && !error && (
        <p className="rounded-2xl bg-white p-6 text-center text-sm text-slate-400 ring-1 ring-slate-100">
          Nenhum marco registrado ainda.
        </p>
      )}

      {Object.entries(groups).map(([month, items]) => (
        <section key={month}>
          <h3 className="mb-3 text-sm font-bold uppercase tracking-wide text-slate-400">{month}</h3>
          <div>
            {items.map((item, i) => (
              <MilestoneCard key={item.id} item={item} index={i} isLast={i === items.length - 1} />
            ))}
          </div>
        </section>
      ))}
    </div>
  )
}
