import { useState, useEffect } from 'react'
import { Clock, MapPin, Check, PartyPopper, Users, Stethoscope, Bus, CalendarDays, Plus, Pencil, Trash2 } from 'lucide-react'
import { getEvents, rsvpEvent, createEvent, updateEvent, deleteEvent } from './api'

const categoryStyles = {
  Festa: { bg: 'bg-pink-100', fg: 'text-pink-600', chip: 'bg-pink-50 text-pink-600', icon: PartyPopper },
  Reunião: { bg: 'bg-indigo-100', fg: 'text-indigo-600', chip: 'bg-indigo-50 text-indigo-600', icon: Users },
  Saúde: { bg: 'bg-emerald-100', fg: 'text-emerald-600', chip: 'bg-emerald-50 text-emerald-600', icon: Stethoscope },
  Passeio: { bg: 'bg-amber-100', fg: 'text-amber-600', chip: 'bg-amber-50 text-amber-600', icon: Bus },
}

// adivinha a categoria do evento (Festa/Reuniao/Saude/Passeio) com base em palavras do titulo
function categoryFor(title = '') {
  const t = title.toLowerCase()
  if (t.includes('festa') || t.includes('junina')) return 'Festa'
  if (t.includes('reunião') || t.includes('pais')) return 'Reunião'
  if (t.includes('saúde') || t.includes('vacina') || t.includes('terapia')) return 'Saúde'
  if (t.includes('passeio') || t.includes('parque')) return 'Passeio'
  return 'Reunião'
}

const filters = ['Todos', 'Festa', 'Reunião', 'Saúde', 'Passeio']

// formata "inicio - fim" em horario local (ex: 14:00 - 16:00)
function formatTimeRange(startsAt, endsAt) {
  const fmt = (d) => new Date(d).toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
  return `${fmt(startsAt)} - ${fmt(endsAt)}`
}

// converte uma data ISO pro formato que o input datetime-local espera, ja no fuso local
function formatDateTimeLocal(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const offset = d.getTimezoneOffset()
  const localDate = new Date(d.getTime() - offset * 60 * 1000)
  return localDate.toISOString().slice(0, 16)
}

// card de um evento da agenda, com confirmacao de presenca e (se puder editar) botoes de editar/excluir
function EventCard({ event, onToggle, busy, index, canEdit, onEdit, onDelete }) {
  const category = categoryFor(event.title)
  const c = categoryStyles[category]
  const Icon = c.icon
  const start = new Date(event.starts_at)

  const isPast = new Date(event.ends_at) < new Date()

  return (
    <article className={`animate-slide-up stagger-${Math.min(index + 1, 6)} flex gap-4 rounded-2xl bg-white p-4 shadow-sm ring-1 ring-slate-100 transition-all duration-200 hover:-translate-y-0.5 hover:shadow-md`}>
      <div className={`flex h-16 w-16 shrink-0 flex-col items-center justify-center rounded-2xl ${c.bg} ${c.fg} transition-transform duration-200 hover:scale-105`}>
        <span className="text-xl font-extrabold leading-none">{start.getDate().toString().padStart(2, '0')}</span>
        <span className="text-xs font-semibold">{start.toLocaleDateString('pt-BR', { month: 'short' }).toUpperCase().replace('.', '')}</span>
      </div>

      <div className="min-w-0 flex-1">
        <div className="flex flex-wrap items-center gap-2">
          <span className={`inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-semibold ${c.chip}`}>
            <Icon className="h-3.5 w-3.5" />
            {category}
          </span>
          <span className="text-xs text-slate-400">{start.toLocaleDateString('pt-BR', { weekday: 'long' })}</span>
        </div>

        <h4 className="mt-1.5 text-base font-bold text-slate-800">{event.title}</h4>
        <p className="mt-0.5 text-sm leading-relaxed text-slate-500">{event.description}</p>

        <div className="mt-3 flex flex-wrap items-center gap-x-4 gap-y-1 text-sm text-slate-500">
          <span className="inline-flex items-center gap-1.5">
            <Clock className="h-4 w-4 text-slate-400" />
            {formatTimeRange(event.starts_at, event.ends_at)}
          </span>
          <span className="inline-flex items-center gap-1.5">
            <MapPin className="h-4 w-4 text-slate-400" />
            {event.location}
          </span>
        </div>

        <div className="mt-3 flex flex-wrap items-center gap-2">
          {isPast ? (
            <span className="inline-flex items-center gap-1.5 rounded-full bg-slate-100 px-3.5 py-1.5 text-xs font-semibold text-slate-500 ring-1 ring-slate-200/50">
              {event.confirmed ? (
                <>
                  <Check className="h-3.5 w-3.5 text-emerald-500" />
                  Presença confirmada (Encerrado)
                </>
              ) : (
                'Evento encerrado'
              )}
            </span>
          ) : (
            <button
              onClick={onToggle}
              disabled={busy}
              className={`inline-flex items-center gap-1.5 rounded-full px-4 py-2 text-sm font-semibold transition-all duration-200 active:scale-95 disabled:opacity-50 ${
                event.confirmed
                  ? 'bg-emerald-500 text-white shadow-md shadow-emerald-200 hover:bg-emerald-600 hover:shadow-lg'
                  : 'border border-slate-200 text-slate-600 hover:border-emerald-300 hover:bg-emerald-50'
              }`}
            >
              {busy ? (
                <span className="h-4 w-4 animate-spin rounded-full border-2 border-current/30 border-t-current" />
              ) : event.confirmed ? (
                <>
                  <Check className="h-4 w-4" />
                  Presença confirmada
                </>
              ) : (
                'Confirmar presença'
              )}
            </button>
          )}

          {canEdit && (
            <div className="flex gap-1 ml-auto">
              <button
                onClick={onEdit}
                aria-label="Editar"
                className="grid h-8 w-8 place-items-center rounded-full text-slate-400 transition hover:bg-slate-100 hover:text-emerald-600"
              >
                <Pencil className="h-4 w-4" />
              </button>
              <button
                onClick={onDelete}
                aria-label="Excluir"
                className="grid h-8 w-8 place-items-center rounded-full text-slate-400 transition hover:bg-rose-50 hover:text-rose-600"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            </div>
          )}
        </div>
      </div>
    </article>
  )
}

// formulario de criar/editar evento (initial preenchido = modo edicao)
function NewEventForm({ onCreate, onClose, initial }) {
  const [title, setTitle] = useState(initial?.title ?? '')
  const [description, setDescription] = useState(initial?.description ?? '')
  const [location, setLocation] = useState(initial?.location ?? '')
  const [startsAt, setStartsAt] = useState(formatDateTimeLocal(initial?.starts_at))
  const [endsAt, setEndsAt] = useState(formatDateTimeLocal(initial?.ends_at))
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  // valida campos obrigatorios e que o fim seja depois do inicio, antes de salvar
  async function submit(e) {
    e.preventDefault()
    if (!title.trim() || !location.trim() || !startsAt || !endsAt) {
      setError('Preencha título, local, data de início e fim.')
      return
    }
    const start = new Date(startsAt)
    const end = new Date(endsAt)
    if (end < start) {
      setError('A data de término deve ser posterior à data de início.')
      return
    }
    setSaving(true)
    setError('')
    try {
      await onCreate({
        title,
        description,
        location,
        starts_at: start.toISOString(),
        ends_at: end.toISOString(),
      })
      onClose()
    } catch (err) {
      setError(err.message)
    } finally {
      setSaving(false)
    }
  }

  return (
    <form onSubmit={submit} className="animate-slide-up space-y-3 rounded-2xl bg-white p-5 shadow-sm ring-1 ring-slate-100">
      <h3 className="font-bold text-slate-800">{initial ? 'Editar evento' : 'Novo evento'}</h3>
      
      <div>
        <label className="mb-1 block text-xs font-medium text-slate-500">Título</label>
        <input
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Título do evento"
          className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
        />
      </div>

      <div>
        <label className="mb-1 block text-xs font-medium text-slate-500">Descrição</label>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Descrição"
          rows={2}
          className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
        />
      </div>

      <div>
        <label className="mb-1 block text-xs font-medium text-slate-500">Local</label>
        <input
          value={location}
          onChange={(e) => setLocation(e.target.value)}
          placeholder="Local"
          className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
        />
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="mb-1 block text-xs font-medium text-slate-500">Início</label>
          <input
            type="datetime-local"
            value={startsAt}
            onChange={(e) => setStartsAt(e.target.value)}
            className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
          />
        </div>
        <div>
          <label className="mb-1 block text-xs font-medium text-slate-500">Término</label>
          <input
            type="datetime-local"
            value={endsAt}
            onChange={(e) => setEndsAt(e.target.value)}
            className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
          />
        </div>
      </div>

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
          {saving ? 'Salvando...' : 'Salvar'}
        </button>
      </div>
    </form>
  )
}

// tela da agenda de eventos: lista, filtra por categoria, confirma presenca e (se profissional) gerencia eventos
export default function Agenda({ canEdit }) {
  const [filter, setFilter] = useState('Todos')
  const [events, setEvents] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [busyId, setBusyId] = useState(null)
  const [showForm, setShowForm] = useState(false)
  const [editingEventId, setEditingEventId] = useState(null)

  useEffect(() => {
    getEvents()
      .then(setEvents)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  // confirma/cancela presenca no evento e atualiza o contador local
  async function toggle(id) {
    setBusyId(id)
    setError('')
    try {
      const { confirmed } = await rsvpEvent(id)
      setEvents((prev) =>
        prev.map((e) => (e.id === id ? { ...e, confirmed, rsvp_count: e.rsvp_count + (confirmed ? 1 : -1) } : e))
      )
    } catch (err) {
      setError(err.message || 'Faça login para confirmar presença.')
    } finally {
      setBusyId(null)
    }
  }

  // cria o evento na API e adiciona no topo da lista
  async function handleCreateEvent(data) {
    const newEvent = await createEvent(data)
    setEvents((prev) => [newEvent, ...prev])
  }

  // atualiza o evento na API e substitui ele na lista local
  async function handleUpdateEvent(id, data) {
    const updated = await updateEvent(id, data)
    setEvents((prev) => prev.map((e) => (e.id === id ? updated : e)))
  }

  // confirma com o usuario e remove o evento
  async function handleDeleteEvent(id) {
    if (!window.confirm('Tem certeza que deseja excluir este evento?')) return
    try {
      await deleteEvent(id)
      setEvents((prev) => prev.filter((e) => e.id !== id))
    } catch (err) {
      setError(err.message)
    }
  }

  const visible = filter === 'Todos' ? events : events.filter((e) => categoryFor(e.title) === filter)

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <div className="flex items-center justify-between gap-3">
        <div>
          <h2 className="flex items-center gap-2 text-xl font-bold text-slate-800">
            <CalendarDays className="h-5 w-5 text-emerald-600" />
            Agenda
          </h2>
          <p className="text-sm text-slate-400">Próximos eventos da escola</p>
        </div>
        {canEdit && !showForm && (
          <button
            onClick={() => setShowForm(true)}
            className="animate-pop-in flex items-center gap-1.5 rounded-full bg-emerald-500 px-4 py-2 text-sm font-semibold text-white shadow-md shadow-emerald-200 transition-all duration-200 hover:bg-emerald-600 hover:shadow-lg active:scale-95"
          >
            <Plus className="h-4 w-4" />
            Novo Evento
          </button>
        )}
      </div>

      {showForm && (
        <NewEventForm
          onCreate={handleCreateEvent}
          onClose={() => setShowForm(false)}
        />
      )}

      {error && <p className="rounded-2xl bg-rose-50 p-4 text-sm text-rose-600">{error}</p>}

      <div className="flex flex-wrap gap-2">
        {filters.map((f) => (
          <button
            key={f}
            onClick={() => setFilter(f)}
            className={`rounded-full px-4 py-1.5 text-sm font-medium transition-all duration-200 active:scale-95 ${
              filter === f
                ? 'bg-emerald-500 text-white shadow-sm'
                : 'bg-white text-slate-500 ring-1 ring-slate-200 hover:bg-slate-50'
            }`}
          >
            {f}
          </button>
        ))}
      </div>

      <div className="space-y-4">
        {loading ? (
          <p className="text-sm text-slate-400">Carregando...</p>
        ) : visible.length === 0 ? (
          <p className="animate-fade-in rounded-2xl bg-white p-6 text-center text-sm text-slate-400 ring-1 ring-slate-100">
            Nenhum evento nesta categoria.
          </p>
        ) : (
          visible.map((event, index) => {
            if (editingEventId === event.id) {
              return (
                <NewEventForm
                  key={event.id}
                  initial={event}
                  onCreate={(data) => handleUpdateEvent(event.id, data)}
                  onClose={() => setEditingEventId(null)}
                />
              )
            }
            return (
              <EventCard
                key={event.id}
                event={event}
                index={index}
                busy={busyId === event.id}
                onToggle={() => toggle(event.id)}
                canEdit={canEdit}
                onEdit={setEditingEventId}
                onDelete={handleDeleteEvent}
              />
            )
          })
        )}
      </div>
    </div>
  )
}
