import { useState } from 'react'
import { ArrowLeft, Eye, EyeOff, GraduationCap, Heart, AlertCircle } from 'lucide-react'
import logoUpEspaco from '../assets/logo_up_espaco.png'
import { login as loginRequest, setToken } from './api'

const roles = [
  {
    id: 'profissional',
    label: 'Profissional',
    desc: 'Professoras e equipe da escola',
    icon: GraduationCap,
  },
  {
    id: 'responsavel',
    label: 'Responsável',
    desc: 'Mãe, pai ou responsável pelo aluno',
    icon: Heart,
  },
]

// tela de login: escolhe o papel (profissional/responsavel) e depois pede email/senha
export default function Login({ onLogin, onForgot, onCreateAccount }) {
  const [role, setRole] = useState(null)
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  // faz login na API e confere se o papel da conta bate com o que foi selecionado na tela
  async function handleSubmit(e) {
    e.preventDefault()
    setError('')
    if (!email || !password) {
      setError('Preencha e-mail e senha.')
      return
    }
    setLoading(true)
    try {
      const { token, user } = await loginRequest(email, password)
      if (user.role !== role) {
        const actualLabel = roles.find((r) => r.id === user.role)?.label ?? user.role
        setError(`Esta conta é de ${actualLabel}. Volte e selecione a opção correta.`)
        return
      }
      setToken(token)
      onLogin(user)
    } catch (err) {
      setError(err.message || 'Não foi possível entrar. Verifique suas credenciais.')
    } finally {
      setLoading(false)
    }
  }

  const selectedRole = roles.find((r) => r.id === role)

  return (
    <div className="flex min-h-screen w-full items-center justify-center overflow-y-auto bg-white p-4 sm:p-6">
      <div className="relative grid w-full max-w-4xl overflow-hidden rounded-[28px] bg-white shadow-2xl shadow-slate-200/60 ring-1 ring-slate-100 md:grid-cols-2 animate-slide-up">
        {/* Painel marca */}
        <div className="relative flex flex-col bg-slate-50 p-6 md:p-8">
          <div className="flex flex-1 items-center justify-center py-2">
            <div className="h-64 w-64 sm:h-72 sm:w-72 md:h-72 md:w-72 transition-transform duration-700 ease-out hover:scale-[1.03]">
              <img src={logoUpEspaco} alt="Up Espaço" className="h-full w-full object-contain" />
            </div>
          </div>
        </div>

        {/* Painel formulário */}
        <div className="relative flex flex-col justify-center overflow-hidden px-7 py-9 sm:px-10">
          {!role ? (
            <div key="role-select" className="animate-slide-in-left">
              <h1 className="mb-2 text-center text-3xl font-extrabold text-slate-800">
                Login
              </h1>
              <p className="mb-7 text-center text-sm text-slate-500">
                Selecione como você vai entrar na Up Espaço.
              </p>

              <div className="space-y-3">
                {roles.map(({ id, label, desc, icon: Icon }, i) => (
                  <button
                    key={id}
                    onClick={() => setRole(id)}
                    className={`animate-pop-in stagger-${i + 1} group flex w-full items-center gap-4 rounded-xl border border-slate-200 bg-white p-4 text-left shadow-sm transition-all duration-200 hover:-translate-y-0.5 hover:border-emerald-400 hover:bg-emerald-50 hover:shadow-md active:scale-[0.98]`}
                  >
                    <span className="grid h-12 w-12 shrink-0 place-items-center rounded-xl bg-slate-100 text-slate-700 transition-colors duration-200 group-hover:bg-emerald-500 group-hover:text-white">
                      <Icon className="h-6 w-6" strokeWidth={1.8} />
                    </span>
                    <div>
                      <p className="font-semibold text-slate-800">{label}</p>
                      <p className="text-xs text-slate-500">{desc}</p>
                    </div>
                  </button>
                ))}
              </div>

              <p className="mt-7 text-center text-sm text-slate-500">
                Não tem conta?{' '}
                <button onClick={onCreateAccount} className="font-semibold text-[#0F5CC0] hover:underline">
                  Criar conta
                </button>
              </p>
            </div>
          ) : (
            <div key="credentials" className="animate-slide-in-right">
              <button
                onClick={() => { setRole(null); setError('') }}
                className="mb-6 inline-flex items-center gap-1.5 self-start text-sm font-medium text-slate-500 transition hover:-translate-x-0.5 hover:text-slate-800"
              >
                <ArrowLeft className="h-4 w-4" />
                Voltar
              </button>

              <h1 className="mb-1 text-center text-3xl font-extrabold text-slate-800">
                Login
              </h1>

              {selectedRole && (
                <div className="animate-pop-in mx-auto mb-6 flex items-center gap-2 rounded-full bg-emerald-50 px-4 py-2 text-emerald-700 ring-1 ring-emerald-100">
                  <selectedRole.icon className="h-4 w-4" />
                  <span className="text-sm font-semibold">Entrando como {selectedRole.label}</span>
                </div>
              )}

              <form onSubmit={handleSubmit} noValidate className="space-y-5">
                {/* E-mail */}
                <div>
                  <label className="mb-1.5 block text-sm font-semibold text-slate-700">
                    E-mail
                  </label>
                  <div className="relative">
                    <input
                      type="email"
                      placeholder="seu@email.com"
                      value={email}
                      onChange={(e) => { setEmail(e.target.value); setError('') }}
                      className="w-full rounded-lg border border-slate-300 bg-white py-3 px-4 text-sm text-slate-800 placeholder-slate-300 shadow-sm outline-none transition-all duration-200 focus:border-emerald-400 focus:ring-2 focus:ring-emerald-100"
                    />
                  </div>
                </div>

                {/* Senha */}
                <div>
                  <label className="mb-1.5 block text-sm font-semibold text-slate-700">
                    Senha
                  </label>
                  <div className="relative">
                    <button
                      type="button"
                      onClick={() => setShowPassword((v) => !v)}
                      className="absolute right-3.5 top-1/2 -translate-y-1/2 text-slate-400 transition hover:text-emerald-600"
                      aria-label={showPassword ? 'Ocultar senha' : 'Mostrar senha'}
                    >
                      {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                    </button>
                    <input
                      type={showPassword ? 'text' : 'password'}
                      placeholder="••••••••"
                      value={password}
                      onChange={(e) => { setPassword(e.target.value); setError('') }}
                      className="w-full rounded-lg border border-slate-300 bg-white py-3 pl-4 pr-10 text-sm text-slate-800 placeholder-slate-300 shadow-sm outline-none transition-all duration-200 focus:border-emerald-400 focus:ring-2 focus:ring-emerald-100"
                    />
                  </div>
                </div>

                {error && (
                  <p className="animate-pop-in flex items-start gap-1.5 rounded-lg bg-rose-50 px-3 py-2 text-xs text-rose-600">
                    <AlertCircle className="mt-0.5 h-3.5 w-3.5 shrink-0" />
                    {error}
                  </p>
                )}

                <div className="flex justify-end">
                  <button type="button" onClick={onForgot} className="text-xs font-semibold text-[#0F5CC0] transition hover:text-[#0F5CC0F] hover:underline">
                    Esqueceu a senha?
                  </button>
                </div>

                <button
                  type="submit"
                  disabled={loading}
                  className="w-full rounded-lg bg-emerald-500 py-3.5 text-sm font-semibold text-white shadow-md shadow-emerald-200 transition-all duration-200 hover:bg-emerald-600 hover:shadow-lg active:scale-[0.98] disabled:opacity-70"
                >
                  {loading ? (
                    <span className="flex items-center justify-center gap-2">
                      <span className="h-4 w-4 animate-spin rounded-full border-2 border-white/40 border-t-white" />
                      Entrando...
                    </span>
                  ) : (
                    'Entrar'
                  )}
                </button>
              </form>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
