# Zopdev Dashboard

This directory contains the frontend application for the Zopdev platform. Built using **React.js**
and **Vite**, it provides the user interface for interacting with the Zopdev API server.

## ⚙️ Tech Stack

- [React.js](https://reactjs.org/)
- [Vite](https://vitejs.dev/)
- [NPM](https://www.npmjs.com/)
- [Atomic Design Architecture](https://bradfrost.com/blog/post/atomic-web-design/)
- [Tailwind CSS](https://tailwindcss.com/)

## ✅ Requirements

Before starting, ensure you have the following installed on your system:

- [Node.js](https://nodejs.org/) (v18 or later recommended)
- [NPM](https://www.npmjs.com/) (comes with Node.js)
- [Git](https://git-scm.com/) (to clone the repository)

## 🧰 Prerequisites

- The backend Zopdev API server should be running locally at `http://localhost:8000` or configured
  accordingly in the `.env` file.
- Basic understanding of React and component-based design.

## 🧩 Project Structure

This project follows the **Atomic Design methodology**, structured as:

```
dashboard/
├── public/
├── src/
│   ├── assets/             # Static assets (images, icons, etc.)
│   ├── components/
│   │   ├── atoms/          # Reusable basic UI elements (buttons, inputs, etc.)
│   │   ├── molecules/      # Groups of atoms (form fields, card layouts)
│   │   ├── organisms/      # Complex UI components (modals, headers)
│   │   └── layouts/        # Page-level layout components
│   ├── pages/              # Route-based views
│   ├── hooks/              # Custom React hooks
│   ├── services/           # API request handlers
│   ├── queries/            # API queries are written here
│   ├── routes/             # Page routes paths are configured here
│   ├── utils/              # Utility functions
│   ├── App.jsx
│   └── main.jsx
├── .env
├── index.html
├── package.json
└── vite.config.js
```

## 🚀 Getting Started

### Installation

1. Navigate to the `dashboard` directory:

   ```bash
   cd dashboard
   ```

2. Install dependencies:

   ```bash
   npm install
   ```

3. Create a `.env` file for local environment variables:

   ```env
   VITE_API_BASE_URL=http://localhost:8000
   ```

### Running the Development Server

```bash
npm run dev
```

The application will start on [http://localhost:3000](http://localhost:3000) by default.

## 🐳 Docker (Optional)

If you'd like to run the frontend in a Docker container:

```bash
docker build -t zopdev/dashboard:v0.0.1 .
docker run -d -p 3000:3000 -e VITE_API_BASE_URL='http://localhost:8000' --name zop-ui zopdev/dashboard:v0.0.1
```

## 🤝 Contributing

Feel free to open issues or submit pull requests to improve the project. Follow the
[contribution guidelines](../CONTRIBUTING.md) for more details.
