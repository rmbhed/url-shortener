# rmbh.me — UI for shortlinks

This is a minimal static frontend for managing shortlinks. It expects a backend exposing two endpoints:

- `GET /api/links` — returns JSON array of `{ shortName, url }` objects
- `POST /api/links` — accepts JSON `{ shortName, url }`; returns `200` on success, `409` if short name already exists

How to run locally (simple):

1. Serve the `ui/` directory as static files. From the project root run:

```bash
cd ui
python -m http.server 8000
```

2. Open `http://localhost:8000` in your browser. The UI will call `/api/links` relative to the page origin.

Notes:
- If you don't have a backend yet, you can mock the endpoints with a tiny server (Express, Flask, etc.) or reverse-proxy requests to your production backend.
-- Shortlink creation expects the backend to return 409 when a short name is already taken.
