# Pruebas dinámicas del auditor (ronda 2) contra la instancia AISLADA aud2 (solo 127.0.0.1).
# Nunca imprime claves ni cookies: se redactan.
import http.client, json, sys, time, os, re, struct, zlib, random, string, base64, statistics

WEB = ("127.0.0.1", 18100)
API = ("127.0.0.1", 18081)
ORIGEN = "http://localhost:18100"
DIR = os.path.dirname(os.path.abspath(__file__))
CLAVE_ADMIN = open(os.path.join(DIR, "aud2", "clave-prueba.txt")).read().strip()
SECRETOS = [CLAVE_ADMIN]

def red(s):
    for x in SECRETOS:
        if x:
            s = s.replace(x, "<REDACTADO>")
    s = re.sub(r"dc_sesion=[A-Za-z0-9_-]+", "dc_sesion=<REDACTADO>", s)
    return s

def out(*a):
    print(red(" ".join(str(x) for x in a)))
    sys.stdout.flush()

def req(metodo, ruta, cuerpo=None, cookie=None, origen=ORIGEN, ct=None, dest=WEB, extra=None, raw=None):
    c = http.client.HTTPConnection(*dest, timeout=30)
    h = {"Accept": "application/json"}
    if origen:
        h["Origin"] = origen
    if cookie:
        h["Cookie"] = cookie
    data = None
    if raw is not None:
        data = raw
        h["Content-Type"] = ct
    elif cuerpo is not None:
        data = json.dumps(cuerpo).encode()
        h["Content-Type"] = ct or "application/json"
    if extra:
        h.update(extra)
    t0 = time.perf_counter()
    c.request(metodo, ruta, body=data, headers=h)
    r = c.getresponse()
    b = r.read()
    dt = (time.perf_counter() - t0) * 1000
    hs = {k.lower(): v for k, v in r.getheaders()}
    sc = [v for k, v in r.getheaders() if k.lower() == "set-cookie"]
    c.close()
    try:
        j = json.loads(b)
    except Exception:
        j = None
    return r.status, hs, b, j, sc, dt

def login(usuario, clave, dest=WEB):
    st, hs, b, j, sc, _ = req("POST", "/api/admin/login", {"usuario": usuario, "password": clave}, dest=dest)
    ck = None
    for v in sc:
        m = re.match(r"(dc_sesion=[^;]+)", v)
        if m and m.group(1) != "dc_sesion=":
            ck = m.group(1)
    return st, ck, sc, j

def seccion(t):
    out("\n## " + t)

# ------------------------------------------------------------------ 1. login admin y cookie
seccion("1. Login admin (instancia aislada) y atributos de la cookie")
st, ADM, sc, j = login("auditor", CLAVE_ADMIN)
out("login admin ->", st, "rol:", j and j.get("rol"))
for v in sc:
    out("Set-Cookie:", re.sub(r"=([^;]+)", lambda m: "=<REDACTADO>" if len(m.group(1)) > 20 else "=" + m.group(1), v, count=1))
assert st == 200 and ADM

# ------------------------------------------------------------------ 2. crear vendedor
seccion("2. Crear vendedor (clave mostrada una vez)")
nombre_formula = "=HYPERLINK(\"http://evil.invalid\",\"x\")"
st, hs, b, j, _, _ = req("POST", "/api/admin/usuarios", {"usuario": "vend.aud", "nombre": nombre_formula, "rol": "vendedor"}, cookie=ADM)
out("crear vendedor ->", st, "cache-control:", hs.get("cache-control"), "| largo clave_temporal:", len(j.get("clave_temporal", "")) if j else None)
CLAVE_V = j["clave_temporal"]; SECRETOS.append(CLAVE_V); VID = j["usuario"]["id"]
st, VEN, _, j = login("vend.aud", CLAVE_V)
out("login vendedor ->", st, "rol:", j and j.get("rol"))
st, hs, b, j, _, _ = req("GET", "/api/admin/usuarios", cookie=ADM)
out("GET usuarios (admin) ->", st, "¿algún campo con hash/clave?:", any(k in json.dumps(j) for k in ["password", "hash", "clave"]))

# ------------------------------------------------------------------ 3. matriz de autorización
seccion("3. Matriz de autorización por rol (sin sesión / vendedor / admin)")
st, hs, b, j, _, _ = req("GET", "/api/admin/pedidos?page_size=100", cookie=ADM)
pedidos = j["items"]; PID = pedidos[0]["id"]
st, hs, b, j, _, _ = req("GET", "/api/admin/productos?page_size=5", cookie=ADM)
PRODS = j["items"] if j and "items" in j else j
PRID = PRODS[0]["id"]
png1 = None
def png(w, h, extra=b""):
    raw = b"".join(b"\x00" + b"\xff\x00\x00" * w for _ in range(h))
    def ch(t, d):
        return struct.pack(">I", len(d)) + t + d + struct.pack(">I", zlib.crc32(t + d) & 0xffffffff)
    return b"\x89PNG\r\n\x1a\n" + ch(b"IHDR", struct.pack(">IIBBBBB", w, h, 8, 2, 0, 0, 0)) + ch(b"IDAT", zlib.compress(raw)) + ch(b"IEND", b"") + extra
def multipart(nombre_campo, contenido, filename="a.png", ctype="image/png"):
    bnd = "----aud" + "".join(random.choice(string.ascii_letters) for _ in range(12))
    body = (f"--{bnd}\r\nContent-Disposition: form-data; name=\"{nombre_campo}\"; filename=\"{filename}\"\r\nContent-Type: {ctype}\r\n\r\n").encode() + contenido + f"\r\n--{bnd}--\r\n".encode()
    return body, "multipart/form-data; boundary=" + bnd
mp, mpct = multipart("imagen", png(4, 4))
rutas = [
    ("GET", "/api/admin/sesion", None),
    ("GET", "/api/admin/hoy", None),
    ("GET", "/api/admin/pedidos", None),
    ("GET", f"/api/admin/pedidos/{PID}", None),
    ("GET", f"/api/admin/pedidos/{PID}/mensajes", None),
    ("GET", "/api/admin/clientes?telefono=%2B56900000090", None),
    ("GET", "/api/admin/productos", None),
    ("GET", f"/api/admin/productos/{PRID}", None),
    ("GET", "/api/admin/pedidos.csv", None),
    ("GET", "/api/admin/resumen", None),
    ("POST", "/api/admin/productos", {}),
    ("PUT", f"/api/admin/productos/{PRID}", {}),
    ("DELETE", "/api/admin/productos/999999", None),
    ("POST", f"/api/admin/productos/{PRID}/imagen", "MP"),
    ("GET", "/api/admin/usuarios", None),
    ("POST", "/api/admin/usuarios", {"usuario": "escalado", "nombre": "Escalado", "rol": "admin"}),
    ("PATCH", f"/api/admin/usuarios/{VID}", {"rol": "admin"}),
    ("POST", f"/api/admin/usuarios/{VID}/clave", None),
    ("POST", "/api/admin/usuarios/1/clave", None),
]
out("| método ruta | sin sesión | vendedor |")
for m, r, cuerpo in rutas:
    res = []
    for ck in (None, VEN):
        if cuerpo == "MP":
            s = req(m, r, cookie=ck, raw=mp, ct=mpct)[0]
        else:
            s = req(m, r, cuerpo, cookie=ck)[0]
        res.append(s)
    out(f"| {m} {r} | {res[0]} | {res[1]} |")
st, hs, b, j, _, _ = req("GET", "/api/admin/sesion", cookie=VEN)
out("vendedor tras intentos de escalar: rol =", j.get("rol"))
st, hs, b, j, _, _ = req("POST", "/api/admin/login", {"usuario": "vend.aud", "password": CLAVE_V, "rol": "admin"})
out("login con campo extra 'rol' ->", st, j and j["error"]["codigo"])

# ------------------------------------------------------------------ 4. IDOR / alcance del vendedor
seccion("4. Alcance del vendedor sobre pedidos (IDOR)")
ids = [p["id"] for p in pedidos]
ok = sum(1 for i in ids if req("GET", f"/api/admin/pedidos/{i}", cookie=VEN)[0] == 200)
out(f"vendedor puede leer {ok}/{len(ids)} pedidos (incluye los registrados por otros usuarios); id inexistente ->", req("GET", "/api/admin/pedidos/999999", cookie=VEN)[0], "; id no numérico ->", req("GET", "/api/admin/pedidos/1%20OR%201=1", cookie=VEN)[0])
# mensajes con datos de transferencia
hay = None
for p in pedidos:
    if p.get("forma_pago") == "transferencia" and p.get("estado") in ("cotizado", "esperando_pago"):
        s, _, b, j, _, _ = req("GET", f"/api/admin/pedidos/{p['id']}/mensajes", cookie=VEN)
        hay = (s, "MARCADOR-DATOS-TRANSF-PRUEBA" in b.decode())
        break
out("vendedor ve DATOS_TRANSFERENCIA en mensajes (esperado, es para enviarlo al cliente):", hay)

# ------------------------------------------------------------------ 5. CSRF / Origin
seccion("5. Origin / CSRF en escrituras")
for etiqueta, o, ex in [("sin Origin", None, None), ("Origin ajeno", "http://evil.invalid", None), ("Origin null", "null", None),
                        ("Origin OK + Sec-Fetch-Site cross-site", ORIGEN, {"Sec-Fetch-Site": "cross-site"})]:
    s = req("PATCH", f"/api/admin/pedidos/{PID}/estado", {"estado": "cancelado", "motivo": "x"}, cookie=ADM, origen=o, extra=ex)[0]
    out(f"PATCH estado {etiqueta} ->", s)
s = req("POST", "/api/admin/usuarios", None, cookie=ADM, raw=b"usuario=x", ct="application/x-www-form-urlencoded")[0]
out("POST usuarios con form-urlencoded ->", s)
s, hs, *_ = req("OPTIONS", "/api/admin/usuarios", origen="http://evil.invalid", extra={"Access-Control-Request-Method": "POST"})
out("preflight OPTIONS desde origen ajeno ->", s, "ACAO:", hs.get("access-control-allow-origin"))

# ------------------------------------------------------------------ 6. pedido nuevo, seguimiento y fugas
seccion("6. Pedido nuevo y página de seguimiento (fugas de datos personales)")
nuevo = {"canal": "regiones", "nombre": "Zoila Prueba Auditor", "telefono": "+56 9 8765 4321", "entrega": "envio",
         "region": "valparaiso", "comuna": "-2-3", "agencia": "starken", "forma_pago": "transferencia",
         "items": [{"producto_id": PRID, "cantidad": 2}, {"descripcion": "=1+1 linea libre", "cantidad": 1, "precio_unitario": 1000}],
         "costo_envio": 4321, "nota_interna": "nota privada MARCADOR-NOTA", "estado_inicial": "esperando_pago"}
s, hs, b, j, _, _ = req("POST", "/api/admin/pedidos", nuevo, cookie=VEN)
out("vendedor crea pedido ->", s, j and j.get("numero"))
NP = j; TOKEN = j["seguimiento_url"].rsplit("/", 1)[-1]
out("seguimiento_url ->", re.sub(r"[A-Za-z0-9_-]{22}$", "<token22>", j["seguimiento_url"]), "| largo token:", len(TOKEN), "| whatsapp_cliente_url empieza con:", (j.get("whatsapp_cliente_url") or "")[:22])
s, hs, b, j, _, _ = req("GET", f"/api/seguimiento/{TOKEN}", origen=None)
out("GET /api/seguimiento/{token} ->", s, "claves:", sorted(j.keys()))
txt = b.decode()
sensibles = ["Zoila", "8765", "4321", "-2-3", "valparaiso", "Valpara", "MARCADOR-NOTA", "total", "subtotal", "items", "telefono", "nombre", str(NP.get("total"))]
out("¿aparece algún dato sensible en el JSON?:", [x for x in sensibles if x in txt])
s, hs, b, j, _, _ = req("GET", f"/seguimiento/{TOKEN}", origen=None, extra={"Accept": "text/html"})
html = b.decode()
out("GET /seguimiento/{token} (HTML SSR) ->", s, "| cache-control:", hs.get("cache-control"), "| x-robots-tag:", hs.get("x-robots-tag"), "| referrer meta no-referrer:", 'name="referrer" content="no-referrer"' in html)
out("¿datos sensibles en el HTML?:", [x for x in ["Zoila", "8765 4321", "98765", "-2-3", "MARCADOR-NOTA", "MARCADOR-DATOS-TRANSF"] if x in html])

# 404 genérico y tiempos
seccion("7. Token de seguimiento: 404 genérico, tiempos y límite")
casos = {"inexistente_22": "A" * 22, "corto": "abc", "largo": "A" * 23, "caracteres": "%2E%2E%2F" + "A" * 13, "numero_pedido": NP["numero"]}
cuerpos = {}
for k, t in casos.items():
    s, hs, b, j, _, _ = req("GET", f"/api/seguimiento/{t}", origen=None, dest=API)
    cuerpos[k] = (s, b)
    out(f"{k} ->", s, b.decode()[:120])
out("¿cuerpos 404 idénticos?:", len(set(v[1] for v in cuerpos.values())) == 1)

# tiempos (directo a la API; se reinicia el contador esperando la ventana si hace falta)
def tiempos(token, n):
    ts = []
    for _ in range(n):
        s, hs, b, j, _, dt = req("GET", f"/api/seguimiento/{token}", origen=None, dest=API)
        if s == 429:
            return ts, True
        ts.append(dt)
    return ts, False
ex, _ = tiempos(TOKEN, 10)
no, _ = tiempos("".join(random.choice(string.ascii_letters) for _ in range(22)), 10)
out(f"tiempos (ms, mediana): existente={statistics.median(ex):.2f} inexistente={statistics.median(no):.2f} (diferencia irrelevante con 128 bits: no permite adivinar)")

# límite por IP: todo el tráfico por la web comparte la IP del contenedor web
n200 = n429 = 0
for i in range(40):
    s = req("GET", f"/api/seguimiento/{'B'*22}", origen=None)[0]
    if s == 429: n429 += 1
    else: n200 += 1
out(f"40 consultas seguidas vía web: no-429={n200} 429={n429}")
s, hs, b, j, _, _ = req("GET", f"/api/seguimiento/{TOKEN}", origen=None, extra={"X-Forwarded-For": "203.0.113.77"})
out("cliente legítimo (otro XFF) con token válido vía web tras el flood ->", s, hs.get("retry-after"))
s, hs, b, _, _, _ = req("GET", f"/seguimiento/{TOKEN}", origen=None, extra={"Accept": "text/html", "X-Forwarded-For": "198.51.100.9"})
out("página /seguimiento/{token} del cliente legítimo tras el flood ->", s, "| ¿muestra el pedido?:", NP["numero"] in b.decode())
s = req("GET", f"/api/seguimiento/{TOKEN}", origen=None, dest=API, extra={"X-Forwarded-For": "203.0.113.99"})[0]
out("directo a la API con XFF inventado (no confiable) ->", s)

# ------------------------------------------------------------------ 8. CSV
seccion("8. Export CSV (solo admin) y CSV/formula injection")
s, hs, b, j, _, _ = req("GET", "/api/admin/pedidos.csv", cookie=ADM, origen=None)
csvt = b.decode("utf-8-sig")
out("admin CSV ->", s, hs.get("content-type"), "|", hs.get("content-disposition"), "| cache-control:", hs.get("cache-control"), "| x-content-type-options:", hs.get("x-content-type-options"))
for linea in csvt.splitlines():
    if NP["numero"] in linea:
        out("fila del pedido de prueba:", linea.replace("Zoila Prueba Auditor", "<nombre>").replace("+56987654321", "<tel>"))
celdas_peligrosas = []
import csv, io
for fila in csv.reader(io.StringIO(csvt), delimiter=";"):
    for c in fila:
        if c[:1] in ("=", "+", "-", "@", "\t", "\r"):
            celdas_peligrosas.append(c[:30])
out("celdas que empiezan con = + - @ TAB CR (sin el apóstrofo):", celdas_peligrosas[:10])
out("¿nota interna o línea libre en el CSV?:", "MARCADOR-NOTA" in csvt, "=1+1" in csvt)

# ------------------------------------------------------------------ 9. imágenes
seccion("9. Subida y entrega de imágenes")
def subir(contenido, filename="a.png", ctype="image/png", ck=ADM, campo="imagen"):
    body, ct = multipart(campo, contenido, filename, ctype)
    return req("POST", f"/api/admin/productos/{PRID}/imagen", cookie=ck, raw=body, ct=ct)
poliglota = png(8, 8, extra=b"<script>alert(1)</script><?php echo 1; ?>")
s, hs, b, j, _, _ = subir(poliglota, "x.html", "text/html")
url = j and j.get("imagen_url")
out("PNG políglota (HTML/PHP al final, filename x.html, CT text/html) ->", s, "imagen_url:", url)
if url:
    s, hs, b, *_ = req("GET", url, origen=None)
    out("GET imagen vía web ->", s, hs.get("content-type"), "| nosniff:", hs.get("x-content-type-options"), "| CSP:", hs.get("content-security-policy"), "| ¿contiene <script>?:", b"<script>" in b, "| ¿empieza con PNG?:", b[:4] == b"\x89PNG")
svg = b'<svg xmlns="http://www.w3.org/2000/svg" onload="alert(1)"/>'
out("SVG ->", subir(svg, "a.svg", "image/svg+xml")[0])
out("GIF ->", subir(b"GIF89a" + b"\x00" * 30, "a.gif", "image/gif")[0])
out("HTML con extensión .png ->", subir(b"<html><script>alert(1)</script></html>", "a.png", "image/png")[0])
bomba = png(1, 1)
bomba = bomba[:16] + struct.pack(">II", 60000, 60000) + bomba[24:]
bomba = bomba[:29] + struct.pack(">I", zlib.crc32(bomba[12:29]) & 0xffffffff) + bomba[33:]
s, hs, b, j, _, _ = subir(bomba)
out("PNG 60000x60000 (cabecera, bomba de descompresión) ->", s, j and j["error"]["codigo"])
out("campo con otro nombre ->", subir(png(4, 4), campo="otro")[0])
out("6 MB ->", subir(b"\x89PNG\r\n\x1a\n" + os.urandom(6 * 1024 * 1024))[0])
for t in ["../../etc/passwd", "..%2f..%2fetc%2fpasswd", "%2e%2e%2f%2e%2e%2fetc%2fpasswd", "....//....//etc/passwd",
          "0123456789abcdef0123456789abcdef.jpg%00.png", "0123456789ABCDEF0123456789ABCDEF.jpg", ".subida-123", "carlitos.db"]:
    sw = req("GET", "/api/imagenes/" + t, origen=None)[0]
    sa = req("GET", "/api/imagenes/" + t, origen=None, dest=API)[0]
    out(f"GET /api/imagenes/{t} -> web {sw} / api {sa}")

# ------------------------------------------------------------------ 10. gestión de usuarios: invalidación de sesiones
seccion("10. Gestión de usuarios: desactivar, regenerar clave, cambio de rol")
s = req("PATCH", f"/api/admin/usuarios/{VID}", {"activo": False}, cookie=ADM)[0]
out("admin desactiva vendedor ->", s, "| sesión vieja del vendedor ->", req("GET", "/api/admin/sesion", cookie=VEN)[0], "| login del desactivado ->", login("vend.aud", CLAVE_V)[0])
s = req("PATCH", f"/api/admin/usuarios/{VID}", {"activo": True}, cookie=ADM)[0]
out("reactiva ->", s, "| sesión vieja sigue inválida ->", req("GET", "/api/admin/sesion", cookie=VEN)[0])
st, VEN2, _, _ = login("vend.aud", CLAVE_V)
s, hs, b, j, _, _ = req("POST", f"/api/admin/usuarios/{VID}/clave", None, cookie=ADM)
NUEVA = j["clave_temporal"]; SECRETOS.append(NUEVA)
out("regenerar clave ->", s, "| sesión previa ->", req("GET", "/api/admin/sesion", cookie=VEN2)[0], "| clave vieja ->", login("vend.aud", CLAVE_V)[0], "| clave nueva ->", login("vend.aud", NUEVA)[0])
# segundo admin degradado
s, hs, b, j, _, _ = req("POST", "/api/admin/usuarios", {"usuario": "admin.dos", "nombre": "Admin Dos", "rol": "admin"}, cookie=ADM)
C2 = j["clave_temporal"]; SECRETOS.append(C2); A2ID = j["usuario"]["id"]
st, A2, _, _ = login("admin.dos", C2)
out("admin.dos ve usuarios ->", req("GET", "/api/admin/usuarios", cookie=A2)[0])
out("admin.dos se quita el rol a sí mismo ->", req("PATCH", f"/api/admin/usuarios/{A2ID}", {"rol": "vendedor"}, cookie=A2)[0])
st, hs, b, j, _, _ = req("GET", "/api/admin/usuarios", cookie=ADM)
AID = [u for u in j["items"] if u["usuario"] == "auditor"][0]["id"]
s, hs, b, j, _, _ = req("POST", f"/api/admin/usuarios/{AID}/clave", None, cookie=A2)
CLAVE_ADMIN2 = j["clave_temporal"]; SECRETOS.append(CLAVE_ADMIN2)
out("admin.dos regenera la clave del admin de api/.env ->", s, "| sesión del admin principal ->", req("GET", "/api/admin/sesion", cookie=ADM)[0], "| su clave original ->", login("auditor", CLAVE_ADMIN)[0])
st, ADM, _, _ = login("auditor", CLAVE_ADMIN2)
out("admin principal entra con la clave que generó admin.dos ->", st)
s = req("PATCH", f"/api/admin/usuarios/{A2ID}", {"rol": "vendedor"}, cookie=ADM)[0] if ADM else None
out("admin principal degrada a admin.dos ->", s, "| admin.dos GET usuarios con su sesión vieja ->", req("GET", "/api/admin/usuarios", cookie=A2)[0])

# ------------------------------------------------------------------ 11. logout
seccion("11. Logout invalida en el servidor")
st, TMP, _, _ = login("vend.aud", NUEVA)
s = req("POST", "/api/admin/logout", None, cookie=TMP)[0]
out("logout ->", s, "| reuso de la cookie ->", req("GET", "/api/admin/sesion", cookie=TMP)[0])

# ------------------------------------------------------------------ 12. paginación y validación
seccion("12. Paginación con tope y validación")
out("page_size=1000 ->", req("GET", "/api/admin/pedidos?page_size=1000", cookie=ADM)[0], "| q de 200 caracteres ->", req("GET", "/api/admin/pedidos?q=" + "a" * 200, cookie=ADM)[0],
    "| q con comodines %_ ->", req("GET", "/api/admin/pedidos?q=%25%5F%25", cookie=ADM)[3].get("total"))
out("resumen 400 días ->", req("GET", "/api/admin/resumen?desde=2025-01-01&hasta=2026-11-14", cookie=ADM)[0])
grande = dict(nuevo); grande["items"] = [{"producto_id": PRID, "cantidad": 10000}] * 51
out("pedido con 51 líneas ->", req("POST", "/api/admin/pedidos", grande, cookie=ADM)[0])
out("body de 70 KB ->", req("POST", "/api/admin/pedidos", None, cookie=ADM, raw=b'{"nota_interna":"' + b"a" * 70000 + b'"}', ct="application/json")[0])

# ------------------------------------------------------------------ 13. DATOS_TRANSFERENCIA fuera del front
time.sleep(62)
seccion("13. DATOS_TRANSFERENCIA no llega a rutas públicas")
pub = ["/", "/api/config-publica", "/api/productos?page_size=50", "/api/opciones-entrega", f"/seguimiento/{TOKEN}", f"/api/seguimiento/{TOKEN}", "/admin"]
for p in pub:
    s, hs, b, *_ = req("GET", p, origen=None, extra={"Accept": "text/html,application/json"})
    out(f"{p if 'seguimiento' not in p else p.replace(TOKEN, '<token>')} -> {s} | marcador presente: {b'MARCADOR-DATOS-TRANSF' in b}")

# ------------------------------------------------------------------ 14. bloqueo de login compartido (al final)
seccion("14. Bloqueo de login: ¿lo comparte todo el tráfico? (TRUSTED_PROXIES vacío)")
res = [req("POST", "/api/admin/login", {"usuario": "no-existe", "password": "x"}, extra={"X-Forwarded-For": f"203.0.113.{i}"})[0] for i in range(6)]
out("atacante: 6 logins fallidos con XFF distintos vía web ->", res)
st, ck, _, j = login("auditor", CLAVE_ADMIN)
out("admin legítimo, clave correcta, inmediatamente después ->", st, j and j.get("error", {}).get("codigo"))
st, ck, _, j = login("auditor", CLAVE_ADMIN, dest=API)
out("admin legítimo directo a la API (otra IP de conexión) ->", st)
