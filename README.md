Sistema empresarial distribuido para la gestión centralizada de imágenes entre dos sedes.

---

Estado
**Fase 1 - Planeación del software.** En esta etapa se definen problema, alcance, actores, requisitos, casos de uso, trazabilidad, arquitectura preliminar, métricas y organización del repositorio. El código funcional se implementará de forma incremental en las siguientes fases.



Problema Organizacional
La empresa trabaja con imágenes de proyectos, diseños, documentación y evidencias desde dos sedes. Cuando los archivos se mantienen en ubicaciones diferentes, se dificulta su localización, clasificación, control de acceso y seguimiento. Además, no existen métricas uniformes para conocer cuántas operaciones se realizan, cuánto tardan o qué errores ocurren.

**ImageFlow** propone un repositorio lógico central con autenticación y trazabilidad para que ambas sedes utilicen las mismas reglas de acceso y administración.



Objetivo General
Desarrollar una aplicación empresarial cliente-servidor en Go para centralizar la gestión segura de imágenes entre dos sedes, integrando autenticación mediante FastAPI, control de acceso, metadatos y registros de operación que permitan validar funcionalidad, rendimiento y uso del sistema.

---

Funcionalidades Previstas
- Iniciar sesión mediante el servicio FastAPI.
- Cargar imágenes autorizadas.
- Registrar metadatos de cada recurso.
- Clasificar imágenes por categoría.
- Listar, buscar y filtrar imágenes.
- Descargar imágenes.
- Eliminar imágenes con rol Administrador.
- Registrar operaciones y métricas.

---

Actores
- **Empleado:** Carga, clasifica, consulta, busca y descarga imágenes.
- **Administrador:** Realiza las operaciones del empleado y puede eliminar recursos y consultar métricas.
- **Servicio FastAPI:** Valida credenciales, emite sesión/token y apoya la verificación de autorización.

---

Arquitectura preliminar

Sede A - Cliente Go ─┐
                     ├── Red entre sedes ── Servidor ImageFlow (Go) ── Almacenamiento
Sede B - Cliente Go ─┘                         │
                                               └── Servicio FastAPI (auth) 

El servidor Go concentra las operaciones de imágenes y consulta el servicio FastAPI antes de completar recursos protegidos. Las dos sedes utilizan el mismo servicio lógico, evitando duplicar reglas de autenticación y gestión.

| Componente | Tecnología | Propósito |
| Cliente y servidor | Go | Operaciones de gestión de imágenes y comunicación cliente-servidor |
| Autenticación | Python + FastAPI | Login, token/sesión y validación de acceso |
| Red | VLAN, VLSM, OSPF | Comunicación controlada entre sedes en fases de diseño e implementación |
| Control de versiones | Git / GitHub | Historial, colaboración y evidencia de contribuciones |
| Estadística | Datos exportados por logs | Tiempos, tamaños, frecuencias, errores y comparaciones |

Requisitos clave
- RF-01: autenticación mediante FastAPI.
- RF-02: carga de imágenes.
- RF-04: consulta y búsqueda.
- RF-06: descarga.
- RF-07: eliminación solo por Administrador.
- RF-08: registro de operaciones.
- RF-10: autorización antes de recursos protegidos.
- RNF-01: contraseñas con hash seguro.
- RNF-02: validación de archivos en el servidor.
- RNF-04: arquitectura modular.
- RNF-06: manejo controlado de errores de red/servidor.

La especificación completa y los criterios de aceptación se encuentran en docs/requirements/.

Datos medibles

Cada caso de uso producirá evidencia cuantitativa:
- tiempo de autenticación;
- intentos de login exitosos y fallidos;
- tamaño de imágenes;
- tiempo de carga y descarga;
- número de búsquedas y descargas;
- operaciones por categoría;
- eliminaciones e intentos denegados;
- tasa de errores.


Estructura del repositorio
proyecto-integrador/
├── app-go/
│ ├── cmd/
│ ├── internal/
│ └── README.md
├── auth-service/
│ ├── app/
│ ├── tests/
│ └── requirements.txt
├── network/
│ ├── diagrams/
│ ├── configs/
│ └── addressing/
├── research/
│ ├── instruments/
│ ├── data/
│ └── analysis/
├── docs/
├── LICENSE
└── README.md


Integrantes y responsabilidades

| Integrante | Área principal | Responsabilidades |
| Sebastián Alexis Cumbal Molineros | Programación y Seguridad | Go cliente-servidor, integración FastAPI, autenticación/autorización, seguridad de archivos y requisitos de software |
| Kevin Cruzate | Redes y Estadística | Red entre sedes, VLAN/VLSM/OSPF en fases posteriores, variables, métricas y análisis estadístico |
| Ambos | Integración y documentación | Casos de uso, trazabilidad, README, pruebas integradas, video y defensa |


Convención de commits
Se utilizarán mensajes breves y descriptivos. Ejemplos:
docs: definir requisitos funcionales de ImageFlow
feat: crear estructura inicial del servidor Go
docs: definir variables y métricas estadísticas
chore: crear estructura inicial de network
docs: completar matriz de trazabilidad

Cada integrante debe realizar commits desde su propia cuenta para que el historial evidencie participación real.

Alcance actual
En Fase 1 no se afirma que el sistema esté implementado. Se han definido los artefactos de planeación necesarios para convertir las necesidades del problema en requisitos verificables. Los resultados de rendimiento y estadísticas se medirán cuando exista un prototipo funcional.

Seguridad prevista
- Contraseñas protegidas con hash adaptativo.
- Tokens/sesiones para acceso a recursos protegidos.
- Verificación de roles en servidor.
- Validación de formato y tamaño de archivos.
- Registro de eventos relevantes sin almacenar secretos.
- Principio de mínimo privilegio.

Referencias técnicas
- Go - Organizing a Go module: https://go.dev/doc/modules/layout
- Go - net/http: https://go.dev/src/net/http/doc.go
- FastAPI - OAuth2, hashing y JWT: https://fastapi.tiangolo.com/tutorial/security/oauth2-jwt/
- OWASP - File Upload Cheat Sheet: https://cheatsheetseries.owasp.org/cheatsheets/File_Upload_Cheat_Sheet.html
- OWASP - Password Storage Cheat Sheet: https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
- RFC 2328 - OSPFv2: https://www.rfc-editor.org/info/rfc2328/
- RFC 1918 - Direcciones privadas: https://www.rfc-editor.org/info/rfc1918/
- NIST/SEMATECH e-Handbook: https://www.itl.nist.gov/div898/handbook/
- WCAG 2.2: https://www.w3.org/TR/WCAG22/

Repositorio
URL: https://github.com/S3bas05/Proyecto-Integrador.git

Licencia
La licencia del proyecto se definirá y agregará al repositorio antes de la entrega definitiva.
Uso de herramientas de IA
Las herramientas de IA pueden emplearse como apoyo para organización, revisión y generación de propuestas. El equipo debe revisar, comprender y validar cualquier contenido incorporado al proyecto y declarar su uso conforme a las indicaciones de la asignatura.
