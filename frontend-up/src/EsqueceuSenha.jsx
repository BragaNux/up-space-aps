import { useState } from 'react'
import { ArrowLeft, MailCheck } from 'lucide-react'
import logoUpEspaco from '../assets/logo_up_espaco.png'
import { forgotPassword } from './api'

// tela de recuperacao de senha: pede o email e pede pra API disparar o link de reset
export default function EsqueceuSenha({ onBack }) {
  const [email, setEmail] = useState('')
  const [loading, setLoading] = useState(false)
  const [sent, setSent] = useState(false)
  const [error, setError] = useState('')

  // valida o email e chama a API de recuperacao de senha
  async function handleSubmit(e) {
    e.preventDefault()
    setError('')
    if (!email) {
      setError('Informe seu e-mail.')
      return
    }
    if (!/^\S+@\S+\.\S+$/.test(email)) {
      setError('E-mail inválido.')
      return
    }
    setLoading(true)
    try {
      await forgotPassword(email)
      setSent(true)
    } catch (err) {
      setError(err.message || 'Não foi possível enviar o link de recuperação.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen w-full items-center justify-center overflow-y-auto bg-white p-4 sm:p-6">
      <div className="relative grid w-full max-w-4xl overflow-hidden rounded-[28px] bg-white shadow-2xl shadow-slate-200/60 ring-1 ring-slate-100 md:grid-cols-2">
        {/* Painel marca */}
        <div className="relative flex flex-col items-center justify-center bg-slate-50 p-6 md:p-8">
          <div className="h-64 w-64 sm:h-72 sm:w-72 md:h-72 md:w-72">
            <img src={logoUpEspaco} alt="Up Espaço" className="h-full w-full object-contain" />
          </div>
        </div>

        {/* Painel formulário */}
        <div className="flex flex-col justify-center px-7 py-9 sm:px-10">
          {sent ? (
            <div className="flex flex-col items-center text-center">
              <span className="grid h-16 w-16 place-items-center rounded-full bg-slate-100 text-slate-900">
                <MailCheck className="h-8 w-8" />
              </span>
              <h1 className="mt-5 text-2xl font-extrabold text-slate-800">
                E-mail enviado!
              </h1>
              <p className="mt-2 text-sm leading-relaxed text-slate-500">
                Enviamos um link de recuperação para{' '}
                <span className="font-semibold text-slate-700">{email}</span>. Verifique
                sua caixa de entrada e o spam.
              </p>

              <button
                onClick={onBack}
                className="mt-7 w-full rounded-lg bg-slate-900 py-3.5 text-sm font-semibold text-white shadow-md shadow-slate-300 transition hover:bg-slate-800 active:scale-[0.98]"
              >
                Voltar ao login
              </button>

              <button
                onClick={() => setSent(false)}
                className="mt-4 text-xs font-medium text-slate-500 transition hover:text-slate-800 hover:underline"
              >
                Não recebeu? Reenviar
              </button>
            </div>
          ) : (
            <>
              <button
                onClick={onBack}
                className="mb-6 inline-flex items-center gap-1.5 self-start text-sm font-medium text-slate-500 transition hover:text-slate-800"
              >
                <ArrowLeft className="h-4 w-4" />
                Voltar
              </button>

              <h1 className="text-3xl font-extrabold text-slate-800">
                Esqueceu a senha?
              </h1>
              <p className="mt-2 text-sm leading-relaxed text-slate-500">
                Sem problemas. Informe o e-mail cadastrado e enviaremos um link para você
                redefinir sua senha.
              </p>

              <form onSubmit={handleSubmit} noValidate className="mt-7 space-y-5">
                <div>
                  <label className="mb-1.5 block text-sm font-semibold text-slate-700">
                    E-mail
                  </label>
                  <input
                    type="email"
                    placeholder="seu@email.com"
                    value={email}
                    onChange={(e) => { setEmail(e.target.value); setError('') }}
                    className="w-full rounded-lg border border-black bg-white px-4 py-3 text-sm text-slate-800 placeholder-slate-300 shadow-sm outline-none transition focus:border-black focus:ring-2 focus:ring-slate-200"
                  />
                  {error && <p className="mt-1.5 text-xs text-red-500">{error}</p>}
                </div>

                <button
                  type="submit"
                  disabled={loading}
                  className="w-full rounded-lg bg-emerald-500 py-3.5 text-sm font-semibold text-white shadow-md shadow-emerald-200 transition hover:bg-emerald-600 active:scale-[0.98] disabled:opacity-70"
                >
                  {loading ? 'Enviando...' : 'Enviar link de recuperação'}
                </button>
              </form>

              <p className="mt-7 text-center text-sm text-slate-500">
                Lembrou a senha?{' '}
                <button onClick={onBack} className="font-semibold text-[#0F5CC0] hover:underline">
                  Entrar
                </button>
              </p>
            </>
          )}
        </div>
      </div>
    </div>
  )
}
