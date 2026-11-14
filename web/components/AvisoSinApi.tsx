import { NEGOCIO } from "@/lib/negocio";

/** Se muestra si el catálogo no responde: la página sigue siendo útil. */
export function AvisoSinApi() {
  return (
    <div className="bg-crema px-4 pt-6">
      <p role="status" className="mx-auto max-w-6xl rounded-[12px] border-2 border-negro bg-blanco p-4 text-[15px]">
        No pudimos cargar el catálogo en este momento. Inténtalo de nuevo en un rato o escríbenos al{" "}
        <strong>{NEGOCIO.telefonoVisible}</strong>.
      </p>
    </div>
  );
}
