import { useState, useEffect, useRef } from 'react'
import {
  Mail,
  Phone,
  MapPin,
  Bell,
  Globe,
  ShieldCheck,
  HelpCircle,
  FileText,
  LogOut,
  ChevronRight,
  Camera,
  Pencil,
} from 'lucide-react'
import defaultAvatar from '../assets/default.png'
import { getMyChildren, updateMe, fileToBase64 } from './api'

const MAX_PHOTO_BYTES = 4 * 1024 * 1024

// linha simples de exibicao "icone + label + valor" (email, telefone, endereco)
function InfoRow({ icon: Icon, label, value }) {
  return (
    <div className="flex items-center gap-3 py-3">
      <span className="grid h-9 w-9 shrink-0 place-items-center rounded-lg bg-slate-50 text-slate-500">
        <Icon className="h-4 w-4" />
      </span>
      <div className="min-w-0">
        <p className="text-xs text-slate-400">{label}</p>
        <p className="truncate text-sm font-medium text-slate-700">{value || '—'}</p>
      </div>
    </div>
  )
}

// linha clicavel de uma opcao de configuracao, com icone e acao opcional (ex: toggle)
function SettingRow({ icon: Icon, label, action, onClick }) {
  return (
    <div
      onClick={onClick}
      className={`flex w-full items-center gap-3 px-4 py-3.5 text-left transition ${onClick ? 'cursor-pointer hover:bg-slate-50' : ''}`}
    >
      <span className="grid h-9 w-9 shrink-0 place-items-center rounded-lg bg-emerald-50 text-emerald-600">
        <Icon className="h-4 w-4" />
      </span>
      <span className="flex-1 text-sm font-medium text-slate-700">{label}</span>
      {action ?? <ChevronRight className="h-4 w-4 text-slate-300" />}
    </div>
  )
}

// switch on/off estilizado, usado nas preferencias
function Toggle({ on, onChange }) {
  return (
    <button
      onClick={onChange}
      role="switch"
      aria-checked={on}
      className={`relative h-6 w-11 shrink-0 rounded-full transition ${on ? 'bg-emerald-500' : 'bg-slate-200'}`}
    >
      <span className={`absolute top-0.5 h-5 w-5 rounded-full bg-white shadow transition ${on ? 'left-[1.375rem]' : 'left-0.5'}`} />
    </button>
  )
}

const roleLabels = { profissional: 'Profissional', responsavel: 'Responsável' }

// tela de perfil do usuario logado: dados pessoais, foto, filhos vinculados e preferencias
export default function Perfil({ user, onLogout, onOpenChild, onUserUpdated }) {
  const [notif, setNotif] = useState(true)
  const [children, setChildren] = useState([])
  const [uploading, setUploading] = useState(false)
  const [photoError, setPhotoError] = useState('')
  const [editing, setEditing] = useState(false)
  const [form, setForm] = useState({ name: '', phone: '', address: '' })
  const [saving, setSaving] = useState(false)
  const [editError, setEditError] = useState('')
  const fileInputRef = useRef(null)

  // preenche o formulario com os dados atuais e entra no modo de edicao
  function startEditing() {
    setForm({
      name: user.name || '',
      phone: user.phone || '',
      address: user.address || '',
    })
    setEditError('')
    setEditing(true)
  }

  // valida o nome e salva nome/telefone/endereco atualizados
  async function handleSave(e) {
    e.preventDefault()
    if (!form.name.trim()) {
      setEditError('O nome não pode estar vazio.')
      return
    }
    setSaving(true)
    setEditError('')
    try {
      const updated = await updateMe({
        name: form.name.trim(),
        phone: form.phone.trim(),
        address: form.address.trim(),
        avatar_url: user.avatar_url,
      })
      onUserUpdated?.(updated)
      setEditing(false)
    } catch (err) {
      setEditError(err.message || 'Não foi possível salvar os dados.')
    } finally {
      setSaving(false)
    }
  }

  useEffect(() => {
    if (user?.role === 'responsavel') {
      getMyChildren().then(setChildren).catch(() => {})
    }
  }, [user])

  // valida tipo/tamanho da foto escolhida e ja salva ela no perfil (converte pra base64 antes)
  async function handlePhotoChange(e) {
    const file = e.target.files?.[0]
    if (!file) return
    setPhotoError('')
    if (!file.type.startsWith('image/')) {
      setPhotoError('Selecione um arquivo de imagem.')
      return
    }
    if (file.size > MAX_PHOTO_BYTES) {
      setPhotoError('Imagem muito grande (máximo 4MB).')
      return
    }
    setUploading(true)
    try {
      const base64 = await fileToBase64(file)
      const updated = await updateMe({ name: user.name, phone: user.phone, address: user.address, avatar_url: base64 })
      onUserUpdated?.(updated)
    } catch (err) {
      setPhotoError(err.message || 'Não foi possível atualizar a foto.')
    } finally {
      setUploading(false)
      e.target.value = ''
    }
  }

  if (!user) return null

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <h2 className="text-xl font-bold text-slate-800">Perfil</h2>

      <section className="animate-slide-up rounded-3xl bg-white p-6 shadow-sm ring-1 ring-slate-100">
        <div className="flex items-start justify-between gap-4">
          <div className="flex items-center gap-4">
            <button
              type="button"
              onClick={() => fileInputRef.current?.click()}
              disabled={uploading}
              className="group relative shrink-0 rounded-full"
              aria-label="Alterar foto de perfil"
            >
              <img
                src={user.avatar_url || defaultAvatar}
                alt={user.name}
                onError={(e) => { e.currentTarget.onerror = null; e.currentTarget.src = defaultAvatar }}
                className="h-20 w-20 rounded-full object-cover ring-4 ring-emerald-50 transition-transform duration-300 group-hover:scale-105"
              />
              <span className="absolute inset-0 grid place-items-center rounded-full bg-slate-900/0 text-white opacity-0 transition-all duration-200 group-hover:bg-slate-900/40 group-hover:opacity-100">
                {uploading ? (
                  <span className="h-5 w-5 animate-spin rounded-full border-2 border-white/40 border-t-white" />
                ) : (
                  <Camera className="h-5 w-5" />
                )}
              </span>
              <input ref={fileInputRef} type="file" accept="image/*" onChange={handlePhotoChange} className="hidden" />
            </button>
            <div className="min-w-0 flex-1">
              <h3 className="text-lg font-extrabold text-slate-800">{user.name}</h3>
              <span className="mt-1 inline-block rounded-full bg-emerald-50 px-2.5 py-0.5 text-xs font-semibold text-emerald-600">
                {roleLabels[user.role] ?? user.role}
              </span>
            </div>
          </div>
          {!editing && (
            <button
              onClick={startEditing}
              className="flex items-center gap-1 text-sm font-medium text-emerald-600 hover:underline"
            >
              <Pencil className="h-3.5 w-3.5" />
              Editar
            </button>
          )}
        </div>
        {photoError && <p className="mt-2 text-xs text-rose-500">{photoError}</p>}

        {editing ? (
          <form onSubmit={handleSave} className="animate-slide-up mt-4 space-y-3">
            <div>
              <label className="mb-1 block text-xs font-medium text-slate-500">Nome completo</label>
              <input
                value={form.name}
                onChange={(e) => setForm((prev) => ({ ...prev, name: e.target.value }))}
                placeholder="Seu nome completo"
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
              />
            </div>

            <div>
              <label className="mb-1 block text-xs font-medium text-slate-500">Telefone</label>
              <input
                value={form.phone}
                onChange={(e) => setForm((prev) => ({ ...prev, phone: e.target.value }))}
                placeholder="Telefone"
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
              />
            </div>

            <div>
              <label className="mb-1 block text-xs font-medium text-slate-500">Endereço</label>
              <input
                value={form.address}
                onChange={(e) => setForm((prev) => ({ ...prev, address: e.target.value }))}
                placeholder="Seu endereço"
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
              />
            </div>

            {editError && <p className="text-xs text-rose-500">{editError}</p>}

            <div className="flex justify-end gap-2 pt-2">
              <button
                type="button"
                onClick={() => setEditing(false)}
                className="rounded-full px-3 py-1.5 text-sm text-slate-500 hover:bg-slate-100"
              >
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
        ) : (
          <>
            <hr className="my-2 border-slate-100" />
            <InfoRow icon={Mail} label="E-mail" value={user.email} />
            <InfoRow icon={Phone} label="Telefone" value={user.phone} />
            <InfoRow icon={MapPin} label="Endereço" value={user.address} />
          </>
        )}
      </section>

      {user.role === 'responsavel' && (
        <section>
          <h3 className="mb-3 text-sm font-bold uppercase tracking-wide text-slate-400">Filhos vinculados</h3>
          {children.length === 0 ? (
            <p className="rounded-2xl bg-white p-4 text-sm text-slate-400 shadow-sm ring-1 ring-slate-100">
              Nenhuma criança vinculada à sua conta.
            </p>
          ) : (
            <div className="space-y-3">
              {children.map((child, i) => (
                <button
                  key={child.id}
                  onClick={onOpenChild}
                  className={`animate-pop-in stagger-${Math.min(i + 1, 6)} flex w-full items-center gap-3 rounded-2xl bg-white p-3 text-left shadow-sm ring-1 ring-slate-100 transition-all duration-200 hover:-translate-y-0.5 hover:shadow-md hover:ring-emerald-200`}
                >
                  <img
                    src={child.photo_url || defaultAvatar}
                    alt={child.name}
                    onError={(e) => { e.currentTarget.onerror = null; e.currentTarget.src = defaultAvatar }}
                    className="h-12 w-12 rounded-xl object-cover"
                  />
                  <div className="flex-1">
                    <p className="font-semibold text-slate-800">{child.name}</p>
                    <p className="text-xs text-slate-400">{child.group_name}</p>
                  </div>
                  <ChevronRight className="h-4 w-4 text-slate-300" />
                </button>
              ))}
            </div>
          )}
        </section>
      )}

      <section>
        <h3 className="mb-3 text-sm font-bold uppercase tracking-wide text-slate-400">Preferências</h3>
        <div className="divide-y divide-slate-100 overflow-hidden rounded-2xl bg-white shadow-sm ring-1 ring-slate-100">
          <SettingRow icon={Bell} label="Notificações" action={<Toggle on={notif} onChange={() => setNotif((v) => !v)} />} />
          <SettingRow icon={Globe} label="Idioma" action={<span className="text-sm text-slate-400">Português</span>} />
          <SettingRow icon={ShieldCheck} label="Privacidade e segurança" />
        </div>
      </section>

      <section>
        <h3 className="mb-3 text-sm font-bold uppercase tracking-wide text-slate-400">Conta</h3>
        <div className="divide-y divide-slate-100 overflow-hidden rounded-2xl bg-white shadow-sm ring-1 ring-slate-100">
          <SettingRow icon={HelpCircle} label="Central de ajuda" />
          <SettingRow icon={FileText} label="Termos e políticas" />
        </div>

        <button
          onClick={onLogout}
          className="mt-4 flex w-full items-center justify-center gap-2 rounded-2xl border border-rose-200 py-3.5 text-sm font-semibold text-rose-600 transition-all duration-200 hover:bg-rose-50 active:scale-[0.98]"
        >
          <LogOut className="h-4 w-4" />
          Sair da conta
        </button>
      </section>

      <p className="pb-2 text-center text-xs text-slate-300">Up - Espaço</p>
    </div>
  )
}
