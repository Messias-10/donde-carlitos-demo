"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { api, ErrorCliente } from "@/lib/api-cliente";
import type { Opcion, Pedido } from "@/lib/tipos";
import { useManejarError } from "./Sesion";
import { Aviso, BTN_PRIMARIO, BTN_SECUNDARIO, Dialogo, INPUT } from "./ui";

/**
 * "Despachar" desde Hoy (QA B-1): agencia, n.º de seguimiento y comuna, y pasa
 * directo a despachado (la API permite pagado → despachado y preparando → despachado).
 */
export function Despachar({
  id,
  numero,
  agencia: agenciaInicial,
  comuna: comunaInicial,
  agencias,
}: {
  id: number;
  numero: string;
  agencia: string | null;
  comuna: string | null;
  agencias: Opcion[];
}) {
  const router = useRouter();
  const manejar = useManejarError();
  const [abierto, setAbierto] = useState(false);
  const [agencia, setAgencia] = useState(agenciaInicial ?? "");
  const [nSeg, setNSeg] = useState("");
  const [comuna, setComuna] = useState(comunaInicial ?? "");
  const [error, setError] = useState("");
  const [guardando, setGuardando] = useState(false);

  async function despachar() {
    if (!agencia || !nSeg.trim() || !comuna.trim()) {
      setError("Para despachar necesitas la agencia, el n.º de seguimiento y la comuna.");
      return;
    }
    setGuardando(true);
    setError("");
    try {
      await api<Pedido>(`/api/admin/pedidos/${id}/estado`, {
        method: "PATCH",
        body: { estado: "despachado", agencia, n_seguimiento: nSeg.trim(), comuna: comuna.trim() },
      });
      setAbierto(false);
      // Hoy es un Server Component: se vuelve a pedir al servidor.
      router.push(`/admin/pedidos/${id}?despachado=1`);
    } catch (err) {
      setError(err instanceof ErrorCliente && err.campos ? Object.values(err.campos).join(" · ") : manejar(err));
    } finally {
      setGuardando(false);
    }
  }

  return (
    <>
      <button
        type="button"
        onClick={() => setAbierto(true)}
        className="inline-flex min-h-11 items-center rounded-[10px] bg-negro px-4 text-[14px] font-bold text-blanco"
      >
        Despachar
      </button>
      <Dialogo abierto={abierto} onCerrar={() => setAbierto(false)} titulo={`Despachar ${numero}`}>
        <div>
          <label htmlFor={`da-${id}`} className="text-[14px] font-bold">
            Agencia
          </label>
          <select id={`da-${id}`} value={agencia} onChange={(e) => setAgencia(e.target.value)} className={INPUT}>
            <option value="">Elige la agencia</option>
            {agencias.map((a) => (
              <option key={a.codigo} value={a.codigo}>
                {a.nombre}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label htmlFor={`dn-${id}`} className="text-[14px] font-bold">
            N.º de seguimiento
          </label>
          <input id={`dn-${id}`} value={nSeg} onChange={(e) => setNSeg(e.target.value)} maxLength={60} className={INPUT} />
        </div>
        <div>
          <label htmlFor={`dc-${id}`} className="text-[14px] font-bold">
            Comuna
          </label>
          <input id={`dc-${id}`} value={comuna} onChange={(e) => setComuna(e.target.value)} maxLength={60} className={INPUT} />
        </div>
        {error ? <Aviso>{error}</Aviso> : null}
        <div className="space-y-2">
          <button type="button" onClick={despachar} disabled={guardando} className={BTN_PRIMARIO}>
            {guardando ? "Guardando…" : "Guardar y marcar despachado"}
          </button>
          <button type="button" onClick={() => setAbierto(false)} className={BTN_SECUNDARIO}>
            Volver
          </button>
        </div>
      </Dialogo>
    </>
  );
}

/** "Actualizar" de Hoy: vuelve a pedir el Server Component. */
export function Actualizar() {
  const router = useRouter();
  return (
    <button type="button" onClick={() => router.refresh()} className="min-h-11 px-2 text-[14px] font-bold underline underline-offset-2">
      Actualizar
    </button>
  );
}
