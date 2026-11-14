# Pendientes — investigador-marca (Fase 1, 2026-11-14)

Lo que no se pudo ver o verificar. **No se debe afirmar en la web ni en la propuesta sin confirmarlo.**

## 1. Instagram @mayte_emporioo: muro de login (pedir a Luckas que inicie sesión una vez)
Sin sesión solo se pudo ver:
- la cabecera (bio, seguidores, insignia verificada): `raw/ig-perfil-02.png`;
- los 12 primeros reels, con fecha y me gusta (metadatos);
- el avatar (`raw/logo-instagram-avatar.png`).

**Quedó sin ver:**
- [ ] **Comentarios** de cualquier post o reel. El modal de registro los tapa: `raw/ig-reel-DfqeJ3CBYN4.png`.
- [ ] **Historias destacadas** (highlights). Ahí suelen estar horario, envíos, medios de pago y "cuidado con estafas".
- [ ] **Grilla completa** más allá de los 12 primeros (el perfil dice 3.653 publicaciones). Falta ver si hay **carruseles o posts con precio en imagen** distintos de TikTok.
- [ ] **Descripciones completas** de los reels: sin sesión salen vacías o recortadas.
- [ ] Si responden **mensajes directos** o mandan todo a WhatsApp, y si hay respuestas fijadas.
- [ ] Los colaboradores de los posts conjuntos ("mayte_emporioo y 3 personas más"): con qué cuentas publican.

**Cómo destrabarlo:** Luckas inicia sesión una vez en la ventana de Playwright. No se guardan cookies ni estado de sesión en el repo. Basta con revisar highlights, 30 posts y comentarios de 10 reels.

## 2. TikTok: lo que no se ve sin sesión
- [ ] **Lives**: si hacen, con qué frecuencia y si venden en vivo. No hay historial público de lives.
- [ ] **Comentarios** más allá de los primeros ~15 por video. TikTok pidió registro después de unos 25 videos.
- [ ] **Videos anteriores a 2024** fuera de los populares (se capturaron 373 videos de muchos más).
- [ ] Texto en pantalla de los videos: solo se revisaron fotogramas sueltos en las capturas. El audio se leyó por los subtítulos automáticos, que pueden traer errores.

## 3. Datos del negocio por confirmar con el dueño
| Dato | Lo que se vio | Fuente | Estado |
|---|---|---|---|
| Dirección Alameda 2432 | Confirmada en las bios de TikTok e IG y en videos de 2025–2026 | `raw/tiktok-perfil-01.png`, `raw/ig-perfil-02.png` | **Confirmado** (redes oficiales) |
| WhatsApp +56 9 2114 5368 | En las bios de TikTok e IG ("pedidos solo a este número") y en tr. de noviembre 2026 | ídem; TT-7707471893363789063 | **Confirmado** (redes oficiales) |
| Horario | Tres versiones: CLAUDE.md L–V 9–17 / sáb 9–15; enero 2026 L–V 9–17:30 / sáb 9–15; noviembre 2026 L–V 9–18 / sáb, dom y feriados 9–16 | TT-7592989161653685522; TT-7707471893363789063 | **Pendiente:** ¿cuál es el vigente? ¿Hay horario de temporada? |
| Envíos a todo Chile | 2025 y abril 2026: "solamente venta presencial". Noviembre 2026: "número para región" y comentario sobre "la niña de los pedidos de regiones". Cuenta secundaria: "Envíos a TODO CHILE" | TT-7505555394513685766, TT-7627590547473059090, TT-7707471893363789063, TT-7706913997378948370 | **Pendiente.** ¿Despachan? ¿Con qué courier? ¿Quién paga? ¿Monto mínimo? |
| Abate Molina 77, local 5 | Usado de marzo 2024 a junio 2025; después ya no aparece | TT-7347392997358669061, TT-7512969200437120262 | **Probablemente cerrado** *(inferencia)*. Confirmar |
| Bascuñán Guerrero 364 ("local amarillo") | 2023 a febrero 2024 | TT-7286474523724107014 | Histórico. Confirmar si sigue |
| Local Temuco, Manuel Rodríguez 1232 | Inaugurado en octubre 2025; Meganoticias dice 17–18 trabajadores | TT-7562937787876445458, TT-7592989161653685522 | **Pendiente:** ¿se incluye en la web? ¿Horario? ¿Mismos precios? |
| "Dos locales" en Alameda 2432 | Uno de hogar, electro y juguetes (efectivo y débito) y otro de Disney y licencias (débito, crédito, transferencia y efectivo) | tr. TT-7707464302034160914 | **Pendiente:** ¿son dos tiendas separadas? ¿Números de local? |
| Medios de pago | Ver fila anterior; preguntan por "link de crédito" | ídem | Confirmar por local |
| Venta por mayor | **Sí**: precio por embalaje o por mayor en casi todos los videos; el hashtag #meiggspormayor | Ver `ofertas.md` | Confirmado como práctica. **Pendiente:** ¿hay mínimo de compra? ¿Lista de precios mayorista? |
| Boleta y factura | "boleta factura lo que usted quiera" | tr. TT-7707464302034160914 | Confirmar |
| Cambios y devoluciones | Un cliente pregunta "hacen cambio si compré un día antes" | TT-7683668603270352136 | **Sin dato** |
| Cómo se canjea "yo quiero mi promoción" | No se explica en los videos | — | **Pendiente.** ¿Muestra el comentario en caja? ¿Hay vigencia? |

## 4. Identidad y cuentas: no afirmar relaciones sin confirmar
- [ ] **@_dondecarlitosemporio7 (TikTok):** se llama "🏆Donde Carlitos Emporio🏆", tiene 538,3 mil seguidores y su bio dice "ÚNICA CUENTA OFICIAL • Evita estafas" y "Envíos a TODO CHILE" (`raw/secundaria-tiktok-_dondecarlitosemporio7.png`). A favor de que sea suya:
  - la oficial usó el hashtag #dondecarlitosemporio7 en 2023–2024;
  - el logo oficial lista "Donde Carlitos Emporio" en TikTok;
  - un cliente etiqueta a "@🏆Donde Carlitos Emporio🏆" en un comentario.

  Pero Luckas dijo que las oficiales son las de Mayte. **¿Es del dueño, de un familiar o de un tercero?** Esto define el aviso antiestafa de la web.
- [ ] **Facebook:** el logo trae el ícono de Facebook con "DONDE CARLITOS EMPORIO", pero Luckas dice que no hay Facebook. ¿Lo hubo? ¿Es la página facebook.com/donde.carlitos.33 (no confirmada)?
- [ ] **"Mayte Emporio" y "Andrés Emporio":** ambos nombres van en el logo. En un video se nombra a un "jefe Andrés" (TT-7707047904761990408). ¿Quiénes son? ¿Se nombran en la web?
- [ ] **Cuentas etiquetadas:** @Xixitan_emporio, @M.alexandraemporio, @la.reinaemporio7, @Andres emporio, @Saiyan, @LavandaHome.cl ("con mi hermanito"). No se afirma ninguna relación.
- [ ] **¿Quién es "Carlitos"?** Los clientes le dicen "Carlitos" al presentador principal. Meganoticias nombra a Carlos Ortega como fundador de "Mayte Emporio". Confirmar si es la misma persona y si quiere aparecer con nombre en la web.
- [ ] **"Familia Emporio":** ¿es el nombre paraguas oficial (Santiago + Temuco)? ¿Cómo quiere el dueño que se llame la web: "Donde Carlitos Emporio", "Mayte Emporio" o "Familia Emporio"?

## 5. Logo
- [ ] **Pedir el archivo original** (vector, o PNG de alta resolución con fondo transparente). Solo hay un raster de 1077 px con fondo negro (`raw/logo-tiktok-avatar.png`).
- [ ] Preguntar si existe una versión **sin los nombres de redes** (texto chico ilegible en tamaño de favicon o header).
- [ ] Confirmar si el 🏆 es parte de la marca o solo de la cuenta secundaria.

## 6. Datos cuantitativos que no se tienen
- [ ] Ventas, tickets promedio y porcentaje de clientes de regiones o mayoristas: no se infieren cifras.
- [ ] Reseñas de Google Maps: no se revisaron, quedan fuera del alcance de las redes. Puede verlas `analista-meiggs`.
- [ ] Calendario de Fiestas Patrias y Cyber: no apareció en la muestra.

## 7. Notas de privacidad y uso
- Las capturas `raw/tt-*.png` muestran nombres de usuario de comentaristas. Son material de investigación interno y **no se usan en la web**.
- En `raw/secundaria-tiktok-_dondecarlitosemporio7.png` se tapó el avatar porque es una foto de personas, incluido un menor.
- En `raw/tt-7441724226857454904.png` se tapó una sugerencia de búsqueda automática de TikTok que no tiene relación con la investigación.
- No se guardaron cookies ni estado de sesión. El navegador se usó sin login.
