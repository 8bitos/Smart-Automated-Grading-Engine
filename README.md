# SAGE (Smart Automated Grading Engine) 🎓
**Sistem Penilaian Esai Otomatis Berbasis RAG & Gemini AI**

SAGE adalah platform *Learning Management System* (LMS) inovatif yang dirancang untuk mengotomatisasi penilaian jawaban esai siswa. Menggunakan metode **Retrieval-Augmented Generation (RAG)**, sistem ini mengevaluasi jawaban secara objektif berdasarkan materi pembelajaran asli, dengan fokus pada level kognitif **C1 (Mengingat) hingga C4 (Menganalisis)**.

---

## 🚀 Fitur Utama

- **RAG-Powered Grading:** Menilai esai menggunakan konteks materi nyata dari database vektor, meminimalisir halusinasi AI.
- **Analitik Kognitif (C1-C4):** Visualisasi penguasaan siswa pada tingkat berpikir Bloom's Taxonomy.
- **Feedback Formatif Otomatis:** Memberikan saran perbaikan instan kepada siswa melalui Gemini AI.
- **Teacher Review Dashboard:** Dashboard bagi guru untuk meninjau dan memvalidasi skor hasil AI.
- **Fast Performance:** Didukung oleh caching Redis untuk respon penilaian yang cepat.

---

## 🛠️ Stack Teknologi

Sistem ini dibangun dengan arsitektur modern:

- **Frontend:** [Next.js 15+](https://nextjs.org/) (React, Tailwind CSS, TypeScript)
- **Backend:** [Golang](https://go.dev/) (Gin Gonic, GORM)
- **AI Model:** [Google Gemini API](https://ai.google.dev/)
- **Database:**
  - **PostgreSQL:** Data relasional (User, Kelas, Rubrik).
  - **Elasticsearch:** Vector database untuk penyimpanan materi (RAG).
  - **Redis:** Caching lapisan kecepatan.
- **Infrastruktur:** Docker & Docker Compose.

---

## 📂 Struktur Proyek

```text
SAGE-System/
├── api-backend/       # Source code Golang (Logic & AI Integration)
├── web-frontend/      # Source code Next.js (Dashboard & UI)
├── docs/              # Dokumentasi Skripsi & Skema Database
├── docker-compose.yml # Konfigurasi Infrastruktur
└── README.md