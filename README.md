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
go run main.go
```

## Endpoint Explanation

`POST /upload (multipart/form-data)` use to upload the CV and Project Report File <br>
Body Request

```json
@cv_file: cv_file.pdf
@report_file: report_file.pdf
```

Response

```json
{
  "message": "Success",
  "status": 200,
  "data": {
    "file_id": "d44411df-e2aa-4915-97f6-69556b20904c"
  }
}
```

---

`POST /evalutate (application/json)` use to evaluate the uploaded file<br>
Body Request

```json
{
  "job_title": "backend engineer",
  "file_id": "d44411df-e2aa-4915-97f6-69556b20904c" // Base on the file_id return from /upload
}
```

Response

```json
{
  "message": "Success",
  "status": 200,
  "data": {
    "job_id": "14dae7fb-ec8b-4a33-a12b-0b00a07d67cc",
    "status": "queued"
  }
}
```

---

`GET /result/{jobId} (application/json)` use to see job result `jobId based on the return from /evaluate`<br>
Body Request

```json
-
```

Response

```json
{
  "message": "Success",
  "status": 200,
  "data": {
    "id": "14dae7fb-ec8b-4a33-a12b-0b00a07d67cc",
    "job_title": "backend engineer",
    "file_id": "d44411df-e2aa-4915-97f6-69556b20904c",
    "status": "completed",
    "result": {
      "cv_match_rate": "0.76/1.0",
      "cv_feedback": "This candidate presents a strong profile for a backend engineer, showcasing solid technical skills in modern backend languages, databases, and cloud technologies (GCP certified). Their experience integrating a machine learning model into a backend service directly addresses",
      "project_score": "1.7",
      "project_feedback": "The provided document is a candidate's CV, not the requested project report for the take-home case study. Therefore, a comprehensive evaluation against the project rubric, particularly for aspects like prompt design, LLM chaining, RAG implementation, and specific error handling, is not possible. Based on the CV, the candidate demonstrates solid backend engineering experience, external API integration, and cloud knowledge (GCP), with a good understanding of testing practices (unit tests).",
      "overall_summary": "This candidate presents a strong profile for a backend engineer, boasting a high CV match rate (0.76) and solid technical skills in modern backend languages, databases, and cloud technologies (GCP certified), including experience integrating ML models. However, a critical gap is the incorrect project submission, as the candidate provided their CV instead of the take-home case study report. This prevents a comprehensive evaluation of their practical skills in prompt design, LLM chaining, or RAG implementation. While their CV demonstrates strong backend experience, external API integration, and testing, their hands-on project execution abilities remain unassessed. Therefore, it is essential to request the correct project report for a complete evaluation."
    }
  }
}
```

---

`PUT /chroma/upsert (application/json)` helper for upsert docs to chroma<br>
Body Request

```json
{
  "collection_name": "test_collection",
  "content_id": "test_content_id",
  "content": "test_content",
  "metadata": {
    "type": "test"
  }
}
```

Response

```json
{
  "message": "Success",
  "status": 200
}
```
