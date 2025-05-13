<p align="center">
  <img src="https://zop.dev/resources/cdn/newsletter/zopdev-transparent-logo.png" alt="zop.dev Logo" width="200">
</p>

<p align="center">ZopDev is a cloud orchestration platform that streamlines cloud management. It automates cloud infrastructure management by optimizing resource allocation, preventing downtime, streamlining deployments, and enabling seamless scaling across AWS, Azure and GCP.</p>

<p align="center">
  </a>
  <a href="./CONTRIBUTING.md">
    <img src="https://img.shields.io/badge/Contribute-Guide-orange?style=for-the-badge" alt="Contributing">
  </a>
</p>

## 💡 **Key Features**

1. **Cloud Agnostic**
2. **Automated Audits**

## 👍 **Contribute**

To contribute please refer to the [CONTRIBUTING.md](./CONTRIBUTING.md) file.

## Installation

### Prerequisites

- Docker installed on your system.
- Node.js version "^18.18.0 || ^19.8.0 || >= 20.0.0" is required.
- [Go](https://golang.org/) (v1.24 or later)

---

### Running Locally

#### zop-api

Run the following command to pull and start the Docker image for the zop-api:

```bash
    docker run -d -p 8000:8000 --name zop-api zopdev/api:v0.2.2
```

#### zop-ui

Run the following command to pull and start the Docker image for the zop-ui:

```bash
    docker run -d -p 3000:8000 -e API_BASE_URL='http://localhost:8000' --name zop-ui zopdev/dashboard:v0.2.2
```

> **Note:** The environment variable `API_BASE_URL` is used by zop-ui to connect to the
> zop-api. Ensure that the value matches the API's running base URL.
