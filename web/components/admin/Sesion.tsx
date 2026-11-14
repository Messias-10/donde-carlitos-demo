"use client";

import Image from "next/image";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { createContext, useCallback, useContext, useEffect, useState } from "react";
import { api, ErrorCliente } from "@/lib/api-cliente";

interface Ctx {
  usuario: string;
  /** Llamar cuando la API responde 401: vuelve al login. */
  expiro: () => void;
}

const SesionCtx = createContext<Ctx | null>(null);

export function useSesion(): Ctx {
  const c = useContext(SesionCtx);
  if (!c) throw new Error("useSesion fuera de AdminShell");
  return c;
}

/** Maneja errores de la API en el panel: 401 manda al login. */
export function useManejarError() {
  const { expiro } = useSesion();
  return useCallback(
    (e: unknown): string => {
      if (e instanceof ErrorCliente) {
        if (e.status === 401) {
          expiro();
          return "Tu sesión terminó. Vuelve a entrar.";
        }
        return e.message;
      }
      return "Ocurrió un error. Intenta de nuevo en un rato.";
    },
    [expiro],
  );
}

type Estado = { tipo: "cargando" } | { tipo: "fuera" } | { tipo: "dentro"; usuario: string };

export function AdminShell({ children }: { children: React.ReactNode }) {
  const [estado, setEstado] = useState<Estado>({ tipo: "cargando" });

  useEffect(() => {
    let vivo = true;
    api<{ usuario: string }>("/api/admin/sesion")
      .then((s) => vivo && setEstado({ tipo: "dentro", usuario: s.usuario }))
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

  return (
    <div className="min-h-dvh bg-crema">
      <header className="sticky top-0 z-30 bg-negro text-blanco">
        <div className="mx-auto flex h-14 max-w-3xl items-center gap-3 px-4">
          <Image src="/marca/logo-donde-carlitos.png" alt="" width={40} height={40} className="size-10 rounded-full" sizes="40px" />
          <span className="titular text-[18px] text-oro">Panel de pedidos</span>
          {estado.tipo === "dentro" ? (
            <button type="button" onClick={salir} className="ml-auto min-h-11 px-2 text-[14px] font-bold text-oro underline underline-offset-2">
              Cerrar sesión
            </button>
          ) : null}
        </div>
        {estado.tipo === "dentro" ? <NavAdmin /> : null}
      </header>
      <main className="mx-auto max-w-3xl px-4 pb-16 pt-4">
        {estado.tipo === "cargando" ? <div className="h-40 animate-pulse rounded-[12px] bg-blanco/70" aria-hidden /> : null}
        {estado.tipo === "fuera" ? <Login onEntrar={(u) => setEstado({ tipo: "dentro", usuario: u })} /> : null}
        {estado.tipo === "dentro" ? (
          <SesionCtx.Provider value={{ usuario: estado.usuario, expiro }}>{children}</SesionCtx.Provider>
        ) : null}
      </main>
    </div>
  );
}

function NavAdmin() {
  const ruta = usePathname();
  const tabs = [
    { href: "/admin", texto: "Pedidos", activo: ruta === "/admin" || ruta.startsWith("/admin/pedidos") },
    { href: "/admin/productos", texto: "Productos", activo: ruta.startsWith("/admin/productos") },
  ];
  return (
    <nav aria-label="Panel" className="border-t border-oro/25">
      <ul className="mx-auto flex max-w-3xl px-2">
        {tabs.map((t) => (
          <li key={t.href} className="flex-1">
            <Link
              href={t.href}
              aria-current={t.activo ? "page" : undefined}
              className={`flex h-11 items-center justify-center text-[15px] font-bold ${t.activo ? "border-b-[3px] border-oro text-oro" : "text-blanco/80"}`}
            >
              {t.texto}
            </Link>
          </li>
        ))}
      </ul>
    </nav>
  );
}

function Login({ onEntrar }: { onEntrar: (usuario: string) => void }) {
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
      const r = await api<{ usuario: string }>("/api/admin/login", { method: "POST", body: { usuario: usuario.trim(), password } });
      setPassword("");
      onEntrar(r.usuario);
    } catch (err) {
      if (err instanceof ErrorCliente) {
        if (err.status === 401) setError("Usuario o contraseña incorrectos.");
        else if (err.status === 429) setError("Demasiados intentos. Espera unos minutos y vuelve a probar.");
        else if (err.codigo === "admin_no_configurado") setError("El panel aún no está configurado.");
        else setError(err.message);
      } else {
        setError("Ocurrió un error. Intenta de nuevo en un rato.");
      }
    } finally {
      setEnviando(false);
    }
  }

  return (
    <form onSubmit={entrar} className="mx-auto mt-6 max-w-sm space-y-4 rounded-[16px] bg-blanco p-5" noValidate>
      <h1 className="titular text-[30px]">Panel de pedidos</h1>
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
        className="flex h-12 w-full items-center justify-center rounded-[12px] bg-negro text-[16px] font-extrabold text-oro disabled:opacity-60"
      >
        {enviando ? "Entrando…" : "Entrar"}
      </button>
    </form>
  );
}
