"use client";

import Image from "next/image";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { createContext, useCallback, useContext, useEffect, useState } from "react";
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
        if (e.status === 0) return "Parece que no tienes internet. Lo que escribiste sigue aquí.";
        if (e.codigo === "conflicto") return "Alguien más cambió este pedido. Recarga para ver la última versión.";
        return e.message;
      }
      return "No pudimos guardar. Revisa tu conexión y vuelve a intentar.";
    },
    [expiro],
  );
}

type Estado = { tipo: "cargando" } | { tipo: "fuera" } | { tipo: "dentro"; sesion: SesionPanel };

export function AdminShell({ children }: { children: React.ReactNode }) {
  const [estado, setEstado] = useState<Estado>({ tipo: "cargando" });

  useEffect(() => {
    let vivo = true;
    api<SesionPanel>("/api/admin/sesion")
      .then((s) => vivo && setEstado({ tipo: "dentro", sesion: s }))
      .catch(() => vivo && setEstado({ tipo: "fuera" }));
    return () => {
      vivo = false;
    };
  }, []);

  const expiro = useCallback(() => setEstado({ tipo: "fuera" }), []);

  async function salir() {
    try {
      await api("/api/admin/logout", { method: "POST" });
    } finally {
      setEstado({ tipo: "fuera" });
    }
  }

  const sesion = estado.tipo === "dentro" ? estado.sesion : null;

  return (
    <div className="min-h-dvh bg-panel text-negro">
      <header className="sticky top-0 z-30 bg-negro text-blanco">
        <div className="mx-auto flex h-14 max-w-3xl items-center gap-3 px-4">
          <Image src="/marca/logo-donde-carlitos.png" alt="" width={36} height={36} className="size-9 rounded-full" sizes="36px" />
          <div className="min-w-0 leading-tight">
            <p className="titular text-[18px] text-oro">Pedidos del WhatsApp</p>
            {sesion ? (
              <p className="truncate text-[12px] text-gris-claro">
                Hola, {sesion.nombre} · {sesion.rol === "admin" ? "Administrador" : "Vendedor"}
              </p>
            ) : null}
          </div>
          {sesion ? (
            <button type="button" onClick={salir} className="ml-auto min-h-11 shrink-0 px-1 text-[14px] font-bold text-oro underline underline-offset-2">
              Cerrar sesión
            </button>
          ) : null}
        </div>
        {sesion ? <NavAdmin esAdmin={sesion.rol === "admin"} /> : null}
      </header>
      <main className="mx-auto max-w-3xl px-4 pb-28 pt-4">
        {estado.tipo === "cargando" ? <div className="h-40 animate-pulse rounded-[12px] bg-blanco" aria-hidden /> : null}
        {estado.tipo === "fuera" ? <Login onEntrar={(s) => setEstado({ tipo: "dentro", sesion: s })} /> : null}
        {sesion ? (
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
        ]
      : []),
  ];
  return (
    <nav aria-label="Panel" className="border-t border-oro/25">
      <ul className="scroll-fila mx-auto flex max-w-3xl gap-1 overflow-x-auto px-2">
        {tabs.map((t) => (
          <li key={t.href} className="shrink-0">
            <Link
              href={t.href}
              aria-current={t.activo ? "page" : undefined}
              className={`flex h-11 items-center px-3 text-[15px] font-bold ${t.activo ? "border-b-[3px] border-oro text-oro" : "text-blanco/80"}`}
            >
              {t.texto}
            </Link>
          </li>
        ))}
      </ul>
    </nav>
  );
}

function Login({ onEntrar }: { onEntrar: (s: SesionPanel) => void }) {
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
      onEntrar(await api<SesionPanel>("/api/admin/sesion"));
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

/** Para pantallas solo de admin: muestra un aviso al vendedor en vez del contenido. */
export function SoloAdmin({ children }: { children: React.ReactNode }) {
  const { esAdmin } = useSesion();
  if (!esAdmin) {
    return <p className="rounded-[12px] bg-blanco p-5 text-[15px]">Esta pantalla es solo para administradores.</p>;
  }
  return <>{children}</>;
}
