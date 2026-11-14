"use client";

import { useCallback, useEffect, useState } from "react";
import { api } from "@/lib/api-cliente";
import { useManejarError } from "./Sesion";

interface Resultado<T> {
  clave: string;
  data?: T;
  error?: string;
}

/**
 * GET a la API del panel. El estado "cargando" se deriva (no hay setState
 * síncrono dentro del efecto). `recargar` vuelve a pedir lo mismo.
 */
export function useApi<T>(url: string) {
  const manejar = useManejarError();
  const [vuelta, setVuelta] = useState(0);
  const [res, setRes] = useState<Resultado<T> | null>(null);
  const clave = `${vuelta}|${url}`;

  useEffect(() => {
    let vivo = true;
    api<T>(url)
      .then((data) => vivo && setRes({ clave, data }))
      .catch((e: unknown) => vivo && setRes({ clave, error: manejar(e) }));
    return () => {
      vivo = false;
    };
  }, [url, clave, manejar]);

  const recargar = useCallback(() => setVuelta((v) => v + 1), []);
  const reemplazar = useCallback((data: T) => setRes((r) => ({ clave: r?.clave ?? clave, data })), [clave]);

  return {
    data: res?.data,
    error: res?.clave === clave ? res.error : undefined,
    cargando: res?.clave !== clave,
    recargar,
    reemplazar,
  };
}
