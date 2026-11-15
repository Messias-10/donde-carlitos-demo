"use client";

import Image from "next/image";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { createContext, useCallback, useContext, useState } from "react";
import { api, ErrorCliente } from "@/lib/api-cliente";
import type { SesionPanel } from "@/lib/tipos";

interface Ctx {
  sesion: SesionPanel;
  esAdmin: boolean;
  /** Llamar cuando la API responde 401: vuelve al login. */
  expiro: () => void;
}

const SesionCtx = createContext<Ctx | null>(null);

export function useSesion(): Ctx {
  const c = useContext(SesionCtx);
  if (!c) throw new Error("useSesion fuera de AdminShell");
  return c;
}

/** Traduce errores de la API del panel a un texto para mostrar. 401 manda al login. */
export function useManejarError() {
  const { expiro } = useSesion();
  return useCallback(
    (e: unknown): string => {
      if (e instanceof ErrorCliente) {
        if (e.status === 401) {
          expiro();
          return "Tu sesión se cerró. Vuelve a entrar.";
        }
        if (e.status === 403 && e.codigo === "sin_permiso") return "Esto lo puede hacer solo un administrador.";
        if (e.status === 403 && e.codigo === "debe_cambiar_clave") {
          expiro();
          return "Primero cambia tu clave.";
        }
        if (e.status === 0) return "Parece que no tienes internet. Lo que escribiste sigue aquí.";
        if (e.codigo === "conflicto") return "Alguien más cambió este pedido. Recarga para ver la última versión.";
        return e.message;
      }
      return "No pudimos guardar. Revisa tu conexión y vuelve a intentar.";
    },
    [expiro],
  );
}

/**
 * La sesión se lee en el servidor (app/admin/layout.tsx) y llega como prop: el panel
 * no espera un fetch del navegador antes de pintar (QA I-8). Con una clave temporal
 * (debe_cambiar_clave) se muestra "Cambia tu clave" antes que nada (seguridad B9).
 */
export function AdminShell({ sesionInicial, children }: { sesionInicial: SesionPanel | null; children: React.ReactNode }) {
  const router = useRouter();
  const [sesion, setSesion] = useState<SesionPanel | null>(sesionInicial);

  const recargarSesion = useCallback(async () => {
    try {
      setSesion(await api<SesionPanel>("/api/admin/sesion"));
    } catch {
      setSesion(null);
    }
    router.refresh();
  }, [router]);

  const expiro = useCallback(() => {
    void recargarSesion();
  }, [recargarSesion]);

  async function salir() {
    try {
      await api("/api/admin/logout", { method: "POST" });
    } finally {
      setSesion(null);
      router.refresh();
    }
  }

  return (
    <div className="min-h-dvh bg-panel text-negro">
      <header className="z-30 bg-negro text-blanco lg:sticky lg:top-0">
        <div className="mx-auto flex h-14 max-w-3xl items-center gap-3 px-4">
          <Image src="/marca/logo-donde-carlitos.png" alt="" width={36} height={36} className="size-9 rounded-full" sizes="36px" />
          <p className="titular min-w-0 text-[18px] text-oro">Pedidos del WhatsApp</p>
          {sesion ? (
            <button type="button" onClick={salir} className="ml-auto min-h-11 shrink-0 px-1 text-[14px] font-bold text-oro underline underline-offset-2">
              Cerrar sesión
            </button>
          ) : null}
        </div>
        {sesion ? (
          <p className="mx-auto max-w-3xl px-4 pb-2 text-[13px] text-gris-claro">
            Hola, {sesion.nombre} · {sesion.rol === "admin" ? "Administrador" : "Vendedor"}
            {sesion.canal_por_defecto ? ` · Canal: ${sesion.canal_por_defecto === "regiones" ? "Regiones" : "Santiago"}` : ""}
          </p>
        ) : null}
        {sesion && !sesion.debe_cambiar_clave ? <NavAdmin esAdmin={sesion.rol === "admin"} /> : null}
      </header>
      <main className="mx-auto max-w-3xl px-4 pb-32 pt-4">
        {!sesion ? <Login onEntrar={recargarSesion} /> : null}
        {sesion && sesion.debe_cambiar_clave ? <CambiarClave onListo={recargarSesion} /> : null}
        {sesion && !sesion.debe_cambiar_clave ? (
          <SesionCtx.Provider value={{ sesion, esAdmin: sesion.rol === "admin", expiro }}>{children}</SesionCtx.Provider>
        ) : null}
      </main>
    </div>
  );
}

function NavAdmin({ esAdmin }: { esAdmin: boolean }) {
  const ruta = usePathname();
  const tabs = [
    { href: "/admin", texto: "Hoy", activo: ruta === "/admin" },
    { href: "/admin/pedidos/nuevo", texto: "Nuevo pedido", activo: ruta === "/admin/pedidos/nuevo" },
    { href: "/admin/buscar", texto: "Buscar", activo: ruta.startsWith("/admin/buscar") },
    ...(esAdmin
      ? [
          { href: "/admin/resumen", texto: "Resumen", activo: ruta.startsWith("/admin/resumen") },
          { href: "/admin/productos", texto: "Productos", activo: ruta.startsWith("/admin/productos") },
          { href: "/admin/usuarios", texto: "Usuarios", activo: ruta.startsWith("/admin/usuarios") },
          { href: "/admin/ajustes", texto: "Ajustes", activo: ruta.startsWith("/admin/ajustes") },
        ]
      : []),
  ];
  // Menú en dos filas si no cabe: nada queda escondido tras un scroll sin señal (QA M-8).
  return (
    <nav aria-label="Panel" className="border-t border-oro/25">
      <ul className="mx-auto flex max-w-3xl flex-wrap px-2">
        {tabs.map((t) => (
          <li key={t.href} className="shrink-0">
            <Link
              href={t.href}
              aria-current={t.activo ? "page" : undefined}
              className={`flex h-11 items-center px-2.5 text-[15px] font-bold ${t.activo ? "border-b-[3px] border-oro text-oro" : "text-blanco/80"}`}
            >
              {t.texto}
            </Link>
          </li>
        ))}
      </ul>
    </nav>
  );
}

function Login({ onEntrar }: { onEntrar: () => Promise<void> | void }) {
  const [usuario, setUsuario] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [enviando, setEnviando] = useState(false);

  async function entrar(e: React.FormEvent) {
    e.preventDefault();
    if (!usuario.trim() || !password) {
      setError("Escribe tu usuario y tu contraseña.");
      return;
    }
    setEnviando(true);
    setError("");
    try {
      await api("/api/admin/login", { method: "POST", body: { usuario: usuario.trim(), password } });
      setPassword("");
      await onEntrar();
    } catch (err) {
      if (err instanceof ErrorCliente) {
        if (err.status === 401) setError("Usuario o contraseña incorrectos.");
        else if (err.status === 429) setError("Demasiados intentos. Espera unos minutos y vuelve a probar.");
        else if (err.codigo === "admin_no_configurado") setError("El panel aún no está configurado.");
        else setError(err.message);
      } else {
        setError("No pudimos entrar. Revisa tu conexión y vuelve a intentar.");
      }
    } finally {
      setEnviando(false);
    }
  }

  return (
    <form onSubmit={entrar} className="mx-auto mt-6 max-w-sm space-y-4 rounded-[16px] bg-blanco p-5" noValidate>
      <h1 className="titular text-[30px]">Pedidos del WhatsApp</h1>
      {error ? (
        <p role="alert" className="rounded-[10px] border-2 border-rojo p-3 text-[15px] font-semibold text-rojo">
          {error}
        </p>
      ) : null}
      <div>
        <label htmlFor="a-usuario" className="text-[15px] font-bold">
          Usuario
        </label>
        <input
          id="a-usuario"
          value={usuario}
          onChange={(e) => setUsuario(e.target.value)}
          autoComplete="username"
          autoCapitalize="none"
          maxLength={64}
          className="mt-1 h-12 w-full rounded-[10px] border-2 border-negro/20 px-3 text-[16px] focus:border-negro"
        />
      </div>
      <div>
        <label htmlFor="a-password" className="text-[15px] font-bold">
          Contraseña
        </label>
        <input
          id="a-password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          autoComplete="current-password"
          maxLength={72}
          className="mt-1 h-12 w-full rounded-[10px] border-2 border-negro/20 px-3 text-[16px] focus:border-negro"
        />
      </div>
      <button
        type="submit"
        disabled={enviando}
        className="flex h-12 w-full items-center justify-center rounded-[12px] bg-negro text-[16px] font-extrabold text-blanco disabled:opacity-60"
      >
        {enviando ? "Entrando…" : "Entrar"}
      </button>
    </form>
  );
}

/** Seguridad B9: con una clave temporal, lo primero es cambiarla. */
function CambiarClave({ onListo }: { onListo: () => Promise<void> | void }) {
  const [actual, setActual] = useState("");
  const [nueva, setNueva] = useState("");
  const [repetir, setRepetir] = useState("");
  const [error, setError] = useState("");
  const [enviando, setEnviando] = useState(false);

  async function guardar(e: React.FormEvent) {
    e.preventDefault();
    if (nueva.length < 12) return setError("La clave nueva tiene que tener al menos 12 caracteres.");
    if (nueva.length > 72) return setError("La clave nueva puede tener hasta 72 caracteres.");
    if (nueva !== repetir) return setError("Las dos claves nuevas no son iguales.");
    if (nueva === actual) return setError("La clave nueva tiene que ser distinta de la actual.");
    setEnviando(true);
    setError("");
    try {
      await api("/api/admin/sesion/clave", { method: "POST", body: { clave_actual: actual, clave_nueva: nueva } });
      setActual("");
      setNueva("");
      setRepetir("");
      await onListo();
    } catch (err) {
      if (err instanceof ErrorCliente && err.campos) setError(Object.values(err.campos).join(" · "));
      else if (err instanceof ErrorCliente && (err.status === 401 || err.codigo === "credenciales_invalidas")) {
        setError("La clave actual no es correcta.");
      } else setError(err instanceof ErrorCliente ? err.message : "No pudimos guardar. Revisa tu conexión y vuelve a intentar.");
    } finally {
      setEnviando(false);
    }
  }

  const campo = "mt-1 h-12 w-full rounded-[10px] border-2 border-negro/20 px-3 text-[16px] focus:border-negro";
  return (
    <form onSubmit={guardar} className="mx-auto mt-6 max-w-sm space-y-4 rounded-[16px] bg-blanco p-5" noValidate>
      <h1 className="titular text-[28px]">Cambia tu clave</h1>
      <p className="text-[15px]">Entraste con una clave temporal. Elige una nueva que solo sepas tú (12 caracteres o más).</p>
      {error ? (
        <p role="alert" className="rounded-[10px] border-2 border-rojo p-3 text-[15px] font-semibold text-rojo">
          {error}
        </p>
      ) : null}
      <div>
        <label htmlFor="c-actual" className="text-[15px] font-bold">
          Clave actual (la temporal)
        </label>
        <input id="c-actual" type="password" value={actual} onChange={(e) => setActual(e.target.value)} autoComplete="current-password" maxLength={72} className={campo} />
      </div>
      <div>
        <label htmlFor="c-nueva" className="text-[15px] font-bold">
          Clave nueva
        </label>
        <input id="c-nueva" type="password" value={nueva} onChange={(e) => setNueva(e.target.value)} autoComplete="new-password" maxLength={72} className={campo} />
      </div>
      <div>
        <label htmlFor="c-repetir" className="text-[15px] font-bold">
          Repite la clave nueva
        </label>
        <input id="c-repetir" type="password" value={repetir} onChange={(e) => setRepetir(e.target.value)} autoComplete="new-password" maxLength={72} className={campo} />
      </div>
      <button
        type="submit"
        disabled={enviando}
        className="flex h-12 w-full items-center justify-center rounded-[12px] bg-negro text-[16px] font-extrabold text-blanco disabled:opacity-60"
      >
        {enviando ? "Guardando…" : "Guardar clave nueva"}
      </button>
    </form>
  );
}

/** Para pantallas solo de admin: muestra un aviso al vendedor en vez del contenido. */
export function SoloAdmin({ children }: { children: React.ReactNode }) {
  const { esAdmin } = useSesion();
  if (!esAdmin) {
    return <p className="rounded-[12px] bg-blanco p-5 text-[15px]">Esta pantalla es solo para administradores.</p>;
  }
  return <>{children}</>;
}
