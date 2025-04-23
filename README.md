# Kavak Document Preprocessor

Este proyecto es un chatbot inteligente para Kavak que procesa documentos y responde preguntas sobre el inventario de autos y la propuesta de valor de la empresa.

## Enlaces útiles

- [Documento de diseño](https://docs.google.com/document/d/1o8E_gJuxCb-lzjaW6KoyrmJyANI31eBje7_gt21VLvM/edit?tab=t.0#heading=h.vf126j1zuuog)

## Características principales

- Procesamiento de documentos CSV con información de autos
- Procesamiento de documentos de texto con información sobre Kavak
  Idealmente estas dos capacidades deberian ser procesadas a través de una lambda con el código que utiliza solo siento triggered por un evento de carga a un bucket o algo similar, sin embargo por efectos de practicidad todo el código quedo en este mismo repositorio
- Almacenamiento de embeddings en Redis para búsqueda semántica
- Chatbot que responde preguntas basadas en el contexto relevante
- API REST para interactuar con el chatbot

## Requisitos previos

- Go 1.23 o superior
- Redis (local o en la nube)
- Cuenta de OpenAI con API key

## Configuración

1. Clona el repositorio:

   ```bash
   git clone https://github.com/DiegoSan99/kavak-document-preprocessor.git
   cd kavak-document-preprocessor
   ```

2. Instala las dependencias:

   ```bash
   go mod download
   ```

3. Crea un archivo `.env` basado en `.env.example`:

   ```bash
   cp .env.example .env
   ```

4. Levanta la instancia de redis vector db

   ```bash
   docker-compose up -d
   ```

5. Configura las variables de entorno en el archivo `.env`:
   ```
   OPENAI_API_KEY=tu_api_key_de_openai
   OPENAI_API_URL=https://api.openai.com/v1
   REDIS_URL=redis://localhost:6379
   LOAD_DATA=true
   ```

## Estructura del proyecto

```
kavak-document-preprocessor/
├── pkg/                    # Código fuente principal
│   ├── config/             # Configuración de la aplicación
│   ├── load/               # Carga de modelos y embeddings
│   ├── openai/             # Cliente de OpenAI
│   ├── prompts/            # Plantillas de prompts para el chatbot
│   ├── services/           # Servicios de la aplicación
│   ├── utils/              # Utilidades
│   ├── vectordb/           # Base de datos vectorial (Redis)
│   └── web/                # Controladores y rutas web
├── sample_caso_ai_engineer.csv  # Datos de ejemplo de autos
├── value_proposal.txt      # Información sobre Kavak
├── .env.example            # Ejemplo de variables de entorno
├── main.go                 # Entrypoint del proyecto
├── go.mod                  # Dependencias de Go
└── README.md               # Este archivo
```

## Cómo ejecutar

### Cargar datos en Redis

Para cargar los datos de ejemplo en Redis, asegúrate de que `LOAD_DATA=true` en tu archivo `.env` y ejecuta:

```bash
go run main.go
```

Esto cargará los datos del CSV y el archivo de texto en Redis.

### Iniciar el servidor

Para iniciar el servidor web:

```bash
go run main.go
```

El servidor estará disponible en `http://localhost:8080`.

## API

### Como utilizarlo con mi cuenta de twilio de Whatsapp

Envía un mensaje de WhatsApp al siguiente número
+14155238886

Con el código

```
join so-bank
```

### Endpoint del chatbot

Si quisieras probar el chatbot en forma local puedes hacer lo siguiente

```
POST /api/v1/chatbot/chat
```

Ejemplo de solicitud:

```json
{
  "From": "whatsapp:+1234567890",
  "To": "whatsapp:+0987654321",
  "Body": "Chevrolet Aveo",
  "NumMedia": "0",
  "ProfileName": "User Name",
  "WaId": "1234567890"
}
```

Ejemplo de respuesta:

La respuesta se retorna de esta forma para que Twilio sea capaz de retornarla al cliente usando TwilioML

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Response>
    <Message>El Chevrolet Aveo del año 2017 tiene un precio de $138,999. Este vehículo está disponible en nuestro inventario y cuenta con todas las garantías de Kavak. Si necesitas más información o deseas agendar una cita para verlo, ¡no dudes en decírmelo!</Message>
</Response>
```

## Cómo funciona

1. **Procesamiento de documentos**: Los documentos CSV y de texto se procesan y se convierten en embeddings vectoriales.
2. **Almacenamiento**: Los embeddings se almacenan en Redis para búsqueda semántica.
3. **Consulta**: Cuando un usuario hace una pregunta, el sistema busca documentos relevantes en Redis.
4. **Generación de respuesta**: El sistema utiliza el contexto relevante para generar una respuesta coherente.
