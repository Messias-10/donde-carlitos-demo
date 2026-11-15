"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useRef, useState } from "react";
import { api, ErrorCliente } from "@/lib/api-cliente";
import { clp } from "@/lib/formato";
import type { Canal, Cliente, Entrega, FormaPago, OpcionesEntrega, Pagina, Pedido, ProductoAdmin } from "@/lib/tipos";
import { IconoFlechaIzq } from "../Iconos";
import { useManejarError, useSesion } from "./Sesion";
import { Aviso, BTN_PRIMARIO, INPUT } from "./ui";
import { useApi } from "./useApi";

// Registro de un pedido mientras se atiende el WhatsApp (arquitectura §C2, copy P3).
// Orden pensado para el pulgar: canal, teléfono, nombre, productos, entrega, pago, total.

interface LineaCatalogo {
  tipo: "catalogo";
  clave: string;
  producto_id: number;
  codigo: string;
  descripcion: string;
  cantidad: string;
  /** Estimado para el total (el precio real lo congela la API). */
  precio: number;
  precio_mayor: number | null;
  mayor_minimo: number | null;
  /** Precio c/u escrito a mano (QA M-5); vacío = el del catálogo. */
  override: string;
}
interface LineaLibre {
  tipo: "libre";
  clave: string;
  descripcion: string;
  cantidad: string;
  precio: string;
}
type Linea = LineaCatalogo | LineaLibre;

let contador = 0;
const nuevaClave = () => `l${++contador}`;
const entero = (s: string): number | null => (/^\d{1,9}$/.test(s.trim()) ? Number(s.trim()) : null);
const CLAVE_CANAL = "dc_panel_canal";

/** Celular chileno: 9 dígitos que empiezan con 9, o 569 + 8. */
function telefonoValido(t: string): boolean {
  const s = t.trim();
  if (!/^\+?[\d\s().-]{8,20}$/.test(s)) return false;
  let d = s.replace(/\D/g, "");
  if (d.length === 11 && d.startsWith("56")) d = d.slice(2);
  return d.length === 9 && d.startsWith("9");
}

function precioLinea(l: Linea): number | null {
  const c = entero(l.cantidad) ?? 0;
  if (l.tipo === "libre") {
    if (!l.precio.trim()) return null;
    const p = entero(l.precio);
    return p === null ? null : p * c;
  }
  if (l.override.trim()) {
    const o = entero(l.override);
    return o === null ? null : o * c;
  }
  const unit = l.precio_mayor && l.mayor_minimo && c >= l.mayor_minimo && l.precio_mayor < l.precio ? l.precio_mayor : l.precio;
  return unit * c;
}

function canalGuardado(): Canal {
  try {
    return window.localStorage.getItem(CLAVE_CANAL) === "regiones" ? "regiones" : "santiago";
  } catch {
    return "santiago";
  }
}

function Segmento<T extends string>({
  valor,
  opciones,
  onCambio,
  etiqueta,
}: {
  valor: T;
  opciones: [T, string][];
  onCambio: (v: T) => void;
  etiqueta: string;
}) {
  return (
    <div role="radiogroup" aria-label={etiqueta} className="grid grid-cols-2 gap-2">
      {opciones.map(([v, t]) => (
        <button
          key={v}
          type="button"
          role="radio"
          aria-checked={valor === v}
          onClick={() => onCambio(v)}
          className={`min-h-12 rounded-[12px] border-2 px-3 text-[16px] font-bold ${
            valor === v ? "border-negro bg-negro text-blanco" : "border-negro/20 bg-blanco text-negro"
          }`}
        >
          {t}
        </button>
      ))}
    </div>
  );
}

function BuscadorProductos({ onElegir }: { onElegir: (p: ProductoAdmin) => void }) {
  const [q, setQ] = useState("");
  const [res, setRes] = useState<ProductoAdmin[]>([]);
  const [buscando, setBuscando] = useState(false);
  const manejar = useManejarError();
  const [error, setError] = useState("");

  useEffect(() => {
    const t = q.trim();
    if (t.length < 2) return;
    let vivo = true;
    const id = window.setTimeout(() => {
      setBuscando(true);
      api<Pagina<ProductoAdmin>>(`/api/admin/productos?q=${encodeURIComponent(t.slice(0, 60))}&page_size=8`)
        .then((r) => {
          if (!vivo) return;
          setRes(r.items.filter((p) => p.activo));
          setError("");
        })
        .catch((e: unknown) => {
          if (vivo) setError(manejar(e));
        })
        .finally(() => {
          if (vivo) setBuscando(false);
        });
    }, 250);
    return () => {
      vivo = false;
      window.clearTimeout(id);
    };
  }, [q, manejar]);

  const visibles = q.trim().length >= 2 ? res : [];

  return (
    <div>
      <label htmlFor="buscar-prod" className="text-[14px] font-bold">
        Busca por nombre o código
      </label>
      <input
        id="buscar-prod"
        type="search"
        value={q}
        onChange={(e) => setQ(e.target.value)}
        className={INPUT}
        placeholder="Ej: squishy, P001"
        autoComplete="off"
        enterKeyHint="search"
      />
      {error ? <p className="mt-1 text-[13px] font-semibold text-rojo">{error}</p> : null}
      {visibles.length > 0 ? (
        <ul className="mt-2 divide-y divide-negro/10 rounded-[10px] bg-blanco ring-1 ring-negro/10" aria-busy={buscando}>
          {visibles.map((p) => (
            <li key={p.id}>
              <button
                type="button"
                onClick={() => {
                  onElegir(p);
                  setQ("");
                  setRes([]);
                }}
                className="flex min-h-12 w-full items-center justify-between gap-3 px-3 py-2 text-left"
              >
                <span className="min-w-0">
                  <span className="block truncate text-[15px] font-semibold">{p.nombre}</span>
                  <span className="text-[13px] text-gris">
                    {p.codigo}
                    {p.precio_mayor && p.mayor_minimo ? ` · ${clp(p.precio_mayor)} desde ${p.mayor_minimo}` : ""}
                  </span>
                </span>
                <span className="shrink-0 font-bold">{clp(p.precio_final)}</span>
              </button>
            </li>
          ))}
        </ul>
      ) : q.trim().length >= 2 && !buscando ? (
        <p className="mt-1 text-[13px] text-gris">Sin resultados. Agrégalo como línea libre.</p>
      ) : null}
    </div>
  );
}

function desdePedido(p: Pedido): Linea[] {
  return p.items.map((it) =>
    it.producto_id
      ? {
          tipo: "catalogo" as const,
          clave: nuevaClave(),
          producto_id: it.producto_id,
          codigo: it.codigo ?? "",
          descripcion: it.descripcion,
          cantidad: String(it.cantidad),
          precio: it.precio_unitario ?? 0,
          precio_mayor: null,
          mayor_minimo: null,
          override: "",
        }
      : {
          tipo: "libre" as const,
          clave: nuevaClave(),
          descripcion: it.descripcion,
          cantidad: String(it.cantidad),
          precio: it.precio_unitario === null ? "" : String(it.precio_unitario),
        },
  );
}

export function FormPedido({ base }: { base?: Pedido }) {
  const router = useRouter();
  const manejar = useManejarError();
  const opciones = useApi<OpcionesEntrega>("/api/opciones-entrega");
  const editando = !!base;

  const { sesion } = useSesion();
  // QA M-15: el canal por defecto del usuario manda; si no tiene, el último usado en este equipo.
  const [canal, setCanal] = useState<Canal>(() => base?.canal ?? sesion.canal_por_defecto ?? canalGuardado());
  const [otraAgencia, setOtraAgencia] = useState("");
  const [telefono, setTelefono] = useState(base?.telefono ?? "");
  const [nombre, setNombre] = useState(base?.nombre ?? "");
  const [lineas, setLineas] = useState<Linea[]>(() => (base ? desdePedido(base) : []));
  const [descuento, setDescuento] = useState(base?.descuento ? String(base.descuento) : "");
  const [entrega, setEntrega] = useState<Entrega>(base?.entrega ?? (canal === "regiones" ? "envio" : "retiro"));
  const [comuna, setComuna] = useState(base?.comuna ?? "");
  const [region, setRegion] = useState(base?.region ?? "");
  const [agencia, setAgencia] = useState(base?.agencia ?? "");
  const [costoEnvio, setCostoEnvio] = useState(base?.costo_envio ? String(base.costo_envio) : "");
  const [porPagar, setPorPagar] = useState(base?.envio_por_pagar ?? false);
  const [formaPago, setFormaPago] = useState<FormaPago>(base?.forma_pago ?? (canal === "regiones" ? "transferencia" : "en_tienda"));
  const [totalManual, setTotalManual] = useState(base?.total_manual ? String(base.total) : "");
  const [nota, setNota] = useState(base?.nota_interna ?? "");
  const [cliente, setCliente] = useState<Cliente | null>(null);
  const [errores, setErrores] = useState<Record<string, string>>({});
  const [error, setError] = useState("");
  const [enviando, setEnviando] = useState(false);
  const telVisto = useRef("");

  function cambiarCanal(c: Canal) {
    setCanal(c);
    try {
      window.localStorage.setItem(CLAVE_CANAL, c);
    } catch {
      /* preferencia opcional */
    }
    // Valores por defecto del canal (mensajes-whatsapp.md §4).
    if (!editando) {
      setEntrega(c === "regiones" ? "envio" : "retiro");
      setFormaPago(c === "regiones" ? "transferencia" : "en_tienda");
    }
  }

  async function buscarCliente() {
    const t = telefono.trim();
    if (!telefonoValido(t) || telVisto.current === t) return;
    telVisto.current = t;
    try {
      const c = await api<Cliente>(`/api/admin/clientes?telefono=${encodeURIComponent(t)}`);
      setCliente(c.total_pedidos > 0 ? c : null);
    } catch {
      setCliente(null);
    }
  }

  function usarCliente() {
    if (!cliente) return;
    if (cliente.nombre && !nombre.trim()) setNombre(cliente.nombre);
    const ultima = cliente.pedidos.find((p) => p.comuna);
    if (ultima?.comuna && !comuna.trim()) setComuna(ultima.comuna);
    if (ultima?.region && !region) setRegion(ultima.region);
  }

  function agregarProducto(p: ProductoAdmin) {
    setLineas((ls) => {
      const i = ls.findIndex((l) => l.tipo === "catalogo" && l.producto_id === p.id);
      if (i >= 0) {
        const l = ls[i];
        const n = [...ls];
        n[i] = { ...l, cantidad: String((entero(l.cantidad) ?? 0) + 1) };
        return n;
      }
      return [
        ...ls,
        {
          tipo: "catalogo",
          clave: nuevaClave(),
          producto_id: p.id,
          codigo: p.codigo,
          descripcion: p.nombre,
          cantidad: "1",
          precio: p.precio_final,
          precio_mayor: p.precio_mayor,
          mayor_minimo: p.mayor_minimo,
          override: "",
        },
      ];
    });
    setErrores((e) => ({ ...e, items: "" }));
  }

  const setLinea = (clave: string, cambios: { cantidad?: string; descripcion?: string; precio?: string; override?: string }) =>
    setLineas((ls) => ls.map((l) => (l.clave === clave ? ({ ...l, ...cambios } as Linea) : l)));

  const calculo = useMemo(() => {
    let sub = 0;
    let sinPrecio = false;
    for (const l of lineas) {
      const v = precioLinea(l);
      if (v === null) sinPrecio = true;
      else sub += v;
    }
    const total = sub - (entero(descuento) ?? 0) + (entrega === "envio" ? (entero(costoEnvio) ?? 0) : 0);
    return { sub, total, sinPrecio };
  }, [lineas, descuento, costoEnvio, entrega]);

  function validar(): Record<string, string> {
    const e: Record<string, string> = {};
    if (!telefonoValido(telefono)) e.telefono = "Revisa el teléfono: tiene que ser un celular chileno (9 dígitos).";
    if (nombre.trim() && !/^[\p{L}\p{M} '’.-]{2,80}$/u.test(nombre.trim())) e.nombre = "Usa solo letras en el nombre.";
    if (lineas.length === 0) e.items = "Agrega al menos un producto.";
    for (const l of lineas) {
      const c = entero(l.cantidad);
      if (l.tipo === "catalogo" && l.override.trim() && entero(l.override) === null) e.items = "El precio tiene que ser un número sin puntos.";
      if (!c || c < 1) e.items = "La cantidad tiene que ser 1 o más.";
      if (l.tipo === "libre" && l.descripcion.trim().length < 2) e.items = "Escribe la descripción de la línea libre.";
      if (l.tipo === "libre" && l.precio.trim() && entero(l.precio) === null) e.items = "El precio tiene que ser un número sin puntos.";
    }
    if (descuento.trim() && entero(descuento) === null) e.descuento = "Escribe el descuento en pesos (solo números).";
    if (entrega === "envio" && costoEnvio.trim() && entero(costoEnvio) === null) e.costo_envio = "Escribe el costo en pesos (solo números).";
    if (calculo.sinPrecio && entero(totalManual) === null) e.total = "Una línea va sin precio: escribe tú el total.";
    return e;
  }

  async function guardar(ev: React.FormEvent) {
    ev.preventDefault();
    const e = validar();
    setErrores(e);
    if (Object.keys(e).length) {
      requestAnimationFrame(() => document.getElementById(`np-${Object.keys(e)[0]}`)?.focus());
      return;
    }
    const cuerpo: Record<string, unknown> = {
      canal,
      telefono: telefono.trim(),
      entrega,
      forma_pago: entrega === "envio" ? "transferencia" : formaPago,
      items: lineas.map((l) =>
        l.tipo === "catalogo"
          ? {
              producto_id: l.producto_id,
              cantidad: entero(l.cantidad),
              ...(l.override.trim() ? { precio_unitario: entero(l.override) } : {}),
            }
          : {
              descripcion: l.descripcion.trim(),
              cantidad: entero(l.cantidad),
              ...(l.precio.trim() ? { precio_unitario: entero(l.precio) } : {}),
            },
      ),
    };
    if (nombre.trim()) cuerpo.nombre = nombre.trim().replace(/\s+/g, " ");
    if (entero(descuento)) cuerpo.descuento = entero(descuento);
    if (entrega === "envio") {
      if (comuna.trim()) cuerpo.comuna = comuna.trim();
      if (region) cuerpo.region = region;
      if (agencia) cuerpo.agencia = agencia;
      if (entero(costoEnvio)) cuerpo.costo_envio = entero(costoEnvio);
      else if (porPagar) cuerpo.envio_por_pagar = true;
    }
    // Si todas las líneas tienen precio, el total lo calcula la API.
    if (calculo.sinPrecio) cuerpo.total = entero(totalManual);
    // QA M-10: la agencia "Otra" se anota en la nota interna (la API no tiene campo para el nombre).
    const notaFinal = [entrega === "envio" && agencia === "otra" && otraAgencia.trim() ? `Agencia: ${otraAgencia.trim()}.` : "", nota.trim()]
      .filter(Boolean)
      .join(" ");
    if (notaFinal) cuerpo.nota_interna = notaFinal.slice(0, 1000);

    setEnviando(true);
    setError("");
    try {
      const p = editando
        ? await api<Pedido>(`/api/admin/pedidos/${base!.id}`, { method: "PUT", body: cuerpo })
        : await api<Pedido>("/api/admin/pedidos", { method: "POST", body: cuerpo });
      router.replace(`/admin/pedidos/${p.id}${editando ? "" : "?nuevo=1"}`);
    } catch (err) {
      setEnviando(false);
      if (err instanceof ErrorCliente && err.campos) {
        const c: Record<string, string> = {};
        for (const [k, v] of Object.entries(err.campos)) c[k.startsWith("items") ? "items" : k] = v;
        if (c.total) c.total = "El total no calza con los productos. Revísalo o deja una línea sin precio.";
        setErrores(c);
      }
      setError(manejar(err));
    }
  }

  const err = (k: string) =>
    errores[k] ? (
      <p id={`npe-${k}`} className="mt-1 text-[13px] font-semibold text-rojo">
        {errores[k]}
      </p>
    ) : null;
  const agencias = opciones.data?.agencias ?? [];

  return (
    <form onSubmit={guardar} noValidate className="space-y-4">
      {editando ? (
        <Link href={`/admin/pedidos/${base.id}`} className="inline-flex min-h-11 items-center gap-1 text-[14px] font-bold">
          <IconoFlechaIzq size={18} /> Volver al pedido
        </Link>
      ) : null}
      <h1 className="titular text-[28px]">{editando ? `Editar ${base.numero}` : "Nuevo pedido"}</h1>

      <section className="space-y-3 rounded-[12px] bg-blanco p-4">
        <Segmento
          etiqueta="Canal"
          valor={canal}
          onCambio={cambiarCanal}
          opciones={[
            ["santiago", "Santiago"],
            ["regiones", "Regiones"],
          ]}
        />
        <div>
          <label htmlFor="np-telefono" className="text-[14px] font-bold">
            Teléfono del cliente
          </label>
          <input
            id="np-telefono"
            type="tel"
            inputMode="tel"
            value={telefono}
            onChange={(e) => {
              setTelefono(e.target.value);
              // QA M-1: si cambia el teléfono, el "cliente conocido" anterior ya no aplica.
              setCliente(null);
              telVisto.current = "";
            }}
            onBlur={buscarCliente}
            placeholder="Pega el número del chat"
            maxLength={20}
            aria-invalid={errores.telefono ? true : undefined}
            className={INPUT}
          />
          {err("telefono")}
          {cliente ? (
            <div className="mt-2 flex flex-wrap items-center gap-2 rounded-[10px] bg-st-prep-bg p-3 text-[14px] text-st-prep-fg">
              <span>
                Cliente conocido: <strong>{cliente.nombre || "sin nombre"}</strong> · {cliente.total_pedidos}{" "}
                {cliente.total_pedidos === 1 ? "pedido" : "pedidos"}
                {cliente.pedidos.find((p) => p.comuna)?.comuna ? ` · última comuna: ${cliente.pedidos.find((p) => p.comuna)!.comuna}` : ""}
              </span>
              <button type="button" onClick={usarCliente} className="min-h-11 rounded-[10px] bg-negro px-3 font-bold text-blanco">
                Usar estos datos
              </button>
            </div>
          ) : null}
        </div>
        <div>
          <label htmlFor="np-nombre" className="text-[14px] font-bold">
            Nombre
          </label>
          <input
            id="np-nombre"
            value={nombre}
            onChange={(e) => setNombre(e.target.value)}
            placeholder="Ej: Pedro Soto"
            maxLength={80}
            autoComplete="off"
            aria-invalid={errores.nombre ? true : undefined}
            className={INPUT}
          />
          {err("nombre")}
        </div>
      </section>

      <section className="space-y-3 rounded-[12px] bg-blanco p-4">
        <h2 className="text-[17px] font-extrabold">Productos</h2>
        <BuscadorProductos onElegir={agregarProducto} />
        {lineas.length > 0 ? (
          <ul className="space-y-2" id="np-items" tabIndex={-1}>
            {lineas.map((l) => {
              const v = precioLinea(l);
              return (
                <li key={l.clave} className="rounded-[10px] bg-panel p-3">
                  <div className="flex gap-2">
                    <div className="w-20 shrink-0">
                      <label htmlFor={`c-${l.clave}`} className="text-[12px] font-bold">
                        Cantidad
                      </label>
                      <input
                        id={`c-${l.clave}`}
                        inputMode="numeric"
                        value={l.cantidad}
                        onChange={(e) => setLinea(l.clave, { cantidad: e.target.value })}
                        className={INPUT}
                      />
                    </div>
                    <div className="min-w-0 flex-1">
                      {l.tipo === "catalogo" ? (
                        <>
                          <p className="text-[12px] font-bold">{l.codigo}</p>
                          <p className="mt-1 text-[15px] font-semibold leading-snug">{l.descripcion}</p>
                          <p className="text-[13px] text-gris">
                            {v === null ? "" : `${clp(v)} ${l.override.trim() ? "(precio escrito a mano)" : "aprox. (precio del catálogo)"}`}
                          </p>
                          <label htmlFor={`o-${l.clave}`} className="mt-2 block text-[12px] font-bold">
                            Precio c/u (opcional, reemplaza el del catálogo)
                          </label>
                          <input
                            id={`o-${l.clave}`}
                            inputMode="numeric"
                            value={l.override}
                            onChange={(e) => setLinea(l.clave, { override: e.target.value })}
                            placeholder={String(l.precio)}
                            className={INPUT}
                          />
                        </>
                      ) : (
                        <>
                          <label htmlFor={`d-${l.clave}`} className="text-[12px] font-bold">
                            Descripción
                          </label>
                          <input
                            id={`d-${l.clave}`}
                            value={l.descripcion}
                            onChange={(e) => setLinea(l.clave, { descripcion: e.target.value })}
                            maxLength={120}
                            className={INPUT}
                          />
                          <label htmlFor={`p-${l.clave}`} className="mt-2 block text-[12px] font-bold">
                            Precio c/u (opcional)
                          </label>
                          <input
                            id={`p-${l.clave}`}
                            inputMode="numeric"
                            value={l.precio}
                            onChange={(e) => setLinea(l.clave, { precio: e.target.value })}
                            className={INPUT}
                          />
                        </>
                      )}
                    </div>
                  </div>
                  <button
                    type="button"
                    onClick={() => setLineas((ls) => ls.filter((x) => x.clave !== l.clave))}
                    className="mt-1 min-h-11 text-[14px] font-bold text-rojo underline underline-offset-2"
                  >
                    Quitar
                  </button>
                </li>
              );
            })}
          </ul>
        ) : null}
        <button
          type="button"
          onClick={() => setLineas((ls) => [...ls, { tipo: "libre", clave: nuevaClave(), descripcion: "", cantidad: "1", precio: "" }])}
          className="min-h-11 rounded-[10px] border-2 border-negro px-4 text-[15px] font-bold"
        >
          + Línea libre
        </button>
        {err("items")}
        <div>
          <label htmlFor="np-descuento" className="text-[14px] font-bold">
            Descuento (opcional)
          </label>
          <input
            id="np-descuento"
            inputMode="numeric"
            value={descuento}
            onChange={(e) => setDescuento(e.target.value)}
            className={INPUT}
          />
          {err("descuento")}
        </div>
      </section>

      <section className="space-y-3 rounded-[12px] bg-blanco p-4">
        <h2 className="text-[17px] font-extrabold">Entrega y pago</h2>
        <Segmento
          etiqueta="Entrega"
          valor={entrega}
          onCambio={(v) => {
            setEntrega(v);
            if (v === "envio") setFormaPago("transferencia");
          }}
          opciones={[
            ["envio", "Despacho"],
            ["retiro", "Retiro en tienda"],
          ]}
        />
        {entrega === "envio" ? (
          <>
            <div>
              <label htmlFor="np-comuna" className="text-[14px] font-bold">
                Comuna
              </label>
              <input id="np-comuna" value={comuna} onChange={(e) => setComuna(e.target.value)} placeholder="Ej: Chillán" maxLength={60} className={INPUT} />
              {err("comuna")}
            </div>
            <div>
              <label htmlFor="np-region" className="text-[14px] font-bold">
                Región (opcional)
              </label>
              <select id="np-region" value={region} onChange={(e) => setRegion(e.target.value)} className={INPUT}>
                <option value="">Sin región</option>
                {opciones.data?.regiones.map((r) => (
                  <option key={r.codigo} value={r.codigo}>
                    {r.nombre}
                  </option>
                ))}
              </select>
            </div>
            <fieldset>
              <legend className="text-[14px] font-bold">
                Agencia {opciones.data?.agencias_es_ejemplo ? <span className="font-normal text-gris">(lista de ejemplo)</span> : null}
              </legend>
              <div className="mt-1 flex flex-wrap gap-2">
                {agencias.map((a) => (
                  <button
                    key={a.codigo}
                    type="button"
                    aria-pressed={agencia === a.codigo}
                    onClick={() => setAgencia(agencia === a.codigo ? "" : a.codigo)}
                    className={`min-h-11 rounded-[10px] border-2 px-3 text-[14px] font-bold ${
                      agencia === a.codigo ? "border-negro bg-negro text-blanco" : "border-negro/20 bg-blanco"
                    }`}
                  >
                    {a.nombre}
                  </button>
                ))}
              </div>
              {agencia === "otra" ? (
                <div className="mt-2">
                  <label htmlFor="np-otra" className="text-[14px] font-bold">
                    ¿Cuál?
                  </label>
                  <input id="np-otra" value={otraAgencia} onChange={(e) => setOtraAgencia(e.target.value)} maxLength={60} className={INPUT} />
                </div>
              ) : null}
              {err("agencia")}
            </fieldset>
            <div>
              <label htmlFor="np-costo_envio" className="text-[14px] font-bold">
                Costo de envío (opcional)
              </label>
              <input
                id="np-costo_envio"
                inputMode="numeric"
                value={costoEnvio}
                onChange={(e) => setCostoEnvio(e.target.value)}
                className={INPUT}
              />
              {err("costo_envio")}
            </div>
            {!costoEnvio.trim() ? (
              <label className="flex min-h-11 items-center gap-3 text-[15px] font-semibold">
                <input type="checkbox" checked={porPagar} onChange={(e) => setPorPagar(e.target.checked)} className="size-6 accent-negro" />
                Envío por pagar
              </label>
            ) : null}
            <p className="text-[13px] text-gris">Con despacho, el pago es por transferencia.</p>
          </>
        ) : (
          <Segmento
            etiqueta="Pago"
            valor={formaPago}
            onCambio={setFormaPago}
            opciones={[
              ["transferencia", "Transferencia"],
              ["en_tienda", "En tienda"],
            ]}
          />
        )}
      </section>

      <section className="space-y-3 rounded-[12px] bg-blanco p-4">
        {calculo.sinPrecio ? (
          <div>
            <label htmlFor="np-total" className="text-[14px] font-bold">
              Total
            </label>
            <input
              id="np-total"
              inputMode="numeric"
              value={totalManual}
              onChange={(e) => setTotalManual(e.target.value)}
              aria-invalid={errores.total ? true : undefined}
              className={INPUT}
            />
            <p className="mt-1 text-[13px] text-gris">Se calcula solo. Si una línea va sin precio, escríbelo tú.</p>
            {err("total")}
          </div>
        ) : (
          <div className="flex items-end justify-between">
            <p className="text-[15px] font-bold">
              Total <span className="block text-[12px] font-normal text-gris">Se calcula solo (la tienda confirma al guardar).</span>
            </p>
            <p className="font-display text-[30px] leading-none">{clp(Math.max(0, calculo.total))}</p>
          </div>
        )}
        {errores.total && !calculo.sinPrecio ? err("total") : null}
        <div>
          <label htmlFor="np-nota" className="text-[14px] font-bold">
            Nota interna (no la ve el cliente)
          </label>
          <textarea id="np-nota" value={nota} onChange={(e) => setNota(e.target.value)} maxLength={1000} className={`${INPUT} h-20 py-2`} />
        </div>
      </section>

      {error ? <Aviso>{error}</Aviso> : null}

      <button type="submit" disabled={enviando} className={`${BTN_PRIMARIO} sticky bottom-3 shadow-lg`}>
        {enviando ? "Guardando…" : editando ? "Guardar cambios" : "Guardar como cotizado"}
      </button>
    </form>
  );
}
