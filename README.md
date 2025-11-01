# Go CV Evaluator with Gemini

## Requirement

- Go
- Chroma DB
- Gemini API Key

## Before start the Application

Input system docs to chroma like job description, scoring rubric, and case study brief with each collection

| Collection Name       | For Document               |
| --------------------- | -------------------------- |
| job_description       | Job Description Docs       |
| case_study_brief      | Case Study Brief Docs      |
| cv_rubric             | CV Scroing Rubric Docs     |
| project_report_rubric | Project Report Rubric Docs |

## Run The App

Copy .env file from .env.example and adjust the env file<br>
Download the package

```bash
go mod download
```

```bash
go mod tidy
```

And then run the app

```bash
go run main.go serve
```

```bash
go run main.go consumer --topic=<consumer_topic>
```

## Repository structure

```
├── api
│   ├── openapi.yaml
│   └── response.go
├── application
│   ├── controllers
│   │   ├── consumer
│   │   │   └── cv_evaluator_controller.go
│   │   ├── hello_controller.go
│   │   ├── job_controller.go
│   │   └── upload_document_controller.go
│   ├── helper
│   │   ├── multipart.go
│   │   ├── parse_json_body.go
│   │   ├── upload_document_mapper.go
│   │   └── validator.go
│   └── services
│       ├── consumer
│       │   └── cv_evaluator_service.go
│       ├── hello_service.go
│       ├── job_service.go
│       ├── kafka_producer.go
│       └── upload_document_service.go
├── bootstrap
│   └── app.go
├── cli
│   ├── consumer.go
│   ├── root.go
│   └── serve.go
├── config
│   └── configuration.go
├── domain
│   ├── models
│   │   ├── dao
│   │   │   └── cv_evaluator_job.go
│   │   ├── chroma_dto.go
│   │   ├── chroma_result.go
│   │   ├── evaluate_dto.go
│   │   ├── job_value.go
│   │   ├── upload_document_dto.go
│   │   └── uploaded_files.go
│   └── repository
│       └── cv_evaluator_job_repository.go
├── handlers
│   ├── consumer.go
│   ├── di.go
│   ├── di_consumer.go
│   └── server.go
├── infrastructure
│   └── middleware
│       └── default_middleware.go
├── internal
│   └── generated
│       └── api.gen.go
├── modules
│   ├── chroma-client
│   │   └── go_chroma_client.go
│   ├── gemini-client
│   │   └── go_gemini_client.go
│   ├── go-mysql
│   │   └── go_mysql.go
│   ├── ingest-document
│   │   └── ingest_document.go
│   ├── job-store
│   │   └── job_store.go
│   └── kafka
│       ├── go_consumer_kafka.go
│       ├── go_kafka_options.go
│       └── go_producer_kafka.go
├── .env.example
├── .gitignore
├── Makefile
├── README.md
├── go.mod
├── go.sum
└── main.go
```
