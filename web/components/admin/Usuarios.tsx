"use client";

import { useState } from "react";
import { api, ErrorCliente } from "@/lib/api-cliente";
import type { Rol, Usuario, UsuarioConClave } from "@/lib/tipos";
import { CopiarTexto } from "../CopiarTexto";
import { SoloAdmin, useManejarError, useSesion } from "./Sesion";
import { Aviso, BTN_PRIMARIO, INPUT } from "./ui";
import { useApi } from "./useApi";

/** La clave temporal se muestra UNA vez: no se guarda en ningún lado del navegador. */
function ClaveUnaVez({ datos, onListo }: { datos: UsuarioConClave; onListo: () => void }) {
  return (
    <section className="space-y-3 rounded-[12px] border-2 border-negro bg-st-espera-bg p-4 text-st-espera-fg">
      <h2 className="text-[17px] font-extrabold">Clave de {datos.usuario.nombre}</h2>
      <p className="text-[15px]">
        Usuario: <strong>{datos.usuario.usuario}</strong>
      </p>
      <p id="clave-temporal" className="break-all rounded-[10px] bg-blanco p-3 font-mono text-[18px] font-bold text-negro">
        {datos.clave_temporal}
      </p>
      <p className="text-[14px] font-semibold">
        Cópiala y pásasela por un canal privado. <strong>No se vuelve a mostrar.</strong>
      </p>
      <div className="flex flex-wrap gap-2">
        <CopiarTexto texto={datos.clave_temporal} etiqueta="Copiar clave" idObjetivo="clave-temporal" className="border-negro bg-blanco text-negro" />
        <button type="button" onClick={onListo} className="min-h-11 rounded-[10px] bg-negro px-4 text-[14px] font-bold text-blanco">
          Ya la copié
        </button>
      </div>
    </section>
  );
}

function Contenido() {
  const { sesion } = useSesion();
  const lista = useApi<{ items: Usuario[] }>("/api/admin/usuarios");
  const manejar = useManejarError();
  const [clave, setClave] = useState<UsuarioConClave | null>(null);
  const [usuario, setUsuario] = useState("");
  const [nombre, setNombre] = useState("");
  const [rol, setRol] = useState<Rol>("vendedor");
  const [error, setError] = useState("");
  const [errores, setErrores] = useState<Record<string, string>>({});
  const [ocupado, setOcupado] = useState(false);

  async function crear(e: React.FormEvent) {
    e.preventDefault();
    const er: Record<string, string> = {};
    if (!/^[a-z][a-z0-9._-]{2,31}$/.test(usuario.trim())) er.usuario = "3 a 32 caracteres: minúsculas, números, punto, guion o guion bajo; empieza con letra.";
    if (nombre.trim().length < 2 || nombre.trim().length > 60) er.nombre = "Escribe el nombre (2 a 60 letras).";
    setErrores(er);
    if (Object.keys(er).length) return;
    setOcupado(true);
    setError("");
    try {
      const r = await api<UsuarioConClave>("/api/admin/usuarios", {
        method: "POST",
        body: { usuario: usuario.trim(), nombre: nombre.trim(), rol },
      });
      setClave(r);
      setUsuario("");
      setNombre("");
      lista.recargar();
    } catch (err) {
      if (err instanceof ErrorCliente && err.campos) setErrores(err.campos);
      setError(err instanceof ErrorCliente && err.codigo === "duplicado" ? "Ya existe un usuario con ese nombre de usuario." : manejar(err));
    } finally {
      setOcupado(false);
    }
  }

  async function activo(u: Usuario, valor: boolean) {
    setError("");
    try {
      await api<Usuario>(`/api/admin/usuarios/${u.id}`, { method: "PATCH", body: { activo: valor } });
      lista.recargar();
    } catch (err) {
      setError(
        err instanceof ErrorCliente && err.codigo === "no_permitido_sobre_si_mismo" ? "No puedes desactivarte a ti mismo." : manejar(err),
      );
    }
  }

  async function nuevaClave(u: Usuario) {
    setError("");
    try {
      setClave(await api<UsuarioConClave>(`/api/admin/usuarios/${u.id}/clave`, { method: "POST" }));
    } catch (err) {
      setError(manejar(err));
    }
  }

  return (
    <div className="space-y-4">
      <h1 className="titular text-[28px]">Usuarios</h1>
      {clave ? <ClaveUnaVez datos={clave} onListo={() => setClave(null)} /> : null}
      {error ? <Aviso>{error}</Aviso> : null}

      <form onSubmit={crear} noValidate className="space-y-3 rounded-[12px] bg-blanco p-4">
        <h2 className="text-[17px] font-extrabold">+ Agregar usuario</h2>
        <div>
          <label htmlFor="u-nombre" className="text-[14px] font-bold">
            Nombre
          </label>
          <input id="u-nombre" value={nombre} onChange={(e) => setNombre(e.target.value)} maxLength={60} placeholder="Ej: Vane (regiones)" className={INPUT} />
          {errores.nombre ? <p className="mt-1 text-[13px] font-semibold text-rojo">{errores.nombre}</p> : null}
        </div>
        <div>
          <label htmlFor="u-usuario" className="text-[14px] font-bold">
            Usuario
          </label>
          <input
            id="u-usuario"
            value={usuario}
            onChange={(e) => setUsuario(e.target.value.toLowerCase())}
            maxLength={32}
            autoCapitalize="none"
            autoComplete="off"
            placeholder="Ej: vane"
            className={INPUT}
          />
          {errores.usuario ? <p className="mt-1 text-[13px] font-semibold text-rojo">{errores.usuario}</p> : null}
        </div>
        <div>
          <label htmlFor="u-rol" className="text-[14px] font-bold">
            Rol
          </label>
          <select id="u-rol" value={rol} onChange={(e) => setRol(e.target.value === "admin" ? "admin" : "vendedor")} className={INPUT}>
            <option value="vendedor">Vendedor (pedidos y clientes)</option>
            <option value="admin">Administrador (todo)</option>
          </select>
        </div>
        <button type="submit" disabled={ocupado} className={BTN_PRIMARIO}>
          {ocupado ? "Creando…" : "Crear usuario y generar clave"}
        </button>
      </form>

      {lista.error ? <Aviso>{lista.error}</Aviso> : null}
      <ul className="space-y-2">
        {lista.data?.items.map((u) => (
          <li key={u.id} className="rounded-[12px] bg-blanco p-4">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <p className="text-[16px] font-bold">
                {u.nombre} <span className="font-normal text-gris">· {u.usuario}</span>
              </p>
              <span className="text-[13px] font-semibold">
                {u.rol === "admin" ? "Administrador" : "Vendedor"} · {u.activo ? "activo" : "desactivado"}
              </span>
            </div>
            {u.id !== sesion.id ? (
              <div className="mt-2 flex flex-wrap gap-2">
                <button
                  type="button"
                  onClick={() => activo(u, !u.activo)}
                  className="min-h-11 rounded-[10px] border-2 border-negro px-3 text-[14px] font-bold"
                >
                  {u.activo ? "Desactivar" : "Activar"}
                </button>
                <button type="button" onClick={() => nuevaClave(u)} className="min-h-11 rounded-[10px] border-2 border-negro px-3 text-[14px] font-bold">
                  Generar clave nueva
                </button>
              </div>
            ) : (
              <p className="mt-1 text-[13px] text-gris">Eres tú.</p>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
}

export function Usuarios() {
  return (
    <SoloAdmin>
      <Contenido />
    </SoloAdmin>
  );
}
