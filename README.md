
![Skypiea AI logo](./web/static/img/skypiea-ai-logo.svg)
<br>

> **Skypiea AI** is a *cloud‑native reference application* that shows how to design, test, ship, and operate a Go service at production scale. Feel free to clone, reuse and feedback.

> It also wraps the Gemini API to expose a built‑in chatbot interface. You can test it, it is free tier though...
>
Built with ❤️ in San Francisco

---
## 🌐 Production URL

*Live demo:* [**https://chat.skypiea-ai.xyz**](https://chat.skypiea-ai.xyz)\
Hosted on an EC2 `t2.micro`, secured by Cloudflare, served through Traefik with default certificates.

---

## ✨ Why does this project exist?

Skypiea AI is **not** a traditional SaaS product that solves a single business problem. Instead, it is a living laboratory for modern backend engineering practices:

- **Scalable** – containerised, stateless by default, ready to autoscale in Kubernetes.
- **Testable** – unit, integration, and functional test layers.
- **Maintainable** – layered architecture, and uniform logging, managed database migrations on startup, run on local easily with little dependency.
- **Secure** – AuthN/Z implementation of JWT and session and similar sensitivity on prod; HTTPS and secrets stored encrypted on k8s/DB etc...
- **Automated** – full CI/CD pipeline that lints → tests → builds → pushes → deploys on every `git push`.
- **Portable** – runs the same on your laptop or in the cloud.

Overall, think of **Skypiea AI** as a live resume: a solid demonstration of senior-level backend practices.

---

## 🔧 Tech Stack

| Layer                  | Tools & Services                            |
| ---------------------- |---------------------------------------------|
| **Language & Runtime** | Go 1.24                                     |
| **Frontend**           | HTMX + Bootstrap + JS                       |
| **Database**           | PostgreSQL                                  |
| **Containerisation**   | Docker                                      |
| **Orchestration**      | Kubernetes (and Docker-Compose if you want) |
| **Packaging**          | Helm chart                                  |
| **Cloud**              | AWS, GHCR, Cloudflare, LetsEncrypt          |
| **CI/CD**              | GitHub Actions workflow chain               | 

---


## 🚀 Quick Start (Local)

```bash
# 1. Clone
$ git clone https://github.com/fukaraca/skypiea && cd skypiea

# 2. Spin up a DB and make migration
$ make migratedb-up

# 3. build images and run the server
$ make dc-build-up

# 4. visit http://localhost:8080 in your browser
```

---


## 📄 License

Skypiea AI is licensed under the **MIT License**. See [`LICENSE`](./LICENSE) for details.

---

## ✉️ Contact

Maintainer: [Linkedin](https://www.linkedin.com/in/fukaraca/)

---


