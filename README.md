# Sistema RAG con Ollama y Golang

Un sistema de Recuperación Aumentada de Generación (RAG) implementado en Go que permite consultar documentos PDF utilizando Ollama y el modelo llama2.

## ¿Qué es RAG?

RAG (Retrieval Augmented Generation) es una técnica de inteligencia artificial que combina la recuperación de información con generación de texto. Permite obtener respuestas precisas sobre documentos específicos utilizando búsqueda vectorial y modelos de lenguaje.

## Características

- **Procesamiento de PDF**: Extrae y procesa texto de archivos PDF
- **Chunking inteligente**: Divide documentos en fragmentos de 500 palabras
- **Base de datos vectorial**: Utiliza ChromeM-Go para búsqueda semántica
- **Integración con Ollama**: Genera respuestas usando el modelo llama2
- **Modo interactivo**: Permite hacer preguntas en tiempo real sobre el documento

## Tecnologías Utilizadas

- **Go 1.24**: Lenguaje de programación principal
- **Ollama**: Plataforma para ejecutar modelos de lenguaje localmente
- **llama2**: Modelo de lenguaje para generar respuestas
- **ChromeM-Go**: Base de datos vectorial en Go puro
- **LangChain Go**: Framework para aplicaciones de LLM
- **pdf**: Librería para procesamiento de archivos PDF

## Instalación

### Prerrequisitos

1. **Instalar Go** (versión 1.24 o superior)
2. **Instalar Ollama**:
   ```bash
   # macOS/Linux
   curl -fsSL https://ollama.ai/install.sh | sh
   
   # Windows
   # Descargar desde https://ollama.ai
   ```

3. **Descargar el modelo llama2**:
   ```bash
   ollama pull llama2
   ```

4. **Ejecutar Ollama**:
   ```bash
   ollama serve
   ```

### Configuración del Proyecto

1. **Clonar el repositorio**:
   ```bash
   git clone <url-del-repositorio>
   cd simple-rag-golang-youtube
   ```

2. **Instalar dependencias**:
   ```bash
   go mod tidy
   ```

## Uso

### Ejecutar el Sistema

```bash
go run main.go <ruta-al-archivo-pdf>
```

**Ejemplo**:
```bash
go run main.go ./monopoly.pdf
```

### Modo Interactivo

Una vez cargado el PDF, el sistema entrará en modo interactivo:

```
Simple RAG AI Agent
====================
RAG System Initialized successfully!
Loading PDF: ./monopoly.pdf
PDF loaded successfully! Found 1 document(s)
Document has 15 chunks
Stored 15 chunks in pure Go vector database

=== Query Mode ===
You can now ask questions about the document.
Type 'quit' to exit.

Question: ¿Cuáles son las reglas básicas del Monopoly?
Generating answer...
Chunk 1 preview: Las reglas básicas del Monopoly incluyen...

Answer: Según el documento, las reglas básicas del Monopoly son...
```

## Arquitectura del Sistema

### Componentes Principales

- **RAGSystem**: Estructura principal que coordina todos los componentes
- **Document**: Representa un documento procesado con sus chunks
- **Vector Database**: ChromeM-Go para almacenamiento y búsqueda vectorial
- **LLM**: Integración con Ollama para generación de respuestas

### Flujo de Procesamiento

1. **Carga de PDF**: Extrae texto página por página
2. **Chunking**: Divide el texto en fragmentos de 500 palabras
3. **Vectorización**: Almacena chunks en la base de datos vectorial
4. **Consulta**: Realiza búsqueda semántica y genera respuesta

### Estructura del Código

```
main.go:17-86    - Función principal e interfaz de usuario
main.go:101-123  - Inicialización del sistema RAG
main.go:125-184  - Carga y procesamiento de PDF
main.go:186-204  - División en chunks
main.go:206-226  - Sistema de consultas semánticas
main.go:228-247  - Generación de respuestas con LLM
```

## Dependencias

```go
require (
    github.com/ledongthuc/pdf v0.0.0-20250511090121-5959a4027728
    github.com/philippgille/chromem-go v0.7.0
    github.com/tmc/langchaingo v0.1.13
)
```

## Configuración Avanzada

### Personalizar Tamaño de Chunks

Modificar la constante en `main.go:188`:
```go
const chunkSize = 500 // Cambiar según necesidades
```

### Cambiar Modelo de Ollama

Modificar en `main.go:112`:
```go
llm, err := ollama.New(ollama.WithModel("llama3")) // Usar otro modelo
```

### Ajustar Cantidad de Chunks Relevantes

Modificar en `main.go:213`:
```go
results, err := r.collection.Query(ctx, question, 5, nil, nil) // Más chunks
```

## Solución de Problemas

### Error: "failed to Initialize Ollama"
- Verificar que Ollama esté ejecutándose: `ollama serve`
- Confirmar que el modelo esté descargado: `ollama list`

### Error: "failed to open PDF"
- Verificar que el archivo PDF existe y es legible
- Comprobar permisos de archivo

### Búsqueda vectorial falla
- El sistema incluye fallback automático a búsqueda por palabras clave
- Verificar que los chunks se almacenaron correctamente

## Ejemplos de Consultas

- "¿Cuál es el tema principal del documento?"
- "Resume los puntos más importantes"
- "¿Qué información específica contiene sobre [tema]?"
- "Explica los conceptos clave mencionados"

## Contribución

1. Fork el repositorio
2. Crear rama para nueva funcionalidad
3. Realizar cambios y pruebas
4. Enviar pull request

## Licencia

Este proyecto está bajo la licencia MIT. Ver archivo LICENSE para más detalles.