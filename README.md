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
