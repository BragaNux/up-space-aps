import { useState, useEffect } from 'react'
import { Megaphone, Paperclip, ChevronDown, AlertTriangle, Info, Pin, Plus, Pencil, Trash2 } from 'lucide-react'
import { getAnnouncements, markAnnouncementRead, createAnnouncement, updateAnnouncement, deleteAnnouncement } from './api'

const priorityStyles = {
  Urgente: { chip: 'bg-rose-50 text-rose-600', icon: AlertTriangle, dot: 'bg-rose-500' },
  Importante: { chip: 'bg-amber-50 text-amber-600', icon: Pin, dot: 'bg-amber-500' },
  Informativo: { chip: 'bg-indigo-50 text-indigo-600', icon: Info, dot: 'bg-indigo-500' },
}

// formata data + hora curta (ex: 17 jun · 14:30)
function formatDate(dateStr) {
  const date = new Date(dateStr)
  return `${date.toLocaleDateString('pt-BR', { day: '2-digit', month: 'short' })} · ${date.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })}`
}

// card de um comunicado, expansivel, com indicador de nao-lido e botoes de editar/excluir se puder
function ComunicadoCard({ item, expanded, onToggle, index, canEdit, onEdit, onDelete }) {
  const p = priorityStyles[item.priority] ?? priorityStyles.Informativo
  const PIcon = p.icon
  return (
    <article
      className={`animate-slide-up stagger-${Math.min(index + 1, 6)} overflow-hidden rounded-2xl bg-white shadow-sm ring-1 transition-all duration-200 hover:shadow-md ${
        item.read ? 'ring-slate-100' : 'ring-emerald-100'
      }`}
    >
      <div className="flex w-full items-start gap-3 p-4">
        <button onClick={onToggle} className="flex flex-1 items-start gap-3 text-left">
          <span className="mt-1.5 flex h-2.5 w-2.5 shrink-0 items-center justify-center">
            {!item.read && <span className={`h-2.5 w-2.5 animate-pulse rounded-full ${p.dot}`} />}
          </span>

          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-2">
              <span className={`inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-semibold ${p.chip}`}>
                <PIcon className="h-3.5 w-3.5" />
                {item.priority}
              </span>
              <span className="text-xs text-slate-400">{item.sender} · {formatDate(item.created_at)}</span>
            </div>

            <h4 className={`mt-1.5 text-base leading-snug ${item.read ? 'font-semibold text-slate-700' : 'font-bold text-slate-900'}`}>
              {item.title}
            </h4>

            {!expanded && <p className="mt-0.5 line-clamp-2 text-sm text-slate-500">{item.preview}</p>}
          </div>
        </button>

        <div className="flex items-center gap-1">
          {canEdit && (
            <>
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
            </>
          )}
          <button onClick={onToggle} aria-label="Expandir comunicado" className="grid h-8 w-8 place-items-center rounded-full hover:bg-slate-50">
            <ChevronDown className={`h-5 w-5 text-slate-400 transition ${expanded ? 'rotate-180' : ''}`} />
          </button>
        </div>
      </div>

      {expanded && (
        <div className="animate-slide-up px-4 pb-4 pl-[2.6rem]">
          <p className="whitespace-pre-line text-sm leading-relaxed text-slate-600">{item.body}</p>

          {item.attachment_name && (
            <button className="mt-4 inline-flex items-center gap-2 rounded-lg border border-slate-200 px-3 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50">
              <Paperclip className="h-4 w-4 text-emerald-600" />
              {item.attachment_name}
            </button>
          )}
        </div>
      )}
    </article>
  )
}

// formulario de criar/editar comunicado (initial preenchido = modo edicao)
function NewAnnouncementForm({ onCreate, onClose, initial }) {
  const [title, setTitle] = useState(initial?.title ?? '')
  const [sender, setSender] = useState(initial?.sender ?? '')
  const [priority, setPriority] = useState(initial?.priority ?? 'Informativo')
  const [preview, setPreview] = useState(initial?.preview ?? '')
  const [body, setBody] = useState(initial?.body ?? '')
  const [attachmentName, setAttachmentName] = useState(initial?.attachment_name ?? '')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  // valida os campos obrigatorios e salva o comunicado
  async function submit(e) {
    e.preventDefault()
    if (!title.trim() || !sender.trim() || !priority.trim() || !preview.trim() || !body.trim()) {
      setError('Preencha todos os campos obrigatórios.')
      return
    }
    setSaving(true)
    setError('')
    try {
      await onCreate({
        title,
        sender,
        priority,
        preview,
        body,
        attachment_name: attachmentName.trim() || null,
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
      <h3 className="font-bold text-slate-800">{initial ? 'Editar comunicado' : 'Novo comunicado'}</h3>
      
      <div>
        <label className="mb-1 block text-xs font-medium text-slate-500">Título</label>
        <input
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Título do comunicado"
          className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
        />
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="mb-1 block text-xs font-medium text-slate-500">Remetente (ex: Direção, Profª Ana)</label>
          <input
            value={sender}
            onChange={(e) => setSender(e.target.value)}
            placeholder="Remetente"
            className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
          />
        </div>
        <div>
          <label className="mb-1 block text-xs font-medium text-slate-500">Prioridade</label>
          <select
            value={priority}
            onChange={(e) => setPriority(e.target.value)}
            className="w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none focus:border-emerald-400"
          >
            <option value="Urgente">Urgente</option>
            <option value="Importante">Importante</option>
            <option value="Informativo">Informativo</option>
          </select>
        </div>
      </div>

      <div>
        <label className="mb-1 block text-xs font-medium text-slate-500">Resumo (visualizado antes de abrir)</label>
        <input
          value={preview}
          onChange={(e) => setPreview(e.target.value)}
          placeholder="Breve resumo informativo"
          className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
        />
      </div>

      <div>
        <label className="mb-1 block text-xs font-medium text-slate-500">Corpo da mensagem</label>
        <textarea
          value={body}
          onChange={(e) => setBody(e.target.value)}
          placeholder="Mensagem detalhada..."
          rows={4}
          className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
        />
      </div>

      <div>
        <label className="mb-1 block text-xs font-medium text-slate-500">Nome do anexo (opcional)</label>
        <input
          value={attachmentName}
          onChange={(e) => setAttachmentName(e.target.value)}
          placeholder="Ex: regulamento_festa.pdf"
          className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
        />
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

// tela de comunicados: lista, marca como lido, filtra nao-lidos e (se profissional) gerencia os avisos
export default function Comunicados({ canEdit }) {
  const [comunicados, setComunicados] = useState([])
  const [openId, setOpenId] = useState(null)
  const [onlyUnread, setOnlyUnread] = useState(false)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [editingId, setEditingId] = useState(null)

  useEffect(() => {
    getAnnouncements()
      .then(setComunicados)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  const unreadCount = comunicados.filter((c) => !c.read).length

  // abre/fecha o comunicado e marca como lido na primeira vez que abre
  function toggle(id) {
    setOpenId((cur) => (cur === id ? null : id))
    const item = comunicados.find((c) => c.id === id)
    if (item && !item.read) {
      markAnnouncementRead(id)
        .then(() => setComunicados((prev) => prev.map((c) => (c.id === id ? { ...c, read: true } : c))))
        .catch(() => {})
    }
  }

  // cria o comunicado na API e adiciona no topo da lista
  async function handleCreateAnnouncement(data) {
    const created = await createAnnouncement(data)
    setComunicados((prev) => [created, ...prev])
  }

  // atualiza o comunicado na API e substitui ele na lista local
  async function handleUpdateAnnouncement(id, data) {
    const updated = await updateAnnouncement(id, data)
    setComunicados((prev) => prev.map((c) => (c.id === id ? updated : c)))
  }

  // confirma com o usuario e remove o comunicado
  async function handleDeleteAnnouncement(id) {
    if (!window.confirm('Tem certeza que deseja remover este comunicado?')) return
    try {
      await deleteAnnouncement(id)
      setComunicados((prev) => prev.filter((c) => c.id !== id))
    } catch (err) {
      setError(err.message)
    }
  }

  const visible = onlyUnread ? comunicados.filter((c) => !c.read) : comunicados

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <div className="flex items-center justify-between gap-4">
        <div>
          <h2 className="flex items-center gap-2 text-xl font-bold text-slate-800">
            <Megaphone className="h-5 w-5 text-emerald-600" />
            Comunicados
          </h2>
          <p className="text-sm text-slate-400">Avisos oficiais da escola</p>
        </div>
        <div className="flex items-center gap-2">
          {unreadCount > 0 && (
            <span className="rounded-full bg-emerald-50 px-3 py-1 text-xs font-semibold text-emerald-600">
              {unreadCount} não {unreadCount === 1 ? 'lido' : 'lidos'}
            </span>
          )}
          {canEdit && !showForm && (
            <button
              onClick={() => setShowForm(true)}
              className="animate-pop-in flex items-center gap-1.5 rounded-full bg-emerald-500 px-4 py-2 text-sm font-semibold text-white shadow-md shadow-emerald-200 transition-all duration-200 hover:bg-emerald-600 active:scale-95"
            >
              <Plus className="h-4 w-4" />
              Novo Comunicado
            </button>
          )}
        </div>
      </div>

      {showForm && (
        <NewAnnouncementForm
          onCreate={handleCreateAnnouncement}
          onClose={() => setShowForm(false)}
        />
      )}

      {error && <p className="rounded-2xl bg-rose-50 p-4 text-sm text-rose-600">{error}</p>}

      <div className="flex gap-2">
        <button
          onClick={() => setOnlyUnread(false)}
          className={`rounded-full px-4 py-1.5 text-sm font-medium transition-all duration-200 active:scale-95 ${
            !onlyUnread ? 'bg-emerald-500 text-white shadow-sm' : 'bg-white text-slate-500 ring-1 ring-slate-200 hover:bg-slate-50'
          }`}
        >
          Todos
        </button>
        <button
          onClick={() => setOnlyUnread(true)}
          className={`rounded-full px-4 py-1.5 text-sm font-medium transition-all duration-200 active:scale-95 ${
            onlyUnread ? 'bg-emerald-500 text-white shadow-sm' : 'bg-white text-slate-500 ring-1 ring-slate-200 hover:bg-slate-50'
          }`}
        >
          Não lidos
        </button>
      </div>

      <div className="space-y-4">
        {loading ? (
          <p className="text-sm text-slate-400">Carregando...</p>
        ) : visible.length === 0 ? (
          <p className="animate-fade-in rounded-2xl bg-white p-6 text-center text-sm text-slate-400 ring-1 ring-slate-100">
            Nenhum comunicado por aqui. 🎉
          </p>
        ) : (
          visible.map((item, index) => {
            if (editingId === item.id) {
              return (
                <NewAnnouncementForm
                  key={item.id}
                  initial={item}
                  onCreate={(data) => handleUpdateAnnouncement(item.id, data)}
                  onClose={() => setEditingId(null)}
                />
              )
            }
            return (
              <ComunicadoCard
                key={item.id}
                item={item}
                index={index}
                expanded={openId === item.id}
                onToggle={() => toggle(item.id)}
                canEdit={canEdit}
                onEdit={() => setEditingId(item.id)}
                onDelete={() => handleDeleteAnnouncement(item.id)}
              />
            )
          })
        )}
      </div>
    </div>
  )
}
