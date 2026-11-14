// Descripción del producto. Las frases marcadas "[EJEMPLO]" en la API no se
// muestran como texto interno: van en gris con la etiqueta "dato de ejemplo" (QA M-7).

const RE = /\[EJEMPLO\]\s*([^\n]*?(?:\.(?=\s|$)|$))/g;

export function Descripcion({ texto }: { texto: string }) {
  const ejemplos: string[] = [];
  const limpio = texto
    .replace(RE, (_m, frase: string) => {
      const f = frase.trim();
      if (f) ejemplos.push(f);
      return "";
    })
    .replace(/[ \t]{2,}/g, " ")
    .trim();

  return (
    <div className="max-w-[60ch] space-y-2">
      {limpio ? <p className="whitespace-pre-line text-[15px] text-negro/85">{limpio}</p> : null}
      {ejemplos.length > 0 ? (
        <ul className="space-y-1">
          {ejemplos.map((e) => (
            <li key={e} className="text-[13px] text-gris">
              <span className="mr-1 inline-block whitespace-nowrap rounded-full border border-dashed border-gris px-2 py-px align-middle text-[11px] font-semibold leading-4">
                dato de ejemplo
              </span>
              {e}
            </li>
          ))}
        </ul>
      ) : null}
    </div>
  );
}
